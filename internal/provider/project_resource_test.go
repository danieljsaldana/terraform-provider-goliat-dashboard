// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccProjectResource(t *testing.T) {
	token := os.Getenv("EXAMPLE_TOKEN")
	if token == "" {
		t.Fatal("EXAMPLE_TOKEN is not set in environment variables")
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"goliatdashboard": func() (*schema.Provider, error) { //nolint:unparam
				return Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: `
provider "goliatdashboard" {
  backend_url = "https://demo.goliat-dashboard.com"
  token       = "` + token + `"
}

resource "goliatdashboard_organization" "test_org" {
  name = "new_provider_org"
  type = "providerOrganizations"
}

resource "goliatdashboard_project" "test" {
  organization = goliatdashboard_organization.test_org.name
  name         = "Test Project"
  description  = "Initial description"
  depends_on   = [goliatdashboard_organization.test_org]
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("goliatdashboard_project.test", "name", "Test Project"),
					resource.TestCheckResourceAttr("goliatdashboard_project.test", "description", "Initial description"),
				),
			},
			{
				Config: `
provider "goliatdashboard" {
  backend_url = "https://demo.goliat-dashboard.com"
  token       = "` + token + `"
}

resource "goliatdashboard_organization" "test_org" {
  name = "new_provider_org"
  type = "providerOrganizations"
}

resource "goliatdashboard_project" "test" {
  organization = goliatdashboard_organization.test_org.name
  name         = "Test Project"
  description  = "Updated description"
  depends_on   = [goliatdashboard_organization.test_org]
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("goliatdashboard_project.test", "description", "Updated description"),
				),
			},
			{
				ResourceName: "goliatdashboard_project.test",
				ImportState:  true,
				Destroy:      true,
			},
		},
	})
}
