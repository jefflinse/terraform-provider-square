package square

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceSquareCatalogItemVariation() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"currency": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"item_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"price": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"pricing_type": {
				Type:     schema.TypeString,
				Required: true,
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
	itemVariation := CatalogItemVariation{
		ItemID:      d.Get("item_id").(string),
		Name:        d.Get("name").(string),
		PricingType: d.Get("pricing_type").(string),
		SKU:         d.Get("sku").(string),
		UPC:         d.Get("upc").(string),
	}

	if itemVariation.PricingType == PricingTypeFixed {
		itemVariation.Price = int64(d.Get("price").(int))
		itemVariation.Currency = d.Get("currency").(string)
	}

	square := meta.(*Client)
	created, err := square.CreateCatalogItemVariation(&itemVariation)
	if err != nil {
		return err
	}

	d.SetId(created.ID)

	return resourceSquareCatalogItemVariationRead(d, meta)
}

func resourceSquareCatalogItemVariationRead(d *schema.ResourceData, meta interface{}) error {
	square := meta.(*Client)
	itemVariation, err := square.RetrieveCatalogItemVariation(d.Id())
	if err != nil {
		return err
	}

	d.Set("item_id", itemVariation.ItemID)
	d.Set("name", itemVariation.Name)
	d.Set("pricing_type", itemVariation.PricingType)
	d.Set("sku", itemVariation.SKU)
	d.Set("upc", itemVariation.UPC)

	if itemVariation.PricingType == PricingTypeFixed {
		d.Set("price", itemVariation.Price)
		d.Set("currency", itemVariation.Currency)
	}

	return nil
}

func resourceSquareCatalogItemVariationUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("item_id") ||
		d.HasChange("name") ||
		d.HasChange("pricing_type") ||
		d.HasChange("price") ||
		d.HasChange("currency") ||
		d.HasChange("sku") ||
		d.HasChange("upc") {
		square := meta.(*Client)

		itemVariation := CatalogItemVariation{
			ID:          d.Id(),
			ItemID:      d.Get("item_id").(string),
			Name:        d.Get("name").(string),
			PricingType: d.Get("pricing_type").(string),
			SKU:         d.Get("sku").(string),
			UPC:         d.Get("upc").(string),
		}

		if itemVariation.PricingType == PricingTypeFixed {
			itemVariation.Price = int64(d.Get("price").(int))
			itemVariation.Currency = d.Get("currency").(string)
		}

		_, err := square.UpdateCatalogItemVariation(&itemVariation)
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceSquareCatalogItemVariationDelete(d *schema.ResourceData, meta interface{}) error {
	square := meta.(*Client)
	_, err := square.DeleteCatalogObject(d.Id())
	return err
}
