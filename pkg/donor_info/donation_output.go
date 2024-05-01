// ----------------------------------------------------------------------------
//
// Donation Output
//
// Author: William Shaffer
// Version: 30-Apr-2024
//
// Copyright (c) William Shaffer
//
// ----------------------------------------------------------------------------

package donor_info

// This file contains code for computing the output related to donations.

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	"math"

	dec "github.com/shopspring/decimal"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type Donations struct {
	DonationsFY2023   dec.Decimal
	DonationsFY2024   dec.Decimal
	DonationsTotal    dec.Decimal
	AvgDonationFY2023 float64
	AvgDonationFY2024 float64
	DonationChange    float64
	CountRepeatDonors int
}

// ----------------------------------------------------------------------------
// Factory Functions
// ----------------------------------------------------------------------------

// NewDonatons creates a Donations structure initializes to zero for each element.
func NewDonations() Donations {
	donations := Donations{
		DonationsFY2023:   dec.Zero,
		DonationsFY2024:   dec.Zero,
		DonationsTotal:    dec.Zero,
		AvgDonationFY2023: 0.0,
		AvgDonationFY2024: 0.0,
		DonationChange:    0.0,
		CountRepeatDonors: 0,
	}
	return donations
}

// ----------------------------------------------------------------------------
// Computational Functions
// ----------------------------------------------------------------------------

// ComputeDonations calculates the breakdown of donations.
func ComputeDonations(donorListPtr *DonorList) Donations {
	var donations = NewDonations()
	var repeatDonationFY2023 = dec.Zero
	var repeatDonationFy2024 = dec.Zero
	var countRepeatDonors = 0

	for _, donor := range *donorListPtr {
		donations.DonationsTotal = donations.DonationsTotal.Add(donor.DonationFY23().Add(donor.DonationFY24()))
		donations.DonationsFY2023 = donations.DonationsFY2023.Add(donor.DonationFY23())
		donations.DonationsFY2024 = donations.DonationsFY2024.Add(donor.donationFY24)
		if donor.IsFY23AndFY24Donor() {
			countRepeatDonors++
			repeatDonationFY2023 = repeatDonationFY2023.Add(donor.DonationFY23())
			repeatDonationFy2024 = repeatDonationFy2024.Add(donor.DonationFY24())
		}
	}

	var repFY2023, _ = repeatDonationFY2023.Float64()
	var repFY2024, _ = repeatDonationFy2024.Float64()
	var repCount = float64(countRepeatDonors)
	donations.CountRepeatDonors = countRepeatDonors
	donations.AvgDonationFY2023 = math.Round(repFY2023 / repCount)
	donations.AvgDonationFY2024 = math.Round(repFY2024 / repCount)
	donations.DonationChange =
		math.Round(100.0 * (donations.AvgDonationFY2024 - donations.AvgDonationFY2023) /
			donations.AvgDonationFY2023)

	return donations
}

// ----------------------------------------------------------------------------
// Methods
// ----------------------------------------------------------------------------

//
