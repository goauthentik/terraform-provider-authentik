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

func boolToPointer(in bool) *bool {
	return &in
}

func httpToDiag(r *http.Response, err error) diag.Diagnostics {
	b, er := ioutil.ReadAll(r.Body)
	if er != nil {
		log.Printf("[DEBUG] authentik: failed to read response: %s", er.Error())
		b = []byte{}
	}
	log.Printf("[DEBUG] authentik: error response: %s", string(b))
	return diag.Errorf("HTTP Error '%s' during request '%s %s': \"%s\"", err.Error(), r.Request.Method, r.Request.URL.Path, string(b))
}
