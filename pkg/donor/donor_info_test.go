// ----------------------------------------------------------------------------
//
// Donor information Test
//
// Author: William Shaffer
// Version: 13-Apr-2024
//
// Copyright (c) William Shaffer
//
// ----------------------------------------------------------------------------

// Package donor_info performs the manipulaton of data for a single donor.
package donor_info

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	"os"
	"testing"

	dec "github.com/shopspring/decimal"
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
/*
func handle(err error, t *testing.T) {
	if err != nil {
		t.Fatalf("%s\n", err)
	}
}
*/

// ----------------------------------------------------------------------------
// Test definitional functions
// ----------------------------------------------------------------------------

// Test_New checks that the New function sets the values of the Donor structure
// properly.
func Test_New(t *testing.T) {
	var donor = New("donor")

	var testFunction = func(t *testing.T) {
		if donor.NameDonor() != "donor" {
			t.Error("Donor name should be 'donor' not: " + donor.NameDonor())
		}
		if !donor.DonationFY23().Equal(ZERO) {
			t.Error("FY23 donation was not set to zero: " + donor.DonationFY23().String())
		}
		if !donor.DonationFY23().Equal(ZERO) {
			t.Error("FY24 donation was not set to zero: " + donor.DonationFY24().String())
		}
	}

	t.Run("Test_New", testFunction)
}

// Test_Total checks that the total of donations is calculated correctly.
func Test_Total(t *testing.T) {
	var donor = New("donor")
	var ten = dec.NewFromInt(10)
	var twenty = dec.NewFromInt(20)
	var thirty = twenty.Add(ten)
	(&donor).AddFY23(ten)
	(&donor).AddFY24(twenty)

	var testFunction = func(t *testing.T) {
		if !donor.TotalDonation().Equal(thirty) {
			t.Error("Total donation does not equal 30: " + donor.TotalDonation().String())
		}
	}

	t.Run("Test_Total", testFunction)
}

// Test_Add checks that donations are added properly to the donor structure.
func Test_Add(t *testing.T) {
	var donor = New("donor")
	var ten = dec.NewFromInt(10)
	var twenty = dec.NewFromInt(20)
	var thirty = twenty.Add(ten)
	var forty = twenty.Add(twenty)
	(&donor).AddFY23(ten)
	(&donor).AddFY23(twenty)
	(&donor).AddFY24(twenty)
	(&donor).AddFY24(twenty)

	var testFunction = func(t *testing.T) {
		if !donor.DonationFY23().Equal(thirty) {
			t.Error("FY2023 donation does not equal 30: " + donor.DonationFY23().String())
		}
		if !donor.DonationFY24().Equal(forty) {
			t.Error("FY2024 donation does not equal 40: " + donor.DonationFY24().String())
		}
	}

	t.Run("Test_Add", testFunction)
}
