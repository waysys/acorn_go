// ----------------------------------------------------------------------------
//
// Grant List
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
	q "acorn_go/pkg/quickbooks"
	"strconv"

	dec "github.com/shopspring/decimal"
	"github.com/waysys/assert/assert"
	d "github.com/waysys/waydate/pkg/date"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type GrantList struct {
	trans []Transaction
	count int
}

// ----------------------------------------------------------------------------
// Factory Function
// ----------------------------------------------------------------------------

// NewGrantList returns an initialized grant list
func NewGrantList() GrantList {
	var grantList = GrantList{
		trans: make([]Transaction, 200),
		count: 0,
	}

	return grantList
}

// AssembleGrantList populates the grant list with bill and AP transaction data.
func AssembleGrantList(billList *q.BillList, tranList *q.TransList) (GrantList, error) {
	var err error = nil
	//
	// Create grant list
	//
	var grantList = NewGrantList()
	//
	// Add grants and transfers
	//
	processBills(billList, &grantList)

	return grantList, err
}

// processBills cycles through bills and populates the grants and transfers
func processBills(billList *q.BillList, grantList *GrantList) {
	var numBills = billList.Size()
	var bill q.EducationBill
	var date d.Date
	var amount dec.Decimal
	var recipient q.Recipient
	var edInst q.Vendor
	var transType TransType
	var transaction Transaction
	//
	// Cycle through list of bills
	//
	for index := 0; index < numBills; index++ {
		//
		// Extract data
		//
		bill = *billList.Get(index)
		date = bill.TransactionDate()
		amount = bill.Amount()
		recipient = bill.Recipient()
		edInst = bill.Vendor()
		transType = convertTranType(bill.BType())
		//
		// Create transaction
		//
		transaction = NewTransaction(
			date,
			transType,
			recipient,
			edInst,
			amount,
		)
		//
		// Add transaction
		//
		grantList.Add(transaction)
	}
}

// convertTranType converts the bill type to the grant transacton type
func convertTranType(billType q.BillType) TransType {
	var tranType TransType
	switch billType {
	case q.Grant:
		tranType = Grant
	case q.Transfer:
		tranType = Transfer
	default:
		assert.Assert(false, "Unknown bill type: "+strconv.Itoa(int(billType)))
	}
	return tranType
}

// ----------------------------------------------------------------------------
// Methods
// ----------------------------------------------------------------------------

// Add adds a grant transaction to the grant list.
func (grantList *GrantList) Add(transaction Transaction) {
	grantList.trans[grantList.count] = transaction
	grantList.count++
}
