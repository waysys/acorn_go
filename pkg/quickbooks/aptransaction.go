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

// The APTransaction models the transactions in the Accounts Payable register.

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	d "github.com/waysys/waydate/pkg/date"

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
	account         string
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

const (
	accountScholarship = "7040"
	accountChecking    = "1010"
)

// Returned check from NC State
const (
	accountCash   = "1010 Cash:Cash in bank - operating"
	accountGrants = "7040 Grants, contracts, & direct assistance:Awards & grants - individuals"
	ncState       = "North Carolina State University"
	mPlazas       = "Plazas, Matthew"
)

var oct31, _ = d.New(10, 31, 2024)

// Refund transactions for Matthew Plazas
var ncDeposit = NewAPTransaction(
	oct31,
	ncState,
	mPlazas,
	Deposit,
	Money(dec.NewFromInt(3000)),
	accountCash)

var ncVendorCredit = NewAPTransaction(
	oct31,
	ncState,
	mPlazas,
	VendorCredit,
	Money(dec.NewFromInt(3000)),
	accountGrants)

// ----------------------------------------------------------------------------
// Factory Functions
// ----------------------------------------------------------------------------

// NewAPTransaction returns a new AP Transaction based on the values
// supplied.
func NewAPTransaction(
	date d.Date,
	vendorName string,
	recipientName string,
	transType QuickbooksTransactionType,
	amount Money,
	account string) APTransaction {
	var vendor = (&APVendorList).Add(vendorName)
	var recipient = (&APRecipientList).Add(recipientName)
	var transaction = APTransaction{
		transactionDate: date,
		vendor:          vendor,
		recipient:       recipient,
		transactionType: transType,
		amount:          amount,
		account:         account,
	}
	return transaction
}

// ----------------------------------------------------------------------------
// Properties
// ----------------------------------------------------------------------------

// TransactionDate returns the transaction date of the transactions.
func (trans *APTransaction) TransactionDate() d.Date {
	return trans.transactionDate
}

// Vendor returns the vendor associated with this transaction.
func (trans *APTransaction) Vendor() *Vendor {
	return trans.vendor
}

// Recipient returns the recipient associated with this transaction.
func (trans *APTransaction) Recipient() *Recipient {
	return trans.recipient
}

// TransactionType returns the type of transaction.
func (trans *APTransaction) TransactionType() QuickbooksTransactionType {
	return trans.transactionType
}

// Amount returns the amount of the transaction.
func (trans *APTransaction) Amount() Money {
	return trans.amount
}

// Account returns the 4 digit account number
func (trans *APTransaction) Account() string {
	var value string

	if len(trans.account) < 4 {
		value = ""
	} else {
		value = trans.account[:4]
	}
	return value
}

// IsScholarshipAccount returns true if the account is the 7050 -
// Grants
func (trans *APTransaction) IsScholarshipAccount() bool {
	return trans.Account() == accountScholarship
}

// IsBankAccount returns true if the account is 1010 - Cash
func (trans *APTransaction) IsBankAccount() bool {
	return trans.Account() == accountChecking
}

// GTZero returns true if the amount is greater zero.
func (trans *APTransaction) GTZero() bool {
	return dec.Decimal(trans.amount).GreaterThan(dec.Zero)
}

// IsBill returns true if the transaction is a valid bill
func (trans *APTransaction) IsBill() bool {
	var result = trans.TransactionType() == Bill
	result = result && trans.IsScholarshipAccount()
	result = result && trans.GTZero()
	return result
}

// IsPayment returns true if the transaction is a valid payment for
// scholarships
func (trans *APTransaction) IsPayment() bool {
	var result = trans.TransactionType() == BillPayment
	result = result && trans.IsBankAccount()
	result = result && trans.GTZero()
	return result
}

// IsVendorCredit returns true if the transaction is a valid
// vendor credit.
func (trans *APTransaction) IsVendorCredit() bool {
	var result = trans.TransactionType() == VendorCredit
	result = result && trans.IsScholarshipAccount()
	result = result && trans.GTZero()
	return result
}

// IsDeposit returns true if the transaction is a deposit.
func (trans *APTransaction) IsDeposit() bool {
	var result = trans.TransactionType() == Deposit
	result = result && trans.IsBankAccount()
	result = result && trans.GTZero()
	return result
}
