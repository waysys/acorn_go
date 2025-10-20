// ----------------------------------------------------------------------------
//
// Donor Count Analysis
//
// Author: William Shaffer
// Version: 24-Sep-2024
//
// Copyright (c) 2024 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

package donations

// The Donor Count List is an array from fiscal year to the associated donor
// count object.  This file contains functions to create the list and
// access the data.

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

type DonorCountAnalysis []DonorCount

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// Factory Functions
// ----------------------------------------------------------------------------

// NewDonorCountAnalysis returns an initialized donoar analysis.
func NewDonorCountAnalysis() DonorCountAnalysis {
	var donorcounts = make(DonorCountAnalysis, a.NumFiscalYears)
	var fiscalYears = a.FYIndicators
	for _, fy := range fiscalYears {
		donorcounts[fy] = NewDonorCount(fy)
	}
	return donorcounts
}

// ----------------------------------------------------------------------------
// Computational Functions
// ----------------------------------------------------------------------------

// ComputeDonorCount calculates the breakdown of donor counts.
func ComputeDonorCount(donationList DonationList) DonorCountAnalysis {
	var dc = NewDonorCountAnalysis()

	for _, donor := range donationList {
		for _, analysisFY := range a.FYIndicators {
			dc[analysisFY].ApplyDonorCount(donor, analysisFY)
		}
	}
	return dc
}

// ----------------------------------------------------------------------------
// Methods
// ----------------------------------------------------------------------------

// DonorCount returns the number of donors for the specified fiscal year
// and year type.
func (dca *DonorCountAnalysis) DonorCount(fy a.FYIndicator, yearType YearType) int {
	var dc = (*dca)[fy]
	return dc.Count(yearType)
}

// DonorCountFiscalYear returns the number of donors for the specified fiscal year.
func (dca *DonorCountAnalysis) DonorCountFiscalYear(fy a.FYIndicator) int {
	assert.Assert(a.IsFYIndicator(fy), "Invalid fiscal year indicator: "+fy.String())
	var dc = (*dca)[fy]
	return dc.TotalDonorCount()
}

// TotalDonations returns the count of all donors.
func (dca *DonorCountAnalysis) TotalDonors() int {
	var total = 0
	for _, dc := range *dca {
		total += dc.Count(CurrentYear)
	}
	return total
}

// Retention returns the percent of donors from the prior year that donate
// in the current year.
func (dca *DonorCountAnalysis) Retention(fy a.FYIndicator) float64 {
	assert.Assert(a.IsFYIndicator(fy), "Invalid fiscal year indicator: "+fy.String())

	var retention float64 = 0.00
	var priorFy = fy.Prior()
	if priorFy == a.OutOfRange {
		return retention
	}
	var currentDC = (*dca)[fy]
	var priorDC = (*dca)[fy]

	if priorDC.TotalDonorCount() == 0 {
		return retention
	}
	var totalPriorDonors = float64(priorDC.TotalDonorCount())
	var currentRepeatDonors = float64(currentDC.Count(PriorYear))
	retention = currentRepeatDonors * 100.00 / totalPriorDonors
	retention = math.Round(retention)

	return retention
}

// Acquisition returns the percent of the new donors in the current year
// compared to the total number of donors from the prior year.
func (dca *DonorCountAnalysis) Acquisition(fy a.FYIndicator) float64 {
	assert.Assert(a.IsFYIndicator(fy), "Invalid fiscal year indicator: "+fy.String())

	var acquisition float64 = 0.0

	var priorFy = fy.Prior()
	if priorFy == a.OutOfRange {
		return acquisition
	}
	var currentDC = (*dca)[fy]
	var priorDC = (*dca)[priorFy]

	if priorDC.TotalDonorCount() == 0 {
		return acquisition
	}

	var totalDonors = float64(priorDC.TotalDonorCount())
	var newDonors = float64(currentDC.Count(CurrentYear))
	acquisition = newDonors * 100.00 / totalDonors
	acquisition = math.Round(acquisition)
	return acquisition
}
