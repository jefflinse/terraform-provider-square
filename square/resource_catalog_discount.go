package square

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jefflinse/terraform-provider-square/square/client"
)

// CatalogDiscountNameMaxLength is the maximum length for a CatalogDiscount's name.
const CatalogDiscountNameMaxLength = 255

func resourceSquareCatalogDiscount() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"amount": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"currency": {
				Type:     schema.TypeString,
				Optional: true,
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
			},
			"percentage": {
				Type:     schema.TypeString,
				Optional: true,
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
	discount := client.CatalogDiscount{
		Amount:         int64(d.Get("amount").(int)),
		Currency:       d.Get("currency").(string),
		LabelColor:     d.Get("label_color").(string),
		ModifyTaxBasis: d.Get("modify_tax_basis").(string),
		Name:           d.Get("name").(string),
		Percentage:     d.Get("percentage").(string),
		PinRequired:    d.Get("pin_required").(bool),
		Type:           d.Get("type").(string),
	}

	square := meta.(*client.Square)
	created, err := square.CreateCatalogDiscount(&discount)
	if err != nil {
		return err
	}

	d.SetId(created.ID)

	return resourceSquareCatalogDiscountRead(d, meta)
}

func resourceSquareCatalogDiscountRead(d *schema.ResourceData, meta interface{}) error {
	square := meta.(*client.Square)
	discount, err := square.RetrieveCatalogDiscount(d.Id())
	if err != nil {
		return err
	}

	d.Set("amount", discount.Amount)
	d.Set("currency", discount.Currency)
	d.Set("label_color", discount.LabelColor)
	d.Set("modify_tax_basis", discount.ModifyTaxBasis)
	d.Set("name", discount.Name)
	d.Set("percentage", discount.Percentage)
	d.Set("pin_required", discount.PinRequired)
	d.Set("type", discount.Type)

	return nil
}

func resourceSquareCatalogDiscountUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("amount") ||
		d.HasChange("currency") ||
		d.HasChange("label_color") ||
		d.HasChange("modify_tax_basis") ||
		d.HasChange("name") ||
		d.HasChange("percentage") ||
		d.HasChange("ping_required") ||
		d.HasChange("type") {
		square := meta.(*client.Square)
		discount := client.CatalogDiscount{
			ID:             d.Id(),
			Amount:         int64(d.Get("amount").(int)),
			Currency:       d.Get("currency").(string),
			LabelColor:     d.Get("label_color").(string),
			ModifyTaxBasis: d.Get("modify_tax_basis").(string),
			Name:           d.Get("name").(string),
			Percentage:     d.Get("percentage").(string),
			PinRequired:    d.Get("pin_required").(bool),
			Type:           d.Get("type").(string),
		}
		_, err := square.UpdateCatalogDiscount(&discount)
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceSquareCatalogDiscountDelete(d *schema.ResourceData, meta interface{}) error {
	square := meta.(*client.Square)
	_, err := square.DeleteCatalogObject(d.Id())
	return err
}
