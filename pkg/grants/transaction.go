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
	a "acorn_go/pkg/accounting"
	q "acorn_go/pkg/quickbooks"
	"strconv"

	dec "github.com/shopspring/decimal"
	"github.com/waysys/assert/assert"
	d "github.com/waysys/waydate/pkg/date"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// The transaction ties together the transaction date, recipient, the educational institution,
// and type of transaction.
type Transaction struct {
	transactionDate d.Date
	transType       TransType
	recipient       *q.Recipient
	edInst          *q.Vendor
	amount          dec.Decimal
}

// ----------------------------------------------------------------------------
// Factory Methods
// ----------------------------------------------------------------------------

func NewTransaction(
	date d.Date,
	transType TransType,
	recipient *q.Recipient,
	edInst *q.Vendor,
	amount dec.Decimal,
) Transaction {
	assert.Assert(IsTransType(transType),
		"invalid transaction type: "+strconv.Itoa(int(transType)))
	var tran = Transaction{
		transactionDate: date,
		transType:       transType,
		recipient:       recipient,
		edInst:          edInst,
		amount:          amount,
	}
	return tran
}

// ----------------------------------------------------------------------------
// Properties
// ----------------------------------------------------------------------------

// TransactionDate returns the date of the transaction
func (trans *Transaction) TransactionDate() d.Date {
	return trans.transactionDate
}

// FiscalYear returns the fiscal year in which the transaction takes place.
func (trans *Transaction) FiscalYear() a.FYIndicator {
	var transactionDate = trans.TransactionDate()
	return a.FiscalYearIndicator(transactionDate)
}

// TransType returns the transaction type
func (trans *Transaction) TransType() TransType {
	return trans.transType
}

// Recipient returns a string with the name of the recipient
func (trans *Transaction) Recipient() string {
	return trans.recipient.Name()
}

// EducationalInstitution returns a string with the name of the
// educational institution
func (trans *Transaction) EducationalInstitution() string {
	return trans.edInst.Name()
}

// Amount returns the amount of the transaction.
func (trans *Transaction) Amount() dec.Decimal {
	return trans.amount
}

// GTZero returns true if the amount is greater than zero.
func (trans *Transaction) GTZero() bool {
	return trans.Amount().GreaterThanOrEqual(dec.Zero)
}

// IsPayment returns true if the transaction type is GrantPayment
func (trans *Transaction) IsPayment() bool {
	var result = trans.TransType() == GrantPayment
	result = result && trans.GTZero()
	return result
}

// IsRefund returns true if the transaction type is Refund
func (trans *Transaction) IsRefund() bool {
	var result = trans.TransType() == Refund
	result = result && trans.GTZero()
	return result
}

// IsGrant returns true if the transaction type is Grant
func (trans *Transaction) IsGrant() bool {
	var result = trans.TransType() == Grant
	result = result && trans.GTZero()
	return result
}

// ----------------------------------------------------------------------------
// Other Functions
// ----------------------------------------------------------------------------

// compare compares two transactions and returns true if the first transaction
// should come before the sedond transaction.
func compare(trans1 *Transaction, trans2 *Transaction) bool {
	var result = false
	var tranDate1 = trans1.TransactionDate()
	var tranDate2 = trans2.TransactionDate()

	if trans1.Recipient() < trans2.Recipient() {
		result = true
	} else if trans1.Recipient() > trans2.Recipient() {
		result = false
	} else {
		if tranDate1.Before(tranDate2) {
			result = true
		} else if tranDate1.After(tranDate2) {
			result = false
		} else {
			result = trans1.TransType().lessTransType(trans2.TransType())
		}
	}
	return result
}
