package main

import (
	"flag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/spectrumwebco/terraform-provider-databricks-ovh/internal/provider"
)

var (
	version string = "dev"
	commit  string = ""
)

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{
		ProviderFunc: provider.New(version),
		Debug:        debugMode,
	}

	plugin.Serve(opts)
}
