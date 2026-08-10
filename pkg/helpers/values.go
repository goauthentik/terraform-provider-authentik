package helpers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// This file is the core fix for H1 (see migration-plan.md): the authentik API returns
// "" / 0 / false / [] for attributes that were never set, but the framework's null and
// empty/zero are distinct values with no semantic-equality bridge between them. Reading
// an API response into state naively writes the zero value where the plan said null,
// which Terraform core rejects with a hard "Provider produced inconsistent result after
// apply" error - hit on every Create, since Create tail-calls Read.
//
// The fix is prior-value-aware mapping: an attribute only becomes null in the new state
// if it was already null before AND the API returned its zero value. Any resource that
// actually sends a zero value keeps seeing it; only "never touched" attributes stay null.

// StringOrNull returns types.StringNull() when v is the empty string and prior was
// already null, and types.StringValue(v) otherwise.
func StringOrNull(prior types.String, v string) types.String {
	if v == "" && prior.IsNull() {
		return types.StringNull()
	}
	return types.StringValue(v)
}

// StringPtrOrNull returns types.StringNull() when v is nil, and types.StringValue(*v)
// otherwise. Unlike StringOrNull, no prior value is needed: a nil pointer from a
// generated *api.Nullable-style field is already an unambiguous null.
func StringPtrOrNull(v *string) types.String {
	if v == nil {
		return types.StringNull()
	}
	return types.StringValue(*v)
}

// Int32OrNull returns types.Int32Null() when v is 0 and prior was already null, and
// types.Int32Value(v) otherwise.
func Int32OrNull(prior types.Int32, v int32) types.Int32 {
	if v == 0 && prior.IsNull() {
		return types.Int32Null()
	}
	return types.Int32Value(v)
}

// BoolOrNull returns types.BoolNull() when v is false and prior was already null, and
// types.BoolValue(v) otherwise.
func BoolOrNull(prior types.Bool, v bool) types.Bool {
	if !v && prior.IsNull() {
		return types.BoolNull()
	}
	return types.BoolValue(v)
}

// ListOrNull returns types.ListNull(elemType) when v is empty and prior was already
// null, and a types.List built from v otherwise.
func ListOrNull[T any](ctx context.Context, prior types.List, elemType attr.Type, v []T) (types.List, diag.Diagnostics) {
	if len(v) == 0 && prior.IsNull() {
		return types.ListNull(elemType), nil
	}
	return types.ListValueFrom(ctx, elemType, v)
}

// StringPtr is the write-direction counterpart of StringOrNull: null config becomes a
// nil pointer, so the field is omitted from the API request and the server-side value
// is left unchanged.
func StringPtr(v types.String) *string {
	if v.IsNull() {
		return nil
	}
	s := v.ValueString()
	return &s
}

// StringPtrEmpty is the write-direction counterpart of StringOrNull for attributes
// where clearing the config must clear the field server-side rather than leaving it
// unchanged (e.g. authentik_application's group/meta_icon/meta_launch_url - see #865).
// Unlike StringPtr, it never returns nil: a null config becomes a pointer to "", which
// the API interprets as "clear this field" rather than "omit this field".
func StringPtrEmpty(v types.String) *string {
	s := v.ValueString()
	return &s
}
