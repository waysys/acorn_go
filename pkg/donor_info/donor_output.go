// ----------------------------------------------------------------------------
//
// Donor Output
//
// Author: William Shaffer
// Version: 17-Apr-2024
//
// Copyright (c) William Shaffer
//
// ----------------------------------------------------------------------------

package donor_info

// This file produces the outputs from the donor list.

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	"acorn_go/pkg/assert"
	"math"
	"sort"

	dec "github.com/shopspring/decimal"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type DonorCount struct {
	TotalDonors           int
	TotalDonorsFY2023     int
	TotalDonorsFY2024     int
	DonorsFY2023Only      int
	DonorsFY2024Only      int
	DonorsFY2023AndFY2024 int
}
type MajorDonor struct {
	MajorDonorsFY2023           int
	MajorDonorsFY2024           int
	DonationsMajorFY2023        dec.Decimal
	DonationsMajorFY2024        dec.Decimal
	AvgMajorDonationFY2023      float64
	AvgMajorDonationFY2024      float64
	DonationChange              float64
	PercentTotalDonationsFY2023 float64
	PercentTotalDonationsFY2024 float64
}

// ----------------------------------------------------------------------------
// Factory Functions
// ----------------------------------------------------------------------------

// NewDonorCount creates a donor count struture initialized to zero for each element.
func NewDonorCount() DonorCount {
	donorCount := DonorCount{
		TotalDonors:           0,
		TotalDonorsFY2023:     0,
		TotalDonorsFY2024:     0,
		DonorsFY2023Only:      0,
		DonorsFY2024Only:      0,
		DonorsFY2023AndFY2024: 0,
	}
	return donorCount
}

// NewMajorDonor creates a MajorDonor structure initialized to zero for each element.
func NewMajorDonor() MajorDonor {
	majorDonor := MajorDonor{
		MajorDonorsFY2023:      0,
		MajorDonorsFY2024:      0,
		DonationsMajorFY2023:   dec.Zero,
		DonationsMajorFY2024:   dec.Zero,
		AvgMajorDonationFY2023: 0.0,
		AvgMajorDonationFY2024: 0.0,
		DonationChange:         0.0,
	}
	return majorDonor
}

// ----------------------------------------------------------------------------
// Operational Functions
// ----------------------------------------------------------------------------

// ComputeDonorCount calculates the number of various categories of donors.
func ComputeDonorCount(donorListPtr *DonorList) DonorCount {
	var donorCount = NewDonorCount()

	//
	// Loop through all the donor entries in the donor list
	//
	for _, donorPtr := range *donorListPtr {
		donorCount.TotalDonors++
		if (*donorPtr).IsFY23DonorOnly() || (*donorPtr).IsFY23AndFY24Donor() {
			donorCount.TotalDonorsFY2023++
		}
		if (*donorPtr).IsFY24DonorOnly() || (*donorPtr).IsFY23AndFY24Donor() {
			donorCount.TotalDonorsFY2024++
		}
		if (*donorPtr).IsFY23DonorOnly() {
			donorCount.DonorsFY2023Only++
		}
		if (*donorPtr).IsFY24DonorOnly() {
			donorCount.DonorsFY2024Only++
		}
		if (*donorPtr).IsFY23AndFY24Donor() {
			donorCount.DonorsFY2023AndFY2024++
		}
	}

	assert.Assert(donorCount.TotalDonors == len(*donorListPtr),
		"Total donors count does not agree with length of donor list")
	var expectedDonorsFY2023 = donorCount.DonorsFY2023Only + donorCount.DonorsFY2023AndFY2024
	assert.Assert(expectedDonorsFY2023 == donorCount.TotalDonorsFY2023,
		"FY2023 donor count does not agree with expected FY2023 donor count")
	var expectedDonorsFY2024 = donorCount.DonorsFY2024Only + donorCount.DonorsFY2023AndFY2024
	assert.Assert(expectedDonorsFY2024 == donorCount.TotalDonorsFY2024,
		"FY2024 donor count does not agree with expected FY2024 donor count")
	var expectedTotalCount = donorCount.DonorsFY2023Only + donorCount.DonorsFY2024Only +
		donorCount.DonorsFY2023AndFY2024
	assert.Assert(expectedTotalCount == donorCount.TotalDonors,
		"Total donor count does not agree with expected donor count")

	return donorCount
}

// NonRepeatDonors returns a list of donors who donated in FY2023, but
// not in FY2024.
func NonRepeatDonors(donorListPtr *DonorList) []string {
	var names = []string{}
	for name, donorPtr := range *donorListPtr {
		if (*donorPtr).IsFY23DonorOnly() {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	return names
}

// ComputeMajorDonors computes the values for the MajorDonor structure
func ComputeMajorDonors(donorListPtr *DonorList) MajorDonor {
	var majorDonor = NewMajorDonor()
	var donationsTotalFY2023 dec.Decimal
	var donationsTotalFY2024 dec.Decimal

	for _, donor := range *donorListPtr {
		donationsTotalFY2023 = donationsTotalFY2023.Add(donor.DonationFY23())
		donationsTotalFY2024 = donationsTotalFY2024.Add(donor.DonationFY24())

		if donor.IsMajorDonorFY23() {
			majorDonor.MajorDonorsFY2023++
			majorDonor.DonationsMajorFY2023 = majorDonor.DonationsMajorFY2023.Add(donor.donationFY23)
		}
		if donor.IsMajorDonorFY24() {
			majorDonor.MajorDonorsFY2024++
			majorDonor.DonationsMajorFY2024 = majorDonor.DonationsMajorFY2024.Add(donor.donationFY24)
		}
	}

	var majCountFY2023 = float64(majorDonor.MajorDonorsFY2023)
	var majFY2023, _ = majorDonor.DonationsMajorFY2023.Float64()
	majorDonor.AvgMajorDonationFY2023 = math.Round(majFY2023 / majCountFY2023)

	var majCountFY2024 = float64(majorDonor.MajorDonorsFY2024)
	var majFY2024, _ = majorDonor.DonationsMajorFY2024.Float64()
	majorDonor.AvgMajorDonationFY2024 = math.Round(majFY2024 / majCountFY2024)

	majorDonor.DonationChange = math.Round(100 *
		(majorDonor.AvgMajorDonationFY2024 - majorDonor.AvgMajorDonationFY2023) /
		majorDonor.AvgMajorDonationFY2023)

	var totFY2023, _ = donationsTotalFY2023.Float64()
	majorDonor.PercentTotalDonationsFY2023 = math.Round(100.0 * majFY2023 / totFY2023)
	var totFY2024, _ = donationsTotalFY2024.Float64()
	majorDonor.PercentTotalDonationsFY2024 = math.Round(100.0 * majFY2024 / totFY2024)
	return majorDonor
}
