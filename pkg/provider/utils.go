package provider

import (
	"encoding/json"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"goauthentik.io/api/v3"
)

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
