// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	assert.True(t, ok, "The provider schema should contain 'backend_url'")

	_, ok = schema["token"]
	assert.True(t, ok, "The provider schema should contain 'token'")
}
