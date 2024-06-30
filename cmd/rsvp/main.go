// ----------------------------------------------------------------------------
//
// Mailing List for Celebration Invitees who have not responded
//
// Author: William Shaffer
// Version: 20-May-2024
//
// Copyright (c) 2024 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

package main

// This program produces a spreadsheet rsvp.xlsx with the names and mailing
// addresses of invitees to the July Celebration who have not responded
// to the email invite.

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	"fmt"
	"os"
	"strconv"

	"acorn_go/pkg/donors"
	g "acorn_go/pkg/guestlist"
	"acorn_go/pkg/spreadsheet"
)

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

const donorListFile = "/home/bozo/golang/acorn_go/data/donors.xlsx"
const donorListTab = "Sheet1"

const guestlistFile = "/home/bozo/golang/acorn_go/data/guestlist.xlsx"
const guestlistTab = "guestlist"

const outputFile = "/home/bozo/Downloads/rsvp.xlsx"
const outputTab = "Donors"

// ----------------------------------------------------------------------------
// Functions
// ----------------------------------------------------------------------------

// main supervises the processing of the invitees.
func main() {
	printHeader()
	//
	// Fetch donors
	//
	var donorList = generateDonorList()
	//
	// Fetch guest list
	//
	var guestlist = generateGuestlist(&donorList)
	//
	// Output name and address information.
	//
	outputAddresses(&guestlist)
}

// generateDonorList creates the collection of donors
func generateDonorList() donors.DonorList {
	var sprdsht spreadsheet.Spreadsheet
	var donorList donors.DonorList
	var err error
	//
	// Obtain spreadsheet data
	//
	sprdsht, err = spreadsheet.ProcessData(donorListFile, donorListTab)
	check(err, "Error generating address list: ")
	//
	// Generate donor list
	//
	donorList, err = donors.NewDonorList(&sprdsht)
	check(err, "Error generating donor list: ")
	return donorList
}

// generateGuestlist creates the list of guests
func generateGuestlist(donorList *donors.DonorList) g.Guestlist {
	var sprdsht spreadsheet.Spreadsheet
	var guestlist g.Guestlist
	var err error
	//
	// Obtain spreadsheet
	//
	sprdsht, err = spreadsheet.ProcessData(guestlistFile, guestlistTab)
	check(err, "Error processing guestlist spreadsheet: ")
	//
	// Generate guest list
	//
	guestlist, err = g.NewGuestlist(&sprdsht, donorList)
	check(err, "Error generating guestlist: ")
	return guestlist
}

// ----------------------------------------------------------------------------
// Output Functions
// ----------------------------------------------------------------------------

// outputAddresses outputs the addresses to a spreadsheet.
// This is a list of guest who have not responeded to the email invitation.
func outputAddresses(guestlist *g.Guestlist) {
	var err error
	var output spreadsheet.SpreadsheetFile
	//
	// selectGuest determines if the guest should be inserted into the spreadsheet
	//
	var selectGuest = func(guest g.Guest) bool {
		var result = guest.NoResponse()
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
	writeCell(&output, "A", row, "Name")
	writeCell(&output, "B", row, "Email")
	writeCell(&output, "C", row, "Street")
	writeCell(&output, "D", row, "City")
	writeCell(&output, "E", row, "State")
	writeCell(&output, "F", row, "Zip")
	row++
	//
	// Process donors
	//
	var keys = guestlist.Keys()
	for _, key := range keys {
		var donor = guestlist.Get(key)
		if guestlist.Contains(key) {
			if selectGuest(*donor) {
				var address = guestlist.Get(key)
				writeCell(&output, "A", row, address.Name())
				writeCell(&output, "B", row, address.Email())
				writeCell(&output, "C", row, address.Street())
				writeCell(&output, "D", row, address.City())
				writeCell(&output, "E", row, address.State())
				writeCell(&output, "F", row, address.Zip())
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
