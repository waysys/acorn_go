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
	WriteOff     TransType = 2
	GrantPayment TransType = 3
)

// ----------------------------------------------------------------------------
// Methods
// ----------------------------------------------------------------------------

// String converts a transaction type to a string value
func (transType TransType) String() string {
	var values = []string{
		"Grant",
		"Transfer",
		"Write-Off",
		"Grant Payment",
	}

	assert.Assert(transType >= 0 && transType <= 3,
		"invalid transaction type: "+strconv.Itoa(int(transType)))
	return values[transType]
}
