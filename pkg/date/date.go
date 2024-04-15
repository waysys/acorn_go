// Package date implements the date algorithms specified by DateBench at
// https://waysysweb.com/waysys/datebench.html.
// The Date type is intended to be invariant.  Use the operations in
// this package to manipulate the dates, rather than directly changing the fields
// in the structure.
package date

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	"acorn_go/pkg/assert"
	"errors"
	"strconv"
	"strings"
	"time"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// Month represents  the number of a month in a date.
//
//	1 <= month <= 12
type Month int

// Day represents the number of the day in a month.
//
//	1 <= day <= 31
type Day int

// Year represents the range of valid years.
//
//	MinYear <= year <= MaxYear
type Year int

// DayOfYear represents the range of days in a year.
//
//	1 <= dayOfYear <= DaysInYear(year)
type DayOfYear int

type Order int

// Date represents a date in the Gregorian calendar
type Date struct {
	month Month
	day   Day
	year  Year
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

const (
	MaxYear = 3999
	MinYear = 1601
)

var daysInMonth = []int{
	31, // January
	28, // February
	31, // March
	30, // April
	31, // May
	30, // June
	31, // July
	31, // August
	30, // September
	31, // October
	30, // November
	31, // December
}

var namesMonth = []string{
	"Jan",
	"Feb",
	"Mar",
	"Apr",
	"May",
	"Jun",
	"Jul",
	"Aug",
	"Sep",
	"Oct",
	"Nov",
	"Dec",
}

var MaxDate, _ = New(12, 31, MaxYear)
var MinDate, _ = New(1, 1, MinYear)

const (
	BEFORE Order = -1
	EQUAL  Order = 0
	AFTER  Order = 1
)

// ----------------------------------------------------------------------------
// Type validation functions
// ----------------------------------------------------------------------------

func isMonth(month Month) error {
	var err error
	switch {
	case month < 1:
		err = errors.New("Month cannot be less than 1: " + strconv.Itoa(int(month)))
	case month > 12:
		err = errors.New("Month cannot be greater than 12: " + strconv.Itoa(int(month)))
	default:
		err = nil
	}
	return err
}

func isYear(year Year) error {
	var err error
	switch {
	case year < MinYear:
		err = errors.New("Year cannot be less than MinYear: " + strconv.Itoa(int(year)))
	case year > MaxYear:
		err = errors.New("Year cannot be greater than MaxYear: " + strconv.Itoa(int(year)))
	default:
		err = nil
	}
	return err
}

// isDay return nil if the day is a valid day of the month (1 <= day <= DaysInMonth(month, year).
// Otherwise, it returns an error.
func isDay(month Month, day Day, year Year) error {
	maxDays, err := DaysInMonth(month, year)
	switch {
	case err != nil:
		break
	case day < 1:
		err = errors.New("Day cannot be less than 1: " + strconv.Itoa(int(day)))
	case int(day) > maxDays:
		err = errors.New("Day cannot be greater than the number of days in the month: " + strconv.Itoa(int(day)))
	default:
		err = nil
	}
	return err
}

// IsDate return nil if the month, day, and year are valid values.  Otherwise, an error is
// returned.
func IsDate(month Month, day Day, year Year) error {
	var err error
	err = isMonth(month)
	if err != nil {
		return err
	}
	err = isYear(year)
	if err != nil {
		return err
	}
	err = isDay(month, day, year)
	return err
}

// IsADate returns nil if the date has valid components.  Otherwise, an error is
// returned.
func IsADate(date Date) error {
	return IsDate(date.month, date.day, date.year)
}

// isDayOfYear return nil if the dayOfYear is a valid day of year.  Otherwise, it
// returns an error.
func isDayOfYear(dayOfYear DayOfYear, year Year) error {
	daysInYear, err := DaysInYear(year)
	switch {
	case err != nil:
		break
	case dayOfYear < 1:
		err = errors.New("Day of year cannot be less than 1: " + strconv.Itoa(int(dayOfYear)))
	case int(dayOfYear) > daysInYear:
		err = errors.New("Day of year cannot be greater than the number of days in the year: " +
			strconv.Itoa(int(dayOfYear)))
	}
	return err
}

// ----------------------------------------------------------------------------
// Date value functions
// ----------------------------------------------------------------------------

// IsLeapYear returns true if the year is a leap year, that is it has 366 days in the year.
// This function is definitional.
func IsLeapYear(year Year) bool {
	assert.Precondition(isYear(year))
	var result bool
	switch {
	case year%400 == 0:
		result = true
	case year%100 == 0:
		result = false
	case year%4 == 0:
		result = true
	default:
		result = false
	}
	return result
}

// DaysInMonth returns the number of days in the specified month and year.
// This function is definitional.
func DaysInMonth(month Month, year Year) (int, error) {
	var err error
	var days int

	// Preconditions
	err = isMonth(month)
	if err != nil {
		return 0, err
	}
	err = isYear(year)
	if err != nil {
		return 0, err
	}

	// Calculations
	days = daysInMonth[month-1]
	if IsLeapYear(year) && (month == 2) {
		// add a day if it is a leap year and the month is February.
		days++
	}

	var postCondition = func() error {
		var err error
		switch {
		case days < 28:
			err = errors.New("days in month amount is less than 28" + strconv.Itoa(days))
		case days > 31:
			err = errors.New("days in month amount is greater than 31" + strconv.Itoa(days))
		default:
			err = nil
		}
		return err
	}
	assert.Postcondition(postCondition())
	return days, nil
}

// DaysInYear returns the number of days in the specified year.
func DaysInYear(year Year) (int, error) {
	var err error
	var days int
	err = isYear(year)
	if err != nil {
		return 0, err
	}
	if IsLeapYear(year) {
		days = 366
	} else {
		days = 365
	}
	// Postcondition:
	//   (isLeapYear(year) and days = 366) or (not isLeapYear(year) and days = 365)
	return days, nil
}

// DayYear computes the day of the year for the specified date.
// This function is definitional.
func DayYear(date Date) DayOfYear {
	assert.Precondition(IsADate(date))
	var days = daysInPriorMonths(date.month, date.year) + int(date.day)
	// Postcondition: 1 <= days <= DaysInYear(date.year)
	return DayOfYear(days)
}

// daysInPriorMonths returns the number of days in all the months prior to
// the specified month for the specified year.
func daysInPriorMonths(month Month, year Year) int {
	assert.Precondition(isMonth(month))
	assert.Precondition(isYear(year))

	// Invariant:
	//   totalDays = for 1 <= m < month : sum(DaysInMonth(m, year))
	// Bound Function: func(m, month) {month - m - 1)
	//
	totalDays := 0
	limit := int(month)
	var daysInMonth int
	// Invariant is true for month = 1
	for m := 1; m < limit; m++ {
		daysInMonth, _ = DaysInMonth(Month(m), year)
		totalDays += daysInMonth
		// Invariant true for m < month
	}
	// PostCondition: totalDays = for 1 <= m < month : sum(DaysInMonth(m, year))
	// PostCondition is true given:
	//    m = month and invariant
	return totalDays
}

// ----------------------------------------------------------------------------
// Date factory functions
// ----------------------------------------------------------------------------

// New returns a date based on the specified month, day, and year
func New(month Month, day Day, year Year) (Date, error) {
	// Precondition
	err := IsDate(month, day, year)
	if err != nil {
		return Date{}, err
	}

	date := Date{
		month: month,
		day:   day,
		year:  year,
	}
	// Postcondition:
	//   date.month = month and date.day = day and date.year = year and IsADate(date)
	return date, nil
}

// Today returns the current date.
func Today() Date {
	now := time.Now()
	month := now.Month()
	day := now.Day()
	year := now.Year()
	date, _ := New(Month(month), Day(day), Year(year))
	assert.Postcondition(IsADate(date))
	return date
}

// FromDayOfYear returns a date based on the day of the year and the specified
// year.
//
// Precondition:
//
//	isDayOfYear(dayOfYear) and isYear(year)
//
// Postcondition:
//
//	date == FromDayOfYear(dayOfYear, year)
//	DayYear(date) == dayOfYear
func FromDayOfYear(dayOfYear DayOfYear, year Year) (Date, error) {
	// Precondition
	err := isDayOfYear(dayOfYear, year)
	if err != nil {
		return Date{}, err
	}
	// Invariant:
	//    isMonth(month) && daysInPastMonths = daysInPastMonth(month, year) &&
	//    remainingDays + daysInPastMonth = dayOfYear &&
	//    remainingDays > 0
	//
	// Bound Function:
	//    func diff(remainingDAys, month) int {
	//        if (remainingDays <= daysInMonth(month)) {
	//           return 0
	//        } else {
	//           return remainingDays
	//        }
	//     }
	//
	var daysInPastMonths = 0
	var remainingDays = int(dayOfYear)
	var month Month = 1
	var daysInMonth int
	// Invariant is true:
	//    remainingDays == dayOfYear && daysInPastMonths = 0 &&
	//    remainingDays + daysInPastMonths == dayOfYear
	for {
		daysInMonth, _ = DaysInMonth(month, year)
		if remainingDays > daysInMonth {
			month++
			daysInPastMonths += daysInMonth
			remainingDays -= daysInMonth
			// Invariant is true:
			//    daysInPastMonths' == daysInPastMonths + daysInMonth &&
			//    remainingDays' == remainingDays - daysInMonth &&
			//    therefore: daysInPastMonths' + remainingDays' == daysInPastMonths + remainingDays == dayOfMonth
		} else {
			break
		}
	}
	// Bound Function: decreases each iteration until 0.
	//
	// Truths
	//   isMonth(month)
	//   isDay(Day(remainingDays))
	//   daysInPastMonths(month) + remainingDays == dayOfYear
	date, _ := New(month, Day(remainingDays), year)
	//
	// Postcondition is true:
	//   dayYear(date) = daysInPastMonth(month) + remainingDays = dayOfYear
	return date, nil
}

// NewFromString converts a string representation of a date in the form MM/DD/YYYY
// into a date
func NewFromString(value string) (Date, error) {
	var err error
	var date = Date{}
	var anInt int
	var month Month
	var day Day
	var year Year
	//
	// Check that format is reasonable
	//
	value = strings.TrimSpace(value)
	if value == "" {
		err = errors.New("String date must not be an empty string")
		return date, err
	}
	var data = strings.Split(value, "/")
	if len(data) != 3 {
		err = errors.New("Invalid date format for string: " + value)
		return date, err
	}
	//
	// Obtain month
	//
	anInt, err = strconv.Atoi(data[0])
	if err != nil {
		return date, err
	}
	month = Month(anInt)
	err = isMonth(month)
	if err != nil {
		return date, err
	}
	//
	// Obtain year
	//
	anInt, err = strconv.Atoi(data[2])
	if err != nil {
		return date, err
	}
	year = Year(anInt)
	err = isYear(year)
	if err != nil {
		return date, err
	}
	//
	// Obtain day
	//
	anInt, err = strconv.Atoi(data[1])
	if err != nil {
		return date, err
	}
	day = Day(anInt)
	err = isDay(month, day, year)
	if err != nil {
		return date, err
	}
	//
	// Create date
	//
	date, err = New(month, day, year)
	return date, err
}

// ----------------------------------------------------------------------------
// Date properties
// ----------------------------------------------------------------------------

// Month returns the month on a date
func (date Date) Month() Month {
	return date.month
}

// Day returns the day on a date
func (date Date) Day() Day {
	return date.day
}

// Year returns the year on a date.
func (date Date) Year() Year {
	return date.year
}

// ----------------------------------------------------------------------------
// Computation
// ----------------------------------------------------------------------------

// Increment return a date one day after the specified date.
// This function is definitional.
func (date Date) Increment() (Date, error) {
	var result Date
	var err error = nil
	var daysInMonth int

	assert.Precondition(IsADate(date))
	// Precondition: date < MaxDate
	if date == MaxDate {
		err = errors.New("cannot increment maximum date")
		return date, err
	}

	daysInMonth, err = DaysInMonth(date.month, date.year)
	if err != nil {
		return date, err
	}
	switch {
	case date.month == 12 && date.day == 31:
		result, err = New(1, 1, date.year+1)
	case int(date.day) == daysInMonth:
		result, err = New(date.month+1, 1, date.year)
	default:
		result, err = New(date.month, date.day+1, date.year)
	}
	// Postcondition: Difference(result, date) = 1
	return result, err
}

// Decrement returns a date one day earlier than the specified date.
func (date Date) Decrement() (Date, error) {
	var err error = nil
	var result Date
	var lastDay int

	assert.Precondition(IsADate(date))
	// Precondition: date > MinDate
	if date == MinDate {
		err = errors.New("cannot decrement minimum date")
		return date, err
	}

	switch {
	case date.day == 1 && date.month == 1:
		result, err = New(12, 31, date.year-1)
	case date.day == 1:
		lastDay, err = DaysInMonth(date.month-1, date.year)
		assert.Assert(err == nil, "Unexpected error from DaysInMonth")
		result, err = New(date.month-1, Day(lastDay), date.year)
	default:
		result, err = New(date.month, date.day-1, date.year)
	}

	var postCondition = func() error {
		var date1, err = result.Increment()
		switch {
		case err != nil:
			break
		case date1 != date:
			var message = "result date " + result.String() + " is not one less than original date " +
				date.String()
			err = errors.New(message)
		default:
			err = nil
		}
		return err
	}

	assert.Postcondition(postCondition())
	return result, err
}

// Add adds the number of days to the date if num > 0. Add subtracts the
// number of days from the date if num < 0.
func Add(date Date, num int) (Date, error) {
	var err error = nil
	var absoluteDate AbsoluteDate
	var resultDate Date

	assert.Precondition(IsADate(date))

	absoluteDate, err = convertToAbsolute(date)
	if err != nil {
		return date, err
	}
	var result = int(absoluteDate) + num
	if result < 1 {
		var message = "value " +
			strconv.Itoa(num) +
			" is too negative to subtract from date " + date.String()
		var err = errors.New(message)
		return date, err
	}
	if result > int(MaxAbsoluteDate) {
		var message = "value " +
			strconv.Itoa(num) +
			" is too large to add to date " + date.String()
		var err = errors.New(message)
		return date, err
	}
	resultDate, err = convertToDate(AbsoluteDate(result))
	if err != nil {
		return date, err
	}
	// Postcondition:
	//    err != nil or
	//    convertToAbsoluteDate(resultDate) = convertToAbsoluteDate(date) + num
	return resultDate, nil
}

// Difference returns the number of days between two dates.
// If date1 is before date2, the number is negative.
// if date1 is after date2, the number is positive.
func Difference(date1 Date, date2 Date) int {
	assert.Precondition(IsADate(date1))
	assert.Postcondition(IsADate(date2))

	var absoluteDate1, _ = convertToAbsolute(date1)
	var absoluteDate2, _ = convertToAbsolute(date2)
	var diff = int(absoluteDate1 - absoluteDate2)

	var postcondition = func() error {
		var date, err = Add(date2, diff)
		switch {
		case err != nil:
			break
		case date1 != date:
			err = errors.New("difference: postcondition failed")
		default:
			err = nil
		}
		return err
	}
	assert.Postcondition(postcondition())
	return diff
}

// ----------------------------------------------------------------------------
// Comparison Functions
// ----------------------------------------------------------------------------

// Compare determines if the date is before, equal, or after the specified
// argument.  This function is definitional.
func (date Date) Compare(anotherDate Date) Order {
	var result Order
	assert.Precondition(IsADate(date))

	switch {
	case date.year < anotherDate.year:
		result = BEFORE
	case date.year == anotherDate.year && date.month < anotherDate.month:
		result = BEFORE
	case date.year == anotherDate.year && date.month == anotherDate.month && date.day < anotherDate.day:
		result = BEFORE
	case date.year == anotherDate.year && date.month == anotherDate.month && date.day == anotherDate.day:
		result = EQUAL
	default:
		result = AFTER
	}
	return result
}

// After returns true if the date is after the specified argument.
// This function is definitional.
func (date Date) After(anotherDate Date) bool {
	var result = date.Compare(anotherDate) == AFTER
	return result
}

// Before returns true if the date is before the specified argument.
// This function is definitional.
func (date Date) Before(anotherDate Date) bool {
	var result = date.Compare(anotherDate) == BEFORE
	return result
}

// Max returns the latter of two dates.
// This function is definitional.
func Max(date1 Date, date2 Date) Date {
	var result Date
	if date1.After(date2) {
		result = date1
	} else {
		result = date2
	}
	return result
}

// Min returns the earlier of two dates.
// This function is definitional.
func Min(date1 Date, date2 Date) Date {
	var result Date
	if date1.Before(date2) {
		result = date1
	} else {
		result = date2
	}
	return result
}

// ----------------------------------------------------------------------------
// Display Functions
// ----------------------------------------------------------------------------

// String displays the date in a format dd-MMM-yyyy.
// This function is definitional.
func (date Date) String() string {
	var day string
	var month string
	var year string
	var dateAsString string
	if date.day < 10 {
		day = "0" + strconv.Itoa(int(date.day))
	} else {
		day = strconv.Itoa(int(date.day))
	}
	month = namesMonth[date.month-1]
	year = strconv.Itoa(int(date.year))
	dateAsString = day + "-" + month + "-" + year
	return dateAsString
}
