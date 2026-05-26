package provider

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "goauthentik.io/api/v3"
)

func TestResourceOutpostDeleteTreatsDeleted405AsSuccess(t *testing.T) {
	const outpostID = "83fcc504-d209-41b8-9840-1883aa2e9640"
	deleteCalled := false
	retrieveCount := 0
	setOutpostDeleteRetryTiming(t, 100*time.Millisecond, time.Millisecond)

	client := testOutpostDeleteClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			deleteCalled = true
			w.WriteHeader(http.StatusMethodNotAllowed)
		case http.MethodGet:
			retrieveCount++
			if retrieveCount == 1 {
				writeOutpostDeleteFixture(w, outpostID)
				return
			}
			w.WriteHeader(http.StatusNotFound)
		default:
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
	})

	d := schema.TestResourceDataRaw(t, resourceOutpost().Schema, map[string]any{})
	d.SetId(outpostID)

	diags := resourceOutpostDelete(t.Context(), d, &APIClient{client: client})
	if diags.HasError() {
		t.Fatalf("expected delete to succeed, got diagnostics: %#v", diags)
	}
	if !deleteCalled {
		t.Fatal("expected DELETE request")
	}
	if retrieveCount < 2 {
		t.Fatalf("expected delete to retry until 404, got %d GET request(s)", retrieveCount)
	}
	if d.Id() != "" {
		t.Fatalf("expected resource ID to be cleared, got %q", d.Id())
	}
}

func TestResourceOutpostDeleteKeeps405WhenOutpostStillExists(t *testing.T) {
	const outpostID = "83fcc504-d209-41b8-9840-1883aa2e9640"
	setOutpostDeleteRetryTiming(t, 5*time.Millisecond, time.Millisecond)

	client := testOutpostDeleteClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			w.WriteHeader(http.StatusMethodNotAllowed)
		case http.MethodGet:
			writeOutpostDeleteFixture(w, outpostID)
		default:
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
	})

	d := schema.TestResourceDataRaw(t, resourceOutpost().Schema, map[string]any{})
	d.SetId(outpostID)

	diags := resourceOutpostDelete(t.Context(), d, &APIClient{client: client})
	if !diags.HasError() {
		t.Fatal("expected delete diagnostics when outpost still exists")
	}
	if d.Id() != outpostID {
		t.Fatalf("expected resource ID to stay set, got %q", d.Id())
	}
}

func setOutpostDeleteRetryTiming(t *testing.T, timeout time.Duration, interval time.Duration) {
	t.Helper()

	previousTimeout := resourceOutpostDeleteRetryTimeout
	previousInterval := resourceOutpostDeleteRetryInterval
	resourceOutpostDeleteRetryTimeout = timeout
	resourceOutpostDeleteRetryInterval = interval
	t.Cleanup(func() {
		resourceOutpostDeleteRetryTimeout = previousTimeout
		resourceOutpostDeleteRetryInterval = previousInterval
	})
}

func writeOutpostDeleteFixture(w http.ResponseWriter, outpostID string) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"pk":"` + outpostID + `","name":"still-here","type":"proxy","providers":[1],"config":{"authentik_host":"http://localhost:9000/"},"managed":null}`))
}

func testOutpostDeleteClient(t *testing.T, handler http.HandlerFunc) *api.APIClient {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	config := api.NewConfiguration()
	config.Servers = api.ServerConfigurations{
		{URL: server.URL + "/api/v3"},
	}
	config.HTTPClient = server.Client()

	return api.NewAPIClient(config)
}
