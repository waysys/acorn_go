// ----------------------------------------------------------------------------
//
// Donor information Test
//
// Author: William Shaffer
//
// Copyright (c) 2024, 2025 William Shaffer All Rights Reserved
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
		if !donor.Donation(a.FY2025).Equal(dec.Zero) {
			t.Error("FY2025 donation was not set to zero: " + donor.Donation(a.FY2025).String())
		}
		if !donor.Donation(a.FY2026).Equal(dec.Zero) {
			t.Error("FY2026 donation was not set to zero: " + donor.Donation(a.FY2026).String())
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
	(&donor).AddDonation(ten, a.FY2025)
	(&donor).AddDonation(twenty, a.FY2026)

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
	(&donor).AddDonation(ten, a.FY2025)
	(&donor).AddDonation(twenty, a.FY2025)
	(&donor).AddDonation(twenty, a.FY2026)
	(&donor).AddDonation(twenty, a.FY2026)

	var testFunction = func(t *testing.T) {
		if !donor.Donation(a.FY2025).Equal(thirty) {
			t.Error("FY2023 donation does not equal 30: " + donor.Donation(a.FY2025).String())
		}
		if !donor.Donation(a.FY2026).Equal(forty) {
			t.Error("FY2024 donation does not equal 40: " + donor.Donation(a.FY2026).String())
		}
	}

	t.Run("Test_Add", testFunction)
}

// Test_DonorStatus checks that donors are identified properly by the fiscal year
// they donated.
func Test_DonorStatus(t *testing.T) {
	var donorFY24 = dn.NewDonorWithDonation("donorFY24")
	(&donorFY24).AddDonation(dec.NewFromInt(100), a.FY2024)
	var donorFY25 = dn.NewDonorWithDonation("donorFY25")
	(&donorFY25).AddDonation(dec.NewFromInt(100), a.FY2025)
	var donorBoth = dn.NewDonorWithDonation("donorBoth")
	(&donorBoth).AddDonation(dec.NewFromInt(90), a.FY2025)
	(&donorBoth).AddDonation(dec.NewFromInt(80), a.FY2026)

	var testFunction = func(t *testing.T) {
		//
		// FY2024 donor
		//
		if !donorFY24.IsDonor(a.FY2024) {
			t.Error("donorFY24 not recognized as a FY2024 donor.")
		}
		if donorFY24.IsDonor(a.FY2026) {
			t.Error("donorFY24 incorrectly identified as FY2024 donor.")
		}
		if donorFY24.IsDonor(a.FY2025) {
			t.Error("donorFY24 incorrectly identified as FY2025 donor.")
		}
		//
		// FY2025 donor
		//
		if donorFY25.IsDonor(a.FY2024) {
			t.Error("donorFY25 incorrectly identified as a FY2024 donor.")
		}
		if !donorFY25.IsDonor(a.FY2025) {
			t.Error("donorFY2025 failed to be identified as FY2025 donor.")
		}
		if donorFY25.IsDonor(a.FY2026) {
			t.Error("donorFY24 incorrectly identified as FY2025 donor.")
		}
		//
		// FY2025 and FY2026 donor
		//
		if !donorBoth.IsDonor(a.FY2025) {
			t.Error("donorBoth incorrectly identified as not a FY2025 donor.")
		}
		if !donorBoth.IsDonor(a.FY2026) {
			t.Error("donorBoth incorrectly identified as not a FY2026 donor.")
		}
		if donorBoth.IsDonor(a.FY2024) {
			t.Error("donorBoth incorrectly identified as FY2024 donor.")
		}
	}

	t.Run("Test_DonorStatus", testFunction)
}

// Test_DonorCountAnalysis tests the creation of the donor count analysis.
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
	var dca = ComputeDonorCount(donationList)
	//
	// Test Function
	//
	var testFunction = func(t *testing.T) {
		var dcfy2024 = dca[a.FY2024]
		var count = dcfy2024.TotalDonorCount()
		if count != 156 {
			t.Error("Incorrect total FY2024 donor count: " + strconv.Itoa(count))
		}
		var dcfy2025 = dca[a.FY2025]
		count = dcfy2025.TotalDonorCount()
		if count != 0 {
			t.Error("Incorrect total FY2025 donor count: " + strconv.Itoa(count))
		}
		var totalDonors = dca.TotalDonors()
		if totalDonors != 178 {
			t.Error("Incorrect total donor count: " + strconv.Itoa(count))
		}
	}

	t.Run("Test_DonorCountAnalysis", testFunction)
}

// Test_DonationAnalysis tests the creation of the donation analysis.
func Test_DonationAnalysis(t *testing.T) {
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
	var dna = ComputeDonations(&donationList)
	//
	// Test Function
	//
	var testFunction = func(t *testing.T) {
		var dcfy2024 = dna[a.FY2024]
		var amount = dcfy2024.TotalDonations()
		if amount != 194708.00 {
			t.Error("Incorrect total FY2024 donations: " + conv(amount))
		}
		var dcfy2025 = dna[a.FY2025]
		amount = dcfy2025.TotalDonations()
		if amount != 0.00 {
			t.Error("Incorrect total FY2025 donations: " + conv(amount))
		}
	}

	t.Run("Test_DonorCountAnalysis", testFunction)
}

// ----------------------------------------------------------------------------
// Support Functions
// ----------------------------------------------------------------------------

// conv converts a float64 to a string.
func conv(value float64) string {
	str := strconv.FormatFloat(value, 'f', -1, 64)
	return str
}
