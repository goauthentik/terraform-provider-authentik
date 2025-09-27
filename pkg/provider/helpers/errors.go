package helpers

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func HTTPToDiag(d *schema.ResourceData, r *http.Response, err error) diag.Diagnostics {
	if r == nil {
		return diag.Errorf("HTTP Error '%s' without http response", err.Error())
	}
	if r.StatusCode == 404 {
		d.SetId("")
		return diag.Diagnostics{}
	}
	buff := &bytes.Buffer{}
	_, er := io.Copy(buff, r.Body)
	if er != nil {
		log.Printf("[DEBUG] authentik: failed to read response: %s", er.Error())
	}
	log.Printf("[DEBUG] authentik: error response: %s", buff.String())
	return diag.Errorf("HTTP Error '%s' during request '%s %s': \"%s\"", err.Error(), r.Request.Method, r.Request.URL.Path, buff.String())
}
