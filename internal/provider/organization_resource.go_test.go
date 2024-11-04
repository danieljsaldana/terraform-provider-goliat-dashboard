// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestAccOrganizationResource(t *testing.T) {
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
				Config: fmt.Sprintf(`
provider "goliatdashboard" {
  backend_url = "https://goliat-dashboard.com"
  token       = "%s"
}

resource "goliatdashboard_organization" "test" {
  name = "New Provider Organization"
  type = "providerOrganizations"
}
`, token),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("goliatdashboard_organization.test", "name", "New Provider Organization"),
					testAccCheckOrganizationCreated("goliatdashboard_organization.test"),
				),
			},
			{
				ImportState:  true,
				Destroy:      true,
				ResourceName: "goliatdashboard_organization.test",
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOrganizationDestroyed("goliatdashboard_organization.test"),
				),
			},
		},
	})
}

func testAccCheckOrganizationCreated(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		token := os.Getenv("EXAMPLE_TOKEN")
		if token == "" {
			return fmt.Errorf("EXAMPLE_TOKEN is not set in environment variables")
		}

		url := "https://goliat-dashboard.com/api/public/provider/organizations"
		time.Sleep(2 * time.Second)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return fmt.Errorf("error creating GET request: %s", err)
		}
		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("error sending GET request: %s", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error reading response: %s", err)
		}

		var result struct {
			ProviderOrganizations []Organization `json:"ProviderOrganizations"`
		}
		if err := json.Unmarshal(body, &result); err != nil {
			return fmt.Errorf("error unmarshalling JSON response: %s", err)
		}

		for _, org := range result.ProviderOrganizations {
			if org.ID == "new_provider_org" && org.Name == "New Provider Organization" {
				fmt.Println("Organization found in ProviderOrganizations.")
				return nil
			}
		}

		return fmt.Errorf("organization with ID: new_provider_org not found in ProviderOrganizations")
	}
}

func testAccCheckOrganizationDestroyed(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		token := os.Getenv("EXAMPLE_TOKEN")
		if token == "" {
			return fmt.Errorf("EXAMPLE_TOKEN is not set in environment variables")
		}

		url := "https://goliat-dashboard.com/api/public/provider/organizations"
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return fmt.Errorf("error creating GET request: %s", err)
		}
		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("error sending GET request: %s", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error reading response: %s", err)
		}

		var result struct {
			ProviderOrganizations []Organization `json:"ProviderOrganizations"`
		}
		if err := json.Unmarshal(body, &result); err != nil {
			return fmt.Errorf("error unmarshalling JSON response: %s", err)
		}

		for _, org := range result.ProviderOrganizations {
			if org.ID == "new_provider_org" {
				return fmt.Errorf("organization with ID: new_provider_org still exists in ProviderOrganizations")
			}
		}

		return nil
	}
}
