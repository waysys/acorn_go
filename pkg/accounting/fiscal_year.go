// ----------------------------------------------------------------------------
//
// Fiscal Year
//
// Author: William Shaffer
// Version: 14-Apr-2024
//
// Copyright (c) 2024 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

package accounting

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	d "acorn_go/pkg/date"
	"acorn_go/pkg/daterange"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type FYIndicator int

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

const (
	FY2023     = FYIndicator(0)
	FY2024     = FYIndicator(1)
	OutOfRange = FYIndicator(2)
)

var IndNames = [3]string{"FY2023", "FY2024", "OutOfRange"}

var FiscalYear2023Begin, _ = d.New(9, 1, 2022)
var FiscalYear2023End, _ = d.New(8, 31, 2023)
var FiscalYear2024Begin, _ = d.New(9, 1, 2023)
var FiscalYear2024End, _ = d.New(8, 31, 2024)

var FiscalYear2023, _ = daterange.New(FiscalYear2023Begin, FiscalYear2023End)
var FiscalYear2024, _ = daterange.New(FiscalYear2024Begin, FiscalYear2024End)

// ----------------------------------------------------------------------------
// Functions
// ----------------------------------------------------------------------------

// FiscalYearIndicator returns a FYIndicator value based on the date provided
// as an argument (in string format MM/DD/YYY)
func FiscalYearIndicator(date d.Date) FYIndicator {
	var indicator = OutOfRange
	//
	// Select Fiscal Year
	//
	if FiscalYear2023.InRange(date) {
		indicator = FY2023
	} else if FiscalYear2024.InRange(date) {
		indicator = FY2024
	} else {
		indicator = OutOfRange
	}
	return indicator
}

// FiscalYearFromYearMonth returns a FYIndicator value based on the YearMonth
func FiscalYearFromYearMonth(yearMonth d.YearMonth) (FYIndicator, error) {
	var err error = nil
	var date d.Date
	var fyIndicator FYIndicator
	date, err = d.New(d.Month(yearMonth.Month), 1, d.Year(yearMonth.Year))
	if err != nil {
		fyIndicator = OutOfRange
	} else {
		fyIndicator = FiscalYearIndicator(date)
	}
	return fyIndicator, err
}

// ----------------------------------------------------------------------------
// Methods
// ----------------------------------------------------------------------------

// String returns the name of the fiscal year indicator
func (ind FYIndicator) String() string {
	return IndNames[int(ind)]
}
