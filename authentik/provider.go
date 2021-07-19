package authentik

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/goauthentik/terraform-provider-authentik/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AUTHENTIK_URL", nil),
			},
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AUTHENTIK_INSECURE", false),
			},
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AUTHENTIK_TOKEN", nil),
				Sensitive:   true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"authentik_application":                   resourceApplication(),
			"authentik_certificate_key_pair":          resourceCertificateKeyPair(),
			"authentik_outpost":                       resourceOutpost(),
			"authentik_policy_binding":                resourcePolicyBinding(),
			"authentik_provider_proxy":                resourceProviderProxy(),
			"authentik_service_connection_docker":     resourceServiceConnectionDocker(),
			"authentik_service_connection_kubernetes": resourceServiceConnectionKubernetes(),
			"authentik_tenant":                        resourceTenant(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"authentik_flow":  dataSourceFlow(),
			"authentik_group": dataSourceGroup(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

type ProviderAPIClient struct {
	client *api.APIClient
}

func httpToDiag(r *http.Response) diag.Diagnostics {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("[DEBUG] authentik: failed to read response: %s", err)
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] authentik: error response: %s", string(b))
	return diag.FromErr(errors.New(string(b)))
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	api_url := d.Get("url").(string)
	token := d.Get("token").(string)
	insecure := d.Get("insecure").(bool)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	akURL, err := url.Parse(api_url)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	config := api.NewConfiguration()
	config.Debug = true
	// TODO versioning
	config.UserAgent = fmt.Sprintf("authentik-terraform@%s", "test")
	config.Host = akURL.Host
	config.Scheme = akURL.Scheme
	config.HTTPClient = &http.Client{
		Transport: GetTLSTransport(insecure),
	}
	config.AddDefaultHeader("Authorization", fmt.Sprintf("Bearer %s", token))
	apiClient := api.NewAPIClient(config)

	return &ProviderAPIClient{
		client: apiClient,
	}, diags
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
