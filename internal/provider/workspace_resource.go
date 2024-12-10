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

type Project struct {
	ID           string `json:"id"`
	Organization string `json:"organization"`
	Name         string `json:"name"`
	Description  string `json:"description"`
}

func resourceProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectCreate,
		Read:   resourceProjectRead,
		Update: resourceProjectUpdate,
		Delete: resourceProjectDelete,
		Importer: &schema.ResourceImporter{
			State: resourceProjectImport,
		},
		Schema: map[string]*schema.Schema{
			"organization": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceProjectCreate(d *schema.ResourceData, meta interface{}) error {
	config, ok := meta.(*Config)
	if !ok {
		return fmt.Errorf("error converting meta to *Config")
	}

	organization, ok := d.Get("organization").(string)
	if !ok {
		return fmt.Errorf("organization must be a string")
	}
	name, ok := d.Get("name").(string)
	if !ok {
		return fmt.Errorf("name must be a string")
	}
	description, ok := d.Get("description").(string)
	if !ok {
		return fmt.Errorf("description must be a string")
	}

	project := Project{
		Organization: organization,
		Name:         name,
		Description:  description,
	}

	body, err := json.Marshal(project)
	if err != nil {
		return fmt.Errorf("error marshalling project to JSON: %s", err)
	}

	url := fmt.Sprintf("%s/api/public/provider/projects", config.BackendURL)
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

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("creation failed, status code: %d, response: %s", resp.StatusCode, string(responseBody))
	}

	var responseData map[string]interface{}
	if err := json.Unmarshal(responseBody, &responseData); err != nil {
		return fmt.Errorf("error unmarshalling response: %s", err)
	}

	if proj, ok := responseData["project"].(map[string]interface{}); ok {
		if id, ok := proj["id"].(string); ok {
			d.SetId(id)
		} else {
			return fmt.Errorf("project ID not found or invalid")
		}
	}

	return resourceProjectRead(d, meta)
}

func resourceProjectRead(d *schema.ResourceData, meta interface{}) error {
	config, ok := meta.(*Config)
	if !ok {
		return fmt.Errorf("error converting meta to *Config")
	}

	url := fmt.Sprintf("%s/api/public/provider/projects", config.BackendURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating GET request: %s", err)
	}
	req.Header.Set("Authorization", "Bearer "+config.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending GET request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %s", err)
	}

	var responseData map[string]interface{}
	if err := json.Unmarshal(body, &responseData); err != nil {
		return fmt.Errorf("error unmarshalling JSON response: %s", err)
	}

	projects, ok := responseData["Projects"].([]interface{})
	if !ok {
		return fmt.Errorf("projects not found in response")
	}

	for _, p := range projects {
		project, ok := p.(map[string]interface{})
		if !ok {
			continue
		}
		if id, ok := project["id"].(string); ok && id == d.Id() {
			if err := d.Set("organization", project["organization"]); err != nil {
				return fmt.Errorf("error setting organization: %s", err)
			}
			if err := d.Set("name", project["name"]); err != nil {
				return fmt.Errorf("error setting name: %s", err)
			}
			if err := d.Set("description", project["description"]); err != nil {
				return fmt.Errorf("error setting description: %s", err)
			}
			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceProjectCreate(d, meta)
}

func resourceProjectDelete(d *schema.ResourceData, meta interface{}) error {
	config, ok := meta.(*Config)
	if !ok {
		return fmt.Errorf("error converting meta to *Config")
	}

	id := d.Id()
	if id == "" {
		return fmt.Errorf("id is not set")
	}

	organization, ok := d.Get("organization").(string)
	if !ok {
		return fmt.Errorf("organization must be a string")
	}

	payload := map[string]string{
		"id":           id,
		"organization": organization,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshalling payload to JSON: %s", err)
	}

	url := fmt.Sprintf("%s/api/public/provider/projects", config.BackendURL)
	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error creating DELETE request: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending DELETE request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("deletion failed, status code: %d", resp.StatusCode)
	}

	d.SetId("")
	return nil
}

func resourceProjectImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	id := d.Id()
	if err := d.Set("organization", id); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
