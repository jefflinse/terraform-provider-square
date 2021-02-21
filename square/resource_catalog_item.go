package square

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	squaremodel "github.com/jefflinse/square-connect/models"
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
				ValidateFunc: func(v interface{}, k string) (wrns []string, errs []error) {
					val := v.(string)
					if len(val) > CatalogItemAbbreviationMaxLength {
						errs = append(errs, fmt.Errorf("item abbreviation '%s' exceeds max length of %d", val, CatalogItemAbbreviationMaxLength))
					}
					return
				},
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
				Type:     schema.TypeSet,
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
	created, err := meta.(*client.Client).UpsertCatalogObject(&squaremodel.CatalogObject{
		ID:       newTempID(),
		Type:     strPtr("ITEM"),
		ItemData: createCatalogItem(d),
	})
	if err != nil {
		return fmt.Errorf("create catalog item: %w", err)
	}

	d.SetId(*created.ID)

	return resourceSquareCatalogItemRead(d, meta)
}

func resourceSquareCatalogItemRead(d *schema.ResourceData, meta interface{}) error {
	obj, err := meta.(*client.Client).RetrieveCatalogObject(d.Id())
	if err != nil {
		return err
	}

	return readCatalogItem(obj.ItemData, d)
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

		client := meta.(*client.Client)
		obj, err := client.RetrieveCatalogObject(d.Id())
		if err != nil {
			return err
		}

		if _, err := client.UpsertCatalogObject(&squaremodel.CatalogObject{
			ID:       strPtr(*obj.ID),
			Type:     strPtr("ITEM"),
			Version:  obj.Version,
			ItemData: createCatalogItem(d),
		}); err != nil {
			return err
		}
	}

	return nil
}

func resourceSquareCatalogItemDelete(d *schema.ResourceData, meta interface{}) error {
	_, err := meta.(*client.Client).DeleteCatalogObject(d.Id())
	return err
}

func createCatalogItem(d *schema.ResourceData) *squaremodel.CatalogItem {
	item := &squaremodel.CatalogItem{
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

	return item
}

func readCatalogItem(item *squaremodel.CatalogItem, d *schema.ResourceData) error {
	d.Set("abbreviation", item.Abbreviation)
	d.Set("available_electronically", item.AvailableElectronically)
	d.Set("available_for_pickup", item.AvailableForPickup)
	d.Set("available_online", item.AvailableOnline)
	d.Set("category_id", item.CategoryID)
	d.Set("description", item.Description)
	d.Set("label_color", item.LabelColor)
	d.Set("name", item.Name)
	d.Set("skip_modifier_screen", item.SkipModifierScreen)
	d.Set("tax_ids", item.TaxIds)

	return nil
}
