// ----------------------------------------------------------------------------
//
// Donor information
//
// Author: William Shaffer
// Version: 12-Apr-2024
//
// Copyright (c) William Shaffer
//
// ----------------------------------------------------------------------------

// Package donor_info performs the manipulaton of data for a single donor.
package donor_info

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	"acorn_go/pkg/assert"

	dec "github.com/shopspring/decimal"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type Donor struct {
	nameDonor    string
	donationFY23 dec.Decimal
	donationFY24 dec.Decimal
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

var ZERO = dec.NewFromInt(0)

// ----------------------------------------------------------------------------
// Factory methods
// ----------------------------------------------------------------------------

func New(name string) Donor {
	assert.Assert(len(name) > 0, "Donor name must not be ")
	donor := Donor{
		nameDonor:    name,
		donationFY23: ZERO,
		donationFY24: ZERO,
	}
	return donor
}

// ----------------------------------------------------------------------------
// Properties
// ----------------------------------------------------------------------------

// Donor returns the donor identifier
func (donor Donor) NameDonor() string {
	return donor.nameDonor
}

func (donor Donor) DonationFY23() dec.Decimal {
	return donor.donationFY23
}

func (donor Donor) DonationFY24() dec.Decimal {
	return donor.donationFY24
}

func (donor Donor) TotalDonation() dec.Decimal {
	var amount = donor.donationFY23
	amount = amount.Add(donor.donationFY24)
	return amount
}
