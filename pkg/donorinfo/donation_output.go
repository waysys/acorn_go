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

package donorinfo

// This file contains code for computing the output related to donations.

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	"math"

	a "acorn_go/pkg/accounting"

	dec "github.com/shopspring/decimal"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type Donations struct {
	donations  [3][2]dec.Decimal
	donorCount [3]int
}

// ----------------------------------------------------------------------------
// Factory Functions
// ----------------------------------------------------------------------------

// NewDonatons creates a Donations structure initializes to zero for each element.
func NewDonations() Donations {
	donations := Donations{
		donations: [3][2]dec.Decimal{
			{ZERO, ZERO}, // FY2023 donor only
			{ZERO, ZERO}, // FY2024 donor only
			{ZERO, ZERO}, // FY2023 and FY2024 donor
		},
		donorCount: [3]int{
			0, 0, 0,
		},
	}
	return donations
}

// ----------------------------------------------------------------------------
// Computational Functions
// ----------------------------------------------------------------------------

// ComputeDonations calculates the breakdown of donations.
func ComputeDonations(donorListPtr *DonorList) Donations {
	var dons = NewDonations()

	for _, donor := range *donorListPtr {
		if donor.IsFY23DonorOnly() {
			dons.donations[DonorFY2023Only][a.FY2023] =
				dons.donations[DonorFY2023Only][a.FY2023].Add(donor.donationFY23)
			dons.donorCount[DonorFY2023Only]++
		} else if donor.IsFY24DonorOnly() {
			dons.donations[DonorFY2024Only][a.FY2024] =
				dons.donations[DonorFY2024Only][a.FY2024].Add(donor.donationFY24)
			dons.donorCount[DonorFY2024Only]++
		} else if donor.IsFY23AndFY24Donor() {
			dons.donations[DonorFY2023AndFY2024][a.FY2023] =
				dons.donations[DonorFY2023AndFY2024][a.FY2023].Add(donor.donationFY23)
			dons.donations[DonorFY2023AndFY2024][a.FY2024] =
				dons.donations[DonorFY2023AndFY2024][a.FY2024].Add(donor.donationFY24)
			dons.donorCount[DonorFY2023AndFY2024]++
		}
	}
	return dons
}

// ----------------------------------------------------------------------------
// Methods
// ----------------------------------------------------------------------------

// Donation returns the donation for a particular type of donor and a
// particular fiscal year
func (dons Donations) Donation(donorType DonorType, fy a.FYIndicator) float64 {
	var value = dons.donations[donorType][fy]
	var donation = value.InexactFloat64()
	return math.Round(donation)
}

// FYDonation returns the total donation for the fiscal year
func (dons Donations) FYDonation(fy a.FYIndicator) float64 {
	var donation float64 = 0
	var donorType DonorType

	for donorType = DonorFY2023Only; donorType <= DonorFY2023AndFY2024; donorType++ {
		donation += dons.Donation(donorType, fy)
	}
	return donation
}

// TotalDonation returns the total donation made for all fiscal years
func (dons Donations) TotalDonation() float64 {
	var donation float64 = 0
	for fy := a.FY2023; fy <= a.FY2024; fy++ {
		donation += dons.FYDonation(fy)
	}
	return donation
}

// DonorCount returns the number of donors of the specified type.
func (dons Donations) DonorCount(donorType DonorType) int {
	var count = dons.donorCount[donorType]
	return count
}

// AvgDonation returns the average donation for a specified donor type and fiscal year
func (dons Donations) AvgDonation(donorType DonorType, fy a.FYIndicator) float64 {
	var donation = dons.Donation(donorType, fy)
	var count = dons.DonorCount(donorType)
	var avg = donation / float64(count)
	return math.Round(avg)
}

// FYAvgDonation returns the average donations for all donors in a fiscal year
func (dons Donations) FYAvgDonation(fy a.FYIndicator) float64 {
	var donation = dons.FYDonation(fy)
	var count = 0
	if fy == a.FY2023 {
		count = dons.DonorCount(DonorFY2023Only) + dons.DonorCount(DonorFY2023AndFY2024)
	} else if fy == a.FY2024 {
		count = dons.DonorCount(DonorFY2024Only) + dons.DonorCount(DonorFY2023AndFY2024)
	}
	var avg = donation / float64(count)
	return math.Round(avg)
}

// DonationChange computes the percent change in donations for a particular donor type.
func (dons Donations) DonationChange(donorType DonorType) float64 {
	var change float64
	var avgFY2023 = dons.AvgDonation(donorType, a.FY2023)
	var avgFY2024 = dons.AvgDonation(donorType, a.FY2024)
	change = (avgFY2024 - avgFY2023) * 100.00 / avgFY2023
	return math.Round(change)
}
