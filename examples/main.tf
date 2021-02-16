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
  name = "My Terraformed Item"
  description = "This was made with Terraform!"
  category_id = square_catalog_category.test.id
}
