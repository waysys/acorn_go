// ----------------------------------------------------------------------------
//
// Non-repeat donor analysis
//
// Author: William Shaffer
// Version: 7-May-2024
//
// Copyright (c) 2024 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

// This program generates a list of FY2024 donors who did not donate in
// FY2025 with the dates and amounts they donated.

package main

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	a "acorn_go/pkg/accounting"
	dn "acorn_go/pkg/donations"
	s "acorn_go/pkg/spreadsheet"
	sp "acorn_go/pkg/support"
	"fmt"
	"os"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

const (
	inputFile      = "/home/bozo/golang/acorn_go/data/donations.xlsx"
	tab            = "Worksheet"
	outputFileName = "/home/bozo/Downloads/nonrepeat.xlsx"
	sheetName      = "Non-Repeat Donors"

	fy = a.FY2025
)

// ----------------------------------------------------------------------------
// Functions
// ----------------------------------------------------------------------------

// main supervises the execution of this program.  It produces a spreadsheet
// with the list of non-repeat donors
func main() {
	var sprdsht s.Spreadsheet
	var donorList dn.DonationList
	var err error

	printHeader()
	//
	// Obtain spreadsheet data
	//
	sprdsht, err = s.ProcessData(inputFile, tab)
	sp.Check(err, "Error processing spreadsheet")
	//
	// Generate donor list
	//
	donorList, err = dn.NewDonationList(&sprdsht)
	sp.Check(err, "Error generating donor list: ")
	//
	// Output non-repeat donor
	//
	outputRetention(&donorList)
	os.Exit(0)
}

// ----------------------------------------------------------------------------
// Print Functions
// ----------------------------------------------------------------------------

// printHeader places the header information at the top of the page
func printHeader() {
	fmt.Println("-----------------------------------------------------------")
	fmt.Println("Acorn Scholarship Fund Retention Analysis")
	fmt.Println("-----------------------------------------------------------")
}

// ----------------------------------------------------------------------------
// Output Functions
// ----------------------------------------------------------------------------

// outputRetention produces a spreadsheet with the non-repeat donors names
// and donation dates and amounts listed.
func outputRetention(donorList *dn.DonationList) {
	var err error
	var output s.SpreadsheetFile
	//
	// Create output spreadsheet
	//
	output, err = s.New(outputFileName, sheetName)
	sp.Check(err, "Error opening output file: ")
	var finish = func() {
		err = output.Save()
		sp.Check(err, "Error saving output file")
		err = output.Close()
		sp.Check(err, "Error closing output file")
	}
	defer finish()
	//
	// Insert Heading
	//
	var row = 1
	s.WriteCell(&output, "A", row, "List of non-repeat donors: "+fy.String())
	row += 2
	s.WriteCell(&output, "A", row, "Donor Name")
	s.WriteCell(&output, "B", row, "Amount of Donation")
	row++
	//
	// Insert donor information
	//
	var names = donorList.DonorKeys()
	for _, name := range names {
		var donor = donorList.Get(name)
		if donor.IsNonRepeatDonor(fy) {
			var amount = donor.Donation(fy.Prior())
			s.WriteCell(&output, "A", row, name)
			s.WriteCellDecimal(&output, "B", row, amount)
			row++
		}
	}
}
