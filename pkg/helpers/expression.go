package helpers

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	_ basetypes.StringTypable                    = ExpressionType{}
	_ basetypes.StringValuable                   = ExpressionValue{}
	_ basetypes.StringValuableWithSemanticEquals = ExpressionValue{}
)

// ExpressionType is the attribute type for authentik expression policy/property-mapping
// source, replacing the SDKv2 DiffSuppressExpression DiffSuppressFunc on 18 attributes.
// authentik always returns expressions with trailing newlines stripped, so a heredoc
// config value ending in "\n" would otherwise show a permanent diff against the API's
// response.
type ExpressionType struct {
	basetypes.StringType
}

func (t ExpressionType) String() string {
	return "ExpressionType"
}

func (t ExpressionType) ValueType(_ context.Context) attr.Value {
	return ExpressionValue{}
}

func (t ExpressionType) Equal(o attr.Type) bool {
	other, ok := o.(ExpressionType)
	if !ok {
		return false
	}
	return t.StringType.Equal(other.StringType)
}

func (t ExpressionType) ValueFromString(_ context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	return ExpressionValue{StringValue: in}, nil
}

func (t ExpressionType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.StringType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err
	}

	stringValue, ok := attrValue.(basetypes.StringValue)
	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	stringValuable, diags := t.ValueFromString(ctx, stringValue)
	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting StringValue to StringValuable: %v", diags)
	}

	return stringValuable, nil
}

// ExpressionValue is the Value type for ExpressionType.
type ExpressionValue struct {
	basetypes.StringValue
}

func (v ExpressionValue) Type(_ context.Context) attr.Type {
	return ExpressionType{}
}

func (v ExpressionValue) Equal(o attr.Value) bool {
	other, ok := o.(ExpressionValue)
	if !ok {
		return false
	}
	return v.StringValue.Equal(other.StringValue)
}

// StringSemanticEquals compares both sides with trailing newlines stripped. This must
// be symmetric: the SDKv2 DiffSuppressExpression this replaces trimmed only the config
// side (`strings.TrimSuffix(new, "\n") == old`), which happened to work there because
// SDKv2 always compares against the API's already-stripped value. The framework instead
// calls proposedNewValuable.StringSemanticEquals(ctx, priorValuable) - i.e. either side
// can be the freshly-read API value or the heredoc-with-newline config value depending
// on direction, so trimming only one side would silently stop suppressing the diff.
func (v ExpressionValue) StringSemanticEquals(_ context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(ExpressionValue)
	if !ok {
		diags.AddError(
			"Semantic Equality Check Error",
			"An unexpected value type was received while performing semantic equality checks. "+
				"Please report this to the provider developers.\n\n"+
				"Expected Value Type: "+fmt.Sprintf("%T", v)+"\n"+
				"Got Value Type: "+fmt.Sprintf("%T", newValuable),
		)
		return false, diags
	}

	return strings.TrimRight(v.ValueString(), "\n") == strings.TrimRight(newValue.ValueString(), "\n"), diags
}
