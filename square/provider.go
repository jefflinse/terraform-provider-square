package square

import (
	"fmt"
	"os"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/jefflinse/terraform-provider-square/square/client"
)

const squarerAPIAccessTokenEnvVar = "SQUARE_API_ACCESS_TOKEN"

// Provider returns the ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"square_catalog_category":       resourceSquareCatalogCategory(),
			"square_catalog_item":           resourceSquareCatalogItem(),
			"square_catalog_item_variation": resourceSquareCatalogItemVariation(),
			"square_catalog_tax":            resourceSquareCatalogTax(),
		},
		ConfigureFunc: configureFn(),
	}
}

func configureFn() func(*schema.ResourceData) (interface{}, error) {
	return func(d *schema.ResourceData) (interface{}, error) {
		token, ok := os.LookupEnv(squarerAPIAccessTokenEnvVar)
		if !ok {
			return nil, fmt.Errorf("%s not set", squarerAPIAccessTokenEnvVar)
		}

		return client.NewClient(token), nil
	}
}
