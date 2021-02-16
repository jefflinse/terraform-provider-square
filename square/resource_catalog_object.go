package square

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jefflinse/terraform-provider-square/square/client"
)

func resourceSquareCatalogCatagory() *schema.Resource {
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
	client := meta.(*client.Square)
	name := d.Get("name").(string)

	id, err := client.CreateCatalogCategory(name)
	if err != nil {
		return err
	}

	d.SetId(id)
	return resourceSquareCatalogCategoryRead(d, meta)
}

func resourceSquareCatalogCategoryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.Square)
	c, err := client.RetrieveCatalogCategory(d.Id())
	if err != nil {
		return err
	}

	d.Set("name", c.Name)

	return nil
}

func resourceSquareCatalogCategoryUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("name") {
		client := meta.(*client.Square)
		name := d.Get("name").(string)
		_, err := client.UpdateCatalogCategory(d.Id(), name)
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceSquareCatalogCategoryDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.Square)
	_, err := client.DeleteCatalogCategory(d.Id())
	if err != nil {
		return err
	}

	return nil
}
