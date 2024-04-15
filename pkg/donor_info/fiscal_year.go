// ----------------------------------------------------------------------------
//
// Fiscal Year
//
// Author: William Shaffer
// Version: 14-Apr-2024
//
// Copyright (c) William Shaffer
//
// ----------------------------------------------------------------------------

package donor_info

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
	FY2023     = FYIndicator(1)
	FY2024     = FYIndicator(2)
	OutOfRange = FYIndicator(3)
)

var IndNames = [4]string{"", "FY2023", "FY2024", "OutOfRange"}

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
func FiscalYearIndicator(dateString string) FYIndicator {
	var indicator = OutOfRange
	//
	// Convert string into date
	//
	var date, err = d.NewFromString(dateString)
	//
	// Select Fiscal Year
	//
	if err != nil {
		indicator = OutOfRange
	} else if FiscalYear2023.InRange(date) {
		indicator = FY2023
	} else if FiscalYear2024.InRange(date) {
		indicator = FY2024
	} else {
		indicator = OutOfRange
	}
	return indicator
}

// ----------------------------------------------------------------------------
// Methods
// ----------------------------------------------------------------------------

// String returns the name of the fiscal year indicator
func (ind FYIndicator) String() string {
	return IndNames[int(ind)]
}
