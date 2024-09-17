// ----------------------------------------------------------------------------
//
// Donations
//
// Author: William Shaffer
// Version: 14-Apr-2024
//
// Copyright (c) 2024 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

package donations

// Donations contains a structure for analyzing repeat donors.

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	a "acorn_go/pkg/accounting"
	"math"

	dec "github.com/shopspring/decimal"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type Donations struct {
	fy        a.FYIndicator
	donations [3]dec.Decimal
}

type YearType int

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

const (
	CurrentYear    = YearType(0)
	PriorYear      = YearType(1)
	PriorPriorYear = YearType(2)
)

// ----------------------------------------------------------------------------
// Factory Functions
// ----------------------------------------------------------------------------

// NewDonatons creates a Donations structure initializes to zero for each element.
func NewDonations(fy a.FYIndicator) Donations {
	donations := Donations{
		fy:        fy,
		donations: [3]dec.Decimal{ZERO, ZERO, ZERO},
	}
	return donations
}

// ----------------------------------------------------------------------------
// Validation Functions
// ----------------------------------------------------------------------------

func IsYearType(yearType YearType) bool {
	var result = CurrentYear <= yearType && yearType <= PriorPriorYear
	return result
}

// ----------------------------------------------------------------------------
// Methods
// ----------------------------------------------------------------------------

// Donation returns the donation for a particular type of donor and a
// particular fiscal year
func (donations Donations) Donation(yearType YearType) float64 {
	var value = donations.donations[yearType]
	var donation = value.InexactFloat64()
	return math.Round(donation)
}

// FYDonation returns the total donation for the fiscal year
func (donations Donations) FYDonation() float64 {
	var donation float64 = 0

	for yearType := CurrentYear; yearType <= PriorPriorYear; yearType++ {
		donation += donations.Donation(yearType)
	}
	return donation
}

// Add adds the donation for the specified year type
func (donations Donations) Add(yearType YearType, amount dec.Decimal) {
	donations.donations[yearType].Add(amount)
}

// Return the string for the associated fiscal year indicator.
func (donations Donations) FiscalYear() string {
	return donations.fy.String()
}

// Return the associated fiscal year indicator
func (donations Donations) FY() a.FYIndicator {
	return donations.fy
}
