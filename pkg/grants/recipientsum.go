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
	grants    []dec.Decimal
}

// ----------------------------------------------------------------------------
// Factory Methods
// ----------------------------------------------------------------------------

// NewRecipientSum returns a new recipient summary initialized to zero.
func NewRecipientSum(name string) RecipientSum {
	assert.Assert(len(name) > 0, "name cannot be an empty string")
	var payments []dec.Decimal
	var refunds []dec.Decimal
	var grants []dec.Decimal
	var recipient = q.NewRecipient(name)
	var recipientSum = RecipientSum{
		recipient: &recipient,
		payments:  zeroArray(payments),
		refunds:   zeroArray(refunds),
		grants:    zeroArray(grants),
	}
	return recipientSum
}

// zeroArray creates a slice with the number of elements with decimal zero
// equal to the number of fiscal years in the reports
func zeroArray(input []dec.Decimal) []dec.Decimal {
	for index := 0; index < a.NumFiscalYears; index++ {
		input = append(input, dec.Zero)
	}
	return input
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

// GrantTotal returns the total grants made to a recipient in a fiscal year.
func (sum *RecipientSum) GrantTotal(fy a.FYIndicator) dec.Decimal {
	assert.Assert(fy != a.OutOfRange, "fiscal year must not be out of range")
	var result = sum.grants[fy]
	return result
}

// IsIndividualRecipient returns true if the recipient received individual
// grants in the specified fiscal year
func (sum *RecipientSum) IsIndividualGrant(fy a.FYIndicator) bool {
	var total = sum.GrantTotal(fy)
	var result = total.GreaterThan(dec.Zero)
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

// AddGrant adds the amount of a grant to the total grants for the fiscal year
func (sum *RecipientSum) AddGrant(fy a.FYIndicator, amount dec.Decimal) {
	assert.Assert(fy != a.OutOfRange, "fiscal year must not be out of range")
	assert.Assert(amount.GreaterThan(dec.Zero), "payment must be greater than zero: "+amount.String())
	sum.grants[fy] = sum.grants[fy].Add(amount)
}
