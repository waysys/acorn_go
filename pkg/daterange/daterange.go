// Package daterange implements date ranges.
// Structures in this package are intended to be invariant.
package daterange

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	d "acorn_go/pkg/date"
	"errors"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type DateRange struct {
	first d.Date
	last  d.Date
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

var errorRange = DateRange{}

// ----------------------------------------------------------------------------
// Functions
// ----------------------------------------------------------------------------

// New create a new date range from two dates.
func New(first d.Date, last d.Date) (DateRange, error) {
	var err error = nil

	// Preconditions
	err = d.IsADate(first)
	if err != nil {
		return errorRange, err
	}
	err = d.IsADate(last)
	if err != nil {
		return errorRange, err
	}
	if first.After(last) {
		var message = "daterange.New: first date " + first.String() +
			" must not be after last date " +
			last.String()
		err = errors.New(message)
		return errorRange, err
	}

	dateRange := DateRange{first, last}
	// Postcondition:
	//   err <> nil or
	//   dateRange.first = first and dateRange.last = last and not dateRange.first.After(dateRange.last)
	return dateRange, err
}

// IsDateRange returns an error is the date range is  not valid.
func IsDateRange(dateRange DateRange) error {
	var err error = nil
	err = d.IsADate(dateRange.first)
	if err == nil {
		err = d.IsADate(dateRange.last)
	}
	if err == nil && dateRange.first.After(dateRange.last) {
		var message = "IsDateRange: first date " + dateRange.first.String() +
			" must not be after last date " +
			dateRange.last.String()
		err = errors.New(message)
	}
	return err
}

// Overlaps returns true if the two date ranges overlap.
func Overlaps(dateRange1 DateRange, dateRange2 DateRange) bool {
	var result bool
	switch {
	case dateRange1.last.Before(dateRange2.first):
		result = false
	case dateRange1.first.After(dateRange2.last):
		result = false
	default:
		result = true
	}
	return result
}

// ----------------------------------------------------------------------------
// Methods
// ----------------------------------------------------------------------------

// First returns the first date in the date range.
func (dateRange DateRange) First() d.Date {
	return dateRange.first
}

// Last returns the last date in the date range.
func (dateRange DateRange) Last() d.Date {
	return dateRange.last
}

// Size returns the number of days in the date range
func (dateRange DateRange) Size() int {
	var size = d.Difference(dateRange.last, dateRange.first)
	return size
}

// InRange returns true if a date in within a date range.
func (dateRange DateRange) InRange(date d.Date) bool {
	var result bool

	switch {
	case date.Before(dateRange.first):
		result = false
	case date.After(dateRange.last):
		result = false
	default:
		result = true
	}
	return result
}

func (dateRange DateRange) String() string {
	var dateRangeString = "(" +
		dateRange.first.String() + "," +
		dateRange.last.String() + ")"
	return dateRangeString
}
