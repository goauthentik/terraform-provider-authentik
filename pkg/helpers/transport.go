package helpers

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"

	"github.com/getsentry/sentry-go"
)

// TestingTransport is a Transport used for testing, always returns a 400 Response
type TestingTransport struct {
	inner http.RoundTripper
}

// NewTestingTransport returns a Transport that fails all requests
func NewTestingTransport(inner http.RoundTripper) *TestingTransport {
	return &TestingTransport{inner}
}

// RoundTrip implements http.RoundTripper
func (tt *TestingTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	body := "mock-failed-request"
	return &http.Response{
		Status:        "400 Bad Request",
		StatusCode:    400,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Body:          io.NopCloser(bytes.NewBufferString(body)),
		ContentLength: int64(len(body)),
		Request:       r,
		Header:        make(http.Header),
	}, nil
}

// GetTLSTransport returns a TLS transport instance, that skips verification if configured to.
func GetTLSTransport(insecure bool) http.RoundTripper {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: insecure,
		},
		Proxy: http.ProxyFromEnvironment,
	}
	return transport
}

type tracingTransport struct {
	inner http.RoundTripper
	ctx   context.Context
}

// NewTracingTransport wraps inner with a Transport that emits a Sentry span per request.
func NewTracingTransport(ctx context.Context, inner http.RoundTripper) *tracingTransport {
	return &tracingTransport{inner, ctx}
}

// RoundTrip implements http.RoundTripper
func (tt *tracingTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	span := sentry.StartSpan(tt.ctx, "authentik.go.http_request")
	r.Header.Set("sentry-trace", span.ToSentryTrace())
	span.Description = fmt.Sprintf("%s %s", r.Method, r.URL.String())
	span.SetTag("url", r.URL.String())
	span.SetTag("method", r.Method)
	defer span.Finish()
	res, err := tt.inner.RoundTrip(r.WithContext(span.Context()))
	return res, err
}
