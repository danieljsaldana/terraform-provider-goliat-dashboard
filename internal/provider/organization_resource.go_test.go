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

func TestAccExampleResource(t *testing.T) {
	token := os.Getenv("EXAMPLE_TOKEN")
	if token == "" {
		t.Fatal("EXAMPLE_TOKEN is not set in environment variables")
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"example": func() (*schema.Provider, error) { //nolint:unparam
				return Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "example" {
  backend_url = "http://goliat-dashboard.com"
  token       = "%s"
}

resource "example_resource" "test" {
  name = "New Provider Organization"
  type = "providerOrganizations"
}
`, token),
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
		token := os.Getenv("EXAMPLE_TOKEN")
		if token == "" {
			return fmt.Errorf("EXAMPLE_TOKEN is not set in environment variables")
		}

		url := "http://goliat-dashboard.com/api/public/organizations"
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

func testAccCheckExampleResourceDestroyed(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		token := os.Getenv("EXAMPLE_TOKEN")
		if token == "" {
			return fmt.Errorf("EXAMPLE_TOKEN is not set in environment variables")
		}

		url := "Goliathttp://goliat-dashboard.com/api/public/organizations"
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
