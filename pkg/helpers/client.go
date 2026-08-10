package helpers

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"sync"

	"github.com/getsentry/sentry-go"
	api "goauthentik.io/api/v3"
)

// ClientOptions configures the shared authentik API client constructor. Two
// ClientOptions that compare equal (via reflect.DeepEqual) get the same *api.APIClient.
type ClientOptions struct {
	URL      string
	Token    string
	Insecure bool
	Headers  map[string]string
	Version  string
	Testing  bool
}

// APIURL normalizes a configured authentik URL to include the /api/v3 suffix,
// preserving any subpath the deployment is served under (e.g. subpath
// installs like https://host/sso/ become https://host/sso/api/v3).
func APIURL(raw string) (string, error) {
	akURL, err := url.Parse(raw)
	if err != nil {
		return "", err
	}
	path := akURL.Path
	if !strings.HasSuffix(path, "/api/v3") {
		path, err = url.JoinPath(path, "/api/v3")
		if err != nil {
			return "", err
		}
	}
	akURL.Path = path
	return akURL.String(), nil
}

var (
	clientCacheMu sync.Mutex
	clientCache   []clientCacheEntry
)

type clientCacheEntry struct {
	opts   ClientOptions
	client *api.APIClient
}

// NewAPIClient builds, or reuses a memoised, authentik API client for the given options.
//
// Memoisation matters, not just tidiness: terraform-plugin-mux invokes ConfigureProvider
// on every underlying server for a single plan/apply. Without memoisation, a muxed
// SDKv2+framework provider pair configured with the same options would each build their
// own client, doubling the RootConfigRetrieve round trip and re-running the global
// sentry.Init.
func NewAPIClient(ctx context.Context, opts ClientOptions) (*api.APIClient, error) {
	clientCacheMu.Lock()
	for _, e := range clientCache {
		if reflect.DeepEqual(e.opts, opts) {
			clientCacheMu.Unlock()
			return e.client, nil
		}
	}
	clientCacheMu.Unlock()

	serverURL, err := APIURL(opts.URL)
	if err != nil {
		return nil, err
	}

	config := api.NewConfiguration()
	config.Debug = true
	config.UserAgent = fmt.Sprintf("authentik-terraform@%s", opts.Version)
	config.Servers = api.ServerConfigurations{
		{
			URL:         serverURL,
			Description: "authentik API Server",
		},
	}
	config.HTTPClient = &http.Client{
		Transport: GetTLSTransport(opts.Insecure),
	}
	if opts.Testing {
		config.HTTPClient = &http.Client{
			Transport: NewTestingTransport(config.HTTPClient.Transport),
		}
	}

	config.AddDefaultHeader("Authorization", fmt.Sprintf("Bearer %s", opts.Token))
	for headerName, headerValue := range opts.Headers {
		config.AddDefaultHeader(headerName, headerValue)
	}
	apiClient := api.NewAPIClient(config)

	rootConfig, _, err := apiClient.RootAPI.RootConfigRetrieve(ctx).Execute()
	if err == nil && rootConfig.ErrorReporting.Enabled {
		dsn := ""
		// Customisable Sentry DSN was added in 2022.11, so only use that DSN when its set
		if rootConfig.ErrorReporting.SentryDsn != "" {
			dsn = rootConfig.ErrorReporting.SentryDsn
		}
		if envDsn, found := os.LookupEnv("SENTRY_DSN"); found {
			dsn = envDsn
		}
		if err := sentry.Init(sentry.ClientOptions{
			Dsn:              dsn,
			EnableTracing:    true,
			Environment:      rootConfig.ErrorReporting.Environment,
			TracesSampleRate: float64(rootConfig.ErrorReporting.TracesSampleRate),
			Release:          fmt.Sprintf("terraform-provider-authentik@%s", opts.Version),
		}); err != nil {
			fmt.Printf("Error during sentry init: %v\n", err)
		} else {
			config.HTTPClient.Transport = NewTracingTransport(context.Background(), config.HTTPClient.Transport)
			apiClient = api.NewAPIClient(config)
		}
	}

	clientCacheMu.Lock()
	clientCache = append(clientCache, clientCacheEntry{opts: opts, client: apiClient})
	clientCacheMu.Unlock()

	return apiClient, nil
}
