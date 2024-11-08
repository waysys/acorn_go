// ----------------------------------------------------------------------------
//
// AP Transaction
//
// Author: William Shaffer
// Version: 25-May-2024
//
// Copyright (c) 2024 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

package quickbooks

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

type QuickbooksTransactionType int

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

const (
	Unknown      QuickbooksTransactionType = 0
	Bill         QuickbooksTransactionType = 1
	BillPayment  QuickbooksTransactionType = 2
	VendorCredit QuickbooksTransactionType = 3
	Deposit      QuickbooksTransactionType = 4
)

const (
	strUnknown      = "Unknown"
	strBill         = "Bill"
	strBillPayment  = "Bill Payment"
	strVendorCredit = "Vendor Credit"
	strDeposit      = "Deposit"
)

var strValues = [5]string{
	strUnknown,
	strBill,
	strBillPayment,
	strVendorCredit,
	strDeposit,
}

// ----------------------------------------------------------------------------
// Factory Functions
// ----------------------------------------------------------------------------

// NewQuickbooksTransactionType returns a transaction type based on the value.
func NewQuickbooksTransactionType(value string) QuickbooksTransactionType {
	var aType = Unknown
	switch value {
	case strBill:
		aType = Bill
	case strBillPayment:
		aType = BillPayment
	case strVendorCredit:
		aType = VendorCredit
	case strDeposit:
		aType = Deposit
	default:
		aType = Unknown
	}
	return aType
}

// ----------------------------------------------------------------------------
// Methods
// ----------------------------------------------------------------------------

// String returns the string description of the transaction type.
func (transType QuickbooksTransactionType) String() string {
	assert.Assert(0 <= transType && transType <= 3,
		"Incorrect value for Quickbooks transaction type: "+strconv.Itoa(int(transType)))
	return strValues[transType]
}
