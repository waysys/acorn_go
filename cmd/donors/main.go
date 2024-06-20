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
// It is currently programmed to create the list for donors in FY2023 and FY2024.

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	"fmt"
	"os"
	"strconv"

	"acorn_go/pkg/donorinfo"
	"acorn_go/pkg/donors"
	"acorn_go/pkg/spreadsheet"
)

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

const donorListFile = "/home/bozo/golang/acorn_go/data/donations.xlsx"
const donorListTab = "Worksheet"

const addressListFile = "/home/bozo/golang/acorn_go/data/donors.xlsx"
const addressListTab = "Sheet1"

const outputFile = "/home/bozo/Downloads/mailing_list.xlsx"
const outputTab = "Donors"

const paperInviteFile = "/home/bozo/Downloads/paper_invite.xlsx"

// ----------------------------------------------------------------------------
// Functions
// ----------------------------------------------------------------------------

// main supervises the processing of the donation data.
func main() {
	var donorList donorinfo.DonorList
	var addressList donors.DonorList

	printHeader()
	//
	// Fetch donor list
	//
	donorList = generateDonorList()
	//
	// Fetch addresses
	//
	addressList = generateAddresses()
	//
	// Output data to spreadsheet
	//
	outputAddresses(&donorList, &addressList)
	outputPaperInviteList(&donorList, &addressList)
}

// generateDonorList creates the donor list
func generateDonorList() donorinfo.DonorList {
	var sprdsht spreadsheet.Spreadsheet
	var donorList donorinfo.DonorList
	var err error
	//
	// Obtain spreadsheet data
	//
	sprdsht, err = spreadsheet.ProcessData(donorListFile, donorListTab)
	check(err, "Error processing spreadsheet: ")
	//
	// Generate donor list
	//
	donorList, err = donorinfo.NewDonorList(&sprdsht)
	check(err, "Error generating donor list: ")
	return donorList
}

// generateAddresses creates the collection of addresses
func generateAddresses() donors.DonorList {
	var sprdsht spreadsheet.Spreadsheet
	var addressList donors.DonorList
	var err error
	//
	// Obtain spreadsheet data
	//
	sprdsht, err = spreadsheet.ProcessData(addressListFile, addressListTab)
	check(err, "Error generating address list: ")
	//
	// Generate address list
	//
	addressList, err = donors.NewDonorList(&sprdsht)
	check(err, "Error generating address list: ")
	return addressList
}

// ----------------------------------------------------------------------------
// Output Functions
// ----------------------------------------------------------------------------

// outputAddresses outputs the addresses to a spreadsheet
// This is the full list of donors
func outputAddresses(donorList *donorinfo.DonorList, addressList *donors.DonorList) {
	var err error
	var output spreadsheet.SpreadsheetFile
	//
	// selectDonor determines if the donor should be inserted into the mailing list
	//
	var selectDonor = func(donor donorinfo.Donor) bool {
		var result = donor.IsFY23AndFY24Donor() || donor.IsFY24DonorOnly() || donor.IsFY23DonorOnly()
		return result
	}
	//
	// Create output spreadsheet
	//
	output, err = spreadsheet.New(outputFile, outputTab)
	check(err, "Error opening output file: ")
	//
	// Insert Heading
	//
	var row = 1
	writeCell(&output, "A", row, "Donor Name")
	writeCell(&output, "B", row, "Email")
	writeCell(&output, "C", row, "Street")
	writeCell(&output, "D", row, "City")
	writeCell(&output, "E", row, "State")
	writeCell(&output, "F", row, "Zip")
	writeCell(&output, "G", row, "Count")
	row++
	//
	// Process donors
	//
	var personCount = 0
	var keys = donorList.DonorKeys()
	for _, key := range keys {
		var donor = donorList.Get(key)
		if addressList.Contains(key) {
			if selectDonor(*donor) {
				var address = addressList.Get(key)
				writeCell(&output, "A", row, address.Name())
				writeCell(&output, "B", row, address.Email())
				writeCell(&output, "C", row, address.Street())
				writeCell(&output, "D", row, address.City())
				writeCell(&output, "E", row, address.State())
				writeCell(&output, "F", row, address.Zip())
				writeCellInt(&output, "G", row, address.NumberInHousehold())
				personCount += address.NumberInHousehold()
				row++
			}
		}
	}
	//
	// Output to person count
	//
	fmt.Println("Number of people donating: " + strconv.Itoa(personCount))
	//
	// Finish
	//
	err = output.Save()
	check(err, "Error saving output file")
	err = output.Close()
	check(err, "Error closing output file")
}

// outputPaperInviteList outputs the names and address of donors without email addresses.
// This list is for FY2023 and FY2024 donors who are not deceased.
func outputPaperInviteList(donorList *donorinfo.DonorList, addressList *donors.DonorList) {
	var err error
	var output spreadsheet.SpreadsheetFile
	//
	// selectDonor determines if the donor should be inserted into the list
	//
	var selectDonor = func(donor donorinfo.Donor, address donors.Donor) bool {
		var result = donor.IsFY23AndFY24Donor() || donor.IsFY24DonorOnly() || donor.IsFY23DonorOnly()
		result = result && !address.HasEmail()
		return result
	}
	//
	// Create output spreadsheet
	//
	output, err = spreadsheet.New(paperInviteFile, outputTab)
	check(err, "Error opening output file: ")
	//
	// Insert Heading
	//
	var row = 1
	writeCell(&output, "A", row, "Donor Name")
	writeCell(&output, "B", row, "")
	writeCell(&output, "C", row, "Street")
	writeCell(&output, "D", row, "City")
	writeCell(&output, "E", row, "State")
	writeCell(&output, "F", row, "Zip")
	writeCell(&output, "G", row, "Count")
	row++
	//
	// Process donors
	//
	var keys = donorList.DonorKeys()
	for _, key := range keys {
		var donor = donorList.Get(key)
		if addressList.Contains(key) {
			var address = addressList.Get(key)
			if selectDonor(*donor, *address) {
				writeCell(&output, "A", row, address.Name())
				writeCell(&output, "C", row, address.Street())
				writeCell(&output, "D", row, address.City())
				writeCell(&output, "E", row, address.State())
				writeCell(&output, "F", row, address.Zip())
				writeCellInt(&output, "G", row, address.NumberInHousehold())
				row++
			}
		}
	}
	//
	// Finish
	//
	err = output.Save()
	check(err, "Error saving output file")
	err = output.Close()
	check(err, "Error closing output file")
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

// ----------------------------------------------------------------------------
// Support Functions
// ----------------------------------------------------------------------------

// check tests an error to see if it is null.  If not, it prints an
// error message and exits the program.
func check(err error, message string) {
	if err != nil {
		fmt.Println(message + err.Error())
		os.Exit(1)
	}
}

// cellName generates a string representing a cell in the spreadsheet.
func cellName(column string, row int) string {
	var cellName = column + strconv.Itoa(row)
	return cellName
}

// writeCell outputs a string value to the specified cell
func writeCell(
	outputPtr *spreadsheet.SpreadsheetFile,
	column string,
	row int,
	value string) {

	var cell = cellName(column, row)
	var err = outputPtr.SetCell(cell, value)
	check(err, "Error writing cell "+cell+": ")
}

// writeCellInt outputs an integer value to the specified cell
func writeCellInt(
	outputPtr *spreadsheet.SpreadsheetFile,
	column string,
	row int,
	value int) {

	var cell = cellName(column, row)
	var err = outputPtr.SetCellInt(cell, value)
	check(err, "Error writing cell "+cell+": ")
}
