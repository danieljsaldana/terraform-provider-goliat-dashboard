// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestAccExampleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"example": func() (*schema.Provider, error) {
				return Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: `
provider "example" {
  backend_url = "http://localhost:4321"
  token       = "123"
}

resource "example_resource" "test" {
  name = "New Provider Organization"
  type = "providerOrganizations"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("example_resource.test", "name", "New Provider Organization"),
					testAccCheckExampleResourceCreated("example_resource.test"),
				),
			},
			{
				ImportState:  true,
				Destroy:      true,
				ResourceName: "example_resource.test",
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExampleResourceDestroyed("example_resource.test"),
				),
			},
		},
	})
}

func testAccCheckExampleResourceCreated(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		url := "http://localhost:4321/api/public/organizations"

		time.Sleep(2 * time.Second)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return fmt.Errorf("error al crear la solicitud GET: %s", err)
		}
		req.Header.Set("Authorization", "Bearer 123")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("error al realizar la solicitud GET: %s", err)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error al leer la respuesta: %s", err)
		}

		var result struct {
			ProviderOrganizations []Organization `json:"ProviderOrganizations"`
		}
		if err := json.Unmarshal(body, &result); err != nil {
			return fmt.Errorf("error al deserializar la respuesta JSON: %s", err)
		}

		for _, org := range result.ProviderOrganizations {
			if org.ID == "new_provider_org" && org.Name == "New Provider Organization" {
				fmt.Println("Organización encontrada en ProviderOrganizations.")
				return nil
			}
		}

		return fmt.Errorf("no se encontró la organización con ID: new_provider_org en ProviderOrganizations")
	}
}

func testAccCheckExampleResourceDestroyed(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		url := "http://localhost:4321/api/public/organizations"
		resp, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("error al realizar la solicitud GET: %s", err)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error al leer la respuesta: %s", err)
		}

		var result struct {
			ProviderOrganizations []Organization `json:"ProviderOrganizations"`
		}
		if err := json.Unmarshal(body, &result); err != nil {
			return fmt.Errorf("error al deserializar la respuesta JSON: %s", err)
		}

		for _, org := range result.ProviderOrganizations {
			if org.ID == "new_provider_org" {
				return fmt.Errorf("la organización aún existe en ProviderOrganizations con ID: new_provider_org")
			}
		}

		return nil
	}
}
