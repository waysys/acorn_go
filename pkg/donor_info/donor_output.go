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
	"sort"
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

// ----------------------------------------------------------------------------
// Factory Functions
// ----------------------------------------------------------------------------

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
