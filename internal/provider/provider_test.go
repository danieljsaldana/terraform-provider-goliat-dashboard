// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("error: %s", err)
	}
}

func TestProvider_Schema(t *testing.T) {
	provider := Provider()
	schema := provider.Schema

	_, ok := schema["backend_url"]
	assert.True(t, ok, "El esquema del provider debe contener 'backend_url'")

	_, ok = schema["token"]
	assert.True(t, ok, "El esquema del provider debe contener 'token'")
}
