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

	dec "github.com/shopspring/decimal"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type Money dec.Decimal

type APTransaction struct {
	transactionDate d.Date
	vendor          *Vendor
	recipient       *Recipient
	transactionType QuickbooksTransactionType
	amount          Money
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// Factory Functions
// ----------------------------------------------------------------------------

// NewAPTransaction returns a new AP Transaction based on the values
// supplied.
func NewAPTransaction(
	date d.Date,
	vendorName string,
	recipientName string,
	transTypeValue string,
	amount Money) APTransaction {
	var vendor = (&APVendorList).Add(vendorName)
	var recipient = (&APRecipientList).Add(recipientName)
	var transType = NewQuickbooksTransactionType(transTypeValue)
	var transaction = APTransaction{
		transactionDate: date,
		vendor:          vendor,
		recipient:       recipient,
		transactionType: transType,
		amount:          amount,
	}
	return transaction
}

// ----------------------------------------------------------------------------
// Properties
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// Methods
// ----------------------------------------------------------------------------