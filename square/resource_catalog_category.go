package square

import (
	"github.com/hashicorp/terraform/helper/schema"
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
	square := meta.(*Client)
	name := d.Get("name").(string)

	id, err := square.CreateCatalogCategory(name)
	if err != nil {
		return err
	}

	d.SetId(id)
	return resourceSquareCatalogCategoryRead(d, meta)
}

func resourceSquareCatalogCategoryRead(d *schema.ResourceData, meta interface{}) error {
	square := meta.(*Client)
	c, err := square.RetrieveCatalogCategory(d.Id())
	if err != nil {
		return err
	}

	d.Set("name", c.Name)

	return nil
}

func resourceSquareCatalogCategoryUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("name") {
		square := meta.(*Client)
		name := d.Get("name").(string)
		_, err := square.UpdateCatalogCategory(d.Id(), name)
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceSquareCatalogCategoryDelete(d *schema.ResourceData, meta interface{}) error {
	square := meta.(*Client)
	_, err := square.DeleteCatalogObject(d.Id())
	return err
}
