package helpers

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"goauthentik.io/api/v3"
)

func ModelSchema(m api.ModelEnum, os map[string]*schema.Schema) map[string]*schema.Schema {
	parts := strings.Split(string(m), ".")
	os["meta_app"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
		Default:  parts[0],
	}
	os["meta_model"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
		Default:  parts[1],
	}
	return os
}
