// This file performs tests on the range package.
package daterange

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	"acorn_go/pkg/date"
	"fmt"
	"os"
	"strconv"
	"testing"
)

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

var date1, _ = date.New(1, 1, 2023)
var date2, _ = date.New(2, 1, 2023)

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

// Test_New checks the creation of date ranges.
func Test_NewRange(t *testing.T) {
	// Case 1: valid range
	var dateRange DateRange
	var err error

	// Case 1: valid date range
	dateRange, err = New(date1, date2)
	if err != nil {
		t.Error("New could not create a valid date range: " + err.Error())
	} else {
		if dateRange.Size() != 31 {
			var message = "Incorrect size " + strconv.Itoa(dateRange.Size()) +
				" for date range: " + dateRange.String()
			t.Error(message)
		}
	}

	// Case 2: invalid date range
	_, err = New(date2, date1)
	if err == nil {
		t.Error("New did not detect invalid date range.")
	} else {
		fmt.Println(err)
	}
}

// Test_InRange checks the determination of a date in a date range.
func Test_InRange(t *testing.T) {
	var dateRange DateRange
	var err error
	var date3 date.Date

	dateRange, err = New(date1, date2)
	handle(err, t)
	var result = dateRange.InRange(date1)
	if !result {
		t.Error("InRange did not detect date1 in date range")
	}

	date3, err = date1.Decrement()
	handle(err, t)
	result = dateRange.InRange(date3)
	if result {
		t.Error("InRange says date3 is in date range")
	}

	date3, err = date2.Increment()
	handle(err, t)
	result = dateRange.InRange(date3)
	if result {
		t.Error("InRange says date3 is in date range")
	}
}
