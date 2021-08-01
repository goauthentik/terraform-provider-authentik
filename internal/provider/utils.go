package provider

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func stringToPointer(in string) *string {
	return &in
}

func stringPointerResolve(in *string) string {
	return *in
}

func intToPointer(in int) *int32 {
	i := int32(in)
	return &i
}

func int32ToPointer(in int32) *int32 {
	return &in
}

func boolToPointer(in bool) *bool {
	return &in
}

func httpToDiag(r *http.Response) diag.Diagnostics {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("[DEBUG] authentik: failed to read response: %s", err)
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] authentik: error response: %s", string(b))
	return diag.Errorf("HTTP Error '%s' during request '%s %s'", string(b), r.Request.Method, r.Request.URL.Path)
}
