package helpers

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// OneOf wraps stringvalidator.OneOf for T ~string enum types (85 call sites in SDKv2,
// via helpers.StringInEnum), so callers can keep passing the generated api.AllowedXxx
// slices directly instead of converting to []string themselves.
func OneOf[T ~string](allowed []T) validator.String {
	values := make([]string, len(allowed))
	for i, v := range allowed {
		values[i] = string(v)
	}
	return stringvalidator.OneOf(values...)
}

var relativeDurationKeys = []string{
	"microseconds",
	"milliseconds",
	"seconds",
	"minutes",
	"hours",
	"days",
	"weeks",
}

// RelativeDuration validates authentik's relative-duration string format
// ("hours=1;minutes=2;seconds=3"), replacing SDKv2's ValidateRelativeDuration.
func RelativeDuration() validator.String {
	return relativeDurationValidator{}
}

type relativeDurationValidator struct{}

func (v relativeDurationValidator) Description(_ context.Context) string {
	return RelativeDurationDescription
}

func (v relativeDurationValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v relativeDurationValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()
	for el := range strings.SplitSeq(value, ";") {
		p := strings.SplitN(el, "=", 2)
		if len(p) < 2 {
			resp.Diagnostics.AddAttributeError(req.Path, "Invalid relative duration", fmt.Sprintf("%q has incorrect amount of elements", el))
			return
		}

		isValid := false
		for _, valid := range relativeDurationKeys {
			if strings.EqualFold(p[0], valid) {
				isValid = true
				break
			}
		}
		if !isValid {
			resp.Diagnostics.AddAttributeError(req.Path, "Invalid relative duration", fmt.Sprintf("%q has incorrect key %q", el, p[0]))
		}
	}
}
