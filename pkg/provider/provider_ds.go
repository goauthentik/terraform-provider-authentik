package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"goauthentik.io/terraform-provider-authentik/pkg/provider/datasources/brand"
	"goauthentik.io/terraform-provider-authentik/pkg/provider/datasources/certificatekeypair"
	"goauthentik.io/terraform-provider-authentik/pkg/provider/datasources/flow"
)

func GetDatasources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"authentik_brand":                            td(brand.DataSource),
		"authentik_certificate_key_pair":             td(certificatekeypair.DataSource),
		"authentik_flow":                             td(flow.DataSource),
		"authentik_group":                            td(dataSourceGroup),
		"authentik_groups":                           td(dataSourceGroups),
		"authentik_outpost":                          td(dataSourceOutpost),
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
		"authentik_user":                             td(dataSourceUser),
		"authentik_users":                            td(dataSourceUsers),
		"authentik_webauthn_device_type":             td(dataSourceWebAuthnDeviceType),
	}
}
