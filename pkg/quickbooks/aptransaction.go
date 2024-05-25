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

// The quickbooks package contains code to read spreadsheets produced by
// Quickbooks and convert the data to arrays or maps of structures.
package quickbooks

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	d "acorn_go/pkg/date"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type APTransaction struct {
	transactionDate d.Date
	vendor          *Vendor
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// Factory Functions
// ----------------------------------------------------------------------------

func NewAPTransaction(date d.Date, vendorName string) APTransaction {
	var vendor = (&APVendorList).Add(vendorName)
	var transaction = APTransaction{
		transactionDate: date,
		vendor:          vendor,
	}
	return transaction
}

// ----------------------------------------------------------------------------
// Properties
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// Methods
// ----------------------------------------------------------------------------
