package helpers

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"goauthentik.io/api/v3"
)

func ModelSchema(m api.ModelEnum, os map[string]*schema.Schema) map[string]*schema.Schema {
	parts := strings.Split(string(m), ".")
	os["meta_app"] = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: fmt.Sprintf("Static value of `%s`", parts[0]),
		DefaultFunc: func() (interface{}, error) {
			return parts[0], nil
		},
	}
	os["meta_model"] = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: fmt.Sprintf("Static value of `%s`", parts[1]),
		DefaultFunc: func() (interface{}, error) {
			return parts[1], nil
		},
	}
	return os
}
