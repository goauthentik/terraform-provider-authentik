package provider

import (
	"io/ioutil"
	"log"
	"net/http"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

// stringOffsetInSlice Return the offset of a matching string in a slice or -1 if not found
func stringOffsetInSlice(s string, list []string) int {
	for offset, entry := range list {
		if entry == s {
			return offset
		}
	}
	return -1
}

// typeListConsistentMerge Consistent merge of TypeList elements, maintaining entries position within the list
// Workaround to TF Plugin SDK issue https://github.com/hashicorp/terraform-plugin-sdk/issues/477
// Taken from https://github.com/alexissavin/terraform-provider-solidserver/blob/master/solidserver/solidserver-helper.go#L62
func typeListConsistentMerge(old []string, new []string) []interface{} {
	// Step 1 Build local list of member indexed by their offset
	old_offsets := make(map[int]string, len(old))
	diff := make([]string, 0, len(new))
	res := make([]interface{}, 0, len(new))

	for _, n := range new {
		if n != "" {
			offset := stringOffsetInSlice(n, old)

			if offset != -1 {
				old_offsets[offset] = n
			} else {
				diff = append(diff, n)
			}
		}
	}

	// Merge sorted entries ordered by their offset with the diff array that contain the new ones
	// Step 2 Sort the index
	keys := make([]int, 0, len(old))
	for k := range old_offsets {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	// Step 3 build the result
	for _, k := range keys {
		res = append(res, old_offsets[k])
	}
	for _, v := range diff {
		res = append(res, v)
	}

	return res
}

func sliceToString(in []interface{}) []string {
	sl := make([]string, len(in))
	for i, m := range in {
		sl[i] = m.(string)
	}
	return sl
}

func stringToPointer(in string) *string {
	return &in
}

func stringPointerResolve(in *string) string {
	return *in
}

func intToPointer(in int) *int32 {
	i := int32(in)
	return &i
}

func boolToPointer(in bool) *bool {
	return &in
}

func httpToDiag(r *http.Response, err error) diag.Diagnostics {
	b, er := ioutil.ReadAll(r.Body)
	if er != nil {
		log.Printf("[DEBUG] authentik: failed to read response: %s", er.Error())
		b = []byte{}
	}
	log.Printf("[DEBUG] authentik: error response: %s", string(b))
	return diag.Errorf("HTTP Error '%s' during request '%s %s': \"%s\"", err.Error(), r.Request.Method, r.Request.URL.Path, string(b))
}
