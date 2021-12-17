package provider

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/getsentry/sentry-go"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"goauthentik.io/api"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
		desc := s.Description
		if s.Default != nil {
			desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
		}
		return strings.TrimSpace(desc)
	}
}

// Provider -
func Provider(version string, testing bool) *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AUTHENTIK_URL", nil),
				Description: "The authentik API endpoint, can optionally be passed as `AUTHENTIK_URL` environmental variable",
			},
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AUTHENTIK_INSECURE", false),
				Description: "Whether to skip TLS verification, can optionally be passed as `AUTHENTIK_INSECURE` environmental variable",
			},
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AUTHENTIK_TOKEN", nil),
				Sensitive:   true,
				Description: "The authentik API token, can optionally be passed as `AUTHENTIK_TOKEN` environmental variable",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"authentik_application":                   resourceApplication(),
			"authentik_certificate_key_pair":          resourceCertificateKeyPair(),
			"authentik_flow_stage_binding":            resourceFlowStageBinding(),
			"authentik_flow":                          resourceFlow(),
			"authentik_group":                         resourceGroup(),
			"authentik_outpost":                       resourceOutpost(),
			"authentik_policy_binding":                resourcePolicyBinding(),
			"authentik_policy_dummy":                  resourcePolicyDummy(),
			"authentik_policy_event_matcher":          resourcePolicyEventMatcher(),
			"authentik_policy_expiry":                 resourcePolicyExpiry(),
			"authentik_policy_expression":             resourcePolicyExpression(),
			"authentik_policy_hibp":                   resourcePolicyHaveIBeenPwend(),
			"authentik_policy_password":               resourcePolicyPassword(),
			"authentik_policy_reputation":             resourcePolicyReputation(),
			"authentik_property_mapping_ldap":         resourceLDAPPropertyMapping(),
			"authentik_property_mapping_notification": resourceNotificationPropertyMapping(),
			"authentik_property_mapping_saml":         resourceSAMLPropertyMapping(),
			"authentik_provider_ldap":                 resourceProviderLDAP(),
			"authentik_provider_oauth2":               resourceProviderOAuth2(),
			"authentik_provider_proxy":                resourceProviderProxy(),
			"authentik_provider_saml":                 resourceProviderSAML(),
			"authentik_scope_mapping":                 resourceScopeMapping(),
			"authentik_service_connection_docker":     resourceServiceConnectionDocker(),
			"authentik_service_connection_kubernetes": resourceServiceConnectionKubernetes(),
			"authentik_source_ldap":                   resourceSourceLDAP(),
			"authentik_source_oauth":                  resourceSourceOAuth(),
			"authentik_source_plex":                   resourceSourcePlex(),
			"authentik_source_saml":                   resourceSourceSAML(),
			"authentik_stage_authenticator_duo":       resourceStageAuthenticatorDuo(),
			"authentik_stage_authenticator_sms":       resourceStageAuthenticatorSms(),
			"authentik_stage_authenticator_static":    resourceStageAuthenticatorStatic(),
			"authentik_stage_authenticator_totp":      resourceStageAuthenticatorTOTP(),
			"authentik_stage_authenticator_validate":  resourceStageAuthenticatorValidate(),
			"authentik_stage_authenticator_webauthn":  resourceStageAuthenticatorWebAuthn(),
			"authentik_stage_captcha":                 resourceStageCaptcha(),
			"authentik_stage_consent":                 resourceStageConsent(),
			"authentik_stage_deny":                    resourceStageDeny(),
			"authentik_stage_dummy":                   resourceStageDummy(),
			"authentik_stage_email":                   resourceStageEmail(),
			"authentik_stage_identification":          resourceStageIdentification(),
			"authentik_stage_invitation":              resourceStageInvitation(),
			"authentik_stage_password":                resourceStagePassword(),
			"authentik_stage_prompt_field":            resourceStagePromptField(),
			"authentik_stage_prompt":                  resourceStagePrompt(),
			"authentik_stage_user_delete":             resourceStageUserDelete(),
			"authentik_stage_user_login":              resourceStageUserLogin(),
			"authentik_stage_user_logout":             resourceStageUserLogout(),
			"authentik_stage_user_write":              resourceStageUserWrite(),
			"authentik_tenant":                        resourceTenant(),
			"authentik_user":                          resourceUser(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"authentik_certificate_key_pair":  dataSourceCertificateKeyPair(),
			"authentik_flow":                  dataSourceFlow(),
			"authentik_group":                 dataSourceGroup(),
			"authentik_property_mapping_ldap": dataSourceLDAPPropertyMapping(),
			"authentik_property_mapping_saml": dataSourceSAMLPropertyMapping(),
			"authentik_scope_mapping":         dataSourceScopeMapping(),
		},
		ConfigureContextFunc: providerConfigure(version, testing),
	}
}

// APIClient Hold the API Client and any relevant configuration
type APIClient struct {
	client *api.APIClient
}

func providerConfigure(version string, testing bool) schema.ConfigureContextFunc {
	return func(c context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		apiURL := d.Get("url").(string)
		token := d.Get("token").(string)
		insecure := d.Get("insecure").(bool)

		// Warning or errors can be collected in a slice type
		var diags diag.Diagnostics

		akURL, err := url.Parse(apiURL)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		config := api.NewConfiguration()
		config.Debug = true
		config.UserAgent = fmt.Sprintf("authentik-terraform@%s", version)
		config.Host = akURL.Host
		config.Scheme = akURL.Scheme
		if testing {
			config.HTTPClient = &http.Client{
				Transport: NewTestingTransport(GetTLSTransport(insecure)),
			}
		} else {
			config.HTTPClient = &http.Client{
				Transport: GetTLSTransport(insecure),
			}
		}

		config.AddDefaultHeader("Authorization", fmt.Sprintf("Bearer %s", token))
		apiClient := api.NewAPIClient(config)

		rootConfig, _, err := apiClient.RootApi.RootConfigRetrieve(context.Background()).Execute()
		if err == nil && rootConfig.ErrorReporting.Enabled {
			dsn := "https://7b485fd979bf48c1acbe38ffe382a541@sentry.beryju.org/14"
			if envDsn, found := os.LookupEnv("SENTRY_DSN"); !found {
				dsn = envDsn
			}
			sentry.Init(sentry.ClientOptions{
				Dsn:              dsn,
				Environment:      rootConfig.ErrorReporting.Environment,
				TracesSampleRate: float64(rootConfig.ErrorReporting.TracesSampleRate),
				Release:          fmt.Sprintf("authentik-terraform-provider@%s", version),
			})
			config.HTTPClient.Transport = NewTracingTransport(context.Background(), config.HTTPClient.Transport)
			apiClient = api.NewAPIClient(config)
		}

		return &APIClient{
			client: apiClient,
		}, diags
	}
}

// TestingTransport Transport used for testing, always returns a 400 Response
type TestingTransport struct {
	inner http.RoundTripper
}

// NewTestingTransport Get a HTTP Transport that fails all requests
func NewTestingTransport(inner http.RoundTripper) *TestingTransport {
	return &TestingTransport{inner}
}

// RoundTrip HTTP Transport
func (tt *TestingTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	body := "mock-failed-request"
	return &http.Response{
		Status:        "400 Bad Request",
		StatusCode:    400,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Body:          ioutil.NopCloser(bytes.NewBufferString(body)),
		ContentLength: int64(len(body)),
		Request:       r,
		Header:        make(http.Header),
	}, nil
}

// GetTLSTransport Get a TLS transport instance, that skips verification if configured via environment variables.
func GetTLSTransport(insecure bool) http.RoundTripper {
	tlsTransport, err := httptransport.TLSTransport(httptransport.TLSClientOptions{
		InsecureSkipVerify: insecure,
	})
	if err != nil {
		panic(err)
	}
	return tlsTransport
}

type tracingTransport struct {
	inner http.RoundTripper
	ctx   context.Context
}

func NewTracingTransport(ctx context.Context, inner http.RoundTripper) *tracingTransport {
	return &tracingTransport{inner, ctx}
}

func (tt *tracingTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	span := sentry.StartSpan(tt.ctx, "authentik.go.http_request")
	r.Header.Set("sentry-trace", span.ToSentryTrace())
	span.Description = fmt.Sprintf("%s %s", r.Method, r.URL.String())
	span.SetTag("url", r.URL.String())
	span.SetTag("method", r.Method)
	defer span.Finish()
	res, err := tt.inner.RoundTrip(r.WithContext(span.Context()))
	return res, err
}
