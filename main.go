package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6/tf6server"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"

	"goauthentik.io/terraform-provider-authentik/pkg/provider"
	"goauthentik.io/terraform-provider-authentik/pkg/sdkprovider"
)

// these will be set by the goreleaser configuration
// to appropriate values for the compiled binary
var version string = "dev"

//go:generate terraform fmt -recursive ./examples/

//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

func main() {
	var debugMode bool
	var versionMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.BoolVar(&versionMode, "version", false, "Show version and exit")
	flag.Parse()

	if versionMode {
		fmt.Println(version)
		return
	}

	ctx := context.Background()

	upgradedSDKServer, err := tf5to6server.UpgradeServer(ctx, sdkprovider.Provider(version, false).GRPCProvider)
	if err != nil {
		log.Fatal(err)
	}

	muxServer, err := tf6muxserver.NewMuxServer(ctx,
		func() tfprotov6.ProviderServer { return upgradedSDKServer },
		providerserver.NewProtocol6(provider.New(version, false)),
	)
	if err != nil {
		log.Fatal(err)
	}

	var serveOpts []tf6server.ServeOpt
	if debugMode {
		serveOpts = append(serveOpts, tf6server.WithManagedDebug())
	}

	serveErr := tf6server.Serve("registry.terraform.io/goauthentik/authentik", func() tfprotov6.ProviderServer {
		return muxServer
	}, serveOpts...)
	if serveErr != nil {
		log.Fatal(serveErr)
	}
}
