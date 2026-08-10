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

// Int32PtrOrNull is StringPtrOrNull's Int32 counterpart, for nullable-int API fields
// (e.g. api.NullableInt32.Get()). ValueInt32Pointer() is NOT a substitute for this in
// the write direction - see Int32Ptr below.
func Int32PtrOrNull(v *int32) types.Int32 {
	if v == nil {
		return types.Int32Null()
	}
	return types.Int32Value(*v)
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

// Int32Ptr is Int32PtrOrNull's write-direction counterpart: null config becomes a nil
// pointer, so a Nullable-typed API field built from it stays unset (omitted).
func Int32Ptr(v types.Int32) *int32 {
	if v.IsNull() {
		return nil
	}
	i := v.ValueInt32()
	return &i
}

// BoolPtr is the write-direction counterpart of BoolOrNull: null config becomes a nil
// pointer so the field is omitted, and an explicit false is sent as false.
//
// This is the fix for SDKv2's GetP[bool], which was built on d.GetOk and therefore
// returned nil for *any* false value - it could not tell "unset" from "explicitly false".
// Combined with omitempty on the request field and the API leaving absent fields unchanged
// on PUT, that meant a boolean could be turned on but never back off: the request omitted
// it, the API kept true, Read returned true, and the plan showed a permanent diff.
// authentik_policy_geoip's check_history_distance/check_impossible_travel are the case
// migration-plan.md flags as "where the null-semantics change is user-visible".
//
// Use this only for Optional booleans with no schema Default. With a Default the value is
// never null, so send it unconditionally with new(v.ValueBool()) instead.
func BoolPtr(v types.Bool) *bool {
	if v.IsNull() {
		return nil
	}
	b := v.ValueBool()
	return &b
}

// Int64Ptr is Int32Ptr for the int64 API fields.
func Int64Ptr(v types.Int64) *int64 {
	if v.IsNull() {
		return nil
	}
	i := v.ValueInt64()
	return &i
}

// Int64OrNull is Int32OrNull for the int64 API fields.
func Int64OrNull(prior types.Int64, v int64) types.Int64 {
	if v == 0 && prior.IsNull() {
		return types.Int64Null()
	}
	return types.Int64Value(v)
}
