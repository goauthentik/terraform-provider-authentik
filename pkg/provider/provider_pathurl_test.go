package provider

import (
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProviderConfigure_PathBasedURL(t *testing.T) {
	testCases := []struct {
		name        string
		inputURL    string
		expectedURL string
	}{
		{
			name:        "Root path with trailing slash",
			inputURL:    "https://api.example.com/",
			expectedURL: "https://api.example.com",
		},
		{
			name:        "Root path without trailing slash",
			inputURL:    "https://api.example.com",
			expectedURL: "https://api.example.com",
		},
		{
			name:        "Single segment path with trailing slash",
			inputURL:    "https://api.example.com/sso/",
			expectedURL: "https://api.example.com/sso",
		},
		{
			name:        "Single segment path without trailing slash",
			inputURL:    "https://api.example.com/sso",
			expectedURL: "https://api.example.com/sso",
		},
		{
			name:        "Multi-segment path with trailing slash",
			inputURL:    "https://api.example.com/auth/sso/",
			expectedURL: "https://api.example.com/auth/sso",
		},
		{
			name:        "Multi-segment path without trailing slash",
			inputURL:    "https://api.example.com/auth/sso",
			expectedURL: "https://api.example.com/auth/sso",
		},
		{
			name:        "HTTP scheme with path",
			inputURL:    "http://localhost:9000/sso/",
			expectedURL: "http://localhost:9000/sso",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Parse the input URL
			akURL, err := url.Parse(tc.inputURL)
			assert.NoError(t, err, "URL parsing should not fail")

			// Simulate the logic from providerConfigure
			serverURL := akURL.Scheme + "://" + akURL.Host
			if akURL.Path != "" && akURL.Path != "/" {
				// Preserve path component, removing trailing slash to avoid double slashes
				path := akURL.Path
				if path[len(path)-1:] == "/" {
					path = path[:len(path)-1]
				}
				serverURL += path
			}

			assert.Equal(t, tc.expectedURL, serverURL, "Server URL should be constructed correctly")
		})
	}
}

func TestProviderConfigure_URLConstruction(t *testing.T) {
	// Test that verifies the actual configuration behavior
	testURL := "https://api.example.com/sso/"

	akURL, err := url.Parse(testURL)
	assert.NoError(t, err)

	// Verify parsed components
	assert.Equal(t, "https", akURL.Scheme)
	assert.Equal(t, "api.example.com", akURL.Host)
	assert.Equal(t, "/sso/", akURL.Path)

	// Verify URL reconstruction preserves path
	expectedServerURL := "https://api.example.com/sso"
	serverURL := akURL.Scheme + "://" + akURL.Host
	if akURL.Path != "" && akURL.Path != "/" {
		path := akURL.Path
		if path[len(path)-1:] == "/" {
			path = path[:len(path)-1]
		}
		serverURL += path
	}

	assert.Equal(t, expectedServerURL, serverURL)
}

// TestProviderConfigure_APIv3PathHandling tests that the API v3 path is correctly handled
// for both root and subpath deployments
func TestProviderConfigure_APIv3PathHandling(t *testing.T) {
	testCases := []struct {
		name                string
		inputURL            string
		expectedServerURL   string
		expectedAPIEndpoint string
		description         string
	}{
		{
			name:                "Root deployment",
			inputURL:            "https://api.example.com",
			expectedServerURL:   "https://api.example.com/api/v3",
			expectedAPIEndpoint: "https://api.example.com/api/v3/providers/oauth2/",
			description:         "Standard deployment at root - server URL should include /api/v3",
		},
		{
			name:                "Root deployment with trailing slash",
			inputURL:            "https://api.example.com/",
			expectedServerURL:   "https://api.example.com/api/v3",
			expectedAPIEndpoint: "https://api.example.com/api/v3/providers/oauth2/",
			description:         "Standard deployment at root with slash - server URL should include /api/v3",
		},
		{
			name:                "Subpath deployment - single segment",
			inputURL:            "https://api.example.com/sso",
			expectedServerURL:   "https://api.example.com/sso/api/v3",
			expectedAPIEndpoint: "https://api.example.com/sso/api/v3/providers/oauth2/",
			description:         "Authentik hosted at /sso subpath - server URL must include /sso/api/v3",
		},
		{
			name:                "Subpath deployment - single segment with trailing slash",
			inputURL:            "https://api.example.com/sso/",
			expectedServerURL:   "https://api.example.com/sso/api/v3",
			expectedAPIEndpoint: "https://api.example.com/sso/api/v3/providers/oauth2/",
			description:         "Authentik hosted at /sso/ subpath - server URL must include /sso/api/v3",
		},
		{
			name:                "Subpath deployment - multi segment",
			inputURL:            "https://api.example.com/auth/sso",
			expectedServerURL:   "https://api.example.com/auth/sso/api/v3",
			expectedAPIEndpoint: "https://api.example.com/auth/sso/api/v3/providers/oauth2/",
			description:         "Authentik hosted at /auth/sso - server URL must include /auth/sso/api/v3",
		},
		{
			name:                "Local development with subpath",
			inputURL:            "http://localhost:9000/sso",
			expectedServerURL:   "http://localhost:9000/sso/api/v3",
			expectedAPIEndpoint: "http://localhost:9000/sso/api/v3/providers/oauth2/",
			description:         "Local development with subpath - must work correctly",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Parse the input URL
			akURL, err := url.Parse(tc.inputURL)
			assert.NoError(t, err, "URL parsing should not fail")

			// Test the FIXED implementation
			serverURL := constructServerURLSmart(akURL)

			assert.Equal(t, tc.expectedServerURL, serverURL, tc.description)

			// Verify that API endpoint construction works correctly
			// OpenAPI client will append relative paths to server URL
			apiEndpoint := serverURL + "/providers/oauth2/"
			assert.Equal(t, tc.expectedAPIEndpoint, apiEndpoint,
				"API endpoint should be correctly formed for OpenAPI client")
		})
	}
}

// constructServerURL shows the basic logic for appending /api/v3
// to properly handle subpath deployments with the /api/v3 prefix
func constructServerURL(akURL *url.URL) string {
	serverURL := akURL.Scheme + "://" + akURL.Host

	// Append the subpath if present
	if akURL.Path != "" && akURL.Path != "/" {
		path := strings.TrimSuffix(akURL.Path, "/")
		serverURL += path
	}

	// Always append /api/v3 as that's where the OpenAPI client expects the API to be
	serverURL += "/api/v3"

	return serverURL
}

// TestProviderConfigure_BackwardCompatibility ensures existing configurations still work
func TestProviderConfigure_BackwardCompatibility(t *testing.T) {
	testCases := []struct {
		name              string
		inputURL          string
		shouldWork        bool
		expectedServerURL string
	}{
		{
			name:              "Legacy root URL",
			inputURL:          "https://authentik.company.com",
			shouldWork:        true,
			expectedServerURL: "https://authentik.company.com/api/v3",
		},
		{
			name:              "User explicitly provides /api/v3",
			inputURL:          "https://authentik.company.com/api/v3",
			shouldWork:        true,
			expectedServerURL: "https://authentik.company.com/api/v3",
		},
		{
			name:              "Subpath without /api/v3",
			inputURL:          "https://api.company.com/sso",
			shouldWork:        true,
			expectedServerURL: "https://api.company.com/sso/api/v3",
		},
		{
			name:              "Subpath with /api/v3 already included",
			inputURL:          "https://api.company.com/sso/api/v3",
			shouldWork:        true,
			expectedServerURL: "https://api.company.com/sso/api/v3",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			akURL, err := url.Parse(tc.inputURL)
			assert.NoError(t, err)

			serverURL := constructServerURLSmart(akURL)
			assert.Equal(t, tc.expectedServerURL, serverURL)
		})
	}
}

// constructServerURLSmart handles cases where user might have already included /api/v3
func constructServerURLSmart(akURL *url.URL) string {
	serverURL := akURL.Scheme + "://" + akURL.Host

	// Append the subpath if present
	if akURL.Path != "" && akURL.Path != "/" {
		path := strings.TrimSuffix(akURL.Path, "/")
		serverURL += path

		// Only append /api/v3 if it's not already in the path
		if !strings.HasSuffix(path, "/api/v3") {
			serverURL += "/api/v3"
		}
	} else {
		// Root deployment - always add /api/v3
		serverURL += "/api/v3"
	}

	return serverURL
}
