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

	d "acorn_go/pkg/date"

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
// Test functions
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

// Test_DonorStatus checks that donors are identified properly by the fiscal year
// they donated.
func Test_DonorStatus(t *testing.T) {
	var donorFY23 = New("donorFY23")
	(&donorFY23).AddFY23(dec.NewFromInt(100))
	var donorFY24 = New("donorFY24")
	(&donorFY24).AddFY24(dec.NewFromInt(100))
	var donorBoth = New("donorBoth")
	(&donorBoth).AddFY23(dec.NewFromInt(90))
	(&donorBoth).AddFY24(dec.NewFromInt(80))

	var testFunction = func(t *testing.T) {
		//
		// FY2023 donor
		//
		if !donorFY23.IsFY23DonorOnly() {
			t.Error("donorFY23 not recognized as a FY2023 donor.")
		}
		if donorFY23.IsFY24DonorOnly() {
			t.Error("donorFY23 incorrectly identified as FY2024 donor.")
		}
		if donorFY23.IsFY23AndFY24Donor() {
			t.Error("donorFY23 incorrectly identified as both FY2023 and FY2024 donor.")
		}
		//
		// FY2024 donor
		//
		if donorFY24.IsFY23DonorOnly() {
			t.Error("donorFY24 incorrectly identified as a FY2023 donor.")
		}
		if !donorFY24.IsFY24DonorOnly() {
			t.Error("donorFY24 failed to be identified as FY2024 donor.")
		}
		if donorFY24.IsFY23AndFY24Donor() {
			t.Error("donorFY24 incorrectly identified as both FY2023 and FY2024 donor.")
		}
		//
		// FY2023 and FY2024 donor
		//
		if donorBoth.IsFY23DonorOnly() {
			t.Error("donorBoth incorrectly identified as a FY2023 donor.")
		}
		if donorBoth.IsFY24DonorOnly() {
			t.Error("donorBoth incorrectly identified as FY2024 donor.")
		}
		if !donorBoth.IsFY23AndFY24Donor() {
			t.Error("donorBoth failed to be identified as both FY2023 and FY2024 donor.")
		}
	}

	t.Run("Test_DonorStatus", testFunction)
}

// Test_FiscalYear checks the determination of the Fiscal Year Indicator
func Test_FiscalYear(t *testing.T) {
	var date d.Date

	date, _ = d.New(1, 1, 1950)
	var indicator = FiscalYearIndicator(date)
	if indicator != OutOfRange {
		t.Error("fiscal year indicator should be OutOfRange, but is: " + indicator.String())
	}

	date, _ = d.New(9, 1, 2022)
	indicator = FiscalYearIndicator(date)
	if indicator != FY2023 {
		t.Error("fiscal year idicator should be FY2023, but is: " + indicator.String())
	}

	date, _ = d.New(8, 31, 2024)
	indicator = FiscalYearIndicator(date)
	if indicator != FY2024 {
		t.Error("fiscal year indiator should be FY2024, but is: " + indicator.String())
	}
}
