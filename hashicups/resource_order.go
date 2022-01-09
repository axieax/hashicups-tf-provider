package hashicups

import (
	"context"
	"strconv"
	"time"

	hc "github.com/hashicorp-demoapp/hashicups-client-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Define the Resource for an Order
func resourceOrder() *schema.Resource {
	return &schema.Resource{
		// Define functionality for the CRUD operations
		CreateContext: resourceOrderCreate,
		ReadContext:   resourceOrderRead,
		UpdateContext: resourceOrderUpdate,
		DeleteContext: resourceOrderDelete,

		// Define Schema for the Resource
		Schema: map[string]*schema.Schema{
			// Keeps track of last updated time to support the Update operation
			"last_updated": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// Define the Schema for Items in an Order
			"items": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					// Nested map using a schema.TypeList with 1 item - chosen because it closely matches the coffee object returned in the response
					// ALT: schema.TypeMap (preferred if you only require a key value map of primitive types, but requires a validation function to enforce required keys)
					Schema: map[string]*schema.Schema{
						"coffee": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									// WARN: below type definition is redundant since the Schema struct defines
									// the types to be keys of string's and values of *schema.Schema's
									// "id": &schema.Schema{
									"id": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"teaser": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"price": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"image": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"quantity": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
		},

		// Support `terraform import` Functionality
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// Function defining the Create operation
// NOTE: parameter m represents meta, which contains the HashiCups API client set by the ConfigureContextFunc (provider.go)
func resourceOrderCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the HashiCups API client from the meta interface
	c := m.(*hc.Client)

	// Warning or errors can be collected in a slice type (dynamic array)
	var diags diag.Diagnostics

	// NOTE: var.(T) is a type assertion
	items := d.Get("items").([]interface{})
	ois := constructOrderItems(items)

	// Client creates the specified Order
	o, err := c.CreateOrder(ois)
	if err != nil {
		return diag.FromErr(err)
	}

	// Sets the resource ID to the order ID
	d.SetId(strconv.Itoa(o.ID))

	// Updates terraform state after resource creation
	resourceOrderRead(ctx, d, m)

	return diags
}

// Function defining the Read operation
func resourceOrderRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*hc.Client)
	var diags diag.Diagnostics

	// Get the order ID provided by the user (Terraform)
	orderID := d.Id()

	// Client retrieves the specified Order
	order, err := c.GetOrder(orderID)
	if err != nil {
		return diag.FromErr(err)
	}

	orderItems := flattenOrderItems(&order.Items)
	if err := d.Set("items", orderItems); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

// Function defining the Update operation
func resourceOrderUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*hc.Client)

	orderID := d.Id()

	// The `hasChange()` function enables you to invoke different APIs or build a
	// targeted request body to update your resource when a specific property changes
	if d.HasChange("items") {
		items := d.Get("items").([]interface{})
		ois := constructOrderItems(items)

		_, err := c.UpdateOrder(orderID, ois)
		if err != nil {
			return diag.FromErr(err)
		}

		// Update last_updated field to current time
		d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	return resourceOrderRead(ctx, d, m)
}

// Function defining the Delete operation
// NOTE: this destroy callback should never update any state on the resource
// if resource already destroyed, this should not return an error
// if the target API doesn't have this functionality, the destroy function should verify the resource exists
func resourceOrderDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	c := m.(*hc.Client)
	var diags diag.Diagnostics

	orderID := d.Id()

	err := c.DeleteOrder(orderID)
	if err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	// if any errors exist, provider assumes resource still exists and all prior state is preserved
	// or else if no errors are returned, provider assumes resource is destroyed and all state is removed
	return diags
}

// Helper function to create OrderItems from maps
func constructOrderItems(items []interface{}) []hc.OrderItem {
	ois := []hc.OrderItem{}

	for _, item := range items {
		i := item.(map[string]interface{})

		// extract coffee from the nested map structure defined above (list with 1 element)
		co := i["coffee"].([]interface{})[0]
		coffee := co.(map[string]interface{})

		oi := hc.OrderItem{
			Coffee: hc.Coffee{
				ID: coffee["id"].(int),
			},
			Quantity: i["quantity"].(int),
		}

		ois = append(ois, oi)
	}

	return ois
}

// Helper function to flatten the OrderItems into a list of maps
func flattenOrderItems(orderItems *[]hc.OrderItem) []interface{} {
	if orderItems == nil {
		return make([]interface{}, 0)
	}

	// `make` creates slices, maps or channels, taking in type, length and capacity
	// unlike `new` which returns a pointer to the type, this returns the same type
	// it seems that capacity can be omitted in this case, since it's the same as length for slices
	ois := make([]interface{}, len(*orderItems), len(*orderItems))

	for i, orderItem := range *orderItems {
		oi := make(map[string]interface{})

		oi["coffee"] = flattenCoffee(orderItem.Coffee)
		oi["quantity"] = orderItem.Quantity
		ois[i] = oi
	}

	return ois
}

// Helper function to flatten the Coffee into a map
func flattenCoffee(coffee hc.Coffee) []interface{} {
	c := make(map[string]interface{})
	c["id"] = coffee.ID
	c["name"] = coffee.Name
	c["teaser"] = coffee.Teaser
	c["description"] = coffee.Description
	c["price"] = coffee.Price
	c["image"] = coffee.Image

	return []interface{}{c}
}
