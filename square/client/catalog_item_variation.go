package client

import (
	"fmt"

	squaremodel "github.com/jefflinse/square-connect/models"
)

const (
	// ItemVariationObjectType designates an object that describes a CatalogItem.
	ItemVariationObjectType = "ITEM_VARIATION"

	// PricingTypeFixed designated a CatalogItemVariation with fixed pricing.
	PricingTypeFixed = "FIXED_PRICING"

	// PricingTypeVariable designated a CatalogItemVariation with variable pricing.
	PricingTypeVariable = "VARIABLE_PRICING"
)

// CatalogItemVariation is a CatalogObject with type ITEM_VARIATION.
type CatalogItemVariation struct {
	ID          string
	ItemID      string
	Name        string
	Price       int64
	Currency    string
	PricingType string
	SKU         string
	UPC         string

	version int64
}

func itemVariationFromObjectData(obj *squaremodel.CatalogObject) *CatalogItemVariation {
	iv := &CatalogItemVariation{
		ID:          *obj.ID,
		ItemID:      obj.ItemVariationData.ItemID,
		Name:        obj.ItemVariationData.Name,
		PricingType: obj.ItemVariationData.PricingType,
		SKU:         obj.ItemVariationData.Sku,
		UPC:         obj.ItemVariationData.Upc,

		version: obj.Version,
	}

	if iv.PricingType == PricingTypeFixed {
		iv.Price = obj.ItemVariationData.PriceMoney.Amount
		iv.Currency = obj.ItemVariationData.PriceMoney.Currency
	}

	return iv
}

func itemVariationDataFromItemVariation(item *CatalogItemVariation) *squaremodel.CatalogItemVariation {
	iv := &squaremodel.CatalogItemVariation{
		ItemID:      item.ItemID,
		Name:        item.Name,
		PricingType: item.PricingType,
		Sku:         item.SKU,
		Upc:         item.UPC,
	}

	if item.PricingType == PricingTypeFixed {
		iv.PriceMoney = &squaremodel.Money{
			Amount:   item.Price,
			Currency: item.Currency,
		}
	}

	return iv
}

// CreateCatalogItemVariation creates a new catalog category.
func (s *Square) CreateCatalogItemVariation(itemVariation *CatalogItemVariation) (*CatalogItemVariation, error) {
	itemID := newTempID()
	created, err := s.upsertCatalogObject(&squaremodel.CatalogObject{
		ID:                &itemID,
		Type:              strPtr(ItemVariationObjectType),
		ItemVariationData: itemVariationDataFromItemVariation(itemVariation),
	})
	if err != nil {
		return nil, fmt.Errorf("create catalog item variation: %w", err)
	}

	return itemVariationFromObjectData(created), nil
}

// RetrieveCatalogItemVariation retrieves a catalog item.
func (s *Square) RetrieveCatalogItemVariation(id string) (*CatalogItemVariation, error) {
	found, err := s.retrieveCatalogObject(id)
	if err != nil {
		return nil, fmt.Errorf("retrieve catalog item variation: %w", err)
	}

	return itemVariationFromObjectData(found), nil
}

// UpdateCatalogItemVariation updates a catalog item.
func (s *Square) UpdateCatalogItemVariation(itemVariation *CatalogItemVariation) (*CatalogItemVariation, error) {
	found, err := s.RetrieveCatalogItemVariation(itemVariation.ID)
	if err != nil {
		return nil, fmt.Errorf("update catalog item variation: %w", err)
	}

	updated, err := s.upsertCatalogObject(&squaremodel.CatalogObject{
		ID:                &found.ID,
		Type:              strPtr(ItemObjectType),
		ItemVariationData: itemVariationDataFromItemVariation(itemVariation),
		Version:           found.version,
	})
	if err != nil {
		return nil, fmt.Errorf("update catalog item variation: %w", err)
	}

	return itemVariationFromObjectData(updated), nil
}
