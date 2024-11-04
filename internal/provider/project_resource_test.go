// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"os"
	"testing"
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
  backend_url = "https://goliat-dashboard.com"
  token       = "` + token + `"
}

resource "goliatdashboard_project" "test" {
  organization = "provider_organization"
  name         = "Test Project"
  description  = "Initial description"
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
  backend_url = "https://goliat-dashboard.com"
  token       = "` + token + `"
}

resource "goliatdashboard_project" "test" {
  organization = "provider_organization"
  name         = "Test Project"
  description  = "Updated description"
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
