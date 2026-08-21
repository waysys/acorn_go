// ----------------------------------------------------------------------------
//
// TermPeriod defines the date ranges for scholarship payments
//
//
// Author: William Shaffer
//
// Copyright (c) 2026 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

package main

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	d "github.com/waysys/waydate/pkg/date"
	r "github.com/waysys/waydate/pkg/daterange"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type TermPeriod string

// termPeriodEntry pairs a date range with the term period it represents.
// Go arrays/slices require every element to share one type, so a range and
// a term period (different types) cannot be combined into a plain array of
// arrays; a small struct is the type-safe equivalent of that pairing.
type termPeriodEntry struct {
	dateRange r.DateRange
	period    TermPeriod
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

const (
	PriorTerm TermPeriod = "Prior Term"
	Spring    TermPeriod = "Spring"
	Summer    TermPeriod = "Summer"
	Fall      TermPeriod = "Fall"
	Outside   TermPeriod = "Outside"
)

var sortOrderTermPeriod = map[TermPeriod]int{
	PriorTerm: 0,
	Spring:    1,
	Summer:    2,
	Fall:      3,
	Outside:   4,
}

// termPeriodTable holds one row per term period, pairing its date range
// with its designation. DetermineTermPeriod walks this table in order,
// so a range and its period always travel together and can't drift out
// of sync the way two separately maintained lists could.
//
// Date ranges indicate periods when scholarships are paid
// not actual seasons in the calendar.
var termPeriodTable = []termPeriodEntry{
	{mustRange(mustDate(9, 1, 2025), mustDate(10, 31, 2025)), PriorTerm},
	{mustRange(mustDate(11, 1, 2025), mustDate(3, 31, 2026)), Spring},
	{mustRange(mustDate(4, 1, 2026), mustDate(6, 30, 2026)), Summer},
	{mustRange(mustDate(7, 1, 2026), mustDate(8, 31, 2026)), Fall},
}

// ----------------------------------------------------------------------------
// Initialization
// ----------------------------------------------------------------------------

// mustDate constructs a Date from a literal, panicking on the
// (only possible) error of an invalid date.
func mustDate(month d.Month, day d.Day, year d.Year) d.Date {
	date, err := d.New(month, day, year)
	if err != nil {
		panic(err)
	}
	return date
}

// mustRange constructs a DateRange, panicking on the
// (only possible) error of an invalid range.
func mustRange(first, last d.Date) r.DateRange {
	dateRange, err := r.New(first, last)
	if err != nil {
		panic(err)
	}
	return dateRange
}

// ----------------------------------------------------------------------------
// Function
// ----------------------------------------------------------------------------

// DetermineTermPeriod returns the term period designation based
// on the bill date.
func DetermineTermPeriod(billDate d.Date) TermPeriod {
	for _, entry := range termPeriodTable {
		if entry.dateRange.InRange(billDate) {
			return entry.period
		}
	}
	return Outside
}

// CompareTermPeriod compares term periods using the sort order
// for term periods
func CompareTermPeriod(
	termPeriod1 TermPeriod,
	termPeriod2 TermPeriod,
) int {
	var order1 = sortOrderTermPeriod[termPeriod1]
	var order2 = sortOrderTermPeriod[termPeriod2]

	var result = 0
	if order1 < order2 {
		result = -1
	} else if order1 > order2 {
		result = 1
	}

	return result
}
