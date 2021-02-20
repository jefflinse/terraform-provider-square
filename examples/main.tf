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
  currency       = "USD"
  item_id        = square_catalog_item.test.id
  pricing_type   = "FIXED_PRICING"
  price          = 350
  sku            = "B123PINT"
  upc            = ""
}

resource "square_catalog_tax" "test" {
	name                      = "My Terraformed Tax"
  applies_to_custom_amounts = true
	calculation_phase         = "TAX_TOTAL_PHASE"
  enabled                   = true
	inclusion_type            = "ADDITIVE"
	percentage                = "4.2"
}

resource "square_catalog_discount" "test1" {
	name                      = "My Terraformed Fixed Discount"
  amount                    = 50
	currency                  = "USD"
	label_color               = "FF0000"
	modify_tax_basis          = "MODIFY_TAX_BASIS"
	pin_required              = false
	type                      = "FIXED_AMOUNT"
}

resource "square_catalog_discount" "test2" {
	name                      = "My Terraformed Variable Discount"
  percentage                = "10.0"
	label_color               = "00FF00"
	modify_tax_basis          = "MODIFY_TAX_BASIS"
	pin_required              = true
	type                      = "FIXED_PERCENTAGE"
}
