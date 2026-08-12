package sdkprovider

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
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
		if s.Computed {
			desc += " Generated."
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
				Optional:    true,
				Description: "The authentik API endpoint, can optionally be passed as `AUTHENTIK_URL` environmental variable",
			},
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to skip TLS verification, can optionally be passed as `AUTHENTIK_INSECURE` environmental variable",
			},
			"token": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The authentik API token, can optionally be passed as `AUTHENTIK_TOKEN` environmental variable",
			},
			"headers": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Sensitive:   true,
				Description: "Optional HTTP headers sent with every request",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"authentik_application_entitlement":                    tr(resourceApplicationEntitlement),
			"authentik_blueprint":                                  tr(resourceBlueprintInstance),
			"authentik_brand":                                      tr(resourceBrand),
			"authentik_certificate_key_pair":                       tr(resourceCertificateKeyPair),
			"authentik_endpoints_connector_agent":                  tr(resourceEndpointsConnectorAgent),
			"authentik_endpoints_connector_agent_enrollment_token": tr(resourceEndpointsEnrollmentToken),
			"authentik_endpoints_device_access_group":              tr(resourceEndpointsDeviceAccessGroup),
			"authentik_endpoints_google_chrome_connector":          tr(resourceEndpointsGoogleChromeConnector),
			"authentik_enterprise_license":                         tr(resourceEnterpriseLicense),
			"authentik_event_rule":                                 tr(resourceEventRule),
			"authentik_event_transport":                            tr(resourceEventTransport),
			"authentik_flow_stage_binding":                         tr(resourceFlowStageBinding),
			"authentik_flow":                                       tr(resourceFlow),
			"authentik_outpost":                                    tr(resourceOutpost),
			"authentik_outpost_provider_attachment":                tr(resourceOutpostProviderAttachment),
			"authentik_provider_proxy":                             tr(resourceProviderProxy),
			"authentik_provider_saml":                              tr(resourceProviderSAML),
			"authentik_rac_endpoint":                               tr(resourceRACEndpoint),
			"authentik_rbac_initial_permissions":                   tr(resourceRBACInitialPermissions),
			"authentik_rbac_permission_role":                       tr(resourceRBACRoleObjectPermission),
			// TODO: Remove in 2026.2 or later
			"authentik_rbac_permission_user":          tr(helpers.MarkDeprecated(resourceRBACUserObjectPermission, "authentik_rbac_permission_role")),
			"authentik_rbac_role":                     tr(resourceRBACRole),
			"authentik_service_connection_docker":     tr(resourceServiceConnectionDocker),
			"authentik_service_connection_kubernetes": tr(resourceServiceConnectionKubernetes),
			"authentik_system_settings":               tr(resourceSystemSettings),
			"authentik_task_schedule":                 tr(resourceTaskSchedule),
			"authentik_token":                         tr(resourceToken),
			"authentik_user":                          tr(resourceUser),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"authentik_application_entitlement":          td(dataSourceApplicationEntitlement),
			"authentik_brand":                            td(dataSourceBrand),
			"authentik_certificate_key_pair":             td(dataSourceCertificateKeyPair),
			"authentik_flow":                             td(dataSourceFlow),
			"authentik_groups":                           td(dataSourceGroups),
			"authentik_outpost":                          td(dataSourceOutpost),
			"authentik_policy_binding":                   td(dataSourcePolicyBinding),
			"authentik_policy_expression":                td(dataSourcePolicyExpression),
			"authentik_property_mapping_provider_rac":    td(dataSourcePropertyMappingProviderRAC),
			"authentik_property_mapping_provider_radius": td(dataSourcePropertyMappingProviderRadius),
			"authentik_property_mapping_provider_saml":   td(dataSourcePropertyMappingProviderSAML),
			"authentik_property_mapping_provider_scim":   td(dataSourcePropertyMappingProviderSCIM),
			"authentik_property_mapping_provider_scope":  td(dataSourcePropertyMappingProviderScope),
			"authentik_property_mapping_source_ldap":     td(dataSourcePropertyMappingSourceLDAP),
			"authentik_provider_oauth2_config":           td(dataSourceProviderOAuth2Config),
			"authentik_provider_saml_metadata":           td(dataSourceProviderSAMLMetadata),
			"authentik_rbac_permission":                  td(dataSourceRBACPermission),
			"authentik_service_connection_kubernetes":    td(dataOutpostServiceConnectionsKubernetes),
			"authentik_source":                           td(dataSourceSource),
			"authentik_stage":                            td(dataSourceStage),
			"authentik_stage_prompt_field":               td(dataSourceStagePromptField),
			"authentik_user":                             td(dataSourceUser),
			"authentik_users":                            td(dataSourceUsers),
			"authentik_webauthn_device_type":             td(dataSourceWebAuthnDeviceType),
		},
		ConfigureContextFunc: providerConfigure(version, testing),
	}
}

// APIClient Hold the API Client and any relevant configuration
type APIClient struct {
	client *api.APIClient
}

// Client exposes the underlying *api.APIClient. The field itself stays unexported so
// in-package resource files keep using the terser c.client, but cross-package callers
// (pkg/acctest's CheckDestroy helper) need a way in.
func (a *APIClient) Client() *api.APIClient {
	return a.client
}

func providerConfigure(version string, testing bool) schema.ConfigureContextFunc {
	return func(c context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
		apiURL := d.Get("url").(string)
		if apiURL == "" {
			apiURL = os.Getenv("AUTHENTIK_URL")
		}
		token := d.Get("token").(string)
		if token == "" {
			token = os.Getenv("AUTHENTIK_TOKEN")
		}
		insecure := d.Get("insecure").(bool)
		if !insecure {
			insecure, _ = strconv.ParseBool(os.Getenv("AUTHENTIK_INSECURE"))
		}

		if apiURL == "" {
			return nil, diag.Errorf("no authentik URL configured, set `url` or the AUTHENTIK_URL environment variable")
		}
		if token == "" {
			return nil, diag.Errorf("no authentik token configured, set `token` or the AUTHENTIK_TOKEN environment variable")
		}

		headers := map[string]string{}
		if _headers, ok := d.GetOk("headers"); ok {
			for headerName, headerValue := range _headers.(map[string]any) {
				headers[headerName] = headerValue.(string)
			}
		}

		apiClient, err := helpers.NewAPIClient(c, helpers.ClientOptions{
			URL:      apiURL,
			Token:    token,
			Insecure: insecure,
			Headers:  headers,
			Version:  version,
			Testing:  testing,
		})
		if err != nil {
			return nil, diag.FromErr(err)
		}

		return &APIClient{
			client: apiClient,
		}, nil
	}
}
