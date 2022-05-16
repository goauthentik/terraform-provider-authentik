package main

import (
	"flag"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"goauthentik.io/terraform-provider-authentik/internal/provider"
)

var (
	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary
	version string = "dev"
)

//go:generate terraform fmt -recursive ./examples/

//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

func main() {
	var debugMode bool
	var versionMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.BoolVar(&versionMode, "version", false, "Show version and exit")
	flag.Parse()

	opts := &plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return provider.Provider(version, false)
		},
		Debug: debugMode,
	}

	if versionMode {
		fmt.Println(version)
		return
	}

	plugin.Serve(opts)
}
