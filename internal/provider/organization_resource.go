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

func resourceExample() *schema.Resource {
	return &schema.Resource{
		Create: resourceExampleCreate,
		Read:   resourceExampleRead,
		Delete: resourceExampleDelete,
		Importer: &schema.ResourceImporter{
			State: resourceExampleImport,
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

func resourceExampleCreate(d *schema.ResourceData, meta interface{}) error {
	config, ok := meta.(*Config)
	if !ok {
		return fmt.Errorf("error al convertir meta a *Config")
	}

	name, ok := d.Get("name").(string)
	if !ok {
		return fmt.Errorf("error al convertir el nombre a string")
	}
	typeVal, ok := d.Get("type").(string)
	if !ok {
		return fmt.Errorf("error al convertir el tipo a string")
	}

	payload := Organization{
		ID:   "new_provider_org",
		Name: name,
		Type: typeVal,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error al convertir los datos a JSON: %s", err)
	}

	url := fmt.Sprintf("%s/api/public/organizations", config.BackendURL)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error al crear la solicitud HTTP: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error al enviar la solicitud al backend: %s", err)
	}
	defer resp.Body.Close()

	responseBody, _ := io.ReadAll(resp.Body)
	fmt.Printf("Response from backend: %s\n", string(responseBody))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("fallo en la creación, código de estado: %d, respuesta: %s", resp.StatusCode, string(responseBody))
	}

	d.SetId(payload.ID)
	fmt.Printf("ID del recurso configurado en Terraform: %s\n", d.Id())
	return nil
}

func resourceExampleRead(d *schema.ResourceData, meta interface{}) error {
	config, ok := meta.(*Config)
	if !ok {
		return fmt.Errorf("error al convertir meta a *Config")
	}

	url := fmt.Sprintf("%s/api/public/organizations", config.BackendURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error al crear la solicitud GET: %s", err)
	}
	req.Header.Set("Authorization", "Bearer "+config.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error al enviar la solicitud GET al backend: %s", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error al leer la respuesta: %s", err)
	}

	var result struct {
		ProviderOrganizations []Organization `json:"ProviderOrganizations"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("error al deserializar la respuesta JSON: %s", err)
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

func resourceExampleDelete(d *schema.ResourceData, meta interface{}) error {
	config, ok := meta.(*Config)
	if !ok {
		return fmt.Errorf("error al convertir meta a *Config")
	}
	resourceID := d.Id()
	resourceName, ok := d.Get("name").(string)
	if !ok {
		return fmt.Errorf("error al convertir el nombre a string")
	}

	payload := map[string]string{
		"id":   resourceID,
		"name": resourceName,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error al convertir los datos a JSON: %s", err)
	}

	url := fmt.Sprintf("%s/api/public/organizations", config.BackendURL)
	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error al crear la solicitud DELETE: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error al enviar la solicitud de eliminación al backend: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("fallo en la eliminación, código de estado: %d", resp.StatusCode)
	}

	d.SetId("")
	return nil
}

func resourceExampleImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	id := d.Id()
	if err := d.Set("name", id); err != nil {
		return nil, err
	}
	return nil, nil
}
