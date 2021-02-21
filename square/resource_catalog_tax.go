package square

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	squaremodel "github.com/jefflinse/square-connect/models"
	"github.com/jefflinse/terraform-provider-square/square/client"
)

const (
	// CatalogTaxNameMaxLength is the maximum length for a CatalogTax's name.
	CatalogTaxNameMaxLength = 255

	// TaxObjectType designates an object that describes a CatalogTax.
	TaxObjectType = "TAX"

	// TaxPhaseSubtotal indicates the fee is calculated based on the payment's subtotal.
	TaxPhaseSubtotal = "TAX_SUBTOTAL_PHASE"

	// TaxPhaseTotal indicates the fee is calculated based on the payment's total.
	TaxPhaseTotal = "TAX_TOTAL_PHASE"
)

func resourceSquareCatalogTax() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"applies_to_custom_amounts": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"calculation_phase": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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
	taxID := newTempID()
	created, err := meta.(*client.Client).UpsertCatalogObject(&squaremodel.CatalogObject{
		ID:      &taxID,
		Type:    strPtr("TAX"),
		TaxData: createCatalogTax(d),
	})
	if err != nil {
		return fmt.Errorf("create catalog tax: %w", err)
	}

	d.SetId(*created.ID)

	return resourceSquareCatalogTaxRead(d, meta)
}

func resourceSquareCatalogTaxRead(d *schema.ResourceData, meta interface{}) error {
	obj, err := meta.(*client.Client).RetrieveCatalogObject(d.Id())
	if err != nil {
		return err
	}

	return readCatalogTax(obj.TaxData, d)
}

func resourceSquareCatalogTaxUpdate(d *schema.ResourceData, meta interface{}) error {
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
			ID:      strPtr(*obj.ID),
			Type:    strPtr("TAX"),
			Version: obj.Version,
			TaxData: createCatalogTax(d),
		}); err != nil {
			return err
		}
	}

	return nil
}

func resourceSquareCatalogTaxDelete(d *schema.ResourceData, meta interface{}) error {
	_, err := meta.(*client.Client).DeleteCatalogObject(d.Id())
	return err
}

func createCatalogTax(d *schema.ResourceData) *squaremodel.CatalogTax {
	tax := &squaremodel.CatalogTax{
		AppliesToCustomAmounts: d.Get("applies_to_custom_amounts").(bool),
		CalculationPhase:       d.Get("calculation_phase").(string),
		Enabled:                d.Get("enabled").(bool),
		InclusionType:          d.Get("inclusion_type").(string),
		Name:                   d.Get("name").(string),
		Percentage:             d.Get("percentage").(string),
	}

	return tax
}

func readCatalogTax(tax *squaremodel.CatalogTax, d *schema.ResourceData) error {
	d.Set("applies_to_custom_amounts", tax.AppliesToCustomAmounts)
	d.Set("calculation_phase", tax.CalculationPhase)
	d.Set("enabled", tax.Enabled)
	d.Set("inclusion_type", tax.InclusionType)
	d.Set("name", tax.Name)
	d.Set("percentage", tax.Percentage)

	return nil
}
