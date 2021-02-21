package square

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	squaremodel "github.com/jefflinse/square-connect/models"
	"github.com/jefflinse/terraform-provider-square/square/client"
)

const (
	// CatalogDiscountNameMaxLength is the maximum length for a discount's name.
	CatalogDiscountNameMaxLength = 255

	// DiscountObjectType is the Square type for a catalog object describing a discount.
	DiscountObjectType = "DISCOUNT"

	// DiscountTypeFixedPercentage applies the discount as a fixed percentage (e.g., 5%) off the item price.
	DiscountTypeFixedPercentage = "FIXED_PERCENTAGE"

	// DiscountTypeFixedAmount applies the discount as a fixed amount (e.g., $1.00) off the item price.
	DiscountTypeFixedAmount = "FIXED_AMOUNT"

	// DiscountTypeVariablePercentage applies the discount as a variable percentage off the item price.
	// The percentage will be specified at the time of sale.
	DiscountTypeVariablePercentage = "VARIABLE_PERCENTAGE"
)

func resourceSquareCatalogDiscount() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"amount": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"percentage"},
			},
			"currency": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"percentage"},
			},
			"label_color": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"modify_tax_basis": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (wrns []string, errs []error) {
					val := v.(string)
					if len(val) > CatalogDiscountNameMaxLength {
						errs = append(errs, fmt.Errorf("discount name '%s' exceeds max length of %d", val, CatalogDiscountNameMaxLength))
					}
					return
				},
			},
			"percentage": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"amount", "currency"},
			},
			"pin_required": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Create: resourceSquareCatalogDiscountCreate,
		Read:   resourceSquareCatalogDiscountRead,
		Update: resourceSquareCatalogDiscountUpdate,
		Delete: resourceSquareCatalogDiscountDelete,
	}
}

func resourceSquareCatalogDiscountCreate(d *schema.ResourceData, meta interface{}) error {
	created, err := meta.(client.SquareAPI).UpsertCatalogObject(&squaremodel.CatalogObject{
		ID:           newTempID(),
		Type:         strPtr(DiscountObjectType),
		DiscountData: expandCatalogDiscount(d),
	})
	if err != nil {
		return err
	}

	d.SetId(*created.ID)

	return resourceSquareCatalogDiscountRead(d, meta)
}

func resourceSquareCatalogDiscountRead(d *schema.ResourceData, meta interface{}) error {
	obj, err := meta.(client.SquareAPI).RetrieveCatalogObject(d.Id())
	if err != nil {
		return err
	}

	return flattenCatalogDiscount(obj.DiscountData, d)
}

func resourceSquareCatalogDiscountUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("amount") ||
		d.HasChange("currency") ||
		d.HasChange("label_color") ||
		d.HasChange("modify_tax_basis") ||
		d.HasChange("name") ||
		d.HasChange("percentage") ||
		d.HasChange("pin_required") ||
		d.HasChange("type") {

		client := meta.(client.SquareAPI)
		obj, err := client.RetrieveCatalogObject(d.Id())
		if err != nil {
			return err
		}

		if _, err := client.UpsertCatalogObject(&squaremodel.CatalogObject{
			ID:           strPtr(*obj.ID),
			Type:         strPtr(DiscountObjectType),
			Version:      obj.Version,
			DiscountData: expandCatalogDiscount(d),
		}); err != nil {
			return err
		}
	}

	return nil
}

func resourceSquareCatalogDiscountDelete(d *schema.ResourceData, meta interface{}) error {
	_, err := meta.(client.SquareAPI).DeleteCatalogObject(d.Id())
	return err
}

func expandCatalogDiscount(d *schema.ResourceData) *squaremodel.CatalogDiscount {
	discount := &squaremodel.CatalogDiscount{
		LabelColor:     d.Get("label_color").(string),
		ModifyTaxBasis: d.Get("modify_tax_basis").(string),
		Name:           d.Get("name").(string),
		PinRequired:    d.Get("pin_required").(bool),
		DiscountType:   d.Get("type").(string),
	}

	switch discount.DiscountType {
	case DiscountTypeFixedAmount:
		discount.AmountMoney = &squaremodel.Money{
			Amount:   int64(d.Get("amount").(int)),
			Currency: d.Get("currency").(string),
		}
	case DiscountTypeFixedPercentage:
		discount.Percentage = d.Get("percentage").(string)
	case DiscountTypeVariablePercentage:
		discount.Percentage = ""
	}

	return discount
}

func flattenCatalogDiscount(discount *squaremodel.CatalogDiscount, d *schema.ResourceData) error {
	d.Set("label_color", discount.LabelColor)
	d.Set("modify_tax_basis", discount.ModifyTaxBasis)
	d.Set("name", discount.Name)
	d.Set("percentage", discount.Percentage)
	d.Set("pin_required", discount.PinRequired)
	d.Set("type", discount.DiscountType)

	switch discount.DiscountType {
	case DiscountTypeFixedAmount:
		d.Set("amount", discount.AmountMoney.Amount)
		d.Set("currency", discount.AmountMoney.Currency)
	case DiscountTypeFixedPercentage:
		d.Set("percentage", discount.Percentage)
	case DiscountTypeVariablePercentage:
		d.Set("percentage", "")
	}

	return nil
}
