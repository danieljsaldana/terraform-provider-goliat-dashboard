---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "Goliat Provider"
subcategory: ""
description: |-
  The Goliat Provider allows you to manage organizational resources in your infrastructure.

---

# Goliat Provider

## Example Usage

```terraform
provider "Goliat" {
  backend_url = "http://localhost:4321"
  token       = "123"
}

resource "organization_resource" "example" {
  name = "New Provider Organization"
  type = "providerOrganizations"
}