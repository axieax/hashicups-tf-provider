package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"terraform-provider-hashicups/hashicups"
)

func main() {
	// Serve the provider defined in provider.go
	plugin.Serve(&plugin.ServeOpts{
		// QUESTION: is the below valid?
		// ProviderFunc: hashicups.Provider,
		ProviderFunc: func() *schema.Provider {
			return hashicups.Provider()
		},
	})
}
