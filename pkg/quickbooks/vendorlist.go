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

type VendorList map[string]*Vendor

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

var APVendorList = make(VendorList)

// ----------------------------------------------------------------------------
// Methods
// ----------------------------------------------------------------------------

// Add creates a vendor structure with the specified name and adds it to
// the VendorList.
func (vendorList *VendorList) Add(name string) *Vendor {
	var vendor *Vendor = nil
	if vendorList.Contains(name) {
		vendor = vendorList.Get(name)
	} else {
		var v = NewVendor(name)
		vendor = &v
		(*vendorList)[name] = vendor
	}
	return vendor
}

// Contains returns true if the vendor already exists in the list.
func (vendorList *VendorList) Contains(name string) bool {
	var _, ok = (*vendorList)[name]
	return ok
}

// Get returns a pointer to a vendor in the venor list
func (vendorList *VendorList) Get(name string) *Vendor {
	assert.Assert(vendorList.Contains(name),
		"the vendor list does not contain a vendor named: "+name)
	return (*vendorList)[name]
}
