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
	q "acorn_go/pkg/quickbooks"
	"strconv"

	dec "github.com/shopspring/decimal"
	"github.com/waysys/assert/assert"
	d "github.com/waysys/waydate/pkg/date"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// The transaction ties together the recipient, the educational institution,
// the award group, fiscal year, and type of transaction.
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
