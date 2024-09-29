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

	dn "acorn_go/pkg/donors"
	g "acorn_go/pkg/guestlist"
	s "acorn_go/pkg/spreadsheet"
	sp "acorn_go/pkg/support"
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
	os.Exit(0)
}

// generateDonorList creates the collection of donors
func generateDonorList() dn.DonorList {
	var sprdsht s.Spreadsheet
	var donorList dn.DonorList
	var err error
	//
	// Obtain spreadsheet data
	//
	sprdsht, err = s.ProcessData(donorListFile, donorListTab)
	sp.Check(err, "Error generating address list: ")
	//
	// Generate donor list
	//
	donorList, err = dn.NewDonorAddressList(&sprdsht)
	sp.Check(err, "Error generating donor list: ")
	return donorList
}

// generateGuestlist creates the list of guests
func generateGuestlist(donorList *dn.DonorList) g.Guestlist {
	var sprdsht s.Spreadsheet
	var guestlist g.Guestlist
	var err error
	//
	// Obtain spreadsheet
	//
	sprdsht, err = s.ProcessData(guestlistFile, guestlistTab)
	sp.Check(err, "Error processing guestlist spreadsheet: ")
	//
	// Generate guest list
	//
	guestlist, err = g.NewGuestlist(&sprdsht, donorList)
	sp.Check(err, "Error generating guestlist: ")
	return guestlist
}

// ----------------------------------------------------------------------------
// Output Functions
// ----------------------------------------------------------------------------

// outputAddresses outputs the addresses to a spreadsheet.
// This is a list of guest who have not responeded to the email invitation.
func outputAddresses(guestlist *g.Guestlist) {
	var err error
	var output s.SpreadsheetFile
	//
	// selectGuest determines if the guest should be inserted into the spreadsheet
	//
	var selectGuest = func(guest *g.Guest) bool {
		var result = guest.NoResponse()
		return result
	}
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
	s.WriteCell(&output, "A", row, "Name")
	s.WriteCell(&output, "B", row, "Email")
	s.WriteCell(&output, "C", row, "Street")
	s.WriteCell(&output, "D", row, "City")
	s.WriteCell(&output, "E", row, "State")
	s.WriteCell(&output, "F", row, "Zip")
	row++
	//
	// Process donors
	//
	var keys = guestlist.Keys()
	for _, key := range keys {
		var guest = guestlist.Get(key)
		if selectGuest(guest) {
			var address = guestlist.Get(key)
			s.WriteCell(&output, "A", row, address.Name())
			s.WriteCell(&output, "B", row, address.Email())
			s.WriteCell(&output, "C", row, address.Street())
			s.WriteCell(&output, "D", row, address.City())
			s.WriteCell(&output, "E", row, address.State())
			s.WriteCell(&output, "F", row, address.Zip())
			row++
		}
	}
}

// ----------------------------------------------------------------------------
// Print Functions
// ----------------------------------------------------------------------------

// printHeader places the header information at the top of the page
func printHeader() {
	fmt.Println("-----------------------------------------------------------")
	fmt.Println("Acorn Scholarship Fund Guests Who Have Not Responded")
	fmt.Println("-----------------------------------------------------------")
}
