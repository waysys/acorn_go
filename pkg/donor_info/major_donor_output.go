// ----------------------------------------------------------------------------
//
// Major Donor Output
//
// Author: William Shaffer
// Version: 03-May-2024
//
// Copyright (c) William Shaffer
//
// ----------------------------------------------------------------------------

package donor_info

// This file computes the output values for major donors

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	"math"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type MajorDonor struct {
	donorCount     [2]int
	donations      [2]float64
	donationsTotal [2]float64
}

// ----------------------------------------------------------------------------
// Factory Functions
// ----------------------------------------------------------------------------

// NewMajorDonor creates a MajorDonor structure initialized to zero for each element.
func NewMajorDonor() MajorDonor {
	majorDonor := MajorDonor{
		donorCount:     [2]int{0, 0},
		donations:      [2]float64{0.0, 0.0},
		donationsTotal: [2]float64{0.0, 0.0},
	}
	return majorDonor
}

// ----------------------------------------------------------------------------
// Operational Functions
// ----------------------------------------------------------------------------

// ComputeMajorDonors computes the values for the MajorDonor structure
func ComputeMajorDonors(donorListPtr *DonorList) MajorDonor {
	var majorDonor = NewMajorDonor()
	var donationFY2023 float64
	var donationFY2024 float64

	for _, donor := range *donorListPtr {
		donationFY2023 = donor.DonationFY23().InexactFloat64()
		donationFY2024 = donor.DonationFY24().InexactFloat64()

		if donor.IsMajorDonorFY23() {
			majorDonor.donorCount[FY2023]++
			majorDonor.donations[FY2023] += donationFY2023
		}
		if donor.IsMajorDonorFY24() {
			majorDonor.donorCount[FY2024]++
			majorDonor.donations[FY2024] += donationFY2024
		}
		majorDonor.donationsTotal[FY2023] += donationFY2023
		majorDonor.donationsTotal[FY2024] += donationFY2024
	}
	return majorDonor
}

// ----------------------------------------------------------------------------
// Methods
// ----------------------------------------------------------------------------

// MajorDonors returns the number of major donors for the specified fiscal year
func (majorDonor MajorDonor) MajorDonorCount(fy FYIndicator) int {
	return majorDonor.donorCount[fy]
}

// DonationsMajor returns the amount of donations by major donors.
func (majorDonor MajorDonor) DonationsMajor(fy FYIndicator) float64 {
	return math.Round(majorDonor.donations[fy])
}

// AverageDonation returns the average donation by major donors for the
// specified fiscal year
func (majorDonor MajorDonor) AvgDonation(fy FYIndicator) float64 {
	var count = float64(majorDonor.MajorDonorCount(fy))
	var donation = majorDonor.donations[fy]
	var avg = math.Round(donation / count)
	return avg
}

// PercentDonation returns the donations by major donors as a percent
// of the total donations for the specified fiscal year.
func (majorDonor MajorDonor) PercentDonation(fy FYIndicator) float64 {
	var donation = majorDonor.donations[fy]
	var donationsTotal = majorDonor.donationsTotal[fy]
	var percent = math.Round(100.0 * donation / donationsTotal)
	return percent
}

// PercentChange returns the percent change in average donations from
// the current fiscal year compared to the prior fiscal year.
func (majorDonor MajorDonor) PercentChange() float64 {
	var avgDonationFY2023 = majorDonor.AvgDonation(FY2023)
	var avgDonationFY2024 = majorDonor.AvgDonation(FY2024)
	var percentChange = 100 * (avgDonationFY2024 - avgDonationFY2023) / avgDonationFY2023
	return math.Round(percentChange)
}
