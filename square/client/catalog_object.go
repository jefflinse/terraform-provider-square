package client

import (
	catalogAPI "github.com/jefflinse/square-connect/client/catalog"
	squaremodel "github.com/jefflinse/square-connect/models"
)

func (s *Square) retrieveCatalogObject(id string) (*squaremodel.CatalogObject, error) {
	params := catalogAPI.NewRetrieveCatalogObjectParams().WithObjectID(id)
	resp, err := s.square.Catalog.RetrieveCatalogObject(params, s.auth())
	if err != nil {
		return nil, err
	}

	return resp.Payload.Object, nil
}

func (s *Square) upsertCatalogObject(obj *squaremodel.CatalogObject) (*squaremodel.CatalogObject, error) {
	params := catalogAPI.NewUpsertCatalogObjectParams().WithBody(&squaremodel.UpsertCatalogObjectRequest{
		IdempotencyKey: newIdempotencyKey(),
		Object:         obj,
	})

	resp, err := s.square.Catalog.UpsertCatalogObject(params, s.auth())
	if err != nil {
		return nil, err
	}

	return resp.Payload.CatalogObject, nil
}
