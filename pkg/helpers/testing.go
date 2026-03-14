package helpers

import (
	"os"
	"testing"
)

func TestAccPreCheck(t *testing.T) {
	testEnvIsSet("AUTHENTIK_URL", t)
	testEnvIsSet("AUTHENTIK_TOKEN", t)
}

func testEnvIsSet(k string, t *testing.T) {
	if v := os.Getenv(k); v == "" {
		t.Fatalf("%[1]s must be set for acceptance tests", k)
	}
}
