package square

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	squaremodel "github.com/jefflinse/square-connect/models"
)

// CatalogItemAbbreviationMaxLength is the maximum length for a CatalogItem's abbreviation.
const CatalogItemAbbreviationMaxLength = 24

func resourceSquareCatalogItem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"abbreviation": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"available_electronically": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"available_for_pickup": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"available_online": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"category_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"label_color": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"skip_modifier_screen": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"tax_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		Create: resourceSquareCatalogItemCreate,
		Read:   resourceSquareCatalogItemRead,
		Update: resourceSquareCatalogItemUpdate,
		Delete: resourceSquareCatalogItemDelete,
	}
}

func resourceSquareCatalogItemCreate(d *schema.ResourceData, meta interface{}) error {
	itemID := newTempID()
	item := squaremodel.CatalogItem{
		Abbreviation:            d.Get("abbreviation").(string),
		AvailableElectronically: d.Get("available_electronically").(bool),
		AvailableForPickup:      d.Get("available_for_pickup").(bool),
		AvailableOnline:         d.Get("available_online").(bool),
		CategoryID:              d.Get("category_id").(string),
		Description:             d.Get("description").(string),
		LabelColor:              d.Get("label_color").(string),
		Name:                    d.Get("name").(string),
		SkipModifierScreen:      d.Get("skip_modifier_screen").(bool),
	}

	taxIDs := d.Get("tax_ids").([]interface{})
	item.TaxIds = []string{}
	for _, tid := range taxIDs {
		item.TaxIds = append(item.TaxIds, tid.(string))
	}

	client := meta.(*Client)
	created, err := client.upsertCatalogObject(&squaremodel.CatalogObject{
		ID:       &itemID,
		Type:     strPtr("ITEM"),
		ItemData: &item,
	})
	if err != nil {
		return fmt.Errorf("create catalog item: %w", err)
	}

	d.SetId(*created.ID)

	return resourceSquareCatalogItemRead(d, meta)
}

func resourceSquareCatalogItemRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	obj, err := client.retrieveCatalogObject(d.Id())
	if err != nil {
		return err
	}

	d.Set("abbreviation", obj.ItemData.Abbreviation)
	d.Set("available_electronically", obj.ItemData.AvailableElectronically)
	d.Set("available_for_pickup", obj.ItemData.AvailableForPickup)
	d.Set("available_online", obj.ItemData.AvailableOnline)
	d.Set("category_id", obj.ItemData.CategoryID)
	d.Set("description", obj.ItemData.Description)
	d.Set("label_colorr", obj.ItemData.LabelColor)
	d.Set("name", obj.ItemData.Name)
	d.Set("skip_modifier_screen", obj.ItemData.SkipModifierScreen)
	d.Set("tax_ids", obj.ItemData.TaxIds)

	return nil
}

func resourceSquareCatalogItemUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("abbreviation") ||
		d.HasChange("available_electronically") ||
		d.HasChange("available_for_pickup") ||
		d.HasChange("available_online") ||
		d.HasChange("category_id") ||
		d.HasChange("description") ||
		d.HasChange("label_color") ||
		d.HasChange("name") ||
		d.HasChange("skip_modifier_screen") ||
		d.HasChange("tax_ids") {

		client := meta.(*Client)
		found, err := client.retrieveCatalogObject(d.Id())
		if err != nil {
			return err
		}

		item := squaremodel.CatalogItem{
			Abbreviation:            d.Get("abbreviation").(string),
			AvailableElectronically: d.Get("available_electronically").(bool),
			AvailableForPickup:      d.Get("available_for_pickup").(bool),
			AvailableOnline:         d.Get("available_online").(bool),
			CategoryID:              d.Get("category_id").(string),
			Description:             d.Get("description").(string),
			LabelColor:              d.Get("label_colorr").(string),
			Name:                    d.Get("name").(string),
			SkipModifierScreen:      d.Get("skip_modifier_screen").(bool),
		}

		taxIDs := d.Get("tax_ids").([]interface{})
		item.TaxIds = []string{}
		for _, tid := range taxIDs {
			item.TaxIds = append(item.TaxIds, tid.(string))
		}

		if _, err := client.upsertCatalogObject(&squaremodel.CatalogObject{
			ID:       strPtr(*found.ID),
			Type:     strPtr("ITEM"),
			Version:  found.Version,
			ItemData: &item,
		}); err != nil {
			return err
		}
	}

	return nil
}

func resourceSquareCatalogItemDelete(d *schema.ResourceData, meta interface{}) error {
	square := meta.(*Client)
	_, err := square.DeleteCatalogObject(d.Id())
	return err
}
