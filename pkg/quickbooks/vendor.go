// ----------------------------------------------------------------------------
//
// AP Transaction
//
// Author: William Shaffer
// Version: 24-May-2024
//
// Copyright (c) 2024 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

package quickbooks

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	"github.com/waysys/assert/assert"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// Vendor is a person or organization that receives scholarship payments.
type Vendor struct {
	name string
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// Factory Functions
// ----------------------------------------------------------------------------

// NewVendor returns a Vendor structure with the specified name.
func NewVendor(name string) Vendor {
	assert.Assert(name != "", "vendor name must not be an empty string")
	var vendor = Vendor{
		name: name,
	}
	return vendor
}

// ----------------------------------------------------------------------------
// Properties
// ----------------------------------------------------------------------------

// Name returns the name of the Vendor
func (vendor Vendor) Name() string {
	return vendor.name
}
