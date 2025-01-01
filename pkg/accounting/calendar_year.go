// ----------------------------------------------------------------------------
//
// Fiscal Year
//
// Author: William Shaffer
// Version: 30-Dec-2024
//
// Copyright (c) 2024 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

package accounting

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	d "github.com/waysys/waydate/pkg/date"
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
	Unknown = YearIndicator(4)
)

var yearIndicators = [5]YearIndicator{
	Y2022,
	Y2023,
	Y2024,
	Y2025,
	Unknown,
}

var indNames = [5]string{"Y2022", "Y2023", "Y2024", "Y2025", "Unknown"}

// ----------------------------------------------------------------------------
// Functions
// ----------------------------------------------------------------------------

// YIndicator returns the year indicator associated with a date
func YIndicator(date d.Date) YearIndicator {
	var year = int(date.Year())
	var indicator = Unknown
	if year >= 2022 && year <= 2025 {
		indicator = yearIndicators[int(date.Year())-2022]
	}
	return indicator
}

// IsYearIndicator returns true if the indicator is between Y2022 and Y2025
func IsYearIndicator(indicator YearIndicator) bool {
	var result = false
	if Y2022 <= indicator && indicator <= Y2025 {
		result = true
	}
	return result
}

// ----------------------------------------------------------------------------
// Methods
// ----------------------------------------------------------------------------

// String returns the name of the fiscal year indicator
func (ind YearIndicator) String() string {
	return indNames[int(ind)]
}
