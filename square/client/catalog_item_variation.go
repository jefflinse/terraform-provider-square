package client

import (
	"fmt"

	catalogAPI "github.com/jefflinse/square-connect/client/catalog"
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
	Price       *Money
	PricingType string
	SKU         string
	UPC         string

	version int64
}

// CreateCatalogItemVariation creates a new catalog category.
func (s *Square) CreateCatalogItemVariation(itemVariation *CatalogItemVariation) (*CatalogItemVariation, error) {
	itemID := newTempID()
	params := catalogAPI.NewUpsertCatalogObjectParams().WithBody(&squaremodel.UpsertCatalogObjectRequest{
		IdempotencyKey: newIdempotencyKey(),
		Object: &squaremodel.CatalogObject{
			ID:   &itemID,
			Type: strPtr(ItemVariationObjectType),
			ItemVariationData: &squaremodel.CatalogItemVariation{
				Name:        itemVariation.Name,
				ItemID:      itemVariation.ItemID,
				PricingType: itemVariation.PricingType,
				Sku:         itemVariation.SKU,
				Upc:         itemVariation.UPC,
			},
		},
	})

	if itemVariation.PricingType == PricingTypeFixed {
		params.Body.Object.ItemVariationData.PriceMoney = &squaremodel.Money{
			Amount:   itemVariation.Price.Amount,
			Currency: itemVariation.Price.Currency,
		}
	}

	resp, err := s.square.Catalog.UpsertCatalogObject(params, s.auth())
	if err != nil {
		return nil, fmt.Errorf("create catalog item variation: %w", err)
	}

	return &CatalogItemVariation{
		ID:   *resp.Payload.CatalogObject.ID,
		Name: resp.Payload.CatalogObject.ItemVariationData.Name,
		Price: &Money{
			Amount:   resp.Payload.CatalogObject.ItemVariationData.PriceMoney.Amount,
			Currency: resp.Payload.CatalogObject.ItemVariationData.PriceMoney.Currency,
		},
		ItemID: resp.Payload.CatalogObject.ItemVariationData.ItemID,

		version: resp.Payload.CatalogObject.Version,
	}, nil
}

// RetrieveCatalogItemVariation retrieves a catalog item.
func (s *Square) RetrieveCatalogItemVariation(id string) (*CatalogItemVariation, error) {
	params := catalogAPI.NewRetrieveCatalogObjectParams().WithObjectID(id)
	resp, err := s.square.Catalog.RetrieveCatalogObject(params, s.auth())
	if err != nil {
		return nil, fmt.Errorf("retrieve catalog item variation: %w", err)
	}

	return &CatalogItemVariation{
		ID:   *resp.Payload.Object.ID,
		Name: resp.Payload.Object.ItemVariationData.Name,
		Price: &Money{
			Amount:   resp.Payload.Object.ItemVariationData.PriceMoney.Amount,
			Currency: resp.Payload.Object.ItemVariationData.PriceMoney.Currency,
		},
		ItemID: resp.Payload.Object.ItemVariationData.ItemID,

		version: resp.Payload.Object.Version,
	}, nil
}

// UpdateCatalogItemVariation updates a catalog item.
func (s *Square) UpdateCatalogItemVariation(itemVariation *CatalogItemVariation) (*CatalogItemVariation, error) {
	foundItemVariation, err := s.RetrieveCatalogItemVariation(itemVariation.ID)
	if err != nil {
		return nil, fmt.Errorf("update catalog item variation: %w", err)
	}

	params := catalogAPI.NewUpsertCatalogObjectParams().WithBody(&squaremodel.UpsertCatalogObjectRequest{
		IdempotencyKey: newIdempotencyKey(),
		Object: &squaremodel.CatalogObject{
			ID:   &foundItemVariation.ID,
			Type: strPtr(ItemObjectType),
			ItemVariationData: &squaremodel.CatalogItemVariation{
				Name:        itemVariation.Name,
				ItemID:      itemVariation.ItemID,
				PricingType: itemVariation.PricingType,
			},
			Version: foundItemVariation.version,
		},
	})

	if itemVariation.PricingType == PricingTypeFixed {
		params.Body.Object.ItemVariationData.PriceMoney = &squaremodel.Money{
			Amount:   itemVariation.Price.Amount,
			Currency: itemVariation.Price.Currency,
		}
	}

	resp, err := s.square.Catalog.UpsertCatalogObject(params, s.auth())
	if err != nil {
		return nil, fmt.Errorf("update catalog item variation: %w", err)
	}

	return &CatalogItemVariation{
		ID:   *resp.Payload.CatalogObject.ID,
		Name: resp.Payload.CatalogObject.ItemVariationData.Name,
		Price: &Money{
			Amount:   resp.Payload.CatalogObject.ItemVariationData.PriceMoney.Amount,
			Currency: resp.Payload.CatalogObject.ItemVariationData.PriceMoney.Currency,
		},
		ItemID: resp.Payload.CatalogObject.ItemVariationData.ItemID,

		version: resp.Payload.CatalogObject.Version,
	}, nil
}

// DeleteCatalogItemVariation deletes a catalog item.
func (s *Square) DeleteCatalogItemVariation(id string) (string, error) {
	params := catalogAPI.NewDeleteCatalogObjectParams().WithObjectID(id)
	resp, err := s.square.Catalog.DeleteCatalogObject(params, s.auth())
	if err != nil {
		return "", fmt.Errorf("delete catalog item variation: %w", err)
	}

	return resp.Payload.DeletedObjectIds[0], nil
}
