package helpers

import "goauthentik.io/api/v3"

// APIClient Hold the API Client and any relevant configuration
type APIClient struct {
	Client *api.APIClient
}
