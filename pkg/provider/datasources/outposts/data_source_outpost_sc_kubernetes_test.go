package outposts_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"goauthentik.io/terraform-provider-authentik/pkg/helpers"
	"goauthentik.io/terraform-provider-authentik/pkg/provider"
)

func TestAccDataSourceOutpostServiceConnectionsKubernetes(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { helpers.TestAccPreCheck(t) },
		ProviderFactories: provider.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOutpostServiceConnectionKubernetesSimple,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.authentik_service_connection_kubernetes.local-cluster", "name", "Local Kubernetes Cluster"),
					resource.TestCheckResourceAttr("data.authentik_service_connection_kubernetes.local-cluster", "local", "true"),
					resource.TestCheckResourceAttrSet("data.authentik_service_connection_kubernetes.local-cluster", "verify_ssl"),
					resource.TestCheckResourceAttrSet("data.authentik_service_connection_kubernetes.local-cluster", "kubeconfig"),
				),
			},
		},
	})
}

func TestAccDataSourceOutpostServiceConnectionsKubernetes_NotFound(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { helpers.TestAccPreCheck(t) },
		ProviderFactories: provider.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceOutpostServiceConnectionKubernetesNotFound,
				ExpectError: regexp.MustCompile(`No Kubernetes Outpost Service Connections found`),
			},
		},
	})
}

const testAccDataSourceOutpostServiceConnectionKubernetesSimple = `
resource "authentik_service_connection_kubernetes" "local" {
  name = "Local Kubernetes Cluster"
  local = true
}

data "authentik_service_connection_kubernetes" "local-cluster" {
  name = authentik_service_connection_kubernetes.local.name
  local = true
}
`

const testAccDataSourceOutpostServiceConnectionKubernetesNotFound = `
data "authentik_service_connection_kubernetes" "missing" {
  name = "definitely-does-not-exist"
}
`
