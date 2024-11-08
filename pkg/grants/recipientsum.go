// ----------------------------------------------------------------------------
//
// Recipients
//
// Author: William Shaffer
// Version: 26-Sep-2024
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
	var values = []dec.Decimal{dec.Zero, dec.Zero, dec.Zero}
	var recipientSum = RecipientSum{
		recipient: &recipient,
		payments:  values,
	}
	return recipientSum
	// Post:
	//   recipientSum.recipient.Name() == name and
	//   recipientSum.payments[0] == dec.Zero and
	//   recipientSum.payments[1] == dec.Zero
}

// ----------------------------------------------------------------------------
// Properties
// ----------------------------------------------------------------------------

// RecipientName returns the name of the recipient
func (sum *RecipientSum) RecipientName() string {
	var name = sum.recipient.Name()
	return name
}

// Payments returns the total payment made to a recipient in a fiscal year.
func (sum *RecipientSum) PaymentTotal(fy a.FYIndicator) dec.Decimal {
	assert.Assert(fy != a.OutOfRange, "fiscal year must not be out of range")
	return sum.payments[fy]
}

// IsRecipient returns true if the recipient benefited from a payment for
// a scholarship in the specified fiscal year.
func (sum *RecipientSum) IsRecipient(fy a.FYIndicator) bool {
	assert.Assert(fy != a.OutOfRange, "fiscal year must not be out of range")
	var result = sum.payments[fy].GreaterThan(dec.Zero)
	return result
}

// ----------------------------------------------------------------------------
// Methods
// ----------------------------------------------------------------------------

// AddPayment adds the amount of a payment to the total payments for the fiscal year.
func (sum *RecipientSum) AddPayment(fy a.FYIndicator, amount dec.Decimal) {
	assert.Assert(fy != a.OutOfRange, "fiscal year must not be out of range")
	assert.Assert(amount.GreaterThan(dec.Zero), "payment must be greater than zero: "+amount.String())
	var originalPayment = sum.PaymentTotal(fy)
	sum.payments[fy] = sum.payments[fy].Add(amount)
	var finalPayment = originalPayment.Add(amount)
	assert.Assert(sum.PaymentTotal(fy).Equal(finalPayment),
		"Payment total is not correct: "+sum.payments[fy].String())
}
