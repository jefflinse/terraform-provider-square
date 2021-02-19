terraform {
  required_providers {
    square = {
      version = "0.0.1"
      source  = "jefflinse.io/com/square"
    }
  }
}

provider "square" {}

resource "square_catalog_category" "test" {
  name = "My Terraformed Category"
}

resource "square_catalog_item" "test" {
  name                     = "My Terraformed Item"
  abbreviation             = "TF"
  available_online         = false
  available_for_pickup     = false
  available_electronically = false
  category_id              = square_catalog_category.test.id
  description              = "This was made with Terraform!"
  label_color              = "0000FF"
  skip_modifier_screen     = false
  tax_ids = [
    square_catalog_tax.test.id
  ]
}

resource "square_catalog_item_variation" "test" {
  name           = "My Terraformed Item Variation"
  item_id        = square_catalog_item.test.id
  pricing_type   = "FIXED_PRICING"
  price_amount   = 350
  price_currency = "USD"
  sku            = "B123PINT"
  upc            = ""
}

resource "square_catalog_tax" "test" {
	name                      = "My Terraformed Tax"
	calculation_phase         = "TAX_TOTAL_PHASE"
	inclusion_type            = "ADDITIVE"
	percentage                = "4.2"
}
