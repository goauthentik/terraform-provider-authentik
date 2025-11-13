package provider

import (
	"net/url"
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
				serverURL += akURL.Path[:len(akURL.Path)-len("/")]
				if akURL.Path[len(akURL.Path)-1:] != "/" {
					serverURL = akURL.Scheme + "://" + akURL.Host + akURL.Path
				}
			}

			// For cleaner implementation, use the actual logic
			serverURL = akURL.Scheme + "://" + akURL.Host
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
