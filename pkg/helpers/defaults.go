package helpers

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
)

// This file exists because defaults.String/Bool/Int32/Float64 (the schema.Attribute.Default
// value) have no accessor for the value they default to - and that value is also needed
// wherever the attribute's description is composed (Desc(..., WithDefault(v))). Without
// these wrappers, every defaulted attribute would state its default twice: once for
// Default, once for the description. Each wrapper embeds the real default.X
// implementation (so it still satisfies schema.XAttribute's Default field directly) and
// adds a Value() accessor.

// StringDefaultValue is a defaults.String that also exposes the value it defaults to.
type StringDefaultValue struct {
	defaults.String
	value string
}

// StringDefault builds a StringDefaultValue for v.
func StringDefault(v string) StringDefaultValue {
	return StringDefaultValue{String: stringdefault.StaticString(v), value: v}
}

// Value returns the default value.
func (d StringDefaultValue) Value() string { return d.value }

// BoolDefaultValue is a defaults.Bool that also exposes the value it defaults to.
type BoolDefaultValue struct {
	defaults.Bool
	value bool
}

// BoolDefault builds a BoolDefaultValue for v.
func BoolDefault(v bool) BoolDefaultValue {
	return BoolDefaultValue{Bool: booldefault.StaticBool(v), value: v}
}

// Value returns the default value.
func (d BoolDefaultValue) Value() bool { return d.value }

// Int32DefaultValue is a defaults.Int32 that also exposes the value it defaults to.
type Int32DefaultValue struct {
	defaults.Int32
	value int32
}

// Int32Default builds an Int32DefaultValue for v.
func Int32Default(v int32) Int32DefaultValue {
	return Int32DefaultValue{Int32: int32default.StaticInt32(v), value: v}
}

// Value returns the default value.
func (d Int32DefaultValue) Value() int32 { return d.value }

// Float64DefaultValue is a defaults.Float64 that also exposes the value it defaults to.
type Float64DefaultValue struct {
	defaults.Float64
	value float64
}

// Float64Default builds a Float64DefaultValue for v.
func Float64Default(v float64) Float64DefaultValue {
	return Float64DefaultValue{Float64: float64default.StaticFloat64(v), value: v}
}

// Value returns the default value.
func (d Float64DefaultValue) Value() float64 { return d.value }
