package square

import (
	"os"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/jefflinse/terraform-provider-square/square/client"
)

// Provider returns the ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			// "access_token": {
			// 	Type:        schema.TypeString,
			// 	Description: "Your Square OAuth Access Token",
			// 	Required:    true,
			// 	// DefaultFunc: schema.EnvDefaultFunc("SQUARE_API_ACCESS_TOKEN", nil),
			// },
		},
		ResourcesMap: map[string]*schema.Resource{
			"square_catalog_category": resourceSquareCatalogCatagory(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			// "todoist_project": dataSourceTodoistProject(),
		},
		ConfigureFunc: configureFn(),
	}
}

func configureFn() func(*schema.ResourceData) (interface{}, error) {
	return func(d *schema.ResourceData) (interface{}, error) {
		// c := client.NewClient(d.Get("access_token").(string))
		c := client.NewClient(os.Getenv("SQUARE_API_ACCESS_TOKEN"))
		return c, nil
	}
}
