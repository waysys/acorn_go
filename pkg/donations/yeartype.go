// ----------------------------------------------------------------------------
//
// Year Types
//
// Author: William Shaffer
// Version: 24-Sep-2024
//
// Copyright (c) 2024 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

package donations

// Year Types identify the current donation year, the prior donation year, and the
// donation year 2 years before the current year.

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

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
// Validation Functions
// ----------------------------------------------------------------------------

func IsYearType(yearType YearType) bool {
	var result = CurrentYear <= yearType && yearType <= PriorPriorYear
	return result
}
