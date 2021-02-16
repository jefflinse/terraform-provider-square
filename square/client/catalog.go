package client

import (
	"fmt"

	catalogAPI "github.com/jefflinse/square-connect/client/catalog"
	squaremodel "github.com/jefflinse/square-connect/models"
)

const (
	// CategoryObjectType designates an object that describes a CatalogCategory.
	CategoryObjectType = "CATEGORY"
)

// CatalogCategory is a CatalogObject with type CATEGORY.
type CatalogCategory struct {
	ID      string
	Name    string
	Version int64
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
		Version: resp.Payload.Object.Version,
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
			Version: category.Version,
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
