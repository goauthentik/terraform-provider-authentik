package helpers

import (
	"encoding/json"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

type ResourceData interface {
	Set(key string, value any) error
	GetOk(key string) (any, bool)
}

func SetWrapper(d ResourceData, key string, data any) {
	err := d.Set(key, data)
	if err != nil {
		panic(err)
	}
}

// Cast string to enum
func CastString[T ~string](raw *string) *T {
	if raw == nil {
		return nil
	}
	t := T(*raw)
	return &t
}

// Get Pointer value from resource data, if castable to generic type
func GetP[T any](d ResourceData, key string) *T {
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
func GetIntP(d ResourceData, key string) *int32 {
	rv, ok := d.GetOk(key)
	if !ok {
		return nil
	}
	if tt, ok := rv.(int); ok {
		return new(int32(tt))
	}
	return nil
}

// Similar to `GetP` however also casts to an int64 as that is what the API prefers
func GetInt64P(d ResourceData, key string) *int64 {
	rv, ok := d.GetOk(key)
	if !ok {
		return nil
	}
	if tt, ok := rv.(int); ok {
		return new(int64(tt))
	}
	return nil
}

func GetJSON[T any](d ResourceData, key string) (T, diag.Diagnostics) {
	var v T
	if sv := GetP[string](d, key); sv != nil {
		err := json.NewDecoder(strings.NewReader(*sv)).Decode(&v)
		if err != nil {
			return v, diag.FromErr(err)
		}
	}
	return v, nil
}

func SetJSON[T any](d ResourceData, key string, v T) diag.Diagnostics {
	var diags diag.Diagnostics
	b, err := json.Marshal(v)
	if err != nil {
		return diag.FromErr(err)
	}
	SetWrapper(d, key, string(b))
	return diags
}
