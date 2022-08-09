package provider

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func setWrapper(d *schema.ResourceData, key string, data interface{}) {
	err := d.Set(key, data)
	if err != nil {
		panic(err)
	}
}

func diffSuppressExpression(k, old, new string, d *schema.ResourceData) bool {
	return strings.TrimSuffix(new, "\n") == old
}

// stringOffsetInSlice Return the offset of a matching string in a slice or -1 if not found
func stringOffsetInSlice(s string, list []string) int {
	for offset, entry := range list {
		if entry == s {
			return offset
		}
	}
	return -1
}

// stringListConsistentMerge Consistent merge of TypeList elements, maintaining entries position within the list
// Workaround to TF Plugin SDK issue https://github.com/hashicorp/terraform-plugin-sdk/issues/477
// Taken from https://github.com/alexissavin/terraform-provider-solidserver/blob/master/solidserver/solidserver-helper.go#L62
func stringListConsistentMerge(old []string, new []string) []interface{} {
	// Step 1 Build local list of member indexed by their offset
	oldOffset := make(map[int]string, len(old))
	diff := make([]string, 0, len(new))
	res := make([]interface{}, 0, len(new))

	for _, n := range new {
		if n != "" {
			offset := stringOffsetInSlice(n, old)

			if offset != -1 {
				oldOffset[offset] = n
			} else {
				diff = append(diff, n)
			}
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
	for _, v := range diff {
		res = append(res, v)
	}

	return res
}

// intOffsetInSlice Return the offset of a matching string in a slice or -1 if not found
func intOffsetInSlice(s int, list []int) int {
	for offset, entry := range list {
		if entry == s {
			return offset
		}
	}
	return -1
}

// intListConsistentMerge Consistent merge of TypeList elements, maintaining entries position within the list
// Workaround to TF Plugin SDK issue https://github.com/hashicorp/terraform-plugin-sdk/issues/477
// Taken from https://github.com/alexissavin/terraform-provider-solidserver/blob/master/solidserver/solidserver-helper.go#L62
func intListConsistentMerge(old []int, new []int) []interface{} {
	// Step 1 Build local list of member indexed by their offset
	oldOffset := make(map[int]int, len(old))
	diff := make([]int, 0, len(new))
	res := make([]interface{}, 0, len(new))

	for _, n := range new {
		offset := intOffsetInSlice(n, old)

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

func sliceToInt(in []interface{}) []int {
	sl := make([]int, len(in))
	for i, m := range in {
		sl[i] = m.(int)
	}
	return sl
}

func slice32ToInt(in []int32) []int {
	sl := make([]int, len(in))
	for i, m := range in {
		sl[i] = int(m)
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

func httpToDiag(d *schema.ResourceData, r *http.Response, err error) diag.Diagnostics {
	if r == nil {
		return diag.Errorf("HTTP Error '%s' without http response", err.Error())
	}
	if r.StatusCode == 404 {
		d.SetId("")
		return diag.Diagnostics{}
	}
	buff := &bytes.Buffer{}
	_, er := io.Copy(buff, r.Body)
	if er != nil {
		log.Printf("[DEBUG] authentik: failed to read response: %s", er.Error())
	}
	log.Printf("[DEBUG] authentik: error response: %s", buff.String())
	return diag.Errorf("HTTP Error '%s' during request '%s %s': \"%s\"", err.Error(), r.Request.Method, r.Request.URL.Path, buff.String())
}
