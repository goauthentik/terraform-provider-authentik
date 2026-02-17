package helpers

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	RelativeDurationDescription = "Format: hours=1;minutes=2;seconds=3."
	JSONDescription             = "JSON format expected. Use `jsonencode()` to pass objects."
)

func MarkDeprecated(resource func() *schema.Resource, newName string) func() *schema.Resource {
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
	sb := &strings.Builder{}
	sb.WriteString("Allowed values:\n")
	for _, v := range allowed {
		fmt.Fprintf(sb, "  - `%s`\n", v)
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

// DiffSuppressExpression Diff suppression for python expressions
func DiffSuppressExpression(k, old, new string, d *schema.ResourceData) bool {
	return strings.TrimSuffix(new, "\n") == old
}

// DiffSuppressJSON Diff suppression for JSON objects
func DiffSuppressJSON(k, old, new string, d *schema.ResourceData) bool {
	var j, j2 interface{}
	if err := json.Unmarshal([]byte(old), &j); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(new), &j2); err != nil {
		return false
	}
	return reflect.DeepEqual(j2, j)
}
