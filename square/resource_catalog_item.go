package square

import (
	"github.com/hashicorp/terraform/helper/schema"
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
	item := CatalogItem{
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
	item.TaxIDs = []string{}
	for _, tid := range taxIDs {
		item.TaxIDs = append(item.TaxIDs, tid.(string))
	}

	square := meta.(*Client)
	created, err := square.CreateCatalogItem(&item)
	if err != nil {
		return err
	}

	d.SetId(created.ID)

	return resourceSquareCatalogItemRead(d, meta)
}

func resourceSquareCatalogItemRead(d *schema.ResourceData, meta interface{}) error {
	square := meta.(*Client)
	c, err := square.RetrieveCatalogItem(d.Id())
	if err != nil {
		return err
	}

	d.Set("abbreviation", c.Abbreviation)
	d.Set("available_electronically", c.AvailableElectronically)
	d.Set("available_for_pickup", c.AvailableForPickup)
	d.Set("available_online", c.AvailableOnline)
	d.Set("category_id", c.CategoryID)
	d.Set("description", c.Description)
	d.Set("label_colorr", c.LabelColor)
	d.Set("name", c.Name)
	d.Set("skip_modifier_screen", c.SkipModifierScreen)
	d.Set("tax_ids", c.TaxIDs)

	return nil
}

func resourceSquareCatalogItemUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("abbreviation") ||
		d.HasChange("available_electronically") ||
		d.HasChange("available_for_pickup") ||
		d.HasChange("available_online") ||
		d.HasChange("category_id") ||
		d.HasChange("description") ||
		d.HasChange("label_colorr") ||
		d.HasChange("name") ||
		d.HasChange("skip_modifier_screen") ||
		d.HasChange("tax_ids") {
		square := meta.(*Client)
		item := CatalogItem{
			ID:                      d.Id(),
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
		item.TaxIDs = []string{}
		for _, tid := range taxIDs {
			item.TaxIDs = append(item.TaxIDs, tid.(string))
		}

		_, err := square.UpdateCatalogItem(&item)
		if err != nil {
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
