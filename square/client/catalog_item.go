package client

import (
	"fmt"

	squaremodel "github.com/jefflinse/square-connect/models"
)

// ItemObjectType designates an object that describes a CatalogItem.
const ItemObjectType = "ITEM"

// CatalogItem is a CatalogObject with type ITEM.
type CatalogItem struct {
	ID string

	Abbreviation            string
	AvailableElectronically bool
	AvailableOnline         bool
	AvailableForPickup      bool
	CategoryID              string
	Description             string
	LabelColor              string
	Name                    string
	SkipModifierScreen      bool
	TaxIDs                  []string

	version int64
}

func itemFromObjectData(obj *squaremodel.CatalogObject) *CatalogItem {
	return &CatalogItem{
		ID:                      *obj.ID,
		Abbreviation:            obj.ItemData.Abbreviation,
		AvailableElectronically: obj.ItemData.AvailableElectronically,
		AvailableForPickup:      obj.ItemData.AvailableElectronically,
		AvailableOnline:         obj.ItemData.AvailableOnline,
		CategoryID:              obj.ItemData.CategoryID,
		Description:             obj.ItemData.Description,
		LabelColor:              obj.ItemData.LabelColor,
		Name:                    obj.ItemData.Name,
		SkipModifierScreen:      obj.ItemData.SkipModifierScreen,
		TaxIDs:                  obj.ItemData.TaxIds,

		version: obj.Version,
	}
}

func itemDataFromItem(item *CatalogItem) *squaremodel.CatalogItem {
	return &squaremodel.CatalogItem{
		Abbreviation:            item.Abbreviation,
		AvailableElectronically: item.AvailableElectronically,
		AvailableForPickup:      item.AvailableElectronically,
		AvailableOnline:         item.AvailableOnline,
		CategoryID:              item.CategoryID,
		Description:             item.Description,
		LabelColor:              item.LabelColor,
		Name:                    item.Name,
		SkipModifierScreen:      item.SkipModifierScreen,
		TaxIds:                  item.TaxIDs,
	}
}

// CreateCatalogItem creates a new catalog category.
func (s *Square) CreateCatalogItem(item *CatalogItem) (*CatalogItem, error) {
	itemID := newTempID()
	created, err := s.upsertCatalogObject(&squaremodel.CatalogObject{
		ID:       &itemID,
		Type:     strPtr(ItemObjectType),
		ItemData: itemDataFromItem(item),
	})
	if err != nil {
		return nil, fmt.Errorf("create catalog item: %w", err)
	}

	return itemFromObjectData(created), nil
}

// RetrieveCatalogItem retrieves a catalog item.
func (s *Square) RetrieveCatalogItem(id string) (*CatalogItem, error) {
	found, err := s.retrieveCatalogObject(id)
	if err != nil {
		return nil, fmt.Errorf("retrieve catalog item: %w", err)
	}

	return itemFromObjectData(found), nil
}

// UpdateCatalogItem updates a catalog item.
func (s *Square) UpdateCatalogItem(item *CatalogItem) (*CatalogItem, error) {
	found, err := s.RetrieveCatalogItem(item.ID)
	if err != nil {
		return nil, fmt.Errorf("update catalog item: %w", err)
	}

	updated, err := s.upsertCatalogObject(&squaremodel.CatalogObject{
		ID:       &found.ID,
		Type:     strPtr(ItemObjectType),
		ItemData: itemDataFromItem(item),
		Version:  item.version,
	})
	if err != nil {
		return nil, fmt.Errorf("update catalog item: %w", err)
	}

	return itemFromObjectData(updated), nil
}
