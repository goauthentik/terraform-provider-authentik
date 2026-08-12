package helpers

import (
	"fmt"
	"strings"

	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// This file replaces the SDKv2 provider's init()-installed schema.SchemaDescriptionBuilder,
// which appended " Defaults to `x`." when a Default was set and " Generated." when an
// attribute was Computed. The framework has no equivalent global post-processing hook,
// so descriptions are composed once, at attribute construction time, instead:
//
//	MarkdownDescription: helpers.Desc(helpers.EnumToDescription(api.AllowedFooEnumValues),
//	                                  helpers.WithDefault(api.FOO_BAR)),
//
// A whole-schema walker that mutates an already-built schema.Schema was considered and
// rejected: schema.Attribute is an interface over six concrete value-type structs, so
// changing a description after the fact means type-switching across all six and
// rebuilding each one - and it would still need these same value-exposing wrappers to
// avoid stating every default twice. Construction-time composition is type-checked, has
// no runtime cost, and lets the handful of attributes that should not say "Generated"
// (despite being Computed) simply not call that option.
//
// Per H2, every attribute with a Default is now also Computed (Default requires
// Computed in the framework), so callers must gate Generated() on "Computed AND no
// Default" themselves - never call both WithDefault and Generated for the same
// attribute, or key Generated purely off Computed, or all 318 defaulted attributes
// wrongly gain a "Generated." suffix they never had in the SDKv2 docs.

type descOptions struct {
	hasDefault bool
	defaultVal any
	generated  bool
}

// DescOption configures Desc's output. See WithDefault and Generated.
type DescOption func(*descOptions)

// WithDefault appends " Defaults to `v`." to the description, matching the SDKv2
// SchemaDescriptionBuilder's handling of Schema.Default.
func WithDefault(v any) DescOption {
	return func(o *descOptions) {
		o.hasDefault = true
		o.defaultVal = v
	}
}

// Generated appends " Generated." to the description, matching the SDKv2
// SchemaDescriptionBuilder's handling of Schema.Computed. Only call this for
// attributes that are Computed with no Default - see the H2 note above.
func Generated() DescOption {
	return func(o *descOptions) {
		o.generated = true
	}
}

// Desc composes an attribute's MarkdownDescription the way the SDKv2 provider's
// SchemaDescriptionBuilder used to, at construction time instead of via a global hook.
func Desc(base string, opts ...DescOption) string {
	o := &descOptions{}
	for _, opt := range opts {
		opt(o)
	}

	desc := base
	if o.hasDefault {
		desc += fmt.Sprintf(" Defaults to `%v`.", o.defaultVal)
	}
	if o.generated {
		desc += " Generated."
	}
	return strings.TrimSpace(desc)
}

// MarkResourceDeprecated appends a deprecation notice to s's description and sets
// DeprecationMessage, matching SDKv2's helpers.MarkDeprecated (kept as-is in
// pkg/sdkprovider until Phase 5, hence the different name here to avoid colliding with
// it). The message must be in both places: DeprecationMessage is what Terraform's own
// deprecation UI reads, but templates/resources.md.tmpl renders the description
// verbatim, so without the "~>" callout appended there too the generated doc page loses
// the deprecation notice.
func MarkResourceDeprecated(s rschema.Schema, newName string) rschema.Schema {
	msg := fmt.Sprintf("This resource is deprecated. Migrate to `%s`.", newName)
	s.DeprecationMessage = msg
	s.MarkdownDescription += fmt.Sprintf("\n\n~> %s", msg)
	return s
}
