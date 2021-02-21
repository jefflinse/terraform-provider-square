package square

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	squaremodel "github.com/jefflinse/square-connect/models"
	"github.com/jefflinse/terraform-provider-square/square/client"
)

const (
	// CategoryObjectType is the Square type for a catalog object describing a category.
	CategoryObjectType = "CATEGORY"

	// CatalogCategoryNameMaxLength is the maximum length for a category's abbreviation.
	CatalogCategoryNameMaxLength = 255
)

func resourceSquareCatalogCategory() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (wrns []string, errs []error) {
					val := v.(string)
					if len(val) > CatalogCategoryNameMaxLength {
						errs = append(errs, fmt.Errorf("category name '%s' exceeds max length of %d", val, CatalogCategoryNameMaxLength))
					}
					return
				},
			},
		},
		Create: resourceSquareCatalogCategoryCreate,
		Read:   resourceSquareCatalogCategoryRead,
		Update: resourceSquareCatalogCategoryUpdate,
		Delete: resourceSquareCatalogCategoryDelete,
	}
}

func resourceSquareCatalogCategoryCreate(d *schema.ResourceData, meta interface{}) error {
	created, err := meta.(client.SquareAPI).UpsertCatalogObject(&squaremodel.CatalogObject{
		ID:           newTempID(),
		Type:         strPtr(CategoryObjectType),
		CategoryData: expandCatalogCategory(d),
	})
	if err != nil {
		return err
	}

	d.SetId(*created.ID)

	return resourceSquareCatalogCategoryRead(d, meta)
}

func resourceSquareCatalogCategoryRead(d *schema.ResourceData, meta interface{}) error {
	obj, err := meta.(client.SquareAPI).RetrieveCatalogObject(d.Id())
	if err != nil {
		return err
	}

	return flattenCatalogCategory(obj.CategoryData, d)
}

func resourceSquareCatalogCategoryUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("name") {
		client := meta.(client.SquareAPI)
		obj, err := client.RetrieveCatalogObject(d.Id())
		if err != nil {
			return err
		}

		if _, err := client.UpsertCatalogObject(&squaremodel.CatalogObject{
			ID:           strPtr(*obj.ID),
			Type:         strPtr(CategoryObjectType),
			Version:      obj.Version,
			CategoryData: expandCatalogCategory(d),
		}); err != nil {
			return err
		}
	}

	return nil
}

func resourceSquareCatalogCategoryDelete(d *schema.ResourceData, meta interface{}) error {
	_, err := meta.(client.SquareAPI).DeleteCatalogObject(d.Id())
	return err
}

func expandCatalogCategory(d *schema.ResourceData) *squaremodel.CatalogCategory {
	return &squaremodel.CatalogCategory{
		Name: d.Get("name").(string),
	}
}

func flattenCatalogCategory(category *squaremodel.CatalogCategory, d *schema.ResourceData) error {
	d.Set("name", category.Name)
	return nil
}
