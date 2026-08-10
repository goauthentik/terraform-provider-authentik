package helpers

import (
	"context"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// offsetInSlice Return the offset of a matching string in a slice or -1 if not found
func offsetInSlice[T comparable](s T, list []T) int {
	for offset, entry := range list {
		if entry == s {
			return offset
		}
	}
	return -1
}

// ListConsistentMerge Consistent merge of TypeList elements, maintaining entries position within the list
// Workaround to TF Plugin SDK issue https://github.com/hashicorp/terraform-plugin-sdk/issues/477
// Taken from https://github.com/alexissavin/terraform-provider-solidserver/blob/master/solidserver/solidserver-helper.go#L62
func ListConsistentMerge[T comparable](old []T, new []T) []T {
	// Step 1 Build local list of member indexed by their offset
	oldOffset := make(map[int]T, len(old))
	diff := make([]T, 0, len(new))
	res := make([]T, 0, len(new))

	for _, n := range new {
		offset := offsetInSlice(n, old)

		if offset != -1 {
			oldOffset[offset] = n
		} else {
			diff = append(diff, n)
		}
	}

	// Merge sorted entries ordered by their offset with the diff array that contain the new ones
	// Step 2 Sort the index
	keys := make([]int, 0, len(old))
	for k := range oldOffset {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	// Step 3 build the result
	for _, k := range keys {
		res = append(res, oldOffset[k])
	}
	res = append(res, diff...)
	return res
}

func CastSlice[T any](d ResourceData, key string) []T {
	sl := make([]T, 0)
	rv, ok := d.GetOk(key)
	if !ok {
		return sl
	}
	in, ok := rv.([]any)
	if !ok {
		return sl
	}
	for _, m := range in {
		sl = append(sl, m.(T))
	}
	return sl
}

// Cast a slice of string-like objects to their respective types
func CastSliceString[T ~string](raw []string) []T {
	nl := make([]T, len(raw))
	for i, r := range raw {
		nl[i] = T(r)
	}
	return nl
}

func CastSliceInt32(in []any) []int32 {
	sl := make([]int32, len(in))
	for i, m := range in {
		sl[i] = int32(m.(int))
	}
	return sl
}

func Slice32ToInt(in []int32) []int {
	sl := make([]int, len(in))
	for i, m := range in {
		sl[i] = int(m)
	}
	return sl
}

// MergeList applies ListConsistentMerge to a framework types.List against a fresh
// slice of API values, staying null-aware per H1: if prior was null and the API
// returned nothing, the result stays null rather than becoming an empty list. The
// algorithm itself is unchanged from ListConsistentMerge; only the null handling and
// the types.List<->[]T conversion are new.
func MergeList[T comparable](ctx context.Context, prior types.List, elemType attr.Type, apiValues []T) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	if prior.IsNull() && len(apiValues) == 0 {
		return types.ListNull(elemType), diags
	}

	var priorValues []T
	if !prior.IsNull() {
		diags.Append(prior.ElementsAs(ctx, &priorValues, false)...)
		if diags.HasError() {
			return types.ListUnknown(elemType), diags
		}
	}

	merged := ListConsistentMerge(priorValues, apiValues)

	listValue, d := types.ListValueFrom(ctx, elemType, merged)
	diags.Append(d...)
	return listValue, diags
}

// SliceOrEmpty is the write-direction counterpart of MergeList, and the list analogue of
// StringPtrEmpty. It converts a types.List to a plain slice, returning a non-nil empty
// slice - never nil - for a null or unknown list.
//
// Both halves of that matter. Every generated API request model gates its list fields on
// `if !IsNil(o.X)`, so a nil slice is omitted from the request body entirely, and for an
// update the server then keeps its current value instead of clearing it. SDKv2's
// CastSlice always returned a non-nil slice, so plain `var x []T; list.ElementsAs(...)`
// is both a behaviour change (removing an optional list no longer clears it) and a hard
// failure: the plan says null, the API keeps the old contents, and the framework raises
// "Provider produced inconsistent result after apply".
//
// The unknown case is separate and louder. ElementsAs cannot write an unknown value into
// a []T target and returns an error diagnostic instead, which means an Optional+Computed
// list (authentik_group's users - unknown in the plan whenever it is unset) fails on
// Create. Treating unknown as "send nothing" is right for those: the whole point of
// Optional+Computed is that the server decides when config is silent.
func SliceOrEmpty[T any](ctx context.Context, l types.List) ([]T, diag.Diagnostics) {
	var diags diag.Diagnostics
	out := []T{}
	if l.IsNull() || l.IsUnknown() {
		return out, diags
	}
	diags.Append(l.ElementsAs(ctx, &out, false)...)
	if out == nil {
		out = []T{}
	}
	return out, diags
}

// MergeStringList is MergeList specialised to types.StringType, the common case for
// slug/PK reference lists (e.g. authentik_group's parents/roles).
func MergeStringList(ctx context.Context, prior types.List, apiValues []string) (types.List, diag.Diagnostics) {
	return MergeList(ctx, prior, types.StringType, apiValues)
}

// MergeInt32List is MergeList specialised to types.Int32Type, for int32 PK reference
// lists (e.g. authentik_group's users).
func MergeInt32List(ctx context.Context, prior types.List, apiValues []int32) (types.List, diag.Diagnostics) {
	return MergeList(ctx, prior, types.Int32Type, apiValues)
}
