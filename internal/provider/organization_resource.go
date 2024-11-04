// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"io"
	"net/http"
)

type Organization struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

func resourceOrganization() *schema.Resource {
	return &schema.Resource{
		Create: resourceOrganizationCreate,
		Read:   resourceOrganizationRead,
		Delete: resourceOrganizationDelete,
		Importer: &schema.ResourceImporter{
			State: resourceOrganizationImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceOrganizationCreate(d *schema.ResourceData, meta interface{}) error {
	config, ok := meta.(*Config)
	if !ok {
		return fmt.Errorf("error converting meta to *Config")
	}

	name, ok := d.Get("name").(string)
	if !ok {
		return fmt.Errorf("error converting name to string")
	}
	typeVal, ok := d.Get("type").(string)
	if !ok {
		return fmt.Errorf("error converting type to string")
	}

	payload := Organization{
		ID:   "new_provider_org",
		Name: name,
		Type: typeVal,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error converting data to JSON: %s", err)
	}

	url := fmt.Sprintf("%s/api/public/provider/organizations", config.BackendURL)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request to backend: %s", err)
	}
	defer resp.Body.Close()

	responseBody, _ := io.ReadAll(resp.Body)
	fmt.Printf("Response from backend: %s\n", string(responseBody))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("creation failed, status code: %d, response: %s", resp.StatusCode, string(responseBody))
	}

	d.SetId(payload.ID)
	fmt.Printf("Resource ID set in Terraform: %s\n", d.Id())
	return nil
}

func resourceOrganizationRead(d *schema.ResourceData, meta interface{}) error {
	config, ok := meta.(*Config)
	if !ok {
		return fmt.Errorf("error converting meta to *Config")
	}

	url := fmt.Sprintf("%s/api/public/provider/organizations", config.BackendURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating GET request: %s", err)
	}
	req.Header.Set("Authorization", "Bearer "+config.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending GET request to backend: %s", err)
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

	found := false
	for _, org := range result.ProviderOrganizations {
		if org.ID == d.Id() {
			found = true
			break
		}
	}

	if !found {
		d.SetId("")
	}

	return nil
}

func resourceOrganizationDelete(d *schema.ResourceData, meta interface{}) error {
	config, ok := meta.(*Config)
	if !ok {
		return fmt.Errorf("error converting meta to *Config")
	}
	resourceID := d.Id()
	resourceName, ok := d.Get("name").(string)
	if !ok {
		return fmt.Errorf("error converting name to string")
	}

	payload := map[string]string{
		"id":   resourceID,
		"name": resourceName,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error converting data to JSON: %s", err)
	}

	url := fmt.Sprintf("%s/api/public/provider/organizations", config.BackendURL)
	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error creating DELETE request: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending DELETE request to backend: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("deletion failed, status code: %d", resp.StatusCode)
	}

	d.SetId("")
	return nil
}

func resourceOrganizationImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	id := d.Id()
	if err := d.Set("name", id); err != nil {
		return nil, err
	}
	return nil, nil
}
