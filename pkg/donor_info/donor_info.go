// ----------------------------------------------------------------------------
//
// Donor information
//
// Author: William Shaffer
// Version: 12-Apr-2024
//
// Copyright (c) 2024 William Shaffer  All Rights Reserved
//
// ----------------------------------------------------------------------------

// Package donor_info performs the manipulaton of data for a single donor.
package donor_info

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	"acorn_go/pkg/assert"
	"strconv"

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

type DonorType int

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

var ZERO = dec.Zero

var MajorDonorLimit = dec.NewFromInt(2000)

const (
	DonorFY2023Only      = 0
	DonorFY2024Only      = 1
	DonorFY2023AndFY2024 = 2
)

var donorTypeNames = [3]string{"Donor FY2023 Only", "Donor FY2024 Only", "Donor FY2023 and FY2024"}

// ----------------------------------------------------------------------------
// Donor type methods
// ----------------------------------------------------------------------------

func (donorType DonorType) String() string {
	assert.Assert(donorType >= DonorFY2023Only && donorType <= DonorFY2023AndFY2024,
		"invalid donor type: "+strconv.Itoa((int(donorType))))
	return donorTypeNames[int(donorType)]
}

// ----------------------------------------------------------------------------
// Factory Functions
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

// Donor returns the donor identifier.
func (donor Donor) NameDonor() string {
	return donor.nameDonor
}

// DonationFY23 returns amount of FY2023 donations for this donor.
func (donor Donor) DonationFY23() dec.Decimal {
	return donor.donationFY23
}

// DonationFY24 returns amount of FY2024 donations for this donor.
func (donor Donor) DonationFY24() dec.Decimal {
	return donor.donationFY24
}

// TotalDonation returns the total donations for this donor.
func (donor Donor) TotalDonation() dec.Decimal {
	var amount = donor.donationFY23
	amount = amount.Add(donor.donationFY24)
	return amount
}

// IsFY23DonorOnly returns true if the donor donated in FY2023 but not
// in FY2024.
func (donor Donor) IsFY23DonorOnly() bool {
	var result = true
	switch {
	case donor.donationFY23.GreaterThan(ZERO) && donor.donationFY24.GreaterThan(ZERO):
		result = false
	case donor.donationFY23.GreaterThan(ZERO):
		result = true
	default:
		result = false
	}
	return result
}

// IsFY24DonorOnly returns true if the donor donated in FY2024 but not
// in FY2023.
func (donor Donor) IsFY24DonorOnly() bool {
	var result = true
	switch {
	case donor.donationFY23.GreaterThan(ZERO) && donor.donationFY24.GreaterThan(ZERO):
		result = false
	case donor.donationFY24.GreaterThan(ZERO):
		result = true
	default:
		result = false
	}
	return result
}

// IsFY23AndFY24Donor returns true if the donor donated in FY2023 and FY2024.
func (donor Donor) IsFY23AndFY24Donor() bool {
	var result = true
	switch {
	case donor.donationFY23.GreaterThan(ZERO) && donor.donationFY24.GreaterThan(ZERO):
		result = true
	default:
		result = false
	}
	return result
}

// IsMajorDonorFY23 returns true if the donor donated $2000 or more in FY2023
func (donor Donor) IsMajorDonorFY23() bool {
	var result = donor.donationFY23.GreaterThanOrEqual(MajorDonorLimit)
	return result
}

// IsMajorDonorFY24 returns true if the donor donated $2000 or more in FY2023
func (donor Donor) IsMajorDonorFY24() bool {
	var result = donor.donationFY24.GreaterThanOrEqual(MajorDonorLimit)
	return result
}

// ----------------------------------------------------------------------------
// Methods
// ----------------------------------------------------------------------------

// AddFY23 adds the payment to the FY2023 donations for this donor.
func (donor *Donor) AddFY23(payment dec.Decimal) {
	donor.donationFY23 = donor.donationFY23.Add(payment)
}

// AddFY24 adds the amount to the FY2024 donations for this donor.
func (donor *Donor) AddFY24(payment dec.Decimal) {
	donor.donationFY24 = donor.donationFY24.Add(payment)
}
