package date

import (
	"errors"
	"slices"
	"strconv"

	"acorn_go/pkg/assert"
)

// This file implements functions related to the day of the week.

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// DayOfWeek is the number of the day of the week.  Sunday is 0.
type DayOfWeek int

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

const (
	SUNDAY    DayOfWeek = 0
	MONDAY    DayOfWeek = 1
	TUESDAY   DayOfWeek = 2
	WEDNESDAY DayOfWeek = 3
	THURSDAY  DayOfWeek = 4
	FRIDAY    DayOfWeek = 5
	SATURDAY  DayOfWeek = 6
)

var weekDays = []DayOfWeek{
	SUNDAY,
	MONDAY,
	TUESDAY,
	WEDNESDAY,
	THURSDAY,
	FRIDAY,
	SATURDAY,
}

// ----------------------------------------------------------------------------
// Functions
// ----------------------------------------------------------------------------

// isDayOfWeek returns true if the argument is a valid day of the week, that is,
// 0 <= dayOfWeek < 7
func isDayOfWeek(dayOfWeek DayOfWeek) error {
	var err error = nil
	if slices.Contains(weekDays, dayOfWeek) {
		err = errors.New("day of week must be between 0 and 6, not " + strconv.Itoa(int(dayOfWeek)))
	}
	return err
}

// WeekDay returns the day of the week for a date.
// This function is definitional.
// This function aligns with actual days of the week for the Gregorian calendar
// because January 1, 1601 (absolute date 1) is a Monday.
func (date Date) WeekDay() (DayOfWeek, error) {
	// Precondition:
	//    IsADate(date)
	var absoluteDate, err = convertToAbsolute(date)
	if err != nil {
		return SUNDAY, err
	}
	var weekDay = DayOfWeek(absoluteDate % 7)
	// Postcondition:
	//   isDayOfWeek(weekDay) and
	//   WeekDay(date.Increment()) = daysOfWeek[(WeekDay(date) + 1) mod 7]
	return weekDay, nil
}

// DayOfWeekAfter returns the closest date after the specified date
// with the specified day of the week.
func (date Date) DayOfWeekAfter(dayOfWeek DayOfWeek) (Date, error) {
	var laterDate Date
	var err error
	var weekday DayOfWeek

	assert.Precondition(isDayOfWeek(dayOfWeek))

	laterDate, err = date.Increment()
	if err != nil {
		return date, err
	}

	for {
		weekday, err = laterDate.WeekDay()
		if err != nil {
			break
		}
		if weekday == dayOfWeek {
			break
		}
		laterDate, err = laterDate.Increment()
		if err != nil {
			break
		}
	}
	return laterDate, err
}

// LastWeekDayOfMonth returns the date of the last specified day of the week
// in a specified month and year.
func LastWeekDayOfMonth(month Month, year Year, dayOfWeek DayOfWeek) (Date, error) {
	var err error
	var lastDay int
	var result Date
	var dayWeek DayOfWeek

	assert.Precondition(isMonth(month))
	assert.Precondition(isYear(year))
	assert.Precondition(isDayOfWeek(dayOfWeek))

	lastDay, err = DaysInMonth(month, year)
	if err != nil {
		return MinDate, err
	}
	result, err = New(month, Day(lastDay), year)
	if err != nil {
		return MinDate, err
	}
	dayWeek, err = result.WeekDay()
	if err != nil {
		return result, err
	}

	for {
		if dayWeek == dayOfWeek {
			break
		}
		result, err = result.Decrement()
		if err != nil {
			break
		}
		dayWeek, err = result.WeekDay()
		if err != nil {
			break
		}
	}
	return result, err
}
