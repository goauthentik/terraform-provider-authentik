package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"goauthentik.io/terraform-provider-authentik/pkg/provider/datasources/brands"
	"goauthentik.io/terraform-provider-authentik/pkg/provider/datasources/certificatekeypairs"
	"goauthentik.io/terraform-provider-authentik/pkg/provider/datasources/flows"
	"goauthentik.io/terraform-provider-authentik/pkg/provider/datasources/groups"
	"goauthentik.io/terraform-provider-authentik/pkg/provider/datasources/outposts"
	"goauthentik.io/terraform-provider-authentik/pkg/provider/datasources/propertymappings"
	"goauthentik.io/terraform-provider-authentik/pkg/provider/datasources/providers"
	"goauthentik.io/terraform-provider-authentik/pkg/provider/datasources/rbac"
	"goauthentik.io/terraform-provider-authentik/pkg/provider/datasources/sources"
	"goauthentik.io/terraform-provider-authentik/pkg/provider/datasources/users"
)

func GetDatasources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"authentik_brand":                            td(brands.DataSource),
		"authentik_certificate_key_pair":             td(certificatekeypairs.DataSource),
		"authentik_flow":                             td(flows.DataSourceFlow),
		"authentik_group":                            td(groups.DataSourceGroup),
		"authentik_groups":                           td(groups.DataSourceGroups),
		"authentik_outpost":                          td(outposts.DataSourceOutpost),
		"authentik_property_mapping_provider_rac":    td(propertymappings.DataSourcePropertyMappingProviderRAC),
		"authentik_property_mapping_provider_radius": td(propertymappings.DataSourcePropertyMappingProviderRadius),
		"authentik_property_mapping_provider_saml":   td(propertymappings.DataSourcePropertyMappingProviderSAML),
		"authentik_property_mapping_provider_scim":   td(propertymappings.DataSourcePropertyMappingProviderSCIM),
		"authentik_property_mapping_provider_scope":  td(propertymappings.DataSourcePropertyMappingProviderScope),
		"authentik_property_mapping_source_ldap":     td(propertymappings.DataSourcePropertyMappingSourceLDAP),
		"authentik_provider_oauth2_config":           td(providers.DataSourceProviderOAuth2Config),
		"authentik_provider_saml_metadata":           td(providers.DataSourceProviderSAMLMetadata),
		"authentik_rbac_permission":                  td(rbac.DataSourceRBACPermission),
		"authentik_service_connection_kubernetes":    td(outposts.DataOutpostServiceConnectionsKubernetes),
		"authentik_source":                           td(sources.DataSourceSource),
		"authentik_stage":                            td(flows.DataSourceStage),
		"authentik_user":                             td(users.DataSourceUser),
		"authentik_users":                            td(users.DataSourceUsers),
		"authentik_webauthn_device_type":             td(dataSourceWebAuthnDeviceType),
	}
}
