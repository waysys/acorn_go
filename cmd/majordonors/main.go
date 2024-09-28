// ----------------------------------------------------------------------------
//
// # Major Donor Analysis
//
// This program produces the major donor spreadsheet.  The tabs in the spreadsheet are:
// -- Major Donor Analysis
// -- Major Donor List by Year
//
// Author: William Shaffer
// Version: 27-Sep-2024
//
// # Copyright (c) 2024 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------
package main

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	a "acorn_go/pkg/accounting"
	dn "acorn_go/pkg/donations"
	md "acorn_go/pkg/majordonor"
	s "acorn_go/pkg/spreadsheet"
	sp "acorn_go/pkg/support"
	"fmt"
	"os"
)

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

const inputFile = "/home/bozo/golang/acorn_go/data/donations.xlsx"
const tab = "Worksheet"
const outputFile = "/home/bozo/Downloads/majordonor.xlsx"

var columns = []string{"B", "C", "D"}

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
	output, err = s.New(outputFile, "Major Donor Count")
	sp.Check(err, "Error opening output file: ")
	var finish = func() {
		err = output.Save()
		sp.Check(err, "Error saving output file: ")
		err = output.Close()
		sp.Check(err, "Error closing output file: ")
		os.Exit(0)
	}
	defer finish()
	//
	// Output the major donor count and donations
	//
	var majorDonor = md.ComputeMajorDonors(&donationList)
	outputMajorDonor(&majorDonor, &output)
	//
	// Output the donations
	//
	output, err = output.AddSheet("Major Donor List")
	sp.Check(err, "Error adding sheet: ")
	outputMajorList(&donationList, &output)
}

// ----------------------------------------------------------------------------
// Print Functions
// ----------------------------------------------------------------------------

// printHeader places the header information at the top of the page
func printHeader() {
	fmt.Println("-----------------------------------------------------------")
	fmt.Println("Acorn Scholarship Fund Major Donor Analysis")
	fmt.Println("-----------------------------------------------------------")
}

// ----------------------------------------------------------------------------
// Output Functions
// ----------------------------------------------------------------------------

// outputMajorDonor inserts the major donor data into the spreadsheet.
func outputMajorDonor(
	md *md.MajorDonor,
	output *s.SpreadsheetFile) {
	var row = 1
	//
	// Place Title
	//
	s.WriteCell(output, "A", row, "Major Donor Analysis")
	//
	// Place Headings
	//
	row += 2
	s.WriteCell(output, "A", row, "")
	for _, fy := range a.FYIndicators {
		s.WriteCell(output, columns[fy], row, fy.String())
	}
	//
	// Ouput Data
	//
	row += 2
	s.WriteCell(output, "A", row, "Major donors")
	for _, fy := range a.FYIndicators {
		s.WriteCellInt(output, columns[fy], row, md.MajorDonorCount(fy))
	}
	row++
	s.WriteCell(output, "A", row, "Major donor donations")
	for _, fy := range a.FYIndicators {
		s.WriteCellFloat(output, columns[fy], row, md.DonationsMajor(fy))
	}
	row++
	s.WriteCell(output, "A", row, "Average major donor donation")
	for _, fy := range a.FYIndicators {
		s.WriteCellFloat(output, columns[fy], row, md.AvgDonation(fy))
	}
	row++
	s.WriteCell(output, "A", row, "Percent of total donations by major donors")
	for _, fy := range a.FYIndicators {
		s.WriteCellFloat(output, columns[fy], row, md.PercentDonation(fy))
	}
	row += 2
	s.WriteCell(output, "A", row, "Percent change in average donations")
	for _, fy := range a.FYIndicators {
		s.WriteCellFloat(output, columns[fy], row, md.PercentChange(fy))
	}
}

// outputMajorList output the list of donors whose combined FY2023 and FY2024
// donations equal or exceed $2000.
func outputMajorList(donorList *dn.DonationList, output *s.SpreadsheetFile) {
	var row = 1
	//
	// Place Headings
	//
	s.WriteCell(output, "A", row, "Donor")
	for _, fy := range a.FYIndicators {
		s.WriteCell(output, columns[fy], row, "Donation "+fy.String())
	}

	s.WriteCell(output, "C", row, "Donation FY24")
	row++
	//
	// Output data
	//
	var names = donorList.DonorKeys()
	for _, name := range names {
		var donor = donorList.Get(name)
		if donor.IsMajorDonorOverall() {
			s.WriteCell(output, "A", row, name)
			for _, fy := range a.FYIndicators {
				s.WriteCellDecimal(output, columns[fy], row, donor.Donation(fy))
			}
			row++
		}
	}
}
