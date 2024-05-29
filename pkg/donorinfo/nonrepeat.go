// ----------------------------------------------------------------------------
//
// Non-Repeat Data
//
// Author: William Shaffer
// Version: 07-May-2024
//
// Copyright (c) William Shaffer
//
// ----------------------------------------------------------------------------

package donorinfo

// This file contains code to generate the retention analysis.

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	d "acorn_go/pkg/date"

	"github.com/waysys/assert/assert"

	dec "github.com/shopspring/decimal"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type NonRepeat struct {
	nameDonor string
	date      d.Date
	amount    dec.Decimal
}

// ----------------------------------------------------------------------------
// Factory Method
// ----------------------------------------------------------------------------

// NewNonRepeat returns an initialized NonRepeat structure.
func NewNonRepeat(name string, dte d.Date, amt dec.Decimal) NonRepeat {
	assert.Assert(amt.GreaterThan(ZERO), "Amount of donation must be greater than 0.00")
	nonRepeat := NonRepeat{
		nameDonor: name,
		date:      dte,
		amount:    amt,
	}
	return nonRepeat
}

// ----------------------------------------------------------------------------
// methods
// ----------------------------------------------------------------------------

// NameDonor returns the name of the donor
func (nr NonRepeat) NameDonor() string {
	return nr.nameDonor
}

// DateDonation returns the date of the donation
func (nr NonRepeat) DateDonation() d.Date {
	return nr.date
}

// AmountDonation returns the amount of the donation
func (nr NonRepeat) AmountDonation() dec.Decimal {
	return nr.amount
}
