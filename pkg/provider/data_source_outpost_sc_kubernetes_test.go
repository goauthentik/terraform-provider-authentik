package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceOutpostServiceConnectionsKubernetes(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
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
			{
				Config: testAccDataSourceOutpostServiceConnectionKubernetesSimple,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.authentik_service_connection_kubernetes.remote-cluster", "name", "Remote Kubernetes Cluster"),
					resource.TestCheckResourceAttr("data.authentik_service_connection_kubernetes.remote-cluster", "local", "false"),
					resource.TestCheckResourceAttrSet("data.authentik_service_connection_kubernetes.remote-cluster", "verify_ssl"),
					resource.TestCheckResourceAttrSet("data.authentik_service_connection_kubernetes.remote-cluster", "kubeconfig"),
				),
			},
		},
	})
}

func TestAccDataSourceOutpostServiceConnectionsKubernetes_NotFound(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceOutpostServiceConnectionKubernetesNotFound,
				ExpectError: regexp.MustCompile(`No Kubernetes Outpost Service Connections found`),
			},
		},
	})
}

const testAccDataSourceOutpostServiceConnectionKubernetesSimple = `
data "authentik_service_connection_kubernetes" "local-cluster" {
  name = "Local Kubernetes Cluster"
  local = true
}

data "authentik_service_connection_kubernetes" "remote-cluster" {
  name = "Remote Kubernetes Cluster"
}
`

const testAccDataSourceOutpostServiceConnectionKubernetesNotFound = `
data "authentik_service_connection_kubernetes" "missing" {
  name = "definitely-does-not-exist"
}
`
