// ----------------------------------------------------------------------------
//
// Mailing Address Program
//
// Author: William Shaffer
// Version: 20-May-2024
//
// Copyright (c) 2024 William Shaffer All Rights Reserved
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

	dn "acorn_go/pkg/donors"
	s "acorn_go/pkg/spreadsheet"
	sp "acorn_go/pkg/support"
)

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

const addressListFile = "/home/bozo/golang/acorn_go/data/donors.xlsx"
const addressListTab = "Sheet1"

const outputFile = "/home/bozo/Downloads/mailing_list.xlsx"
const outputTab = "Donors"

// ----------------------------------------------------------------------------
// Functions
// ----------------------------------------------------------------------------

// main supervises the processing of the donation data.
func main() {
	var donorList dn.DonorList

	printHeader()
	//
	// Fetch addresses
	//
	donorList = generateDonorList()
	//
	// Output data to spreadsheet
	//
	outputAddresses(&donorList)
}

// generateDonorList creates the donor list
func generateDonorList() dn.DonorList {
	var sprdsht s.Spreadsheet
	var donorList dn.DonorList
	var err error
	//
	// Obtain spreadsheet data
	//
	sprdsht, err = s.ProcessData(addressListFile, addressListTab)
	sp.Check(err, "Error processing spreadsheet: ")
	//
	// Generate donor list
	//
	donorList, err = dn.NewDonorAddressList(&sprdsht)
	sp.Check(err, "Error generating donor list: ")
	return donorList
}

// ----------------------------------------------------------------------------
// Output Functions
// ----------------------------------------------------------------------------

// outputAddresses outputs the addresses to a spreadsheet
// This is the full list of donors
func outputAddresses(donorList *dn.DonorList) {
	var err error
	var output s.SpreadsheetFile
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
	// Insert Heading
	//
	var row = 1
	s.WriteCell(&output, "A", row, "Donor Name")
	s.WriteCell(&output, "B", row, "Email")
	s.WriteCell(&output, "C", row, "Street")
	s.WriteCell(&output, "D", row, "City")
	s.WriteCell(&output, "E", row, "State")
	s.WriteCell(&output, "F", row, "Zip")
	s.WriteCell(&output, "G", row, "Count")
	row++
	//
	// Process donors
	//
	var personCount = 0
	var keys = donorList.Keys()
	for _, key := range keys {
		var donor = donorList.Get(key)
		s.WriteCell(&output, "A", row, donor.Name())
		s.WriteCell(&output, "B", row, donor.Email())
		s.WriteCell(&output, "C", row, donor.Street())
		s.WriteCell(&output, "D", row, donor.City())
		s.WriteCell(&output, "E", row, donor.State())
		s.WriteCell(&output, "F", row, donor.Zip())
		s.WriteCellInt(&output, "G", row, donor.NumberInHousehold())
		personCount += donor.NumberInHousehold()
		row++
	}
	//
	// Output to person count
	//
	fmt.Println("Number of people donating: " + strconv.Itoa(personCount))
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
