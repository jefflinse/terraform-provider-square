package square

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	squaremodel "github.com/jefflinse/square-connect/models"
)

func resourceSquareCatalogCategory() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Create: resourceSquareCatalogCategoryCreate,
		Read:   resourceSquareCatalogCategoryRead,
		Update: resourceSquareCatalogCategoryUpdate,
		Delete: resourceSquareCatalogCategoryDelete,
	}
}

func resourceSquareCatalogCategoryCreate(d *schema.ResourceData, meta interface{}) error {
	categoryID := newTempID()
	created, err := meta.(*Client).upsertCatalogObject(&squaremodel.CatalogObject{
		ID:           &categoryID,
		Type:         strPtr("CATEGORY"),
		CategoryData: createCatalogCategory(d),
	})
	if err != nil {
		return fmt.Errorf("create catalog category: %w", err)
	}

	d.SetId(*created.ID)

	return resourceSquareCatalogCategoryRead(d, meta)
}

func resourceSquareCatalogCategoryRead(d *schema.ResourceData, meta interface{}) error {
	obj, err := meta.(*Client).retrieveCatalogObject(d.Id())
	if err != nil {
		return err
	}

	return readCatalogCategory(obj.CategoryData, d)
}

func resourceSquareCatalogCategoryUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("name") {
		client := meta.(*Client)
		obj, err := client.retrieveCatalogObject(d.Id())
		if err != nil {
			return err
		}

		if _, err := client.upsertCatalogObject(&squaremodel.CatalogObject{
			ID:           strPtr(*obj.ID),
			Type:         strPtr("CATEGOORY"),
			Version:      obj.Version,
			CategoryData: createCatalogCategory(d),
		}); err != nil {
			return err
		}
	}

	return nil
}

func resourceSquareCatalogCategoryDelete(d *schema.ResourceData, meta interface{}) error {
	_, err := meta.(*Client).DeleteCatalogObject(d.Id())
	return err
}

func createCatalogCategory(d *schema.ResourceData) *squaremodel.CatalogCategory {
	return &squaremodel.CatalogCategory{
		Name: d.Get("name").(string),
	}
}

func readCatalogCategory(category *squaremodel.CatalogCategory, d *schema.ResourceData) error {
	d.Set("name", category.Name)
	return nil
}
