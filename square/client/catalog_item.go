package client

import (
	"fmt"

	catalogAPI "github.com/jefflinse/square-connect/client/catalog"
	squaremodel "github.com/jefflinse/square-connect/models"
)

// ItemObjectType designates an object that describes a CatalogItem.
const ItemObjectType = "ITEM"

// CatalogItem is a CatalogObject with type ITEM.
type CatalogItem struct {
	ID          string
	Name        string
	Description string
	CategoryID  string

	version int64
}

// CreateCatalogItem creates a new catalog category.
func (s *Square) CreateCatalogItem(item *CatalogItem) (*CatalogItem, error) {
	itemID := newTempID()
	params := catalogAPI.NewBatchUpsertCatalogObjectsParams().WithBody(&squaremodel.BatchUpsertCatalogObjectsRequest{
		IdempotencyKey: newIdempotencyKey(),
		Batches: []*squaremodel.CatalogObjectBatch{
			{
				Objects: []*squaremodel.CatalogObject{
					{
						ID:   &itemID,
						Type: strPtr(ItemObjectType),
						ItemData: &squaremodel.CatalogItem{
							Name:        item.Name,
							Description: item.Description,
							CategoryID:  item.CategoryID,
						},
					},
				},
			},
		},
	})

	resp, err := s.square.Catalog.BatchUpsertCatalogObjects(params, s.auth())
	if err != nil {
		return nil, fmt.Errorf("failed to create Square catalog item: %w", err)
	}

	return &CatalogItem{
		ID:          *resp.Payload.Objects[0].ID,
		Name:        resp.Payload.Objects[0].ItemData.Name,
		Description: resp.Payload.Objects[0].ItemData.Description,
		CategoryID:  resp.Payload.Objects[0].ItemData.CategoryID,
		version:     resp.Payload.Objects[0].Version,
	}, nil
}

// RetrieveCatalogItem retrieves a catalog item.
func (s *Square) RetrieveCatalogItem(id string) (*CatalogItem, error) {
	params := catalogAPI.NewRetrieveCatalogObjectParams().WithObjectID(id)
	resp, err := s.square.Catalog.RetrieveCatalogObject(params, s.auth())
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve Square catalog item: %w", err)
	}

	return &CatalogItem{
		ID:          *resp.Payload.Object.ID,
		Name:        resp.Payload.Object.ItemData.Name,
		Description: resp.Payload.Object.ItemData.Description,
		CategoryID:  resp.Payload.Object.ItemData.CategoryID,
		version:     resp.Payload.Object.Version,
	}, nil
}

// UpdateCatalogItem updates a catalog item.
func (s *Square) UpdateCatalogItem(item *CatalogItem) (*CatalogItem, error) {
	foundItem, err := s.RetrieveCatalogItem(item.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to find Square catalog item: %w", err)
	}

	params := catalogAPI.NewUpsertCatalogObjectParams().WithBody(&squaremodel.UpsertCatalogObjectRequest{
		IdempotencyKey: newIdempotencyKey(),
		Object: &squaremodel.CatalogObject{
			ID:   &foundItem.ID,
			Type: strPtr(ItemObjectType),
			ItemData: &squaremodel.CatalogItem{
				Name:        item.Name,
				Description: item.Description,
				CategoryID:  item.CategoryID,
			},
			Version: foundItem.version,
		},
	})

	resp, err := s.square.Catalog.UpsertCatalogObject(params, s.auth())
	if err != nil {
		return nil, fmt.Errorf("failed to update Square catalog item: %w", err)
	}

	return &CatalogItem{
		ID:          *resp.Payload.CatalogObject.ID,
		Name:        resp.Payload.CatalogObject.ItemData.Name,
		Description: resp.Payload.CatalogObject.ItemData.Description,
		CategoryID:  resp.Payload.CatalogObject.ItemData.CategoryID,
		version:     resp.Payload.CatalogObject.Version,
	}, nil
}

// DeleteCatalogItem deletes a catalog item.
func (s *Square) DeleteCatalogItem(id string) (string, error) {
	params := catalogAPI.NewDeleteCatalogObjectParams().WithObjectID(id)
	resp, err := s.square.Catalog.DeleteCatalogObject(params, s.auth())
	if err != nil {
		return "", fmt.Errorf("failed to delete Square catalog item: %w", err)
	}

	return resp.Payload.DeletedObjectIds[0], nil
}
