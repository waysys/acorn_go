// ----------------------------------------------------------------------------
//
// Transactions
//
// Author: William Shaffer
// Version: 30-May-2024
//
// Copyright (c) 2024 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

package grants

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

// TransType represents the type of transaction associated with scholarships.
type TransType int

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

const (
	Grant        TransType = 0
	Transfer     TransType = 1
	GrantPayment TransType = 2
	WriteOff     TransType = 3
	Refund       TransType = 4
)

var TransTypes = []TransType{Grant, Transfer, GrantPayment, WriteOff, Refund}

// ----------------------------------------------------------------------------
// Validation
// ----------------------------------------------------------------------------

// IsTransType returns true if the argument is a valid value for the traansaction
// type
func IsTransType(value TransType) bool {
	var result = true
	if value < Grant || value > Refund {
		result = false
	}
	return result
}

// ----------------------------------------------------------------------------
// Methods
// ----------------------------------------------------------------------------

// String converts a transaction type to a string value
func (transType TransType) String() string {
	var values = []string{
		"Grant",
		"Transfer",
		"Grant Payment",
		"Write-Off",
		"Refund",
	}

	assert.Assert(transType >= 0 && transType <= 4,
		"invalid transaction type: "+strconv.Itoa(int(transType)))
	return values[transType]
}

// lessTransType returns true if this transaction type is less
// than another tranaction type.
func (transType TransType) lessTransType(anotherTransType TransType) bool {
	return transType < anotherTransType
}
