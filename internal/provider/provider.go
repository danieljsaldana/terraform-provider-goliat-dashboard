// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Config struct {
	BackendURL string
	Token      string
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"backend_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"token": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"goliatdashboard_organization": resourceOrganization(),
			"goliatdashboard_project":      resourceProject(),
		},
		ConfigureFunc: configureProvider,
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	backendURL, ok := d.Get("backend_url").(string)
	if !ok {
		return nil, fmt.Errorf("backend_url must be a string")
	}
	token, ok := d.Get("token").(string)
	if !ok {
		return nil, fmt.Errorf("token must be a string")
	}
	return &Config{
		BackendURL: backendURL,
		Token:      token,
	}, nil
}
