package square

import (
	"fmt"

	squaremodel "github.com/jefflinse/square-connect/models"
)

const (
	// DiscountObjectType designates an object that describes a CatalogDiscount.
	DiscountObjectType = "DISCOUNT"

	// DiscountTypeFixedPercentage applies the discount as a fixed percentage (e.g., 5%) off the item price.
	DiscountTypeFixedPercentage = "FIXED_PERCENTAGE"

	// DiscountTypeFixedAmount applies the discount as a fixed amount (e.g., $1.00) off the item price.
	DiscountTypeFixedAmount = "FIXED_AMOUNT"

	// DiscountTypeVariablePercentage applies the discount as a variable percentage off the item price.
	// The percentage will be specified at the time of sale.
	DiscountTypeVariablePercentage = "VARIABLE_PERCENTAGE"
)

// CatalogDiscount is a CatalogObject with type DISCOUNT.
type CatalogDiscount struct {
	ID             string
	Amount         int64
	Currency       string
	Type           string
	LabelColor     string
	ModifyTaxBasis string
	Name           string
	Percentage     string
	PinRequired    bool

	version int64
}

func discountFromObjectData(obj *squaremodel.CatalogObject) *CatalogDiscount {
	discount := &CatalogDiscount{
		ID:             *obj.ID,
		Type:           obj.DiscountData.DiscountType,
		LabelColor:     obj.DiscountData.LabelColor,
		ModifyTaxBasis: obj.DiscountData.ModifyTaxBasis,
		Name:           obj.DiscountData.Name,
		Percentage:     obj.DiscountData.Percentage,
		PinRequired:    obj.DiscountData.PinRequired,

		version: obj.Version,
	}

	if discount.Type == DiscountTypeFixedAmount {
		discount.Amount = obj.DiscountData.AmountMoney.Amount
		discount.Currency = obj.DiscountData.AmountMoney.Currency
	}

	return discount
}

func discountDataFromDiscount(discount *CatalogDiscount) *squaremodel.CatalogDiscount {
	obj := &squaremodel.CatalogDiscount{
		DiscountType:   discount.Type,
		LabelColor:     discount.LabelColor,
		ModifyTaxBasis: discount.ModifyTaxBasis,
		Name:           discount.Name,
		Percentage:     discount.Percentage,
		PinRequired:    discount.PinRequired,
	}

	if discount.Type == DiscountTypeFixedAmount {
		obj.AmountMoney = &squaremodel.Money{
			Amount:   discount.Amount,
			Currency: discount.Currency,
		}
	}

	return obj
}

// CreateCatalogDiscount creates a new catalog discount.
func (c *Client) CreateCatalogDiscount(discount *CatalogDiscount) (*CatalogDiscount, error) {
	discountID := newTempID()
	created, err := c.upsertCatalogObject(&squaremodel.CatalogObject{
		ID:           &discountID,
		Type:         strPtr(DiscountObjectType),
		DiscountData: discountDataFromDiscount(discount),
	})
	if err != nil {
		return nil, fmt.Errorf("create catalog discount: %w", err)
	}

	return discountFromObjectData(created), nil
}

// RetrieveCatalogDiscount retrieves a catalog discount.
func (c *Client) RetrieveCatalogDiscount(id string) (*CatalogDiscount, error) {
	found, err := c.retrieveCatalogObject(id)
	if err != nil {
		return nil, fmt.Errorf("retrieve catalog discount: %w", err)
	}

	return discountFromObjectData(found), nil
}

// UpdateCatalogDiscount updates a catalog discount.
func (c *Client) UpdateCatalogDiscount(discount *CatalogDiscount) (*CatalogDiscount, error) {
	found, err := c.RetrieveCatalogDiscount(discount.ID)
	if err != nil {
		return nil, err
	}

	updated, err := c.upsertCatalogObject(&squaremodel.CatalogObject{
		ID:           &found.ID,
		Type:         strPtr(DiscountObjectType),
		DiscountData: discountDataFromDiscount(discount),
		Version:      found.version,
	})
	if err != nil {
		return nil, fmt.Errorf("update catalog discount: %w", err)
	}

	return discountFromObjectData(updated), nil
}
