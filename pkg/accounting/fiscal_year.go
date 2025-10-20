// ----------------------------------------------------------------------------
//
// Fiscal Year
//
// Author: William Shaffer
//
// Copyright (c) 2024, 2025 William Shaffer All Rights Reserved
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

const NumFiscalYears = 4

const (
	FY2023     = FYIndicator(0)
	FY2024     = FYIndicator(1)
	FY2025     = FYIndicator(2)
	FY2026     = FYIndicator(3)
	OutOfRange = FYIndicator(4)
)

var FYIndicators = [NumFiscalYears]FYIndicator{
	FY2023,
	FY2024,
	FY2025,
	FY2026,
}

var IndNames = [NumFiscalYears + 1]string{"FY2023", "FY2024", "FY2025", "FY2026", "OutOfRange"}
var FiscalYear2023Begin d.Date
var FiscalYear2023End d.Date
var FiscalYear2024Begin d.Date
var FiscalYear2024End d.Date
var FiscalYear2025Begin d.Date
var FiscalYear2025End d.Date
var FiscalYear2026Begin d.Date
var FiscalYear2026End d.Date

var FiscalYears [NumFiscalYears]daterange.DateRange

var fiscalYear2023 daterange.DateRange
var fiscalYear2024 daterange.DateRange
var fiscalYear2025 daterange.DateRange
var fiscalYear2026 daterange.DateRange

// Identify fiscal years for analysis
const (
	CurrentFiscalYear = FY2026
)

// ----------------------------------------------------------------------------
// Initialization
// ----------------------------------------------------------------------------

func init() {
	var err error

	FiscalYear2023Begin, err = d.New(9, 1, 2022)
	if err != nil {
		panic(err)
	}
	FiscalYear2023End, err = d.New(8, 31, 2023)
	if err != nil {
		panic(err)
	}

	FiscalYear2024Begin, err = d.New(9, 1, 2023)
	if err != nil {
		panic(err)
	}
	FiscalYear2024End, err = d.New(8, 31, 2024)
	if err != nil {
		panic(err)
	}

	FiscalYear2025Begin, err = d.New(9, 1, 2024)
	if err != nil {
		panic(err)
	}
	FiscalYear2025End, err = d.New(8, 31, 2025)
	if err != nil {
		panic(err)
	}

	FiscalYear2026Begin, err = d.New(9, 1, 2025)
	if err != nil {
		panic(err)
	}
	FiscalYear2026End, err = d.New(8, 31, 2026)
	if err != nil {
		panic(err)
	}

	fiscalYear2023, err = daterange.New(FiscalYear2023Begin, FiscalYear2023End)
	if err != nil {
		panic(err)
	}
	fiscalYear2024, err = daterange.New(FiscalYear2024Begin, FiscalYear2024End)
	if err != nil {
		panic(err)
	}
	fiscalYear2025, err = daterange.New(FiscalYear2025Begin, FiscalYear2025End)
	if err != nil {
		panic(err)
	}
	fiscalYear2026, err = daterange.New(FiscalYear2026Begin, FiscalYear2026End)
	if err != nil {
		panic(err)
	}

	FiscalYears = [NumFiscalYears]daterange.DateRange{
		fiscalYear2023,
		fiscalYear2024,
		fiscalYear2025,
		fiscalYear2026,
	}
}

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
		var rng = FiscalYears[int(fy)]
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

// IsFYIndicator returns true if the FYIndicator is a valid value
func IsFYIndicator(fy FYIndicator) bool {
	var result = false
	// OutOfRange is not a valid FYIndicator
	if fy >= FY2023 && fy < OutOfRange {
		result = true
	}
	return result
}

// ----------------------------------------------------------------------------
// Methods
// ----------------------------------------------------------------------------

// String returns the name of the fiscal year indicator
func (ind FYIndicator) String() string {
	var index = int(ind)
	if index < 0 || index >= len(IndNames) {
		index = len(IndNames) - 1 // OutOfRange
	}
	return IndNames[int(index)]
}

// Prior returns the fiscal year indicator before the specified fiscal year.
func (ind FYIndicator) Prior() FYIndicator {
	var result = OutOfRange
	var index = int(ind) - 1
	if ind == OutOfRange {
		result = OutOfRange
	} else if index >= 0 {
		result = FYIndicators[index]
	}
	return result
}
