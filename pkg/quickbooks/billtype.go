// ----------------------------------------------------------------------------
//
// Bill Type
//
// Author: William Shaffer
// Version: 27-May-2024
//
// Copyright (c) 2024 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

package quickbooks

// BillType distinguishes between a bill associated with a grant versus a
// bill resulting in a transfer from one educational institution to another.

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	"strconv"

	"github.com/waysys/assert/assert"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type BillType int

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

const (
	Grant        BillType = 0
	Transfer     BillType = 1
	Individual   BillType = 2
	Unidentified BillType = 3
)

const (
	strGrant        = "Grant"
	strTransfer     = "Transfer"
	strIndividual   = "Individual"
	strUnidentified = "Unidentified"
)

var billTypeStrings = []string{
	strGrant,
	strTransfer,
	strIndividual,
	strUnidentified,
}

// ----------------------------------------------------------------------------
// Factory Function
// ----------------------------------------------------------------------------

// NewBillType returns a bill type based on the value provided.
func NewBillType(value string) BillType {
	var result BillType
	switch value {
	case "Transfer":
		result = Transfer
	case strIndividual:
		result = Individual
	default:
		result = Grant
	}
	return result
}

// ----------------------------------------------------------------------------
// Methods
// ----------------------------------------------------------------------------

// String returns the string description of the bill type
func (billType BillType) String() string {
	assert.Assert(0 <= billType && billType <= 3,
		"Incorrect value for bill type: "+strconv.Itoa(int(billType)))
	return billTypeStrings[billType]
}
