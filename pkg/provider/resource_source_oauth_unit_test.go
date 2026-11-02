package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// When `oidc_well_known_url` is set, authentik fills in the OAuth2 endpoint URLs
// from the discovery document. Those server-populated values must not show up as
// a diff on the next plan when they are not set in the configuration.
func TestResourceSourceOAuthWellKnownURLsNoDiff(t *testing.T) {
	r := resourceSourceOAuth()
	discovered := []string{"authorization_url", "access_token_url", "profile_url", "oidc_jwks_url"}

	for _, k := range discovered {
		s := r.Schema[k]
		require.NotNil(t, s, k)
		assert.True(t, s.Optional, "%s must be optional", k)
		assert.True(t, s.Computed, "%s must be computed", k)
	}

	state := &terraform.InstanceState{
		ID: "test",
		Attributes: map[string]string{
			"id":                  "test",
			"name":                "test",
			"slug":                "test",
			"oidc_well_known_url": "https://idp.example.com/.well-known/openid-configuration",
			"authorization_url":   "https://idp.example.com/authorize",
			"access_token_url":    "https://idp.example.com/token",
			"profile_url":         "https://idp.example.com/userinfo",
			"oidc_jwks_url":       "https://idp.example.com/jwks",
		},
	}
	cfg := terraform.NewResourceConfigRaw(map[string]any{
		"name":                "test",
		"slug":                "test",
		"oidc_well_known_url": "https://idp.example.com/.well-known/openid-configuration",
	})

	diff, err := r.Diff(context.Background(), state, cfg, nil)
	require.NoError(t, err)
	if diff != nil {
		for _, k := range discovered {
			_, ok := diff.Attributes[k]
			assert.False(t, ok, "unexpected diff on %s", k)
		}
	}
}
