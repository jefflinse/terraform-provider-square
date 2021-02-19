package client

import (
	"fmt"

	squaremodel "github.com/jefflinse/square-connect/models"
)

const (
	// TaxObjectType designates an object that describes a CatalogTax.
	TaxObjectType = "TAX"

	// TaxPhaseSubtotal indicates the fee is calculated based on the payment's subtotal.
	TaxPhaseSubtotal = "TAX_SUBTOTAL_PHASE"

	// TaxPhaseTotal indicates the fee is calculated based on the payment's total.
	TaxPhaseTotal = "TAX_TOTAL_PHASE"
)

// CatalogTax is a CatalogObject with type TAX.
type CatalogTax struct {
	ID                     string
	AppliesToCustomAmounts bool
	CalculationPhase       string
	Enabled                bool
	InclusionType          string
	Name                   string
	Percentage             string

	version int64
}

func taxFromObjectData(obj *squaremodel.CatalogObject) *CatalogTax {
	return &CatalogTax{
		ID:                     *obj.ID,
		AppliesToCustomAmounts: obj.TaxData.AppliesToCustomAmounts,
		CalculationPhase:       obj.TaxData.CalculationPhase,
		Enabled:                obj.TaxData.Enabled,
		InclusionType:          obj.TaxData.InclusionType,
		Name:                   obj.TaxData.Name,
		Percentage:             obj.TaxData.Percentage,

		version: obj.Version,
	}
}

func taxDataFromTax(tax *CatalogTax) *squaremodel.CatalogTax {
	return &squaremodel.CatalogTax{
		AppliesToCustomAmounts: tax.AppliesToCustomAmounts,
		CalculationPhase:       tax.CalculationPhase,
		Enabled:                tax.Enabled,
		InclusionType:          tax.InclusionType,
		Name:                   tax.Name,
		Percentage:             tax.Percentage,
	}
}

// CreateCatalogTax creates a new catalog tax.
func (s *Square) CreateCatalogTax(tax *CatalogTax) (*CatalogTax, error) {
	taxID := newTempID()
	created, err := s.upsertCatalogObject(&squaremodel.CatalogObject{
		ID:      &taxID,
		Type:    strPtr(TaxObjectType),
		TaxData: taxDataFromTax(tax),
	})
	if err != nil {
		return nil, fmt.Errorf("create catalog tax: %w", err)
	}

	return taxFromObjectData(created), nil
}

// RetrieveCatalogTax retrieves a catalog tax.
func (s *Square) RetrieveCatalogTax(id string) (*CatalogTax, error) {
	found, err := s.retrieveCatalogObject(id)
	if err != nil {
		return nil, fmt.Errorf("retrieve catalog tax: %w", err)
	}

	return taxFromObjectData(found), nil
}

// UpdateCatalogTax updates a catalog tax.
func (s *Square) UpdateCatalogTax(tax *CatalogTax) (*CatalogTax, error) {
	found, err := s.RetrieveCatalogTax(tax.ID)
	if err != nil {
		return nil, err
	}

	updated, err := s.upsertCatalogObject(&squaremodel.CatalogObject{
		ID:      &found.ID,
		Type:    strPtr(TaxObjectType),
		TaxData: taxDataFromTax(tax),
		Version: found.version,
	})
	if err != nil {
		return nil, fmt.Errorf("update catalog tax: %w", err)
	}

	return taxFromObjectData(updated), nil
}
