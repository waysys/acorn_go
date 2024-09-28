// ----------------------------------------------------------------------------
//
// Major Donor Output
//
// Author: William Shaffer
// Version: 03-May-2024
//
// Copyright (c) 2024 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

package majordonor

// This file computes the output values for major donors

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	a "acorn_go/pkg/accounting"
	dn "acorn_go/pkg/donations"
	"math"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type MajorDonor struct {
	donorCount     [3]int
	donations      [3]float64
	donationsTotal [3]float64
}

// ----------------------------------------------------------------------------
// Factory Functions
// ----------------------------------------------------------------------------

// NewMajorDonor creates a MajorDonor structure initialized to zero for each element.
func NewMajorDonor() MajorDonor {
	majorDonor := MajorDonor{
		donorCount:     [3]int{0, 0, 0},
		donations:      [3]float64{0.0, 0.0, 0.0},
		donationsTotal: [3]float64{0.0, 0.0, 0.0},
	}
	return majorDonor
}

// ----------------------------------------------------------------------------
// Operational Functions
// ----------------------------------------------------------------------------

// ComputeMajorDonors computes the values for the MajorDonor structure
func ComputeMajorDonors(donationList *dn.DonationList) MajorDonor {
	var majorDonor = NewMajorDonor()
	var donations [3]float64
	//
	// Loop through donors with annual donation information.
	//
	for _, donor := range *donationList {
		//
		// Loop through fiscal years
		//
		for _, fy := range a.FYIndicators {
			donations[fy] = donor.Donation(fy).InexactFloat64()
			if donor.IsMajorDonor(fy) {
				majorDonor.donorCount[fy]++
				majorDonor.donations[fy] += donations[fy]
			}
			majorDonor.donationsTotal[fy] += donations[fy]
		}
	}
	return majorDonor
}

// ----------------------------------------------------------------------------
// Methods
// ----------------------------------------------------------------------------

// MajorDonors returns the number of major donors for the specified fiscal year
func (majorDonor MajorDonor) MajorDonorCount(fy a.FYIndicator) int {
	return majorDonor.donorCount[fy]
}

// DonationsMajor returns the amount of donations by major donors.
func (majorDonor MajorDonor) DonationsMajor(fy a.FYIndicator) float64 {
	return math.Round(majorDonor.donations[fy])
}

// AverageDonation returns the average donation by major donors for the
// specified fiscal year
func (majorDonor MajorDonor) AvgDonation(fy a.FYIndicator) float64 {
	var count = float64(majorDonor.MajorDonorCount(fy))
	var donation = majorDonor.donations[fy]
	var avg = 0.00
	if count > 0.00 {
		avg = math.Round(donation / count)
	}
	return avg
}

// PercentDonation returns the donations by major donors as a percent
// of the total donations for the specified fiscal year.
func (majorDonor MajorDonor) PercentDonation(fy a.FYIndicator) float64 {
	var donation = majorDonor.donations[fy]
	var donationsTotal = majorDonor.donationsTotal[fy]
	var percent = 0.00
	if donationsTotal > 0.00 {
		percent = math.Round(100.0 * donation / donationsTotal)
	}
	return percent
}

// PercentChange returns the percent change in average donations from
// the current fiscal year compared to the prior fiscal year.
func (majorDonor MajorDonor) PercentChange(fy a.FYIndicator) float64 {
	var percentChange = 0.00
	switch fy {
	case a.FY2023:
		percentChange = 0.00
	case a.FY2024:
		var avgDonationFY2023 = majorDonor.AvgDonation(a.FY2023)
		var avgDonationFY2024 = majorDonor.AvgDonation(a.FY2024)
		if avgDonationFY2023 > 0.00 {
			percentChange = 100.00 * (avgDonationFY2024 - avgDonationFY2023) / avgDonationFY2023
		}
	case a.FY2025:
		var avgDonationFY2024 = majorDonor.AvgDonation(a.FY2024)
		var avgDonationFY2025 = majorDonor.AvgDonation(a.FY2025)
		if avgDonationFY2024 > 0.00 {
			percentChange = 100 * (avgDonationFY2025 - avgDonationFY2024) / avgDonationFY2024
		}
	}
	return math.Round(percentChange)
}
