package provider

import (
	"context"
	"fmt"
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
func Provider(version string) *schema.Provider {
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
			"authentik_provider_oauth2":               resourceProviderOAuth2(),
			"authentik_provider_proxy":                resourceProviderProxy(),
			"authentik_service_connection_docker":     resourceServiceConnectionDocker(),
			"authentik_service_connection_kubernetes": resourceServiceConnectionKubernetes(),
			"authentik_stage_authenticator_duo":       resourceStageAuthenticatorDuo(),
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
		},
		DataSourcesMap: map[string]*schema.Resource{
			"authentik_flow":  dataSourceFlow(),
			"authentik_group": dataSourceGroup(),
		},
		ConfigureContextFunc: providerConfigure(version),
	}
}

// APIClient Hold the API Client and any relevant configuration
type APIClient struct {
	client *api.APIClient
}

func providerConfigure(version string) schema.ConfigureContextFunc {
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
		config.HTTPClient = &http.Client{
			Transport: GetTLSTransport(insecure),
		}
		config.AddDefaultHeader("Authorization", fmt.Sprintf("Bearer %s", token))
		apiClient := api.NewAPIClient(config)

		return &APIClient{
			client: apiClient,
		}, diags
	}
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
