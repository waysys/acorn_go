// ----------------------------------------------------------------------------
//
// AP Transaction
//
// Author: William Shaffer
// Version: 27-May-2024
//
// Copyright (c) 2027 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

package quickbooks

// An EducationBill is a subset of AP Transactions

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	dec "github.com/shopspring/decimal"
	"github.com/waysys/assert/assert"
	d "github.com/waysys/waydate/pkg/date"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type EducationBill struct {
	trans    *APTransaction
	billType BillType
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// Factory Function
// ----------------------------------------------------------------------------

// NewEducationBill returns a new education bill based on the inputs.
func NewEducationBill(trans *APTransaction, value string) EducationBill {
	assert.Assert(trans != nil, "transaction must not be nil")
	var bill = EducationBill{
		trans:    trans,
		billType: NewBillType(value),
	}
	return bill
}

// ----------------------------------------------------------------------------
// Methods
// ----------------------------------------------------------------------------

// Trans returns the transaction associated with the bill
func (bill *EducationBill) Trans() *APTransaction {
	return bill.trans
}

// TransactionDate returns the date of the transaction
func (bill *EducationBill) TransactionDate() d.Date {
	return bill.trans.TransactionDate()
}

// BillType returns the bill type for the bill
func (bill *EducationBill) BType() BillType {
	return bill.billType
}

// Ammount returns the amount of the bill
func (bill *EducationBill) Amount() dec.Decimal {
	return dec.Decimal(bill.trans.Amount())
}

// Recipient returns the recipient associated with this bill
func (bill *EducationBill) Recipient() Recipient {
	return *bill.trans.recipient
}

// Vendor returns the vendor associated with this bill
func (bill *EducationBill) Vendor() Vendor {
	return *bill.trans.vendor
}
