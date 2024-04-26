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

import (
	"strconv"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type YearMonth struct {
	Year  int
	Month int
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

const (
	startYear = 2022
	endYear   = 2024
)

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
		Year:  year,
		Month: month,
	}
	return yearMonth, err
}

func NewYearMonthFromDate(date Date) (YearMonth, error) {
	var year = int(date.Year())
	var month = int(date.Month())
	return NewYearMonth(year, month)
}

// Keys returns a slice with the year month structures for 2023 and 2024.
func Keys() ([]YearMonth, error) {
	var keys []YearMonth
	var err error
	var yearMonth YearMonth

	for year := startYear; year <= endYear; year++ {
		for month := 1; month <= 12; month++ {
			switch {
			case year == startYear && month < 9:
				// do nothing
			case year == endYear && month > 5:
				// do nothing
			default:
				yearMonth, err = NewYearMonth(year, month)
				if err != nil {
					return keys, err
				}
				keys = append(keys, yearMonth)
			}
		}
	}
	return keys, nil
}

// ----------------------------------------------------------------------------
// Methods
// ----------------------------------------------------------------------------

// String converts a year month to a string in the format MM/YYYY.
func (yearMonth YearMonth) String() string {
	var value = strconv.Itoa(yearMonth.Month) + "/" + strconv.Itoa(yearMonth.Year)
	return value
}
