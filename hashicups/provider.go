package hashicups

import (
	"context"

	"github.com/hashicorp-demoapp/hashicups-client-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Define the Provider
func Provider() *schema.Provider {
	return &schema.Provider{
		// Define the Schema for the Provider (fields to be passed in from terraform)
		Schema: map[string]*schema.Schema{
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HASHICUPS_USERNAME", nil),
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("HASHICUPS_PASSWORD", nil),
			},
		},
		// Support Authentication Functionality
		ConfigureContextFunc: providerConfigure,

		// Define Resources
		ResourcesMap: map[string]*schema.Resource{
			"hashicups_order": resourceOrder(),
		},
		// Define Data Sources
		DataSourcesMap: map[string]*schema.Resource{
			"hashicups_coffees": dataSourceCoffees(),
			"hashicups_order":   dataSourceOrder(),
		},
	}
}

// Authentication Functionality
func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	// Warning or errors can be collected in a slice type (Diagnostics is defined as []Diagnostic)
	var diags diag.Diagnostics

	// Apparently warnings only appear when the provider also errors.. not sure if this is true in practice..
	/* diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Warning Message Summary",
		Detail:   "This is the detailed warning message from providerConfigure",
	}) */

	// NOTE: cannot declare untyped nil in golang, e.g. up := nil
	up := &username
	if username == "" {
		up = nil
	}

	pp := &password
	if password == "" {
		pp = nil
	}

	// API connects a new HTTP client
	c, err := hashicups.NewClient(nil, up, pp)
	if err != nil {
		// return nil, diag.FromErr(err)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create HashiCups client",
			Detail:   "Unable to auth user for authenticated HashiCups client",
		})

		return nil, diags
	}

	return c, diags
}
