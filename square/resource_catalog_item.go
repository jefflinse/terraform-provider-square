package square

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jefflinse/terraform-provider-square/square/client"
)

func resourceSquareCatalogItem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"category_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		Create: resourceSquareCatalogItemCreate,
		Read:   resourceSquareCatalogItemRead,
		Update: resourceSquareCatalogItemUpdate,
		Delete: resourceSquareCatalogItemDelete,
	}
}

func resourceSquareCatalogItemCreate(d *schema.ResourceData, meta interface{}) error {
	item := client.CatalogItem{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		CategoryID:  d.Get("category_id").(string),
	}

	square := meta.(*client.Square)
	created, err := square.CreateCatalogItem(&item)
	if err != nil {
		return err
	}

	d.SetId(created.ID)

	return resourceSquareCatalogItemRead(d, meta)
}

func resourceSquareCatalogItemRead(d *schema.ResourceData, meta interface{}) error {
	square := meta.(*client.Square)
	c, err := square.RetrieveCatalogItem(d.Id())
	if err != nil {
		return err
	}

	d.Set("name", c.Name)
	d.Set("description", c.Description)
	d.Set("category_id", c.CategoryID)

	return nil
}

func resourceSquareCatalogItemUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("name") || d.HasChange("description") || d.HasChange("category_id") {
		square := meta.(*client.Square)
		item := client.CatalogItem{
			ID:          d.Id(),
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			CategoryID:  d.Get("category_id").(string),
		}
		_, err := square.UpdateCatalogItem(&item)
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceSquareCatalogItemDelete(d *schema.ResourceData, meta interface{}) error {
	square := meta.(*client.Square)
	_, err := square.DeleteCatalogItem(d.Id())
	if err != nil {
		return err
	}

	return nil
}
