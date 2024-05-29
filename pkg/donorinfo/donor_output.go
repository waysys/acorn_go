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

package donorinfo

// This file produces the outputs from the donor list.

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	"math"

	"github.com/waysys/assert/assert"
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

// ----------------------------------------------------------------------------
// Methods
// ----------------------------------------------------------------------------

// RetentionRate returns the percent of the total donors in the prior year
// that donate in the current year
func (dc DonorCount) RetentionRate() float64 {
	var priorYearDonors = float64(dc.TotalDonorsFY2023)
	var currentYearRepeatDonors = float64(dc.DonorsFY2023AndFY2024)
	var retention = 100 * currentYearRepeatDonors / priorYearDonors
	return math.Round(retention)
}

// Acquisition rate returns the percent of the total donors in the prior year
// that are new donors in the current year
func (dc DonorCount) AcquisitionRate() float64 {
	var priorYearDonors = float64(dc.TotalDonorsFY2023)
	var newDonors = float64(dc.DonorsFY2024Only)
	var acquisition = 100 * newDonors / priorYearDonors
	return math.Round(acquisition)
}
