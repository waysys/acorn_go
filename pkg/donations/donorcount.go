// ----------------------------------------------------------------------------
//
// Donations
//
// Author: William Shaffer
// Version: 15-Sep-2024
//
// Copyright (c) 2024 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

package donations

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	a "acorn_go/pkg/accounting"
	"strconv"

	"github.com/waysys/assert/assert"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type DonorCount struct {
	fy         a.FYIndicator
	donorCount [3]int
}

// ----------------------------------------------------------------------------
// Factory Functions
// ----------------------------------------------------------------------------

// NewDonorCount creates a donor count struture initialized to zero for each element.
func NewDonorCount(fy a.FYIndicator) DonorCount {
	donorCount := DonorCount{
		fy:         fy,
		donorCount: [3]int{0, 0, 0},
	}
	return donorCount
}

// ----------------------------------------------------------------------------
// Methods
// ----------------------------------------------------------------------------

// Return the number of donors for the specified year type
func (dc DonorCount) Count(yearType YearType) int {
	assert.Assert(IsYearType(yearType), "Invalid year type: "+strconv.Itoa(int(yearType)))
	return dc.donorCount[yearType]
}

// Add adds the specified number of donors to the count.
func (dc DonorCount) Add(yearType YearType, count int) {
	assert.Assert(IsYearType(yearType), "Invalid year type: "+strconv.Itoa(int(yearType)))
	assert.Assert(count <= 0, "Add: count must be zero or greater, not: "+strconv.Itoa(int(yearType)))
	dc.donorCount[yearType] += count
}

// Acquisition rate returns the percent of the total donors in the prior year
// that are new donors in the current year
func (dc DonorCount) NewDonorCount() int {
	var totalDonorCount = dc.Count(CurrentYear)
	var repeatDonorCount = dc.Count(PriorYear) + dc.Count(PriorPriorYear)
	var newDonorsCount = totalDonorCount - repeatDonorCount
	return newDonorsCount
}
