package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
)

func TestProvider(t *testing.T) {
	p := Provider("testing", false)
	if err := p.InternalValidate(); err != nil {
		t.Fatalf("err: %[1]s", err)
	}
}

func TestProviderConfigure_PathBasedURL(t *testing.T) {
	testCases := []struct {
		name        string
		inputURL    string
		expectedURL string
	}{
		{
			name:        "Root path with trailing slash",
			inputURL:    "https://api.example.com/",
			expectedURL: "https://api.example.com/api/v3",
		},
		{
			name:        "Root path without trailing slash",
			inputURL:    "https://api.example.com",
			expectedURL: "https://api.example.com/api/v3",
		},
		{
			name:        "Single segment path with trailing slash",
			inputURL:    "https://api.example.com/sso/",
			expectedURL: "https://api.example.com/sso/api/v3",
		},
		{
			name:        "Single segment path without trailing slash",
			inputURL:    "https://api.example.com/sso",
			expectedURL: "https://api.example.com/sso/api/v3",
		},
		{
			name:        "Multi-segment path with trailing slash",
			inputURL:    "https://api.example.com/auth/sso/",
			expectedURL: "https://api.example.com/auth/sso/api/v3",
		},
		{
			name:        "Multi-segment path without trailing slash",
			inputURL:    "https://api.example.com/auth/sso",
			expectedURL: "https://api.example.com/auth/sso/api/v3",
		},
		{
			name:        "HTTP scheme with path",
			inputURL:    "http://localhost:9000/sso/",
			expectedURL: "http://localhost:9000/sso/api/v3",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := Provider("testing", true)

			_ac, diag := p.ConfigureContextFunc(t.Context(), schema.TestResourceDataRaw(t, p.Schema, map[string]any{
				"url":      tc.inputURL,
				"token":    "",
				"insecure": false,
			}))
			assert.Nil(t, diag)
			ac := _ac.(*helpers.APIClient)
			assert.Equal(t, tc.expectedURL, ac.Client.GetConfig().Servers[0].URL, "Server URL should be constructed correctly")
		})
	}
}
