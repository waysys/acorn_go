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
	refunds   []dec.Decimal
}

// ----------------------------------------------------------------------------
// Factory Methods
// ----------------------------------------------------------------------------

// NewRecipientSum returns a new recipient summary initialized to zero.
func NewRecipientSum(name string) RecipientSum {
	assert.Assert(len(name) > 0, "name cannot be an empty string")
	var recipient = q.NewRecipient(name)
	var recipientSum = RecipientSum{
		recipient: &recipient,
		payments:  []dec.Decimal{dec.Zero, dec.Zero, dec.Zero},
		refunds:   []dec.Decimal{dec.Zero, dec.Zero, dec.Zero},
	}
	return recipientSum
}

// ----------------------------------------------------------------------------
// Properties
// ----------------------------------------------------------------------------

// RecipientName returns the name of the recipient
func (sum *RecipientSum) RecipientName() string {
	var name = sum.recipient.Name()
	return name
}

// PaymentTotal returns the total payment made to a recipient in a fiscal year.
func (sum *RecipientSum) PaymentTotal(fy a.FYIndicator) dec.Decimal {
	assert.Assert(fy != a.OutOfRange, "fiscal year must not be out of range")
	return sum.payments[fy]
}

// RefundTotal returns the total refunds made to a recipient in a fiscal year.
func (sum *RecipientSum) RefundTotal(fy a.FYIndicator) dec.Decimal {
	assert.Assert(fy != a.OutOfRange, "fiscal year must not be out of range")
	return sum.refunds[fy]
}

// IsRecipient returns true if the recipient benefited from a payment for
// a scholarship in the specified fiscal year.
func (sum *RecipientSum) IsRecipient(fy a.FYIndicator) bool {
	assert.Assert(fy != a.OutOfRange, "fiscal year must not be out of range")
	var result = sum.payments[fy].GreaterThan(dec.Zero)
	return result
}

// NetPaymentTotal returns the total payments minus the total refunds.
func (sum *RecipientSum) NetPaymentTotal(fy a.FYIndicator) dec.Decimal {
	assert.Assert(fy != a.OutOfRange, "fiscal year must not be out of range")
	var result = sum.PaymentTotal(fy).Sub(sum.RefundTotal(fy))
	return result
}

// ----------------------------------------------------------------------------
// Methods
// ----------------------------------------------------------------------------

// AddPayment adds the amount of a payment to the total payments for the fiscal year.
func (sum *RecipientSum) AddPayment(fy a.FYIndicator, amount dec.Decimal) {
	assert.Assert(fy != a.OutOfRange, "fiscal year must not be out of range")
	assert.Assert(amount.GreaterThan(dec.Zero), "payment must be greater than zero: "+amount.String())
	sum.payments[fy] = sum.payments[fy].Add(amount)
}

// AddRefund adds the amount of a refund to the total refunds for the fiscal year.
func (sum *RecipientSum) AddRefund(fy a.FYIndicator, amount dec.Decimal) {
	assert.Assert(fy != a.OutOfRange, "fiscal year must not be out of range")
	assert.Assert(amount.GreaterThan(dec.Zero), "payment must be greater than zero: "+amount.String())
	sum.refunds[fy] = sum.refunds[fy].Add(amount)
}
