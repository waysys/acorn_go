// ----------------------------------------------------------------------------
//
// Donor Count
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
	dn "acorn_go/pkg/donors"
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
func (dc *DonorCount) Count(yearType YearType) int {
	assert.Assert(IsYearType(yearType), "Invalid year type: "+strconv.Itoa(int(yearType)))
	return dc.donorCount[yearType]
}

// Add adds the specified number of donors to the count.
func (dc *DonorCount) Add(yearType YearType, count int) {
	assert.Assert(IsYearType(yearType), "Invalid year type: "+strconv.Itoa(int(yearType)))
	assert.Assert(count >= 0, "Add: count must be zero or greater, not: "+strconv.Itoa(int(yearType)))
	dc.donorCount[yearType] += count
}

// TotalDonorCount returns the total number of donors for the fiscal year.
func (dc *DonorCount) TotalDonorCount() int {
	var total = 0
	for _, count := range dc.donorCount {
		total += count
	}
	return total
}

// FiscalYear returns the fiscal year as a string
func (dc *DonorCount) FiscalYear() string {
	return dc.fy.String()
}

// FY returns the fiscal year
func (dc *DonorCount) FY() a.FYIndicator {
	return dc.fy
}

// ApplyDonorCount increments the proper donor count for the donor.
func (dc *DonorCount) ApplyDonorCount(donor *dn.Donor, analysisFy a.FYIndicator) {
	assert.Assert(a.IsFYIndicator(analysisFy), "invalid fiscal year indicator: "+analysisFy.String())
	var priorFy = analysisFy.Prior()
	var priorPriorFy = priorFy.Prior()

	if donor.IsDonor(analysisFy) {
		if (priorFy != a.OutOfRange) && donor.IsDonor(priorFy) {
			dc.Add(PriorYear, 1)
		} else if (priorPriorFy != a.OutOfRange) && donor.IsDonor(priorPriorFy) {
			dc.Add(PriorPriorYear, 1)
		} else {
			dc.Add(CurrentYear, 1)
		}
	}
}
