// ----------------------------------------------------------------------------
//
// Year Month
//
// Author: William Shaffer
// Version: 25-Apr-2024
//
// Copyright (c) William Shaffer
//
// ----------------------------------------------------------------------------

package date

// This file manages year month structures

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type YearMonth struct {
	Year  int
	Month int
}

// ----------------------------------------------------------------------------
// Factory Functions
// ----------------------------------------------------------------------------

func NewYearMonth(year int, month int) (YearMonth, error) {
	var err error = nil
	//
	// Preconditions
	//
	err = isYear(Year(year))
	if err != nil {
		return YearMonth{}, err
	}
	err = isMonth(Month(month))
	if err != nil {
		return YearMonth{}, err
	}
	//
	// Set the values
	//
	yearMonth := YearMonth{
		year,
		month,
	}
	return yearMonth, nil
}

func NewYearMonthFromDate(date Date) (YearMonth, error) {
	var year = int(date.Day())
	var month = int(date.Month())
	return NewYearMonth(year, month)
}
