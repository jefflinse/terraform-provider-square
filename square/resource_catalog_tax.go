package square

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jefflinse/terraform-provider-square/square/client"
)

// CatalogTaxNameMaxLength is the maximum length for a CatalogTax's name.
const CatalogTaxNameMaxLength = 255

func resourceSquareCatalogTax() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"applies_to_custom_amounts": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"calculation_phase": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"inclusion_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"percentage": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Create: resourceSquareCatalogTaxCreate,
		Read:   resourceSquareCatalogTaxRead,
		Update: resourceSquareCatalogTaxUpdate,
		Delete: resourceSquareCatalogTaxDelete,
	}
}

func resourceSquareCatalogTaxCreate(d *schema.ResourceData, meta interface{}) error {
	tax := client.CatalogTax{
		AppliesToCustomAmounts: d.Get("applies_to_custom_amounts").(bool),
		CalculationPhase:       d.Get("calculation_phase").(string),
		Enabled:                d.Get("enabled").(bool),
		InclusionType:          d.Get("inclusion_type").(string),
		Name:                   d.Get("name").(string),
		Percentage:             d.Get("percentage").(string),
	}

	square := meta.(*client.Square)
	created, err := square.CreateCatalogTax(&tax)
	if err != nil {
		return err
	}

	d.SetId(created.ID)

	return resourceSquareCatalogTaxRead(d, meta)
}

func resourceSquareCatalogTaxRead(d *schema.ResourceData, meta interface{}) error {
	square := meta.(*client.Square)
	t, err := square.RetrieveCatalogTax(d.Id())
	if err != nil {
		return err
	}

	d.Set("applies_to_custom_amounts", t.AppliesToCustomAmounts)
	d.Set("calculation_phase", t.CalculationPhase)
	d.Set("enabled", t.Enabled)
	d.Set("inclusion_type", t.InclusionType)
	d.Set("name", t.Name)
	d.Set("percentage", t.Percentage)

	return nil
}

func resourceSquareCatalogTaxUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("applies_to_custom_amounts") ||
		d.HasChange("calculation_phase") ||
		d.HasChange("enabled") ||
		d.HasChange("inclusion_type") ||
		d.HasChange("name") ||
		d.HasChange("percentage") {
		square := meta.(*client.Square)
		tax := client.CatalogTax{
			ID:                     d.Id(),
			AppliesToCustomAmounts: d.Get("applies_to_custom_amounts").(bool),
			CalculationPhase:       d.Get("calculation_phase").(string),
			Enabled:                d.Get("enabled").(bool),
			InclusionType:          d.Get("inclusion_type").(string),
			Name:                   d.Get("name").(string),
			Percentage:             d.Get("percentage").(string),
		}
		_, err := square.UpdateCatalogTax(&tax)
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceSquareCatalogTaxDelete(d *schema.ResourceData, meta interface{}) error {
	square := meta.(*client.Square)
	_, err := square.DeleteCatalogObject(d.Id())
	return err
}
