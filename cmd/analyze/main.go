// ----------------------------------------------------------------------------
//
// Donor analysis program
//
// Author: William Shaffer
// Version: 12-Apr-2024
//
// Copyright (c) 2024 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

// This package processes donation information based on the payment transactons
// from the donations.xlsx spreadsheet obtained from QuickBooks.  The package will then output
// the results in another spreadsheet.  These are the tabs in the spreadsheet:
package main

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	"fmt"
	"os"

	dn "acorn_go/pkg/donations"
	"acorn_go/pkg/spreadsheet"
	s "acorn_go/pkg/support"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

const inputFile = "/home/bozo/golang/acorn_go/data/donations.xlsx"
const tab = "Worksheet"
const outputFile = "/home/bozo/Downloads/analysis.xlsx"

// ----------------------------------------------------------------------------
// Functions
// ----------------------------------------------------------------------------

// main supervises the processing of the donation data.
func main() {
	var sprdsht spreadsheet.Spreadsheet
	var donationList dn.DonationList
	var err error
	var output spreadsheet.SpreadsheetFile

	printHeader()
	//
	// Obtain spreadsheet data
	//
	sprdsht, err = spreadsheet.ProcessData(inputFile, tab)
	s.Check(err, "Error processing spreadsheet: ")
	//
	// Obtain donation list
	//
	donationList, err = dn.NewDonationList(&sprdsht)
	s.Check(err, "Error generating donor list: ")
	//
	// Obtain donation analysis
	//
	var donationAnalysis = dn.ComputeDonations(&donationList)
	//
	// Output donation analysis
	//
	output, err = spreadsheet.New(outputFile, "Donor Count")
	s.Check(err, "Error opening output file: ")
	var finish = func() {
		err = output.Save()
		s.Check(err, "Error saving output file: ")
		err = output.Close()
		s.Check(err, "Error closing output file: ")
		os.Exit(0)
	}
	defer finish()

}

// ----------------------------------------------------------------------------
// Print Functions
// ----------------------------------------------------------------------------

// printHeader places the header information at the top of the page
func printHeader() {
	fmt.Println("-----------------------------------------------------------")
	fmt.Println("Acorn Scholarship Fund Donation Analysis")
	fmt.Println("-----------------------------------------------------------")
}
