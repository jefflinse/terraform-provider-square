package square

import (
	catalogAPI "github.com/jefflinse/square-connect/client/catalog"
	squaremodel "github.com/jefflinse/square-connect/models"
)

func (c *Client) retrieveCatalogObject(id string) (*squaremodel.CatalogObject, error) {
	params := catalogAPI.NewRetrieveCatalogObjectParams().WithObjectID(id)
	resp, err := c.square.Catalog.RetrieveCatalogObject(params, c.auth())
	if err != nil {
		return nil, err
	}

	return resp.Payload.Object, nil
}

func (c *Client) upsertCatalogObject(obj *squaremodel.CatalogObject) (*squaremodel.CatalogObject, error) {
	params := catalogAPI.NewUpsertCatalogObjectParams().WithBody(&squaremodel.UpsertCatalogObjectRequest{
		IdempotencyKey: newIdempotencyKey(),
		Object:         obj,
	})

	resp, err := c.square.Catalog.UpsertCatalogObject(params, c.auth())
	if err != nil {
		return nil, err
	}

	return resp.Payload.CatalogObject, nil
}

// DeleteCatalogObject deletes a catalog object with the specified ID.
func (c *Client) DeleteCatalogObject(id string) (string, error) {
	params := catalogAPI.NewDeleteCatalogObjectParams().WithObjectID(id)
	resp, err := c.square.Catalog.DeleteCatalogObject(params, c.auth())
	if err != nil {
		return "", err
	}

	return resp.Payload.DeletedObjectIds[0], nil
}
