package provider

import (
	"context"
	"testing"

	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

// TestRegistryExpressionAttributesUseExpressionType is a registry-wide sweep asserting that
// every attribute named "expression" declares CustomType: helpers.ExpressionType{}.
//
// It exists because a missing CustomType is the one mistake in this family that no other
// gate catches. ExpressionType serialises to a plain protocol `string`, so the docs-drift
// job renders the page identically either way, ValidateImplementation is happy, and the
// unit tests for a given resource still pass - the only symptom is a permanent diff for any
// user whose expression is a heredoc (which is the common case, since these are multi-line
// Python snippets). With 18 such attributes spread across 15 resources and 14 of them
// landing in a single batch, checking them one file at a time is exactly the sort of review
// that misses one.
//
// Sweeping the registry rather than listing files also means resources migrated later are
// covered automatically: the remaining expression attributes arrive with the policies batch
// (resource_policy_expression.go and friends), and this test will fail if any of them
// forgets the custom type.
func TestRegistryExpressionAttributesUseExpressionType(t *testing.T) {
	ctx := context.Background()
	p := New("testing", true)

	seen := 0
	for _, newResource := range p.Resources(ctx) {
		r := newResource()

		var metaResp fwresource.MetadataResponse
		r.Metadata(ctx, fwresource.MetadataRequest{ProviderTypeName: "authentik"}, &metaResp)

		var schemaResp fwresource.SchemaResponse
		r.Schema(ctx, fwresource.SchemaRequest{}, &schemaResp)
		require.Empty(t, schemaResp.Diagnostics)

		attr, ok := schemaResp.Schema.Attributes["expression"]
		if !ok {
			continue
		}
		seen++

		t.Run(metaResp.TypeName, func(t *testing.T) {
			stringAttr, ok := attr.(schema.StringAttribute)
			require.Truef(t, ok, "expression should be a StringAttribute, got %T", attr)
			assert.Equal(t, helpers.ExpressionType{}, stringAttr.CustomType,
				"expression must use helpers.ExpressionType or heredoc values show a permanent diff")
		})
	}

	// Guards against the sweep silently passing because it found nothing - e.g. if the
	// attribute were ever renamed, or the registry stopped being reachable from here.
	assert.GreaterOrEqual(t, seen, 14, "expected at least the 14 property mappings to carry an expression attribute")
}
