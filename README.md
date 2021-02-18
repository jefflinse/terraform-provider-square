A Terraform provider for the Square Connect V2 API (squareup.com)

# terraform-provider-square

Manage your Square POS catalog using Terraform:

```hcl
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

resource "square_catalog_item_variation" "test" {
  name = "My Terraformed Item Variation"
  price = 350
  item_id = square_catalog_item.test.id
}

```

To run locally:

1. Clone this repo
2. Run `make` to build and install the provider
3. Set `SQUARE_API_ACCESS_TOKEN` in your environment to your Square OAUTH access token.
4. Copy the example above into a Terraform configuration file
5. Run `terraform init`
6. Run `terraform apply` to create the example resources

## Project Status

This projects is very much in its infancy. The feature set is limited to my own original needs, but I am actively developing this provider. Contributions are absolutely welcomed.

