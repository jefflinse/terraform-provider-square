package client

import (
	"fmt"

	catalogAPI "github.com/jefflinse/square-connect/client/catalog"
	squaremodel "github.com/jefflinse/square-connect/models"
)

const (
	// CategoryObjectType designates an object that describes a CatalogCategory.
	CategoryObjectType = "CATEGORY"

	// ItemObjectType designates an object that describes a CatalogItem.
	ItemObjectType = "ITEM"
)

// CatalogCategory is a CatalogObject with type CATEGORY.
type CatalogCategory struct {
	ID   string
	Name string

	version int64
}

// CatalogItem is a CatalogObject with type ITEM.
type CatalogItem struct {
	ID          string
	Name        string
	Description string
	CategoryID  string

	version int64
}

// CreateCatalogCategory creates a new catalog category.
func (s *Square) CreateCatalogCategory(name string) (string, error) {
	params := catalogAPI.NewUpsertCatalogObjectParams().WithBody(&squaremodel.UpsertCatalogObjectRequest{
		IdempotencyKey: newIdempotencyKey(),
		Object: &squaremodel.CatalogObject{
			ID:   newTempID(),
			Type: strPtr(CategoryObjectType),
			CategoryData: &squaremodel.CatalogCategory{
				Name: name,
			},
		},
	})

	resp, err := s.square.Catalog.UpsertCatalogObject(params, s.auth())
	if err != nil {
		return "", fmt.Errorf("failed to create Square catalog category: %w", err)
	}

	return *resp.Payload.CatalogObject.ID, nil
}

// RetrieveCatalogCategory retrieves a catalog category.
func (s *Square) RetrieveCatalogCategory(id string) (*CatalogCategory, error) {
	params := catalogAPI.NewRetrieveCatalogObjectParams().WithObjectID(id)
	resp, err := s.square.Catalog.RetrieveCatalogObject(params, s.auth())
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve Square catalog category: %w", err)
	}

	return &CatalogCategory{
		ID:      *resp.Payload.Object.ID,
		Name:    *&resp.Payload.Object.CategoryData.Name,
		version: resp.Payload.Object.Version,
	}, nil
}

// UpdateCatalogCategory updates a catalog category.
func (s *Square) UpdateCatalogCategory(id string, name string) (string, error) {
	category, err := s.RetrieveCatalogCategory(id)
	if err != nil {
		return "", err
	}

	params := catalogAPI.NewUpsertCatalogObjectParams().WithBody(&squaremodel.UpsertCatalogObjectRequest{
		IdempotencyKey: newIdempotencyKey(),
		Object: &squaremodel.CatalogObject{
			ID:   &category.ID,
			Type: strPtr(CategoryObjectType),
			CategoryData: &squaremodel.CatalogCategory{
				Name: name,
			},
			Version: category.version,
		},
	})

	resp, err := s.square.Catalog.UpsertCatalogObject(params, s.auth())
	if err != nil {
		return "", fmt.Errorf("failed to update Square catalog category: %w", err)
	}

	return *resp.Payload.CatalogObject.ID, nil
}

// DeleteCatalogCategory deletes a catalog category.
func (s *Square) DeleteCatalogCategory(id string) (string, error) {
	params := catalogAPI.NewDeleteCatalogObjectParams().WithObjectID(id)
	resp, err := s.square.Catalog.DeleteCatalogObject(params, s.auth())
	if err != nil {
		return "", fmt.Errorf("failed to delete Square catalog category: %w", err)
	}

	return resp.Payload.DeletedObjectIds[0], nil
}

// CreateCatalogItem creates a new catalog category.
func (s *Square) CreateCatalogItem(item *CatalogItem) (*CatalogItem, error) {
	itemID := newTempID()
	params := catalogAPI.NewUpsertCatalogObjectParams().WithBody(&squaremodel.UpsertCatalogObjectRequest{
		IdempotencyKey: newIdempotencyKey(),
		Object: &squaremodel.CatalogObject{
			ID:   itemID,
			Type: strPtr(ItemObjectType),
			ItemData: &squaremodel.CatalogItem{
				Name:        item.Name,
				Description: item.Description,
				CategoryID:  item.CategoryID,
			},
		},
	})

	resp, err := s.square.Catalog.UpsertCatalogObject(params, s.auth())
	if err != nil {
		return nil, fmt.Errorf("failed to create Square catalog item: %w", err)
	}

	return &CatalogItem{
		ID:          *resp.Payload.CatalogObject.ID,
		Name:        resp.Payload.CatalogObject.ItemData.Name,
		Description: resp.Payload.CatalogObject.ItemData.Description,
		CategoryID:  resp.Payload.CatalogObject.ItemData.CategoryID,
		version:     resp.Payload.CatalogObject.Version,
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
