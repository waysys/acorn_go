// ----------------------------------------------------------------------------
//
// Mailing Address Program
//
// Author: William Shaffer
// Version: 20-May-2024
//
// Copyright (c) 2024, 2025 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

package main

// This program produces a spreadsheet of donor names and addresses.

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	"fmt"
	"strconv"

	a "acorn_go/pkg/accounting"
	dn "acorn_go/pkg/donors"
	s "acorn_go/pkg/spreadsheet"
	sp "acorn_go/pkg/support"
)

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

const addressListFile = "/home/bozo/golang/acorn_go/data/donors.xlsx"
const addressListTab = "Sheet1"

const donationFile = "/home/bozo/golang/acorn_go/data/donations.xlsx"
const donationTab = "Worksheet"

const outputFile = "/home/bozo/Downloads/mailing_list.xlsx"
const outputTab = "Donors"
const inviteTab = "Invitees"

// ----------------------------------------------------------------------------
// Functions
// ----------------------------------------------------------------------------

// main supervises the processing of the donation data.
func main() {
	var donorList dn.DonorList
	var output s.SpreadsheetFile
	var err error

	printHeader()
	//
	// Create output spreadsheet
	//
	output, err = s.New(outputFile, outputTab)
	sp.Check(err, "Error opening output file: ")
	defer func() {
		err = output.Save()
		sp.Check(err, "Error saving output file")
		err = output.Close()
		sp.Check(err, "Error closing output file")
	}()
	//
	// Fetch addresses
	//
	donorList = generateDonorList()
	//
	// Output mailing list
	//
	outputAddresses(&donorList, &output)
	//
	// Output the invite list
	//
	output, err = output.AddSheet(inviteTab)
	sp.Check(err, "Error adding sheet: ")
	outputInviteList(&donorList, &output)

	printFooter()
}

// generateDonorList creates the donor list
func generateDonorList() dn.DonorList {
	var sprdsht s.Spreadsheet
	var donorList dn.DonorList
	var err error
	//
	// Obtain donor spreadsheet data
	//
	sprdsht, err = s.ProcessData(addressListFile, addressListTab)
	sp.Check(err, "Error processing spreadsheet: ")
	//
	// Generate donor list
	//
	donorList, err = dn.NewDonorAddressList(&sprdsht)
	sp.Check(err, "Error generating donor list: ")
	//
	// Obtain donation spreadsheet data
	//
	sprdsht, err = s.ProcessData(donationFile, donationTab)
	sp.Check(err, "Error processing spreadsheet: ")
	//
	// Add donation data to donor list
	//
	dn.AddDonations(&sprdsht, &donorList)
	return donorList
}

// ----------------------------------------------------------------------------
// Output Functions
// ----------------------------------------------------------------------------

// outputAddresses outputs the addresses to a spreadsheet
// This is the full list of donors
func outputAddresses(donorList *dn.DonorList, output *s.SpreadsheetFile) {
	//
	// Insert Heading
	//
	var row = 1
	s.WriteCell(output, "A", row, "Donor Name")
	s.WriteCell(output, "B", row, "Email")
	s.WriteCell(output, "C", row, "Street")
	s.WriteCell(output, "D", row, "City")
	s.WriteCell(output, "E", row, "State")
	s.WriteCell(output, "F", row, "Zip")
	s.WriteCell(output, "G", row, "Count")
	row++
	//
	// Process donors
	//
	var personCount = 0
	var keys = donorList.Keys()
	for _, key := range keys {
		var donor = donorList.Get(key)
		s.WriteCell(output, "A", row, donor.Name())
		s.WriteCell(output, "B", row, donor.Email())
		s.WriteCell(output, "C", row, donor.Street())
		s.WriteCell(output, "D", row, donor.City())
		s.WriteCell(output, "E", row, donor.State())
		s.WriteCell(output, "F", row, donor.Zip())
		s.WriteCellInt(output, "G", row, donor.NumberInHousehold())
		personCount += donor.NumberInHousehold()
		row++
	}
	//
	// Output person count
	//
	fmt.Println("Number of people donating: " + strconv.Itoa(personCount))
}

// outputAddresses outputs the addresses to a spreadsheet
// This is the full list of donors
func outputInviteList(donorList *dn.DonorList, output *s.SpreadsheetFile) {
	//
	// Insert Heading
	//
	var row = 1
	s.WriteCell(output, "A", row, "Donor Name")
	s.WriteCell(output, "B", row, "Email")
	s.WriteCell(output, "C", row, "Street")
	s.WriteCell(output, "D", row, "City")
	s.WriteCell(output, "E", row, "State")
	s.WriteCell(output, "F", row, "Zip")
	s.WriteCell(output, "G", row, "Count")
	row++
	//
	// Process donors
	//
	var personCount = 0
	var keys = donorList.Keys()
	for _, key := range keys {
		var donor = donorList.Get(key)
		if donor.IsCalDonor(a.Y2024) || donor.IsCalDonor(a.Y2025) {
			s.WriteCell(output, "A", row, donor.Name())
			s.WriteCell(output, "B", row, donor.Email())
			s.WriteCell(output, "C", row, donor.Street())
			s.WriteCell(output, "D", row, donor.City())
			s.WriteCell(output, "E", row, donor.State())
			s.WriteCell(output, "F", row, donor.Zip())
			s.WriteCellInt(output, "G", row, donor.NumberInHousehold())
			personCount += donor.NumberInHousehold()
			row++
		}
	}
	//
	// Output person count
	//
	fmt.Println("Number of people invited: " + strconv.Itoa(personCount))
}

// ----------------------------------------------------------------------------
// Print Functions
// ----------------------------------------------------------------------------

// printHeader places the header information at the top of the page
func printHeader() {
	fmt.Println("-----------------------------------------------------------")
	fmt.Println("Acorn Scholarship Fund Mailing List")
	fmt.Println("-----------------------------------------------------------")
}

// printFooter prints the notice that the program is finished.
func printFooter() {
	fmt.Println("-----------------------------------------------------------")
	fmt.Println("Program is finished")
	fmt.Println("-----------------------------------------------------------")
}
