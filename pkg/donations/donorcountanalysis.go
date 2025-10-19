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
	dn "acorn_go/pkg/donors"
	"math"
	"strconv"

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
	donorcounts[a.FY2024] = NewDonorCount(a.FY2024)
	donorcounts[a.FY2025] = NewDonorCount(a.FY2025)
	donorcounts[a.FY2026] = NewDonorCount(a.FY2026)
	return donorcounts
}

// ----------------------------------------------------------------------------
// Computational Functions
// ----------------------------------------------------------------------------

// ComputeDonorCount calculates the breakdown of donor counts.
func ComputeDonorCount(donationList DonationList) DonorCountAnalysis {
	var dc = NewDonorCountAnalysis()

	for _, donor := range donationList {
		for _, fy := range a.FYIndicators {
			if fy != a.OutOfRange {
				if donor.IsDonor(fy) {
					var donorCount = dc[fy]
					applyDonorCount(&donorCount, donor)
					dc[fy] = donorCount
				}
			}
		}
	}
	return dc
}

// ApplyDonorCount increments the proper donor count for the donor.
func applyDonorCount(donorCount *DonorCount, donor *dn.Donor) {
	var fy = donorCount.fy
	switch fy {
	case a.FY2024:
		donorCount.Add(CurrentYear, 1)
	case a.FY2025:
		if donor.IsDonor(a.FY2024) {
			donorCount.Add(PriorYear, 1)
		} else {
			donorCount.Add(CurrentYear, 1)
		}
	case a.FY2026:
		if donor.IsDonor(a.FY2025) {
			donorCount.Add(PriorYear, 1)
		} else if donor.IsDonor(a.FY2024) {
			donorCount.Add(PriorPriorYear, 1)
		} else {
			donorCount.Add(CurrentYear, 1)
		}
	}
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
	var retention float64 = 0.00
	switch fy {
	case a.FY2024:
		retention = 0.00
	case a.FY2025:
		var repeatDonors = float64((*dca)[a.FY2025].Count(PriorYear))
		var totalDonors = float64((*dca)[a.FY2024].TotalDonorCount())
		retention = math.Round(repeatDonors * 100 / totalDonors)
	case a.FY2026:
		var repeatDonors = float64((*dca)[a.FY2026].Count(PriorYear))
		var totalDonors = float64((*dca)[a.FY2025].TotalDonorCount())
		retention = math.Round(repeatDonors * 100 / totalDonors)
	default:
		assert.Assert(false, "Invalid fiscal year indicator: "+strconv.Itoa(int(fy)))
	}
	return retention
}

// Acquisition returns the percent of the new donors in the current year
// compared to the total number of donors from the prior year.
func (dca *DonorCountAnalysis) Acquisition(fy a.FYIndicator) float64 {
	var acquisition float64 = 0.0
	switch fy {
	case a.FY2024:
		acquisition = 0.0
	case a.FY2025:
		var newDonors = float64((*dca)[a.FY2025].Count(CurrentYear))
		var totalDonors = float64((*dca)[a.FY2024].TotalDonorCount())
		acquisition = math.Round(newDonors * 100 / totalDonors)
	case a.FY2026:
		var newDonors = float64((*dca)[a.FY2026].Count(CurrentYear))
		var totalDonors = float64((*dca)[a.FY2025].TotalDonorCount())
		acquisition = math.Round(newDonors * 100 / totalDonors)
	default:
		assert.Assert(false, "Invalid fiscal year indicator: "+strconv.Itoa(int(fy)))
	}
	return acquisition
}
