package date

import (
	"fmt"
	"os"
	"strconv"
	"testing"
)

// ----------------------------------------------------------------------------
// Test Main
// ----------------------------------------------------------------------------

func TestMain(m *testing.M) {
	exitVal := m.Run()
	os.Exit(exitVal)
}

// ----------------------------------------------------------------------------
// Support functions
// ----------------------------------------------------------------------------

// handle checks an error return.  If it is not nil, it calls t.Fatalf to
// fail the test and print the error.
func handle(err error, t *testing.T) {
	if err != nil {
		t.Fatalf("%s\n", err)
	}
}

// ----------------------------------------------------------------------------
// Test definitional functions
// ----------------------------------------------------------------------------

// Test_isMonth checks that isMonth is validating integers representing months
// correctly.
func Test_isMonth(t *testing.T) {
	type aTest struct {
		name  string
		value Month
		error bool
	}
	var data = []aTest{
		{"too low value", 0, true},
		{"too high value", 13, true},
		{"valid value", 6, false},
	}

	var tt aTest
	var testFunction = func(t *testing.T) {
		err := isMonth(tt.value)
		if tt.error && (err == nil) {
			t.Error("isMonth should have reported error for " + strconv.Itoa(int(tt.value)))
		} else if !tt.error && err != nil {
			t.Error("isMonth incorrectly reported an error for " + strconv.Itoa(int(tt.value)))
		}
	}
	for _, d := range data {
		tt = d
		t.Run(d.name, testFunction)
	}
	return
}

// Test_IsLeapYear tests that IsLeapYear returns a valid answer for these conditions:
//
//	-- The year is divisible by 400
//	-- The year is divisible by 100 but not 400
//	-- The year is divisible by 4 but not 100
//	-- The year is not divisible by 4
func Test_IsLeapYear(t *testing.T) {
	type aTest struct {
		name       string
		year       Year
		isLeapYear bool
	}

	var data = []aTest{
		{name: "Divisible by 400", year: 2000, isLeapYear: true},
		{name: "Divisible by 100 but not 400", year: 1900, isLeapYear: false},
		{name: "Divisible by 4 but not 100", year: 2004, isLeapYear: true},
		{name: "Not divisible by 4", year: 2023, isLeapYear: false},
	}

	var tt aTest
	var testFunction = func(t *testing.T) {
		result := IsLeapYear(tt.year)
		if result != tt.isLeapYear {
			if tt.isLeapYear {
				t.Error("Year should have been identified as a leap year but was not: " +
					strconv.Itoa(int(tt.year)))
			} else {
				t.Error("Year should not have been identified as a leap year but was: " +
					strconv.Itoa(int(tt.year)))
			}
		}
		return
	}

	for _, d := range data {
		tt = d
		t.Run(d.name, testFunction)
	}
	return
}

// Test_DaysInYear checks the final results in DaysInYear with the sum of
// days in month.
func Test_DaysInYear(t *testing.T) {
	var month int
	var sum = 0
	for month = 1; month <= 12; month++ {
		sum += daysInMonth[month-1]
	}
	// Non-leap year case
	var actualDaysInYear, _ = DaysInYear(2023)
	if actualDaysInYear != sum {
		var message = "Actual days in year " + strconv.Itoa(actualDaysInYear) +
			" does not agree with sum of days in months " + strconv.Itoa(sum)
		t.Error(message)
	}

	// Leap year case
	actualDaysInYear, _ = DaysInYear(2024)
	if actualDaysInYear != 366 {
		var message = "Actual days in year " + strconv.Itoa(actualDaysInYear) +
			" does not agree with days in leap year " + strconv.Itoa(366)
		t.Error(message)
	}
	return
}

// Test_DayYear tests the computation of the day of the year.
func Test_DayYear(t *testing.T) {
	type aTest struct {
		name      string
		date      Date
		dayOfYear DayOfYear
	}

	// Test dates
	var date1, _ = New(6, 1, 2023)

	// Test data
	var data = []aTest{
		{"MinDate", MinDate, 1},
		{"MaxDate", MaxDate, 365},
		{"Other Date", date1, 152},
	}

	var tt aTest
	var testFunction = func(t *testing.T) {
		var dayOfYear = DayYear(tt.date)
		if dayOfYear != tt.dayOfYear {
			var message = "day of year for date " + tt.date.String() +
				" does not equal expected day of year " + strconv.Itoa(int(tt.dayOfYear))
			t.Error(message)
		}
		return
	}

	for _, tt = range data {
		t.Run(tt.name, testFunction)
	}
	return
}

// ----------------------------------------------------------------------------
// Test factory methods
// ----------------------------------------------------------------------------

// Test_FromDayOfYear tests that this factory method returns the right result.
func Test_FromDayOfYear(t *testing.T) {
	var year Year = 2024
	totalDays, _ := DaysInYear(year)

	var dayOfYear DayOfYear

	var testFunction = func(t *testing.T) {
		date, err := FromDayOfYear(dayOfYear, year)
		handle(err, t)
		dayOfYearOfDate := DayYear(date)
		if dayOfYear != dayOfYearOfDate {
			message := "days of year do not agree for dayOfYear: " + strconv.Itoa(int(dayOfYear))
			t.Error(message)
		}
	}

	for index := 1; index <= totalDays; index++ {
		dayOfYear = DayOfYear(index)
		name := "test for day " + strconv.Itoa(index)
		t.Run(name, testFunction)
	}
}

// TestToday tests that the Today function returns a valid date.
func Test_Today(t *testing.T) {
	var testFunction = func(t *testing.T) {
		date := Today()
		month := date.Month()
		day := date.Day()
		year := date.Year()
		fmt.Printf("Current date is: %d / %d / %d \n", month, day, year)
	}
	t.Run("Today Test", testFunction)
}

// ----------------------------------------------------------------------------
// Test date calculations
// ----------------------------------------------------------------------------

func Test_AddDate(t *testing.T) {
	var date, err1 = New(12, 22, 2023)
	handle(err1, t)
	var expected1, err2 = New(1, 1, 2024)
	handle(err2, t)
	var actual1, err3 = Add(date, 10)
	handle(err3, t)

	var expected2, err4 = New(12, 12, 2023)
	handle(err4, t)
	var actual2, err5 = Add(date, -10)
	handle(err5, t)

	var actual3, err6 = Add(date, 0)
	handle(err6, t)

	type aTest struct {
		name string
		want Date
		got  Date
	}

	var data = []aTest{
		{"Add 10", expected1, actual1},
		{"Subtract 10", expected2, actual2},
		{"Add 0", date, actual3},
	}

	var tt aTest
	var testFunction = func(t *testing.T) {
		if tt.want != tt.got {
			t.Error("Expected date " + tt.want.String() + " but actual date " + tt.got.String())
		}
	}

	for _, d := range data {
		tt = d
		t.Run(d.name, testFunction)
	}
}

// Test_AddError tests the response of the Add function when an amount is
// added to a date which would make it later than MaxDate.
func Test_AddError(t *testing.T) {
	var _, err = Add(MaxDate, 10)
	if err == nil {
		t.Fatalf("Add() failed to detect date overflow")
	} else {
		fmt.Println(err)
	}
	_, err = Add(MinDate, -10)
	if err == nil {
		t.Fatalf("Add() failed to detect date underflow")
	} else {
		fmt.Println(err)
	}
}

// Test_Difference tests the computation of the number of days between two
// dates.
func Test_Difference(t *testing.T) {
	var date = MinDate

	var testFunction = func(t *testing.T) {
		for i := 0; i < int(MaxAbsoluteDate); i++ {
			var newDate, err0 = Add(date, i)
			handle(err0, t)
			var diff1 = Difference(newDate, date)
			if diff1 != i {
				t.Fatalf("Diff %d does not equal i %d\n", diff1, i)
			}
			var diff2 = Difference(date, newDate)
			if -diff2 != i {
				t.Fatalf("Diff %d dones not equal -%d\n", diff2, i)
			}
		}
	}

	t.Run("Computation of difference", testFunction)
}

// ----------------------------------------------------------------------------
// Test comparison functions
// ----------------------------------------------------------------------------

// Test_Compare test the comparison functions.
func Test_Compare(t *testing.T) {
	var date1, err1 = New(12, 22, 2023)
	handle(err1, t)
	var date2, err2 = New(12, 23, 2023)
	handle(err2, t)
	if date1.Compare(date2) != BEFORE {
		t.Fatalf("Comparison: date1 < date2 not BEFORE")
	}
	if date2.Compare(date1) != AFTER {
		t.Fatalf("Compareison: date2 > date1 not AFTER")
	}
	if date2.Compare(date2) != EQUAL {
		t.Fatalf("Comoparison: date2 == date2 not EQUAL")
	}
	if !date1.Before(date2) {
		t.Fatalf("Before: date1 < date2 not true")
	}
	if !date2.After(date1) {
		t.Fatalf("After: date2 > date1 not true")
	}
	if Max(date1, date2) != date2 {
		t.Fatalf("Max: date2 was not returned")
	}
	if Min(date1, date2) != date1 {
		t.Fatalf("Min: date1 was not selected")
	}
}

// ----------------------------------------------------------------------------
// Test absolute date functions
// ----------------------------------------------------------------------------

// Test_ConvertToDate tests the conversion from absolute date to WayDate.
func Test_ConvertToDate(t *testing.T) {
	type aTest struct {
		name         string
		absoluteDate AbsoluteDate
		date         Date
	}
	var dec19, _ = New(12, 19, 2023)
	var data = []aTest{
		{"MinDate", MinAbsoluteDate, MinDate},
		{"MaxDate", MaxAbsoluteDate, MaxDate},
		{"19-Dec-2023", 154485, dec19},
	}

	var tt aTest
	var testFunction = func(t *testing.T) {
		var absoluteDate = tt.absoluteDate
		var expectedDate = tt.date
		var actualDate, err = convertToDate(absoluteDate)
		if err != nil {
			t.Fatalf("Unexpecte error: %s", err)
		}
		if expectedDate != actualDate {
			t.Error("Expected date not equal to actual date")
		}
	}

	for _, d := range data {
		tt = d
		t.Run(d.name, testFunction)
	}
}

// Test_ConvertToAbsolute tests the conversion of a date to an absolute date.
func Test_ConvertToAbsolute(t *testing.T) {
	type aTest struct {
		name         string
		absoluteDate AbsoluteDate
		date         Date
	}
	var dec19, _ = New(12, 19, 2023)
	var data = []aTest{
		{"MinDate", MinAbsoluteDate, MinDate},
		{"MaxDate", MaxAbsoluteDate, MaxDate},
		{"19-Dec-2023", 154485, dec19},
	}

	var tt aTest
	var testFunction = func(t *testing.T) {
		var date = tt.date
		var expectedAbsoluteDate = tt.absoluteDate
		var actualAbsoluteDate, err = convertToAbsolute(date)
		if err != nil {
			t.Fatalf("Unexpecte error: %s", err)
		}
		if expectedAbsoluteDate != actualAbsoluteDate {
			t.Error("Expected absolute date not equal to actual absolute date")
		}
	}

	for _, d := range data {
		tt = d
		t.Run(d.name, testFunction)
	}
}

// ----------------------------------------------------------------------------
// Test day of week
// ----------------------------------------------------------------------------

// Test_WeekDay tests the calculation of a week day from a date.
func Test_WeekDay(t *testing.T) {
	type aTest struct {
		name      string
		date      Date
		dayOfWeek DayOfWeek
	}

	var sunday, err1 = New(12, 24, 2023)
	handle(err1, t)
	var saturday, err2 = New(12, 30, 2023)
	handle(err2, t)

	var data = []aTest{
		{"Sunday", sunday, SUNDAY},
		{"Saturday", saturday, SATURDAY},
	}

	var tt aTest
	var testFunction = func(t *testing.T) {
		var weekDay, err = tt.date.WeekDay()
		handle(err, t)
		if weekDay != tt.dayOfWeek {
			t.Fatalf("Date %s has incorrect day of week: %d\n", tt.date.String(), weekDay)
		}
	}

	for _, d := range data {
		tt = d
		t.Run(d.name, testFunction)
	}
	return
}
