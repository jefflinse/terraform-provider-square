package square

import (
	"fmt"

	catalogAPI "github.com/jefflinse/square-connect/client/catalog"
	squaremodel "github.com/jefflinse/square-connect/models"
)

// CategoryObjectType designates an object that describes a CatalogCategory.
const CategoryObjectType = "CATEGORY"

// CatalogCategory is a CatalogObject with type CATEGORY.
type CatalogCategory struct {
	ID   string
	Name string

	version int64
}

// CreateCatalogCategory creates a new catalog category.
func (c *Client) CreateCatalogCategory(name string) (string, error) {
	categoryID := newTempID()
	params := catalogAPI.NewUpsertCatalogObjectParams().WithBody(&squaremodel.UpsertCatalogObjectRequest{
		IdempotencyKey: newIdempotencyKey(),
		Object: &squaremodel.CatalogObject{
			ID:   &categoryID,
			Type: strPtr(CategoryObjectType),
			CategoryData: &squaremodel.CatalogCategory{
				Name: name,
			},
		},
	})

	resp, err := c.square.Catalog.UpsertCatalogObject(params, c.auth())
	if err != nil {
		return "", fmt.Errorf("create catalog category: %w", err)
	}

	return *resp.Payload.CatalogObject.ID, nil
}

// RetrieveCatalogCategory retrieves a catalog category.
func (c *Client) RetrieveCatalogCategory(id string) (*CatalogCategory, error) {
	params := catalogAPI.NewRetrieveCatalogObjectParams().WithObjectID(id)
	resp, err := c.square.Catalog.RetrieveCatalogObject(params, c.auth())
	if err != nil {
		return nil, fmt.Errorf("retrieve catalog category: %w", err)
	}

	return &CatalogCategory{
		ID:      *resp.Payload.Object.ID,
		Name:    resp.Payload.Object.CategoryData.Name,
		version: resp.Payload.Object.Version,
	}, nil
}

// UpdateCatalogCategory updates a catalog category.
func (c *Client) UpdateCatalogCategory(id string, name string) (string, error) {
	category, err := c.RetrieveCatalogCategory(id)
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

	resp, err := c.square.Catalog.UpsertCatalogObject(params, c.auth())
	if err != nil {
		return "", fmt.Errorf("update catalog category: %w", err)
	}

	return *resp.Payload.CatalogObject.ID, nil
}