package square

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	squaremodel "github.com/jefflinse/square-connect/models"
	"github.com/jefflinse/terraform-provider-square/square/client"
)

const (
	// CatalogTaxNameMaxLength is the maximum length for a tax's name.
	CatalogTaxNameMaxLength = 255

	// TaxObjectType is the Square type for a catalog object describing a tax.
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
				ValidateFunc: func(v interface{}, k string) (wrns []string, errs []error) {
					val := v.(string)
					if len(val) > CatalogTaxNameMaxLength {
						errs = append(errs, fmt.Errorf("tax name '%s' exceeds max length of %d", val, CatalogTaxNameMaxLength))
					}
					return
				},
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
	created, err := meta.(client.SquareAPI).UpsertCatalogObject(&squaremodel.CatalogObject{
		ID:      newTempID(),
		Type:    strPtr(TaxObjectType),
		TaxData: expandCatalogTax(d),
	})
	if err != nil {
		return err
	}

	d.SetId(*created.ID)

	return resourceSquareCatalogTaxRead(d, meta)
}

func resourceSquareCatalogTaxRead(d *schema.ResourceData, meta interface{}) error {
	obj, err := meta.(client.SquareAPI).RetrieveCatalogObject(d.Id())
	if err != nil {
		return err
	}

	return flattenCatalogTax(obj.TaxData, d)
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

		client := meta.(client.SquareAPI)
		obj, err := client.RetrieveCatalogObject(d.Id())
		if err != nil {
			return err
		}

		if _, err := client.UpsertCatalogObject(&squaremodel.CatalogObject{
			ID:      strPtr(*obj.ID),
			Type:    strPtr(TaxObjectType),
			Version: obj.Version,
			TaxData: expandCatalogTax(d),
		}); err != nil {
			return err
		}
	}

	return nil
}

func resourceSquareCatalogTaxDelete(d *schema.ResourceData, meta interface{}) error {
	_, err := meta.(client.SquareAPI).DeleteCatalogObject(d.Id())
	return err
}

func expandCatalogTax(d *schema.ResourceData) *squaremodel.CatalogTax {
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

func flattenCatalogTax(tax *squaremodel.CatalogTax, d *schema.ResourceData) error {
	d.Set("applies_to_custom_amounts", tax.AppliesToCustomAmounts)
	d.Set("calculation_phase", tax.CalculationPhase)
	d.Set("enabled", tax.Enabled)
	d.Set("inclusion_type", tax.InclusionType)
	d.Set("name", tax.Name)
	d.Set("percentage", tax.Percentage)

	return nil
}
