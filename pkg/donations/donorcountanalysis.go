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
	var donorcounts = make(DonorCountAnalysis, 3)
	donorcounts[a.FY2023] = NewDonorCount(a.FY2023)
	donorcounts[a.FY2024] = NewDonorCount(a.FY2024)
	donorcounts[a.FY2025] = NewDonorCount(a.FY2025)
	return donorcounts
}

// ----------------------------------------------------------------------------
// Computational Functions
// ----------------------------------------------------------------------------

// ComputeDonorCount calculates the breakdown of donor counts.
func ComputeDonorCount(donationListPtr *DonationList) DonorCountAnalysis {
	var dc = NewDonorCountAnalysis()

	for _, donor := range *donationListPtr {
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
	case a.FY2023:
		donorCount.Add(CurrentYear, 1)
	case a.FY2024:
		if donor.IsDonor(a.FY2023) {
			donorCount.Add(PriorYear, 1)
		} else {
			donorCount.Add(CurrentYear, 1)
		}
	case a.FY2025:
		if donor.IsDonor(a.FY2024) {
			donorCount.Add(PriorYear, 1)
		} else if donor.IsDonor(a.FY2023) {
			donorCount.Add(PriorPriorYear, 1)
		} else {
			donorCount.Add(CurrentYear, 1)
		}
	}
}

// ----------------------------------------------------------------------------
// Methods
// ----------------------------------------------------------------------------

// TotalDonations returns the count of all donors.
func (dca *DonorCountAnalysis) TotalDonors() int {
	var total = 0
	for _, dc := range *dca {
		total += dc.TotalDonorCount()
	}
	return total
}

// Retention returns the percent of donors from the prior year that donate
// in the current year.
func (dca *DonorCountAnalysis) Retention(fy a.FYIndicator) float64 {
	var retention float64 = 0.00
	switch fy {
	case a.FY2023:
		retention = 0.00
	case a.FY2024:
		var repeatDonors = float64((*dca)[a.FY2024].Count(PriorYear))
		var totalDonors = float64((*dca)[a.FY2023].TotalDonorCount())
		retention = math.Round(repeatDonors * 100 / totalDonors)
	case a.FY2025:
		var repeatDonors = float64((*dca)[a.FY2025].Count(PriorYear))
		var totalDonors = float64((*dca)[a.FY2024].TotalDonorCount())
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
	case a.FY2023:
		acquisition = 0.0
	case a.FY2024:
		var newDonors = float64((*dca)[a.FY2024].Count(CurrentYear))
		var totalDonors = float64((*dca)[a.FY2023].TotalDonorCount())
		acquisition = math.Round(newDonors * 100 / totalDonors)
	case a.FY2025:
		var newDonors = float64((*dca)[a.FY2025].Count(CurrentYear))
		var totalDonors = float64((*dca)[a.FY2024].TotalDonorCount())
		acquisition = math.Round(newDonors * 100 / totalDonors)
	default:
		assert.Assert(false, "Invalid fiscal year indicator: "+strconv.Itoa(int(fy)))
	}
	return acquisition
}
