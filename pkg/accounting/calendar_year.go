// ----------------------------------------------------------------------------
//
// Calendar Year
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
	d "github.com/waysys/waydate/pkg/date"
	"github.com/waysys/waydate/pkg/daterange"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// YearIndicator represents a calendar year identifier used by the accounting
// package. Valid values map to specific years (Y2022..Y2026). Unknown is used
// for out-of-range values.
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

var yearIndicators = []YearIndicator{
	Y2022,
	Y2023,
	Y2024,
	Y2025,
	Y2026,
}

var indNames = []string{"Y2022", "Y2023", "Y2024", "Y2025", "Y2026", "Unknown"}

var NumYears = len(yearIndicators)
var MinYear = 2022
var MaxYear = 2026

var currentMonthStart d.Date
var currentMonthEnd d.Date
var currentMonth daterange.DateRange

// ----------------------------------------------------------------------------
// Initialization
// ----------------------------------------------------------------------------

func init() {
	var err error
	currentMonthStart, err = d.New(11, 1, 2025)
	if err != nil {
		panic(err)
	}
	currentMonthEnd, err = d.New(11, 30, 2025)
	if err != nil {
		panic(err)
	}
	currentMonth, err = daterange.New(currentMonthStart, currentMonthEnd)
	if err != nil {
		panic(err)
	}
}

// ----------------------------------------------------------------------------
// Functions
// ----------------------------------------------------------------------------

// YIndicator returns the year indicator associated with a date
func YIndicator(date d.Date) YearIndicator {
	var year = int(date.Year())
	var indicator = Unknown
	if year >= MinYear && year <= MaxYear {
		indicator = yearIndicators[year-MinYear]
	}
	return indicator
}

// IsYearIndicator reports whether the indicator is between Y2022 and Y2026.
func IsYearIndicator(indicator YearIndicator) bool {
	var result = false
	if Y2022 <= indicator && indicator <= Y2026 {
		result = true
	}
	return result
}

// IsCurrentMonth reports whether the date is in the defined current month.
func IsCurrentMonth(date d.Date) bool {
	var result = currentMonth.InRange(date)
	return result
}

// CurrentMonth returns the current month range as a string.
func CurrentMonth() string {
	var result = currentMonth.String()
	return result
}

// ----------------------------------------------------------------------------
// Methods
// ----------------------------------------------------------------------------

// String returns the name of the YearIndicator. It returns "Unknown" for out-of-range values.
func (ind YearIndicator) String() string {
	var index = int(ind)
	if index < 0 || index >= NumYears {
		index = NumYears - 1 // Unknown
	}
	return indNames[index]
}
