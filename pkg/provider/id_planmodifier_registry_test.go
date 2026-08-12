package provider

import (
	"context"
	"testing"

	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// h3MutableIDResources are the resources whose id is derived from a *mutable* attribute, per
// H3 in migration-plan.md: nine take res.Slug and authentik_token takes res.Identifier.
// UseStateForUnknown() on their id is wrong - it would pin the old value in the plan while
// the API returns a new one, failing the apply the moment the slug or identifier changes.
//
// Every other resource keys its id on a stable server-generated UUID (res.Pk and friends)
// and should pin it, or users see "(known after apply)" on every plan where SDKv2 showed the
// prior value.
var h3MutableIDResources = map[string]bool{
	"authentik_application":     true,
	"authentik_flow":            true,
	"authentik_source_kerberos": true,
	"authentik_source_ldap":     true,
	"authentik_source_oauth":    true,
	"authentik_source_plex":     true,
	"authentik_source_saml":     true,
	"authentik_source_scim":     true,
	"authentik_source_telegram": true,
	"authentik_token":           true,
}

// TestRegistryIDPlanModifiers checks both directions of the H3 rule across the whole
// registry, and is the analogue of TestRegistryExpressionAttributesUseExpressionType.
//
// Neither mistake is visible to any other gate. UseStateForUnknown is not part of the
// protocol schema, so the docs-drift job and ValidateImplementation are both blind to it,
// and a resource's own unit tests pass either way - getting it wrong on an H3 resource only
// shows up as a failed apply when a user renames a slug, and getting it wrong on a stable-id
// resource only shows up as plan noise. As later batches add the remaining H3 resources
// (authentik_flow, authentik_token and the four other sources) this test covers them
// automatically.
func TestRegistryIDPlanModifiers(t *testing.T) {
	ctx := context.Background()
	p := New("testing", true)

	seenH3 := 0
	for _, newResource := range p.Resources(ctx) {
		r := newResource()

		var metaResp fwresource.MetadataResponse
		r.Metadata(ctx, fwresource.MetadataRequest{ProviderTypeName: "authentik"}, &metaResp)

		var schemaResp fwresource.SchemaResponse
		r.Schema(ctx, fwresource.SchemaRequest{}, &schemaResp)
		require.Empty(t, schemaResp.Diagnostics)

		attr, ok := schemaResp.Schema.Attributes["id"]
		require.Truef(t, ok, "%s must declare an explicit id: the framework does not inject one (H3)", metaResp.TypeName)

		stringAttr, ok := attr.(schema.StringAttribute)
		require.Truef(t, ok, "%s id should be a StringAttribute, got %T", metaResp.TypeName, attr)
		assert.Truef(t, stringAttr.Computed, "%s id must be Computed", metaResp.TypeName)

		if h3MutableIDResources[metaResp.TypeName] {
			seenH3++
			assert.Emptyf(t, stringAttr.PlanModifiers,
				"%s derives its id from a mutable attribute, so id must have no plan modifiers - "+
					"UseStateForUnknown would fail the apply when that attribute changes", metaResp.TypeName)
			continue
		}

		assert.NotEmptyf(t, stringAttr.PlanModifiers,
			"%s has a stable id, so it should pin it with UseStateForUnknown to avoid "+
				"\"(known after apply)\" on every plan", metaResp.TypeName)
	}

	// Guard against the H3 half passing vacuously.
	assert.GreaterOrEqual(t, seenH3, 4,
		"expected at least authentik_application plus the three sources migrated so far")
}
