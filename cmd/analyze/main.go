// ----------------------------------------------------------------------------
//
// # Donor analysis program
//
// This program produces the analysis spreadsheet.  The tabs in the spreadsheet are:
// -- Donor Count
// -- Donation Analysis
//
// Author: William Shaffer
// Version: 24-Sep-2024
//
// # Copyright (c) 2024 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------
package main

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	"fmt"
	"os"

	dn "acorn_go/pkg/donations"
	s "acorn_go/pkg/spreadsheet"
	sp "acorn_go/pkg/support"
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
	var sprdsht s.Spreadsheet
	var donationList dn.DonationList
	var err error
	var output s.SpreadsheetFile

	printHeader()
	//
	// Obtain spreadsheet data
	//
	sprdsht, err = s.ProcessData(inputFile, tab)
	sp.Check(err, "Error processing spreadsheet: ")
	//
	// Obtain donation list
	//
	donationList, err = dn.NewDonationList(&sprdsht)
	sp.Check(err, "Error generating donor list: ")
	//
	// Obtain donation analysis
	//
	// var donationAnalysis = dn.ComputeDonations(&donationList)
	//
	// Output donor count
	//
	output, err = s.New(outputFile, "Donor Count")
	sp.Check(err, "Error opening output file: ")
	var finish = func() {
		err = output.Save()
		sp.Check(err, "Error saving output file: ")
		err = output.Close()
		sp.Check(err, "Error closing output file: ")
		os.Exit(0)
	}
	defer finish()
	var donorCountAnalysis = dn.ComputeDonorCount(&donationList)
	outputDonorCount(&donorCountAnalysis, &output)
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

// ----------------------------------------------------------------------------
// Output Functions
// ----------------------------------------------------------------------------

// outputDonorCount produces the donor count tab.
func outputDonorCount(
	donorCountAnalysis *dn.DonorCountAnalysis,
	output *s.SpreadsheetFile) {
	var row int = 1
	//
	// Place Title
	//
	s.WriteCell(output, "A", row, "Donor Count Analysis")
	//
	// Place Headings
	//
	row += 2
	s.WriteCell(output, "A", row, "Prior Years")
	s.WriteCell(output, "B", row, "Prior Prior Year")
	s.WriteCell(output, "C", row, "Prior Year")
	s.WriteCell(output, "D", row, "Current Year")
	s.WriteCell(output, "E", row, "Total Count for Year")
	row++
	s.WriteCell(output, "A", row, "Year of Donation")
	row++
	for _, dc := range *donorCountAnalysis {
		s.WriteCell(output, "A", row, dc.FiscalYear())
		s.WriteCellInt(output, "B", row, dc.Count(dn.PriorPriorYear))
		s.WriteCellInt(output, "C", row, dc.Count(dn.PriorYear))
		s.WriteCellInt(output, "D", row, dc.Count(dn.CurrentYear))
		s.WriteCellInt(output, "E", row, dc.TotalDonorCount())
		row++
	}
	row += 2
	s.WriteCell(output, "A", row, "Total Donors")
}
