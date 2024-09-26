// ----------------------------------------------------------------------------
//
// Donation Output
//
// Author: William Shaffer
// Version: 30-Apr-2024
//
// Copyright (c) 2024 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

package donations

// This file contains code for computing the output related to donations.

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	a "acorn_go/pkg/accounting"
	dn "acorn_go/pkg/donors"

	dec "github.com/shopspring/decimal"
	"github.com/waysys/assert/assert"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type DonationAnalysis []Donations

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

var ZERO = dec.Zero

// ----------------------------------------------------------------------------
// Factory Functions
// ----------------------------------------------------------------------------

// NewDonationAnalysis returns the array of Donations
func NewDonationAnalysis() DonationAnalysis {
	var donations = make(DonationAnalysis, 3)
	donations[a.FY2023] = NewDonations(a.FY2023)
	donations[a.FY2024] = NewDonations(a.FY2024)
	donations[a.FY2025] = NewDonations(a.FY2025)
	return donations
}

// ----------------------------------------------------------------------------
// Computational Functions
// ----------------------------------------------------------------------------

// ComputeDonations calculates the breakdown of donations.
func ComputeDonations(donationListPtr *DonationList) DonationAnalysis {
	var da = NewDonationAnalysis()
	//
	// Loop through the list of donations
	//
	for _, donor := range *donationListPtr {
		for _, fy := range a.FYIndicators {
			if fy != a.OutOfRange {
				var amount = donor.Donation(fy)
				da.applyAmount(fy, donor, amount)
			}
		}
	}
	return da
}

// ----------------------------------------------------------------------------
// Methods
// ----------------------------------------------------------------------------

// applyAmount applies the amount to the proper elements
func (da *DonationAnalysis) applyAmount(
	fy a.FYIndicator,
	donor *dn.Donor,
	amount dec.Decimal) {
	switch fy {
	case a.FY2023:
		da.add(fy, CurrentYear, amount)
	case a.FY2024:
		if donor.IsDonor(a.FY2023) {
			da.add(fy, PriorYear, amount)
		} else {
			da.add(fy, CurrentYear, amount)
		}
	case a.FY2025:
		if donor.IsDonor(a.FY2024) {
			da.add(fy, PriorYear, amount)
		} else if donor.IsDonor(a.FY2023) {
			da.add(fy, PriorPriorYear, amount)
		} else {
			da.add(fy, CurrentYear, amount)
		}
	default:
		assert.Assert(false, "Invalid fiscal year")
	}
}

// Add adds a donation amount to the donations for the fiscal year
func (da *DonationAnalysis) add(fy a.FYIndicator, yearType YearType, amount dec.Decimal) {
	var donations = (*da)[fy]
	assert.Assert(donations.FY() == fy, "Incorrect fiscal year indicator")
	donations.Add(yearType, amount)
	(*da)[fy] = donations
}
