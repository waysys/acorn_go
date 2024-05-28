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

// BillType returns the bill type for the bill
func (bill *EducationBill) BType() BillType {
	return bill.billType
}
