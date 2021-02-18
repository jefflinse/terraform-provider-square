package square

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jefflinse/terraform-provider-square/square/client"
)

func resourceSquareCatalogItemVariation() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"price": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"item_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Create: resourceSquareCatalogItemVariationCreate,
		Read:   resourceSquareCatalogItemVariationRead,
		Update: resourceSquareCatalogItemVariationUpdate,
		Delete: resourceSquareCatalogItemVariationDelete,
	}
}

func resourceSquareCatalogItemVariationCreate(d *schema.ResourceData, meta interface{}) error {
	itemVariation := client.CatalogItemVariation{
		Name:   d.Get("name").(string),
		Price:  int64(d.Get("price").(int)),
		ItemID: d.Get("item_id").(string),
	}

	square := meta.(*client.Square)
	created, err := square.CreateCatalogItemVariation(&itemVariation)
	if err != nil {
		return err
	}

	d.SetId(created.ID)

	return resourceSquareCatalogItemVariationRead(d, meta)
}

func resourceSquareCatalogItemVariationRead(d *schema.ResourceData, meta interface{}) error {
	square := meta.(*client.Square)
	c, err := square.RetrieveCatalogItemVariation(d.Id())
	if err != nil {
		return err
	}

	d.Set("name", c.Name)
	d.Set("price", c.Price)
	d.Set("item_id", c.ItemID)

	return nil
}

func resourceSquareCatalogItemVariationUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("name") || d.HasChange("price") || d.HasChange("item_id") {
		square := meta.(*client.Square)
		itemVariation := client.CatalogItemVariation{
			ID:     d.Id(),
			Name:   d.Get("name").(string),
			Price:  d.Get("price").(int64),
			ItemID: d.Get("item_id").(string),
		}
		_, err := square.UpdateCatalogItemVariation(&itemVariation)
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceSquareCatalogItemVariationDelete(d *schema.ResourceData, meta interface{}) error {
	square := meta.(*client.Square)
	_, err := square.DeleteCatalogItemVariation(d.Id())
	if err != nil {
		return err
	}

	return nil
}
