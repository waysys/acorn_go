// ----------------------------------------------------------------------------
//
// Recipients
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

	dec "github.com/shopspring/decimal"
	"github.com/waysys/assert/assert"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// RecipientSum contains a summary of payments made to recipients by fiscal year.
type RecipientSum struct {
	recipient *q.Recipient
	payments  []dec.Decimal
}

// ----------------------------------------------------------------------------
// Factory Methods
// ----------------------------------------------------------------------------

// NewRecipientSum returns a new recipient summary initialized to zero.
func NewRecipientSum(name string) RecipientSum {
	assert.Assert(len(name) > 0, "name cannot be an empty string")
	var recipient = q.NewRecipient(name)
	var values = []dec.Decimal{dec.Zero, dec.Zero}
	var recipientSum = RecipientSum{
		recipient: &recipient,
		payments:  values,
	}
	return recipientSum
}

// ----------------------------------------------------------------------------
// Properties
// ----------------------------------------------------------------------------

// RecipientName returns the name of the recipient
func (sum RecipientSum) RecipientName() string {
	var name = sum.recipient.Name()
	return name
}

// Payments returns the total payment made to a recipient in a fiscal year.
func (sum RecipientSum) PaymentTotal(fy a.FYIndicator) dec.Decimal {
	assert.Assert(fy != a.OutOfRange, "fiscal year must not be out of range")
	return sum.payments[fy]
}

// ----------------------------------------------------------------------------
// Methods
// ----------------------------------------------------------------------------

// AddPayment adds the amount of a payment to the total payments for the fiscal year.
func (sum RecipientSum) AddPayment(fy a.FYIndicator, amount dec.Decimal) {
	assert.Assert(fy != a.OutOfRange, "fiscal year must not be out of range")
	assert.Assert(amount.GreaterThan(dec.Zero), "payment must be greater than zero: "+amount.String())
	sum.payments[fy] = sum.payments[fy].Add(amount)
}