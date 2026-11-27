package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"goauthentik.io/api/v3"
)

func redirectURI(url string, t api.RedirectURITypeEnum) api.RedirectURI {
	return api.RedirectURI{
		MatchingMode:    api.MATCHINGMODEENUM_STRICT,
		Url:             url,
		RedirectUriType: &t,
	}
}

func redirectURIConfig(url string, t string) map[string]any {
	m := map[string]any{"matching_mode": "strict", "url": url}
	if t != "" {
		m["redirect_uri_type"] = t
	}
	return m
}

// authentik returns `redirect_uri_type = "authorization"` for entries created
// without a type. When the configuration doesn't set it, that server default
// must not end up in state, or every plan shows a diff and every apply rewrites
// the provider (#956).
func TestMergeRedirectURIs_DefaultTypeNotInConfig(t *testing.T) {
	config := []any{
		redirectURIConfig("https://a.example.com/cb", ""),
		redirectURIConfig("https://b.example.com/cb", ""),
	}
	server := []api.RedirectURI{
		redirectURI("https://a.example.com/cb", api.REDIRECTURITYPEENUM_AUTHORIZATION),
		redirectURI("https://b.example.com/cb", api.REDIRECTURITYPEENUM_AUTHORIZATION),
	}

	assert.Equal(t, []map[string]any{
		redirectURIConfig("https://a.example.com/cb", ""),
		redirectURIConfig("https://b.example.com/cb", ""),
	}, mergeRedirectURIs(config, server))
}

func TestMergeRedirectURIs_ExplicitTypesKept(t *testing.T) {
	config := []any{
		redirectURIConfig("https://a.example.com/cb", "authorization"),
		redirectURIConfig("https://a.example.com/logout", "logout"),
	}
	server := []api.RedirectURI{
		redirectURI("https://a.example.com/cb", api.REDIRECTURITYPEENUM_AUTHORIZATION),
		redirectURI("https://a.example.com/logout", api.REDIRECTURITYPEENUM_LOGOUT),
	}

	assert.Equal(t, []map[string]any{
		redirectURIConfig("https://a.example.com/cb", "authorization"),
		redirectURIConfig("https://a.example.com/logout", "logout"),
	}, mergeRedirectURIs(config, server))
}

// A non-default type on the server must still show up, so a change made
// outside Terraform is detected.
func TestMergeRedirectURIs_NonDefaultServerTypeShown(t *testing.T) {
	config := []any{
		redirectURIConfig("https://a.example.com/cb", ""),
	}
	server := []api.RedirectURI{
		redirectURI("https://a.example.com/cb", api.REDIRECTURITYPEENUM_LOGOUT),
	}

	assert.Equal(t, []map[string]any{
		redirectURIConfig("https://a.example.com/cb", "logout"),
	}, mergeRedirectURIs(config, server))
}

// On import (no configuration yet) the server's types are kept as they are.
func TestMergeRedirectURIs_Import(t *testing.T) {
	server := []api.RedirectURI{
		redirectURI("https://a.example.com/cb", api.REDIRECTURITYPEENUM_AUTHORIZATION),
		redirectURI("https://a.example.com/logout", api.REDIRECTURITYPEENUM_LOGOUT),
	}

	assert.Equal(t, []map[string]any{
		redirectURIConfig("https://a.example.com/cb", "authorization"),
		redirectURIConfig("https://a.example.com/logout", "logout"),
	}, mergeRedirectURIs(nil, server))
}

// Config order is kept, and entries only on the server are appended.
func TestMergeRedirectURIs_OrderAndServerOnlyEntries(t *testing.T) {
	config := []any{
		redirectURIConfig("https://b.example.com/cb", ""),
		redirectURIConfig("https://a.example.com/cb", ""),
	}
	server := []api.RedirectURI{
		redirectURI("https://a.example.com/cb", api.REDIRECTURITYPEENUM_AUTHORIZATION),
		redirectURI("https://c.example.com/cb", api.REDIRECTURITYPEENUM_AUTHORIZATION),
		redirectURI("https://b.example.com/cb", api.REDIRECTURITYPEENUM_AUTHORIZATION),
	}

	assert.Equal(t, []map[string]any{
		redirectURIConfig("https://b.example.com/cb", ""),
		redirectURIConfig("https://a.example.com/cb", ""),
		redirectURIConfig("https://c.example.com/cb", "authorization"),
	}, mergeRedirectURIs(config, server))
}
