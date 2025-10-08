// ----------------------------------------------------------------------------
//
// Fiscal Year
//
// Author: William Shaffer
// Version: 15-Sep-2024
//
// Copyright (c) 2024 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

package accounting

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	"github.com/waysys/waydate/pkg/daterange"

	d "github.com/waysys/waydate/pkg/date"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type FYIndicator int

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

const (
	FY2024     = FYIndicator(0)
	FY2025     = FYIndicator(1)
	FY2026     = FYIndicator(2)
	OutOfRange = FYIndicator(3)
)

var FYIndicators = [3]FYIndicator{
	FY2024,
	FY2025,
	FY2026,
}

var IndNames = [4]string{"FY2024", "FY2025", "FY2026", "OutOfRange"}

var FiscalYear2024Begin, _ = d.New(9, 1, 2023)
var FiscalYear2024End, _ = d.New(8, 31, 2024)
var FiscalYear2025Begin, _ = d.New(9, 1, 2024)
var FiscalYear2025End, _ = d.New(8, 31, 2025)
var FiscalYear2026Begin, _ = d.New(9, 1, 2025)
var FiscalYear2026End, _ = d.New(8, 31, 2026)

var FiscalYears = [3]daterange.DateRange{
	fiscalYear2024,
	fiscalYear2025,
	fiscalYear2026,
}

var fiscalYear2024, _ = daterange.New(FiscalYear2024Begin, FiscalYear2024End)
var fiscalYear2025, _ = daterange.New(FiscalYear2025Begin, FiscalYear2025End)
var fiscalYear2026, _ = daterange.New(FiscalYear2026Begin, FiscalYear2026End)

// ----------------------------------------------------------------------------
// Functions
// ----------------------------------------------------------------------------

// FiscalYearIndicator returns a FYIndicator value based on the date provided
// as an argument.
func FiscalYearIndicator(date d.Date) FYIndicator {
	var indicator = OutOfRange
	//
	// Select Fiscal Year
	//
	for _, fy := range FYIndicators {
		if fy == OutOfRange {
			indicator = OutOfRange
			break
		}
		var rng = FiscalYears[fy]
		if rng.InRange(date) {
			indicator = fy
			break
		}
	}
	return indicator
}

// FiscalYearFromYearMonth returns a FYIndicator value based on the YearMonth.
func FiscalYearFromYearMonth(yearMonth d.YearMonth) (FYIndicator, error) {
	var err error = nil
	var date d.Date
	var indicator = OutOfRange
	date, err = d.New(d.Month(yearMonth.Month), 1, d.Year(yearMonth.Year))
	if err == nil {
		indicator = FiscalYearIndicator(date)
	}
	return indicator, err
}

// IsFYIndicator return true if the FYIndicator is a valid value
func IsFYIndicator(fy FYIndicator) bool {
	var result = false
	if fy >= FY2024 && fy <= OutOfRange {
		result = true
	}
	return result
}

// ----------------------------------------------------------------------------
// Methods
// ----------------------------------------------------------------------------

// String returns the name of the fiscal year indicator
func (ind FYIndicator) String() string {
	return IndNames[int(ind)]
}

// Prior returns the fiscal year indicator before the specified fiscal year.
func (ind FYIndicator) Prior() FYIndicator {
	var result = OutOfRange
	switch ind {
	case FY2024:
		result = OutOfRange
	case FY2025:
		result = FY2024
	case FY2026:
		result = FY2025
	default:
		result = OutOfRange
	}
	return result
}
