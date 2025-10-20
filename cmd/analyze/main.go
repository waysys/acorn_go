// ----------------------------------------------------------------------------
//
// # Donor analysis program
//
// This program produces the analysis spreadsheet.  The tabs in the spreadsheet are:
// -- Donor Count
// -- Donation Analysis
// -- Average Donations
//
// Author: William Shaffer
//
//	Copyright (c) 2024, 2025 Acorn Scholarship Fund All Rights Reserved
//
// ----------------------------------------------------------------------------
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

const inputFile = "/home/bozo/golang/acorn_go/data/donations.xlsx"
const tab = "Worksheet"
const outputFile = "/home/bozo/Downloads/analysis.xlsx"

// Output tabs
const (
	donorCount       = "Donor Count"
	donationAnalysis = "Donation Analysis"
	averageDonations = "Average Donations"
)

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
	// Open the output spreadsheet
	//
	output, err = s.New(outputFile, donorCount)
	sp.Check(err, "Error opening output file: ")
	//
	// Defer function to close spreadsheet
	//
	var finish = func() {
		err = output.Save()
		sp.Check(err, "Error saving output file: ")
		err = output.Close()
		sp.Check(err, "Error closing output file: ")
		os.Exit(0)
	}
	defer finish()
	//
	// Output the donor count
	//
	var donorCountAnalysis = dn.ComputeDonorCount(donationList)
	outputDonorCount(donorCountAnalysis, &output)
	//
	// Output the donations
	//
	output, err = output.AddSheet(donationAnalysis)
	sp.Check(err, "Error adding sheet: ")
	var donationAnalysis = dn.ComputeDonations(donationList)
	outputDonations(donationAnalysis, &output)
	//
	// Output average donations
	//
	output, err = output.AddSheet(averageDonations)
	sp.Check(err, "Error adding sheet: ")
	var dac = dn.NewDonationsAndCounts(donationAnalysis, donorCountAnalysis)
	outputAvgDonations(&dac, &output)
	//
	// Print completion notice
	//
	printFooter()
}

// ----------------------------------------------------------------------------
// Print Functions
// ----------------------------------------------------------------------------

// printHeader outputs a notice that the program has started
func printHeader() {
	fmt.Println("-----------------------------------------------------------")
	fmt.Println("Acorn Scholarship Fund Donation Analysis")
	fmt.Println("-----------------------------------------------------------")
}

// printFooter outputs a notice that the program has completed
func printFooter() {
	fmt.Println("-----------------------------------------------------------")
	fmt.Println("Program is completed")
	fmt.Println("-----------------------------------------------------------")
}

// ----------------------------------------------------------------------------
// Output Functions
// ----------------------------------------------------------------------------

// outputDonorCount produces the donor count tab.
func outputDonorCount(
	donorCountAnalysis dn.DonorCountAnalysis,
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
	s.WriteCell(output, "F", row, "Retention Percent")
	s.WriteCell(output, "G", row, "Acquisition Percent")
	row++
	s.WriteCell(output, "A", row, "Year of Donation")
	row++
	for _, dc := range donorCountAnalysis {
		s.WriteCell(output, "A", row, dc.FiscalYear())
		s.WriteCellInt(output, "B", row, dc.Count(dn.PriorPriorYear))
		s.WriteCellInt(output, "C", row, dc.Count(dn.PriorYear))
		s.WriteCellInt(output, "D", row, dc.Count(dn.CurrentYear))
		s.WriteCellInt(output, "E", row, dc.TotalDonorCount())
		s.WriteCellFloat(output, "F", row, donorCountAnalysis.Retention(dc.FY()))
		s.WriteCellFloat(output, "G", row, donorCountAnalysis.Acquisition(dc.FY()))
		row++
	}
	row += 2
	s.WriteCell(output, "A", row, "Total Donors")
	s.WriteCellInt(output, "B", row, donorCountAnalysis.TotalDonors())
}

// outputDonations produces the donation analysis tab
func outputDonations(
	donationAnalysis dn.DonationAnalysis,
	output *s.SpreadsheetFile) {
	var row int = 1
	//
	// Place Title
	//
	s.WriteCell(output, "A", row, "Donation Analysis")
	//
	// Place Headings
	//
	row += 2
	s.WriteCell(output, "A", row, "Fiscal Year")
	s.WriteCell(output, "B", row, "Prior Prior Year")
	s.WriteCell(output, "C", row, "Prior Year")
	s.WriteCell(output, "D", row, "Current Year")
	s.WriteCell(output, "E", row, "Total Donation for Year")
	row++
	s.WriteCell(output, "A", row, "Year of Donation")
	row++
	for _, don := range donationAnalysis {
		s.WriteCell(output, "A", row, don.FiscalYear())
		s.WriteCellFloat(output, "B", row, don.Donation(dn.PriorPriorYear))
		s.WriteCellFloat(output, "C", row, don.Donation(dn.PriorYear))
		s.WriteCellFloat(output, "D", row, don.Donation(dn.CurrentYear))
		s.WriteCellFloat(output, "E", row, don.TotalDonations())
		row++
	}
	row += 2
	s.WriteCell(output, "A", row, "Total Donations")
	s.WriteCellFloat(output, "B", row, donationAnalysis.TotalDonations())
}

// outputAvgDonations produces the average donation tab.
func outputAvgDonations(
	dac *dn.DonationsAndCounts,
	output *s.SpreadsheetFile) {
	var row int = 1
	//
	// Place Title
	//
	s.WriteCell(output, "A", row, "Average Donation Analysis")
	//
	// Place Headings
	//
	row += 2
	s.WriteCell(output, "A", row, "Prior Years")
	s.WriteCell(output, "B", row, "Prior Prior Year")
	s.WriteCell(output, "C", row, "Prior Year")
	s.WriteCell(output, "D", row, "Current Year")
	s.WriteCell(output, "E", row, "Total Average for Year")
	row++
	s.WriteCell(output, "A", row, "Year of Donation")
	row++
	for fy := a.FY2024; fy <= a.FY2026; fy++ {
		s.WriteCell(output, "A", row, fy.String())
		s.WriteCellFloat(output, "B", row, dac.AvgDonation(fy, dn.PriorPriorYear))
		s.WriteCellFloat(output, "C", row, dac.AvgDonation(fy, dn.PriorYear))
		s.WriteCellFloat(output, "D", row, dac.AvgDonation(fy, dn.CurrentYear))
		s.WriteCellFloat(output, "E", row, dac.AvgTotalDonationFiscalYear(fy))
		row++
	}
	row += 2
	s.WriteCell(output, "A", row, "Total Average Donations")
	s.WriteCellFloat(output, "B", row, dac.TotalAvgDonation())
}
