// ----------------------------------------------------------------------------
//
// Donation Output
//
// Author: William Shaffer
//
// Copyright (c) 2024, 2025 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

package donations

// This file contains code for computing the output related to donations.

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	a "acorn_go/pkg/accounting"

	"github.com/waysys/assert/assert"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type DonationAnalysis []Donations

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// Factory Functions
// ----------------------------------------------------------------------------

// NewDonationAnalysis returns the array of Donations
func NewDonationAnalysis() DonationAnalysis {
	var donations []Donations
	for _, fy := range a.FYIndicators {
		var donation = NewDonations(fy)
		donations = append(donations, donation)
	}
	return donations
}

// ----------------------------------------------------------------------------
// Computational Functions
// ----------------------------------------------------------------------------

// ComputeDonations calculates the breakdown of donations.
func ComputeDonations(donationList DonationList) DonationAnalysis {
	var da = NewDonationAnalysis()

	for _, donor := range donationList {
		for _, fy := range a.FYIndicators {
			var donation = da[fy]
			var amount = donor.Donation(fy)
			donation.ApplyAmount(donor, fy, amount)
		}
	}
	return da
}

// ----------------------------------------------------------------------------
// Methods
// ----------------------------------------------------------------------------

// Donation returns the donations for a specified fiscal year and year type.
func (da *DonationAnalysis) Donation(fy a.FYIndicator, yearType YearType) float64 {
	assert.Assert(a.IsFYIndicator(fy), "Invalid fiscal year indicator: "+fy.String())
	assert.Assert(IsYearType(yearType), "Invalid year type: "+yearType.String())
	var don = (*da)[fy]
	return don.Donation(yearType)
}

// DonationFiscalYear returns the total donations for the fiscal year.
func (da *DonationAnalysis) DonationFiscalYear(fy a.FYIndicator) float64 {
	assert.Assert(a.IsFYIndicator(fy), "Invalid fiscal year indicator: "+fy.String())
	var don = (*da)[fy]
	return don.TotalDonations()
}

// TotalDonations returns the total of donations for all specified years.
func (da *DonationAnalysis) TotalDonations() float64 {
	var total = 0.00
	for _, don := range *da {
		total += don.TotalDonations()
	}
	return total
}
