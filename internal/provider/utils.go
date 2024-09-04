package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func markDeprecated(resource func() *schema.Resource, newName string) func() *schema.Resource {
	return func() *schema.Resource {
		res := resource()
		res.DeprecationMessage = fmt.Sprintf("This resource is deprecated. Migrate to `%s`.", newName)
		res.Description += fmt.Sprintf("\n\n~> %s", res.DeprecationMessage)
		return res
	}
}

func StringInEnum[T ~string](items []T) schema.SchemaValidateDiagFunc {
	nv := make([]string, len(items))
	for i, v := range items {
		nv[i] = string(v)
	}
	return validation.ToDiagFunc(validation.StringInSlice(nv, false))
}

func EnumToDescription[T ~string](allowed []T) string {
	sb := strings.Builder{}
	sb.WriteString("Allowed values:\n")
	for _, v := range allowed {
		_, _ = sb.WriteString(fmt.Sprintf("  - `%s`\n", v))
	}
	return sb.String()
}

func setWrapper(d *schema.ResourceData, key string, data interface{}) {
	err := d.Set(key, data)
	if err != nil {
		panic(err)
	}
}

// diffSuppressExpression Diff suppression for python expressions
func diffSuppressExpression(k, old, new string, d *schema.ResourceData) bool {
	return strings.TrimSuffix(new, "\n") == old
}

// diffSuppressJSON Diff suppression for JSON objects
func diffSuppressJSON(k, old, new string, d *schema.ResourceData) bool {
	var j, j2 interface{}
	if err := json.Unmarshal([]byte(old), &j); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(new), &j2); err != nil {
		return false
	}
	return reflect.DeepEqual(j2, j)
}

// offsetInSlice Return the offset of a matching string in a slice or -1 if not found
func offsetInSlice[T comparable](s T, list []T) int {
	for offset, entry := range list {
		if entry == s {
			return offset
		}
	}
	return -1
}

// listConsistentMerge Consistent merge of TypeList elements, maintaining entries position within the list
// Workaround to TF Plugin SDK issue https://github.com/hashicorp/terraform-plugin-sdk/issues/477
// Taken from https://github.com/alexissavin/terraform-provider-solidserver/blob/master/solidserver/solidserver-helper.go#L62
func listConsistentMerge[T comparable](old []T, new []T) []interface{} {
	// Step 1 Build local list of member indexed by their offset
	oldOffset := make(map[int]T, len(old))
	diff := make([]T, 0, len(new))
	res := make([]interface{}, 0, len(new))

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
	for _, v := range diff {
		res = append(res, v)
	}

	return res
}

func castSlice[T any](in []interface{}) []T {
	sl := make([]T, len(in))
	for i, m := range in {
		sl[i] = m.(T)
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
