package helpers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// IsNotFound reports whether hr is a 404 response. Deliberately narrow: callers decide
// what a 404 means for the operation they're in, rather than this function deciding for
// them. In particular this is NOT a drop-in port of pkg/helpers.HTTPToDiag's implicit
// "404 always clears state" behaviour - only Read should act on it by calling
// resp.State.RemoveResource(ctx). A 404 during Update is a real error (the SDKv2 version
// silently dropped the resource and reported success on a 404 there; that is a bug, not
// behaviour to preserve). A 404 during Delete is a success as usual.
func IsNotFound(hr *http.Response) bool {
	return hr != nil && hr.StatusCode == http.StatusNotFound
}

// HTTPError converts a failed authentik API call into diagnostics. It does not special
// case 404 - check IsNotFound(hr) first if a 404 should be handled differently.
func HTTPError(hr *http.Response, err error) diag.Diagnostics {
	var diags diag.Diagnostics

	if hr == nil {
		diags.AddError("authentik API request failed", err.Error())
		return diags
	}

	body, readErr := io.ReadAll(hr.Body)
	if readErr != nil {
		body = []byte(fmt.Sprintf("failed to read response body: %s", readErr.Error()))
	}

	summary := "authentik API request failed"
	if hr.Request != nil {
		summary = fmt.Sprintf("authentik API request '%s %s' failed", hr.Request.Method, hr.Request.URL.Path)
	}
	diags.AddError(summary, fmt.Sprintf("%s: %s", err.Error(), string(body)))
	return diags
}
