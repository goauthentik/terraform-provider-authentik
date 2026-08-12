package helpers

import (
	"context"

	"github.com/getsentry/sentry-go"
)

// Span starts a Sentry span and returns a func to defer to close it, for the
// `defer helpers.Span(ctx, ...)()` pattern. It does not propagate the span's context
// into ctx-taking calls made inside the deferred block - this is intentionally the same
// behaviour as pkg/sdkprovider/tracing.go's tr()/td() wrappers, which pass the original
// ctx through unchanged too.
func Span(ctx context.Context, transaction, operation, description string) func() {
	span := sentry.StartSpan(ctx, operation, sentry.WithTransactionName(transaction))
	span.Description = description
	return span.Finish
}
