// ----------------------------------------------------------------------------
//
// Average Donations
//
// Author: William Shaffer
// Version: 26-Sep-2024
//
// Copyright (c) 2024 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

package donations

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	a "acorn_go/pkg/accounting"

	"math"

	"github.com/waysys/assert/assert"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type DonationsAndCounts struct {
	da *DonationAnalysis
	ca *DonorCountAnalysis
}

// ----------------------------------------------------------------------------
// Factory Functions
// ----------------------------------------------------------------------------

// NewDonationsAndCounts
func NewDonationsAndCounts(da *DonationAnalysis, ca *DonorCountAnalysis) DonationsAndCounts {
	var dac = DonationsAndCounts{
		da: da,
		ca: ca,
	}
	return dac
}

// ----------------------------------------------------------------------------
// Methods
// ----------------------------------------------------------------------------

// AvgDonation returns the average donation for a specified fiscal year and
// year type.
func (dac *DonationsAndCounts) AvgDonation(fy a.FYIndicator, yearType YearType) float64 {
	assert.Assert(IsYearType(yearType), "Invalid year type: "+yearType.String())
	assert.Assert(a.IsFYIndicator(fy), "Invalid fiscal year indicator: "+fy.String())
	var average = 0.00
	var donation = dac.da.Donation(fy, yearType)
	var count = float64(dac.ca.DonorCount(fy, yearType))
	if count > 0.00 {
		average = math.Round(donation / count)
	}
	return average
}

// AvgTotalDonationFiscalYear returns the average donation for a specified year.
func (dac *DonationsAndCounts) AvgTotalDonationFiscalYear(fy a.FYIndicator) float64 {
	assert.Assert(a.IsFYIndicator(fy), "Invalid fiscal year indicator: "+fy.String())
	var average = 0.00
	var donation = dac.da.DonationFiscalYear(fy)
	var count = float64(dac.ca.DonorCountFiscalYear(fy))
	if count > 0.00 {
		average = math.Round(donation / count)
	}
	return average
}

// TotalAvgDonation returns the average donation for all donors.
func (dac *DonationsAndCounts) TotalAvgDonation() float64 {
	var average = 0.00
	var donation = dac.da.TotalDonations()
	var count = float64(dac.ca.TotalDonors())
	if count > 0.00 {
		average = math.Round(donation / count)
	}
	return average
}
