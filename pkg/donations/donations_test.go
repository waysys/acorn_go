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

// Package donorinfo performs the manipulaton of data for a single donor.
package donations

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	"os"
	"strconv"
	"testing"

	a "acorn_go/pkg/accounting"

	dn "acorn_go/pkg/donors"

	d "github.com/waysys/waydate/pkg/date"

	dec "github.com/shopspring/decimal"

	"acorn_go/pkg/spreadsheet"
)

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

const inputFile = "/home/bozo/golang/acorn_go/data/donations.xlsx"
const tab = "Worksheet"

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
	var donor = dn.NewDonorWithDonation("donor")

	var testFunction = func(t *testing.T) {
		if donor.Name() != "donor" {
			t.Error("Donor name should be 'donor' not: " + donor.Name())
		}
		if !donor.Donation(a.FY2023).Equal(dec.Zero) {
			t.Error("FY23 donation was not set to zero: " + donor.Donation(a.FY2023).String())
		}
		if !donor.Donation(a.FY2023).Equal(dec.Zero) {
			t.Error("FY24 donation was not set to zero: " + donor.Donation(a.FY2023).String())
		}
	}

	t.Run("Test_New", testFunction)
}

// Test_Total checks that the total of donations is calculated correctly.
func Test_Total(t *testing.T) {
	var donor = dn.NewDonorWithDonation("donor")
	var ten = dec.NewFromInt(10)
	var twenty = dec.NewFromInt(20)
	var thirty = twenty.Add(ten)
	(&donor).AddDonation(ten, a.FY2023)
	(&donor).AddDonation(twenty, a.FY2024)

	var testFunction = func(t *testing.T) {
		if !donor.TotalDonation().Equal(thirty) {
			t.Error("Total donation does not equal 30: " + donor.TotalDonation().String())
		}
	}

	t.Run("Test_Total", testFunction)
}

// Test_Add checks that donations are added properly to the donor structure.
func Test_Add(t *testing.T) {
	var donor = dn.NewDonorWithDonation("donor")
	var ten = dec.NewFromInt(10)
	var twenty = dec.NewFromInt(20)
	var thirty = twenty.Add(ten)
	var forty = twenty.Add(twenty)
	(&donor).AddDonation(ten, a.FY2023)
	(&donor).AddDonation(twenty, a.FY2023)
	(&donor).AddDonation(twenty, a.FY2024)
	(&donor).AddDonation(twenty, a.FY2024)

	var testFunction = func(t *testing.T) {
		if !donor.Donation(a.FY2023).Equal(thirty) {
			t.Error("FY2023 donation does not equal 30: " + donor.Donation(a.FY2023).String())
		}
		if !donor.Donation(a.FY2024).Equal(forty) {
			t.Error("FY2024 donation does not equal 40: " + donor.Donation(a.FY2024).String())
		}
	}

	t.Run("Test_Add", testFunction)
}

// Test_DonorStatus checks that donors are identified properly by the fiscal year
// they donated.
func Test_DonorStatus(t *testing.T) {
	var donorFY23 = dn.NewDonorWithDonation("donorFY23")
	(&donorFY23).AddDonation(dec.NewFromInt(100), a.FY2023)
	var donorFY24 = dn.NewDonorWithDonation("donorFY24")
	(&donorFY24).AddDonation(dec.NewFromInt(100), a.FY2024)
	var donorBoth = dn.NewDonorWithDonation("donorBoth")
	(&donorBoth).AddDonation(dec.NewFromInt(90), a.FY2023)
	(&donorBoth).AddDonation(dec.NewFromInt(80), a.FY2024)

	var testFunction = func(t *testing.T) {
		//
		// FY2023 donor
		//
		if !donorFY23.IsDonor(a.FY2023) {
			t.Error("donorFY23 not recognized as a FY2023 donor.")
		}
		if donorFY23.IsDonor(a.FY2024) {
			t.Error("donorFY23 incorrectly identified as FY2024 donor.")
		}
		if donorFY23.IsDonor(a.FY2025) {
			t.Error("donorFY23 incorrectly identified as FY2025 donor.")
		}
		//
		// FY2024 donor
		//
		if donorFY24.IsDonor(a.FY2023) {
			t.Error("donorFY24 incorrectly identified as a FY2023 donor.")
		}
		if !donorFY24.IsDonor(a.FY2024) {
			t.Error("donorFY24 failed to be identified as FY2024 donor.")
		}
		if donorFY24.IsDonor(a.FY2025) {
			t.Error("donorFY24 incorrectly identified as FY2025 donor.")
		}
		//
		// FY2023 and FY2024 donor
		//
		if !donorBoth.IsDonor(a.FY2023) {
			t.Error("donorBoth incorrectly identified as not a FY2023 donor.")
		}
		if !donorBoth.IsDonor(a.FY2024) {
			t.Error("donorBoth incorrectly identified as not a FY2024 donor.")
		}
		if donorBoth.IsDonor(a.FY2025) {
			t.Error("donorBoth incorrectly identified as FY2025 donor.")
		}
	}

	t.Run("Test_DonorStatus", testFunction)
}

// Test_FiscalYear checks the determination of the Fiscal Year Indicator
func Test_FiscalYear(t *testing.T) {
	var date d.Date

	date, _ = d.New(1, 1, 1950)
	var indicator = a.FiscalYearIndicator(date)
	if indicator != a.OutOfRange {
		t.Error("fiscal year indicator should be OutOfRange, but is: " + indicator.String())
	}

	date, _ = d.New(9, 1, 2022)
	indicator = a.FiscalYearIndicator(date)
	if indicator != a.FY2023 {
		t.Error("fiscal year idicator should be FY2023, but is: " + indicator.String())
	}

	date, _ = d.New(8, 31, 2024)
	indicator = a.FiscalYearIndicator(date)
	if indicator != a.FY2024 {
		t.Error("fiscal year indiator should be FY2024, but is: " + indicator.String())
	}
}

// Test_FiscalYearFromYearMonth checks the determination of the Fiscal Year Indicator based
// on a specified year and month
func Test_FiscalYearFromYearMonth(t *testing.T) {
	var indicator a.FYIndicator

	var evaluate = func(month int, year int) a.FYIndicator {
		var err error = nil
		var yearMonth d.YearMonth

		yearMonth, err = d.NewYearMonth(year, month)
		if err != nil {
			t.Error("year month not created: " + err.Error())
			return a.OutOfRange
		}
		indicator, err = a.FiscalYearFromYearMonth(yearMonth)
		if err != nil {
			t.Error("fiscal year indicator not create: " + err.Error())
			return a.OutOfRange
		}
		return indicator
	}
	//
	// Beginning of FY2023
	//
	indicator = evaluate(9, 2022)
	if indicator != a.FY2023 {
		t.Error("indicator is not correct value: " + indicator.String())
	}
	//
	// End of FY2023
	//
	indicator = evaluate(8, 2023)
	if indicator != a.FY2023 {
		t.Error("indicator is not correct valaue: " + indicator.String())
	}
	//
	// Beginning of FY2024
	//
	indicator = evaluate(9, 2023)
	if indicator != a.FY2024 {
		t.Error("indicator is not correct value: " + indicator.String())
	}
	//
	// End of FY2024
	//
	indicator = evaluate(8, 2024)
	if indicator != a.FY2024 {
		t.Error("indicator is not correct value: " + indicator.String())
	}

}

// Test_DonorCountAnalysis tests the create of the donor count analysis.
func Test_DonorCountAnalysis(t *testing.T) {
	var err error = nil
	var sprdsht spreadsheet.Spreadsheet
	var donationList DonationList
	//
	// Obtain spreadsheet data
	//
	sprdsht, err = spreadsheet.ProcessData(inputFile, tab)
	if err != nil {
		t.Error("Error reading spreadsheet: " + err.Error())
	}
	//
	// Obtain donation list
	//
	donationList, err = NewDonationList(&sprdsht)
	if err != nil {
		t.Error("Error creating donation list")
	}
	//
	// Create donor count analysis
	//
	var dca = ComputeDonorCount(&donationList)
	//
	// Test Function
	//
	var testFunction = func(t *testing.T) {
		var dcfy2023 = dca[a.FY2023]
		var count = dcfy2023.TotalDonorCount()
		if count != 125 {
			t.Error("Incorrect total FY2023 donor count: " + strconv.Itoa(count))
		}
		var dcfy2024 = dca[a.FY2024]
		count = dcfy2024.TotalDonorCount()
		if count != 156 {
			t.Error("Incorrect total FY2024 donor count: " + strconv.Itoa(count))
		}
		var dcfy2025 = dca[a.FY2025]
		count = dcfy2025.TotalDonorCount()
		if count != 0 {
			t.Error("Incorrect total FY2025 donor count: " + strconv.Itoa(count))
		}
	}

	t.Run("Test_DonorCountAnalysis", testFunction)
}
