package helpers

import (
	"testing"

	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/stretchr/testify/assert"
)

func TestDesc(t *testing.T) {
	assert.Equal(t, "some text", Desc("some text"))
	assert.Equal(t, "some text Defaults to `foo`.", Desc("some text", WithDefault("foo")))
	assert.Equal(t, "some text Generated.", Desc("some text", Generated()))
	// Never both in practice (H2: every Default implies Computed), but Desc itself
	// composes whatever it's given rather than enforcing that - the enforcement is
	// "don't call both" at each call site, documented in describe.go.
	assert.Equal(t, "some text Defaults to `1`. Generated.", Desc("some text", WithDefault(1), Generated()))
}

func TestMarkResourceDeprecated(t *testing.T) {
	s := rschema.Schema{
		MarkdownDescription: "Manage a widget.",
	}
	out := MarkResourceDeprecated(s, "authentik_new_widget")
	assert.Equal(t, "This resource is deprecated. Migrate to `authentik_new_widget`.", out.DeprecationMessage)
	assert.Contains(t, out.MarkdownDescription, "Manage a widget.")
	assert.Contains(t, out.MarkdownDescription, "~> This resource is deprecated. Migrate to `authentik_new_widget`.")
}
