package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	api "goauthentik.io/api/v3"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

// authentik_stage_prompt_field is the first resource to migrate an attribute that carried
// DiffSuppressExpression (placeholder and initial_value; 16 more follow in the property
// mappings batch), so these tests cover both halves of what replaced it.

// TestStagePromptFieldExpression_SemanticEquality is the check migration-plan.md asks for
// on helpers.Expression: a heredoc value ending in "\n" against the API's stripped value
// must compare equal, so the plan is empty rather than showing a permanent diff. It runs
// the comparison in both directions, because the framework calls
// proposedNewValuable.StringSemanticEquals(ctx, priorValuable) and which side holds the
// trailing newline depends on whether the fresh value came from config or from the API.
func TestStagePromptFieldExpression_SemanticEquality(t *testing.T) {
	ctx := context.Background()

	heredoc := helpers.NewExpressionValue("return True\n")
	stripped := helpers.NewExpressionValue("return True")

	eq, diags := heredoc.StringSemanticEquals(ctx, stripped)
	require.False(t, diags.HasError(), diags)
	assert.True(t, eq, "config value with a trailing newline must equal the API's stripped value")

	eq, diags = stripped.StringSemanticEquals(ctx, heredoc)
	require.False(t, diags.HasError(), diags)
	assert.True(t, eq, "the comparison must be symmetric - SDKv2's DiffSuppressExpression trimmed only one side")

	// A real difference must still be reported as one.
	eq, diags = stripped.StringSemanticEquals(ctx, helpers.NewExpressionValue("return False"))
	require.False(t, diags.HasError(), diags)
	assert.False(t, eq)
}

// TestStagePromptFieldFromAPI_ExpressionNullVsEmpty is H1 for the custom type. Semantic
// equality cannot cover this case - value_semantic_equality.go returns early when the prior
// value is null - so ExpressionOrNull is still needed on top of ExpressionType.
func TestStagePromptFieldFromAPI_ExpressionNullVsEmpty(t *testing.T) {
	r := &stagePromptFieldResource{}

	data := &stagePromptFieldModel{
		Placeholder:  helpers.NewExpressionNull(),
		InitialValue: helpers.NewExpressionNull(),
	}

	empty := ""
	res := &api.Prompt{
		Pk:           "uuid-1",
		Name:         "field",
		FieldKey:     "key",
		Label:        "Label",
		Type:         api.PROMPTTYPEENUM_TEXT,
		Placeholder:  &empty,
		InitialValue: &empty,
	}

	r.fromAPI(data, res)

	assert.True(t, data.Placeholder.IsNull(), "placeholder was null and the API returned \"\", so it stays null")
	assert.True(t, data.InitialValue.IsNull(), "initial_value was null and the API returned \"\", so it stays null")
	assert.True(t, data.Order.IsNull(), "order has no Default and the API returned 0, so it stays null")

	// The four defaulted attributes must come back concrete even though every API value
	// here is the zero value - discovery #5.
	require.False(t, data.Required.IsNull())
	require.False(t, data.PlaceholderExpression.IsNull())
	require.False(t, data.InitialValueExpression.IsNull())
	require.False(t, data.SubText.IsNull(), "sub_text defaults to \"\", which is exactly the value the API returns")
	assert.Equal(t, "", data.SubText.ValueString())
}

// TestStagePromptFieldToRequest_ExpressionOmittedWhenNull pins the write direction:
// placeholder/initial_value used GetP in SDKv2, so a null config must omit the field rather
// than send "".
func TestStagePromptFieldToRequest_ExpressionOmittedWhenNull(t *testing.T) {
	r := &stagePromptFieldResource{}

	data := &stagePromptFieldModel{
		Name:                   types.StringValue("field"),
		FieldKey:               types.StringValue("key"),
		Label:                  types.StringValue("Label"),
		Type:                   types.StringValue(string(api.PROMPTTYPEENUM_TEXT)),
		Required:               types.BoolValue(false),
		Placeholder:            helpers.NewExpressionNull(),
		InitialValue:           helpers.NewExpressionNull(),
		PlaceholderExpression:  types.BoolValue(false),
		InitialValueExpression: types.BoolValue(false),
		Order:                  types.Int32Null(),
		SubText:                types.StringValue(""),
	}

	body := r.toRequest(data)

	assert.Nil(t, body.Placeholder, "a null placeholder must be omitted, not sent as \"\"")
	assert.Nil(t, body.InitialValue, "a null initial_value must be omitted, not sent as \"\"")
	assert.Nil(t, body.Order, "a null order must be omitted")

	// sub_text has a Default, so unlike the two expressions it is always sent.
	require.NotNil(t, body.SubText)
	assert.Equal(t, "", *body.SubText)

	// The trailing newline must survive the write direction untouched: semantic equality is
	// a state-comparison mechanism, not a normalisation step, so what the user wrote is
	// what gets sent.
	data.Placeholder = helpers.NewExpressionValue("return True\n")
	body = r.toRequest(data)
	require.NotNil(t, body.Placeholder)
	assert.Equal(t, "return True\n", *body.Placeholder)
}
