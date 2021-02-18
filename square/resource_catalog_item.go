package square

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jefflinse/terraform-provider-square/square/client"
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
		},
		Create: resourceSquareCatalogItemCreate,
		Read:   resourceSquareCatalogItemRead,
		Update: resourceSquareCatalogItemUpdate,
		Delete: resourceSquareCatalogItemDelete,
	}
}

func resourceSquareCatalogItemCreate(d *schema.ResourceData, meta interface{}) error {
	item := client.CatalogItem{
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

	square := meta.(*client.Square)
	created, err := square.CreateCatalogItem(&item)
	if err != nil {
		return err
	}

	d.SetId(created.ID)

	return resourceSquareCatalogItemRead(d, meta)
}

func resourceSquareCatalogItemRead(d *schema.ResourceData, meta interface{}) error {
	square := meta.(*client.Square)
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

	return nil
}

func resourceSquareCatalogItemUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("name") || d.HasChange("description") || d.HasChange("category_id") {
		square := meta.(*client.Square)
		item := client.CatalogItem{
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
		_, err := square.UpdateCatalogItem(&item)
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceSquareCatalogItemDelete(d *schema.ResourceData, meta interface{}) error {
	square := meta.(*client.Square)
	_, err := square.DeleteCatalogItem(d.Id())
	if err != nil {
		return err
	}

	return nil
}
