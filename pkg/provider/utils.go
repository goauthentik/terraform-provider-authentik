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

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"goauthentik.io/api/v3"
)

const (
	RelativeDurationDescription = "Format: hours=1;minutes=2;seconds=3."
	JSONDescription             = "JSON format expected. Use `jsonencode()` to pass objects."
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

func ValidateRelativeDuration(i interface{}, p cty.Path) diag.Diagnostics {
	validKV := []string{
		"microseconds",
		"milliseconds",
		"seconds",
		"minutes",
		"hours",
		"days",
		"weeks",
	}
	return validation.ToDiagFunc(func(i interface{}, s string) (warnings []string, errors []error) {
		v, ok := i.(string)
		if !ok {
			errors = append(errors, fmt.Errorf("expected type of %s to be string", s))
			return warnings, errors
		}
		for _, el := range strings.Split(v, ";") {
			p := strings.Split(el, "=")
			if len(p) < 2 {
				errors = append(errors, fmt.Errorf("%s has incorrect amount of elements", el))
				return warnings, errors
			}
			isValid := false
			for _, valid := range validKV {
				if strings.EqualFold(p[0], valid) {
					isValid = true
				}
			}
			if !isValid {
				errors = append(errors, fmt.Errorf("%s has incorrect key %s", el, p[0]))
			}
		}
		return warnings, errors
	})(i, p)
}

func ValidateJSON(i interface{}, p cty.Path) diag.Diagnostics {
	return validation.ToDiagFunc(func(i interface{}, s string) (warnings []string, errors []error) {
		v, ok := i.(string)
		if !ok {
			errors = append(errors, fmt.Errorf("expected type of %s to be string", s))
			return warnings, errors
		}
		var j interface{}
		err := json.Unmarshal([]byte(v), &j)
		if err != nil {
			errors = append(errors, err)
			return warnings, errors
		}
		return warnings, errors
	})(i, p)
}

func setWrapper(d *schema.ResourceData, key string, data interface{}) {
	err := d.Set(key, data)
	if err != nil {
		panic(err)
	}
}

// Get Pointer value from resource data, if castable to generic type
func getP[T any](d *schema.ResourceData, key string) *T {
	rv, ok := d.GetOk(key)
	if !ok {
		return nil
	}
	if tt, ok := rv.(T); ok {
		return &tt
	}
	return nil
}

// Similar to `getP` however also casts to an int32 as that is what the API prefers
func getIntP(d *schema.ResourceData, key string) *int32 {
	rv, ok := d.GetOk(key)
	if !ok {
		return nil
	}
	if tt, ok := rv.(int); ok {
		return api.PtrInt32(int32(tt))
	}
	return nil
}

// Similar to `getP` however also casts to an int64 as that is what the API prefers
func getInt64P(d *schema.ResourceData, key string) *int64 {
	rv, ok := d.GetOk(key)
	if !ok {
		return nil
	}
	if tt, ok := rv.(int); ok {
		return api.PtrInt64(int64(tt))
	}
	return nil
}

func getJSON[T any](d *schema.ResourceData, key string) (T, diag.Diagnostics) {
	var v T
	if sv := getP[string](d, key); sv != nil {
		err := json.NewDecoder(strings.NewReader(*sv)).Decode(&v)
		if err != nil {
			return v, diag.FromErr(err)
		}
	}
	return v, nil
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

func castSliceInt32(in []interface{}) []int32 {
	sl := make([]int32, len(in))
	for i, m := range in {
		sl[i] = int32(m.(int))
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
