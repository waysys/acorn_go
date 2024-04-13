package date

import (
	"acorn_go/pkg/assert"
	"errors"
	"strconv"
)

// This file implements the calculations for absolute dates.

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// AbsoluteDate is the number of days from 31-Dec-1600.
//
// absoluteDate >= 1 and absoluteDate <= 876216
type AbsoluteDate int

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// MinAbsoluteDate is the absolute date for 1-Jan-1601
const MinAbsoluteDate AbsoluteDate = 1

// MaxAbsoluteDate is the absolute date for 31-Dec-3999
const MaxAbsoluteDate AbsoluteDate = 876216

const (
	daysIn400YearCycle = 400*365 + 100 - 3
	daysIn100YearCycle = 365*100 + 25 - 1
	daysIn4YearCycle   = 365*4 + 1
	daysIn1YearCycle   = 365
)

// ----------------------------------------------------------------------------
// Functions
// ----------------------------------------------------------------------------

// isAbsoluteDate returns an error if the absolute date is not valid.
func isAbsoluteDate(absoluteDate AbsoluteDate) error {
	var err error = nil
	if absoluteDate < MinAbsoluteDate || absoluteDate > MaxAbsoluteDate {
		err = errors.New("isAbsoluteDate: invalid absolute date: " + strconv.Itoa(int(absoluteDate)))
	}
	return err
}

// convertToAbsoluteDate converts a date into an absolute date.
// This function is definitional
func convertToAbsolute(date Date) (AbsoluteDate, error) {
	// Precondition
	var err = IsADate(date)
	if err != nil {
		return MinAbsoluteDate, err
	}

	// Definition
	var dayYear = int(DayYear(date))
	var pastDays = daysInPastYears(date.year)
	var absoluteDate = AbsoluteDate(pastDays + dayYear)

	// Postcondition:
	assert.Assert(1 <= absoluteDate && absoluteDate <= MaxAbsoluteDate,
		"convertToAbsolute: Absolute date is outside of bounds: "+strconv.Itoa(int(absoluteDate)))
	return absoluteDate, nil
}

// ConvertToDate converts an absolute date to a date (month, day, year)
func convertToDate(absoluteDate AbsoluteDate) (Date, error) {
	assert.Precondition(isAbsoluteDate(absoluteDate))

	var year = yearFromAbsolute(absoluteDate)
	var dayOfYear = DayOfYear(int(absoluteDate) - daysInPastYears(year))
	var date, err = FromDayOfYear(dayOfYear, year)

	var postCondition = func() error {
		var err error
		var converted AbsoluteDate

		converted, err = convertToAbsolute(date)
		if converted != absoluteDate {
			var message = "convertToDate: absolute date of converted date " + strconv.Itoa(int(converted)) +
				" is not equal original absolute date " + strconv.Itoa(int(absoluteDate))
			err = errors.New(message)
		}
		return err
	}

	assert.Postcondition(postCondition())
	return date, err
}

// yearFromAbsolute determines the year of the specified absolute date.
// This algorithm is from:
// Edward M. Reingold and Nachum Dershowitz, Calendrical Calculations:
// The Millennium Edition (Cambridge, UK: Cambridge University Press, 2001)
func yearFromAbsolute(absoluteDate AbsoluteDate) Year {
	var year Year
	var num400YearCycles = (absoluteDate - 1) / daysIn400YearCycle
	var remainder100 = (absoluteDate - 1) % daysIn400YearCycle
	var num100YearCycles = remainder100 / daysIn100YearCycle
	var remainder4 = remainder100 % daysIn100YearCycle
	var num4YearCycles = remainder4 / daysIn4YearCycle
	var remainder1 = remainder4 % daysIn4YearCycle
	var num1YearCycles = remainder1 / daysIn1YearCycle
	if (num100YearCycles == 4) || (num1YearCycles == 4) {
		year = Year(400*num400YearCycles + 100*num100YearCycles + 4*num4YearCycles + num1YearCycles + 1600)
	} else {
		year = Year(400*num400YearCycles + 100*num100YearCycles + 4*num4YearCycles + num1YearCycles + 1601)
	}
	return year
}

// daysInPastYears computes the number of days starting in 1-Jan-1601 and ending
// in 31-Dec-year-1.
//
// Note that daysInPastYears(MinYear) = 0.
//
// This algorithm is from:
// Edward M. Reingold and Nachum Dershowitz, Calendrical Calculations:
// The Millennium Edition (Cambridge, UK: Cambridge University Press, 2001) p. 53
func daysInPastYears(year Year) int {
	var y = int(year) - 1601
	var result = 365*y + // days in prior years if all years had 365 days
		y/4 - // plus julian leap days in prior years if all years divided by 4 were leap years
		y/100 + // minus prior century years if all years divisible by 100 were not leap years
		y/400 // plus years divisible by 400 which have leap days
	return result
}
