package square

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jefflinse/terraform-provider-square/square/client"
)

func resourceSquareCatalogItemVariation() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"item_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"price_amount": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"price_currency": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"pricing_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  client.PricingTypeVariable,
			},
			"sku": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"upc": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		Create: resourceSquareCatalogItemVariationCreate,
		Read:   resourceSquareCatalogItemVariationRead,
		Update: resourceSquareCatalogItemVariationUpdate,
		Delete: resourceSquareCatalogItemVariationDelete,
	}
}

func resourceSquareCatalogItemVariationCreate(d *schema.ResourceData, meta interface{}) error {
	itemVariation := client.CatalogItemVariation{
		ItemID:      d.Get("item_id").(string),
		Name:        d.Get("name").(string),
		PricingType: d.Get("pricing_type").(string),
		SKU:         d.Get("sku").(string),
		UPC:         d.Get("upc").(string),
	}

	if itemVariation.PricingType == client.PricingTypeFixed {
		itemVariation.Price = &client.Money{
			Amount:   int64(d.Get("price_amount").(int)),
			Currency: d.Get("price_currency").(string),
		}
	}

	square := meta.(*client.Square)
	created, err := square.CreateCatalogItemVariation(&itemVariation)
	if err != nil {
		return err
	}

	d.SetId(created.ID)

	return resourceSquareCatalogItemVariationRead(d, meta)
}

func resourceSquareCatalogItemVariationRead(d *schema.ResourceData, meta interface{}) error {
	square := meta.(*client.Square)
	itemVariation, err := square.RetrieveCatalogItemVariation(d.Id())
	if err != nil {
		return err
	}

	d.Set("item_id", itemVariation.ItemID)
	d.Set("name", itemVariation.Name)
	d.Set("pricing_type", itemVariation.PricingType)
	d.Set("sku", itemVariation.PricingType)
	d.Set("upc", itemVariation.PricingType)

	if itemVariation.PricingType == client.PricingTypeFixed {
		d.Set("price_amount", itemVariation.Price.Amount)
		d.Set("price_currency", itemVariation.Price.Currency)
	}

	return nil
}

func resourceSquareCatalogItemVariationUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("item_id") ||
		d.HasChange("name") ||
		d.HasChange("pricing_type") ||
		d.HasChange("sku") ||
		d.HasChange("upc") ||
		d.HasChange("price_amount") ||
		d.HasChange("price_currency") {
		square := meta.(*client.Square)

		itemVariation := client.CatalogItemVariation{
			ID:          d.Id(),
			ItemID:      d.Get("item_id").(string),
			Name:        d.Get("name").(string),
			PricingType: d.Get("pricing_type").(string),
			SKU:         d.Get("sku").(string),
			UPC:         d.Get("upc").(string),
		}

		if itemVariation.PricingType == client.PricingTypeFixed {
			itemVariation.Price = &client.Money{
				Amount:   int64(d.Get("price_amount").(int)),
				Currency: d.Get("price_currency").(string),
			}
		}

		_, err := square.UpdateCatalogItemVariation(&itemVariation)
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceSquareCatalogItemVariationDelete(d *schema.ResourceData, meta interface{}) error {
	square := meta.(*client.Square)
	_, err := square.DeleteCatalogItemVariation(d.Id())
	if err != nil {
		return err
	}

	return nil
}
