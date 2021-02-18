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
}

resource "square_catalog_item_variation" "test" {
  name = "My Terraformed Item Variation"
  item_id = square_catalog_item.test.id
  price = 350
}
