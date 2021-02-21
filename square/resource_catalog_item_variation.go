package square

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	squaremodel "github.com/jefflinse/square-connect/models"
	"github.com/jefflinse/terraform-provider-square/square/client"
)

const (
	// ItemVariationObjectType designates an object that describes a CatalogItemVariation.
	ItemVariationObjectType = "ITEM_VARIATION"

	// PricingTypeFixed designated a CatalogItemVariation with fixed pricing.
	PricingTypeFixed = "FIXED_PRICING"

	// PricingTypeVariable designated a CatalogItemVariation with variable pricing.
	PricingTypeVariable = "VARIABLE_PRICING"
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
	created, err := meta.(*client.Client).UpsertCatalogObject(&squaremodel.CatalogObject{
		ID:                newTempID(),
		Type:              strPtr("ITEM_VARIATION"),
		ItemVariationData: createCatalogItemVariation(d),
	})
	if err != nil {
		return fmt.Errorf("create catalog item: %w", err)
	}

	d.SetId(*created.ID)

	return resourceSquareCatalogItemVariationRead(d, meta)
}

func resourceSquareCatalogItemVariationRead(d *schema.ResourceData, meta interface{}) error {
	obj, err := meta.(*client.Client).RetrieveCatalogObject(d.Id())
	if err != nil {
		return err
	}

	return readCatalogItemVariation(obj.ItemVariationData, d)
}

func resourceSquareCatalogItemVariationUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("item_id") ||
		d.HasChange("name") ||
		d.HasChange("pricing_type") ||
		d.HasChange("price") ||
		d.HasChange("currency") ||
		d.HasChange("sku") ||
		d.HasChange("upc") {

		client := meta.(*client.Client)
		obj, err := client.RetrieveCatalogObject(d.Id())
		if err != nil {
			return err
		}

		if _, err := client.UpsertCatalogObject(&squaremodel.CatalogObject{
			ID:                strPtr(*obj.ID),
			Type:              strPtr("ITEM_VARIATION"),
			Version:           obj.Version,
			ItemVariationData: createCatalogItemVariation(d),
		}); err != nil {
			return err
		}
	}

	return nil
}

func resourceSquareCatalogItemVariationDelete(d *schema.ResourceData, meta interface{}) error {
	_, err := meta.(*client.Client).DeleteCatalogObject(d.Id())
	return err
}

func createCatalogItemVariation(d *schema.ResourceData) *squaremodel.CatalogItemVariation {
	itemVariation := &squaremodel.CatalogItemVariation{
		ItemID:      d.Get("item_id").(string),
		Name:        d.Get("name").(string),
		PricingType: d.Get("pricing_type").(string),
		Sku:         d.Get("sku").(string),
		Upc:         d.Get("upc").(string),
	}

	if itemVariation.PricingType == PricingTypeFixed {
		itemVariation.PriceMoney = &squaremodel.Money{
			Amount:   int64(d.Get("price").(int)),
			Currency: d.Get("currency").(string),
		}
	}

	return itemVariation
}

func readCatalogItemVariation(itemVariation *squaremodel.CatalogItemVariation, d *schema.ResourceData) error {
	d.Set("item_id", itemVariation.ItemID)
	d.Set("name", itemVariation.Name)
	d.Set("pricing_type", itemVariation.PricingType)
	d.Set("sku", itemVariation.Sku)
	d.Set("upc", itemVariation.Upc)

	if itemVariation.PricingType == PricingTypeFixed {
		d.Set("price", itemVariation.PriceMoney.Amount)
		d.Set("currency", itemVariation.PriceMoney.Currency)
	}

	return nil
}
