package provider

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	api "goauthentik.io/api/v3"
)

func TestAccResourceStagePromptField(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceStagePromptField(rName, "radio-button-group", "return ['true', 'false']", true, "false"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_prompt_field.field", "name", rName),
					resource.TestCheckResourceAttr("authentik_stage_prompt_field.field", "type", "radio-button-group"),
					resource.TestCheckResourceAttr("authentik_stage_prompt_field.field", "placeholder", "return ['true', 'false']"),
					resource.TestCheckResourceAttr("authentik_stage_prompt_field.field", "initial_value", "false"),
				),
			},
			// Verify that updating the field clears placeholder and initial_value
			{
				Config: testAccResourceStagePromptField(rName, "checkbox", "", false, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("authentik_stage_prompt_field.field", "type", "checkbox"),
					resource.TestCheckResourceAttr("authentik_stage_prompt_field.field", "placeholder", ""),
					resource.TestCheckResourceAttr("authentik_stage_prompt_field.field", "initial_value", ""),
				),
			},
		},
	})
}

func testAccResourceStagePromptField(name string, fieldType string, placeholder string, placeholderExpression bool, initialValue string) string {
	return fmt.Sprintf(`
resource "authentik_stage_prompt_field" "field" {
  name                   = "%[1]s"
  field_key              = "%[1]s"
  label                  = "%[1]s"
  type                   = "%[2]s"
  required               = false
  placeholder            = "%[3]s"
  placeholder_expression = %[4]t
  initial_value          = "%[5]s"
}
`, name, fieldType, placeholder, placeholderExpression, initialValue)
}

func TestResourceStagePromptFieldUpdateSendsEmptyStrings(t *testing.T) {
	var putBody map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v3/stages/prompt/prompts/prompt-1/", r.URL.Path)
		if r.Method == http.MethodPut {
			raw, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			require.NoError(t, json.Unmarshal(raw, &putBody))
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"pk": "prompt-1",
			"name": "consent",
			"field_key": "consent",
			"label": "consent",
			"type": "checkbox",
			"required": false,
			"placeholder": "",
			"initial_value": "",
			"order": 0,
			"prompt_stages_obj": [],
			"sub_text": "",
			"placeholder_expression": false,
			"initial_value_expression": false
		}`))
	}))
	t.Cleanup(server.Close)

	config := api.NewConfiguration()
	config.Servers = api.ServerConfigurations{{
		URL: server.URL + "/api/v3",
	}}
	client := &APIClient{
		client: api.NewAPIClient(config),
	}

	d := schema.TestResourceDataRaw(t, resourceStagePromptField().Schema, map[string]any{
		"name":          "consent",
		"field_key":     "consent",
		"label":         "consent",
		"type":          "checkbox",
		"required":      false,
		"placeholder":   "",
		"initial_value": "",
	})
	d.SetId("prompt-1")

	diags := resourceStagePromptFieldUpdate(t.Context(), d, client)

	require.False(t, diags.HasError(), diags)
	// Omitting the keys would make authentik keep the previous values
	require.Contains(t, putBody, "placeholder")
	require.Contains(t, putBody, "initial_value")
	assert.Equal(t, "", putBody["placeholder"])
	assert.Equal(t, "", putBody["initial_value"])
}
