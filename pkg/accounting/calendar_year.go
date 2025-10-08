// ----------------------------------------------------------------------------
//
// Calendar Year
//
// Author: William Shaffer
// Version: 02-Oct-2025
//
// Copyright (c) 2024, 2025 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

package accounting

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	d "github.com/waysys/waydate/pkg/date"
	"github.com/waysys/waydate/pkg/daterange"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type YearIndicator int

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

const (
	Y2022   = YearIndicator(0)
	Y2023   = YearIndicator(1)
	Y2024   = YearIndicator(2)
	Y2025   = YearIndicator(3)
	Y2026   = YearIndicator(4)
	Unknown = YearIndicator(5)
)

var yearIndicators = [6]YearIndicator{
	Y2022,
	Y2023,
	Y2024,
	Y2025,
	Y2026,
	Unknown,
}

var indNames = [6]string{"Y2022", "Y2023", "Y2024", "Y2025", "Y2026", "Unknown"}

var currentMonthStart, _ = d.New(9, 1, 2025)
var currentMonthEnd, _ = d.New(9, 30, 2025)
var currentMonth, _ = daterange.New(currentMonthStart, currentMonthEnd)

// ----------------------------------------------------------------------------
// Functions
// ----------------------------------------------------------------------------

// YIndicator returns the year indicator associated with a date
func YIndicator(date d.Date) YearIndicator {
	var year = int(date.Year())
	var indicator = Unknown
	if year >= 2022 && year <= 2026 {
		indicator = yearIndicators[int(date.Year())-2022]
	}
	return indicator
}

// IsYearIndicator returns true if the indicator is between Y2022 and Y2026
func IsYearIndicator(indicator YearIndicator) bool {
	var result = false
	if Y2022 <= indicator && indicator <= Y2026 {
		result = true
	}
	return result
}

// IsCurrentMonth returns true if the date is in the defined current month
func IsCurrentMonth(date d.Date) bool {
	var result = currentMonth.InRange(date)
	return result
}

// CurrentMonth returns the current month range as a string
func CurrentMonth() string {
	var result = currentMonth.String()
	return result
}

// ----------------------------------------------------------------------------
// Methods
// ----------------------------------------------------------------------------

// String returns the name of the fiscal year indicator
func (ind YearIndicator) String() string {
	return indNames[int(ind)]
}
