package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/shztki/terraform-provider-nifcloud/nifcloud"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: nifcloud.Provider})
}
