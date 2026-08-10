package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPIURL_PathBased(t *testing.T) {
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
			got, err := APIURL(tc.inputURL)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedURL, got)
		})
	}
}
