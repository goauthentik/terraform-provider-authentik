package helpers

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsNotFound(t *testing.T) {
	assert.False(t, IsNotFound(nil))
	assert.False(t, IsNotFound(&http.Response{StatusCode: 500}))
	assert.True(t, IsNotFound(&http.Response{StatusCode: 404}))
}

func TestHTTPError(t *testing.T) {
	diags := HTTPError(nil, errors.New("boom"))
	assert.True(t, diags.HasError())

	resp := &http.Response{
		StatusCode: 500,
		Body:       io.NopCloser(bytes.NewBufferString(`{"detail":"failed"}`)),
		Request:    &http.Request{Method: "GET", URL: &url.URL{Path: "/api/v3/core/groups/x/"}},
	}
	diags = HTTPError(resp, errors.New("request failed"))
	assert.True(t, diags.HasError())
}
