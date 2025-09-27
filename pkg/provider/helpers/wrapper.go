package helpers

import (
	"encoding/json"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"goauthentik.io/api/v3"
)

func SetWrapper(d *schema.ResourceData, key string, data interface{}) {
	err := d.Set(key, data)
	if err != nil {
		panic(err)
	}
}

// Get Pointer value from resource data, if castable to generic type
func GetP[T any](d *schema.ResourceData, key string) *T {
	rv, ok := d.GetOk(key)
	if !ok {
		return nil
	}
	if tt, ok := rv.(T); ok {
		return &tt
	}
	return nil
}

// Similar to `GetP` however also casts to an int32 as that is what the API prefers
func GetIntP(d *schema.ResourceData, key string) *int32 {
	rv, ok := d.GetOk(key)
	if !ok {
		return nil
	}
	if tt, ok := rv.(int); ok {
		return api.PtrInt32(int32(tt))
	}
	return nil
}

// Similar to `GetP` however also casts to an int64 as that is what the API prefers
func GetInt64P(d *schema.ResourceData, key string) *int64 {
	rv, ok := d.GetOk(key)
	if !ok {
		return nil
	}
	if tt, ok := rv.(int); ok {
		return api.PtrInt64(int64(tt))
	}
	return nil
}

func GetJSON[T any](d *schema.ResourceData, key string) (T, diag.Diagnostics) {
	var v T
	if sv := GetP[string](d, key); sv != nil {
		err := json.NewDecoder(strings.NewReader(*sv)).Decode(&v)
		if err != nil {
			return v, diag.FromErr(err)
		}
	}
	return v, nil
}

func CastSlice_New[T any](d *schema.ResourceData, key string) []T {
	sl := make([]T, 0)
	rv, ok := d.GetOk(key)
	if !ok {
		return sl
	}
	in, ok := rv.([]interface{})
	if !ok {
		return sl
	}
	for _, m := range in {
		sl = append(sl, m.(T))
	}
	return sl
}

func CastSlice[T any](in []interface{}) []T {
	sl := make([]T, len(in))
	for i, m := range in {
		sl[i] = m.(T)
	}
	return sl
}

func CastSliceInt32(in []interface{}) []int32 {
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
