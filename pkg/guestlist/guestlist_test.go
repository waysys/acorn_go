// ----------------------------------------------------------------------------
//
// Guestlist
//
// Author: William Shaffer
// Version: 29-Jun-2024
//
// Copyright (c) 2024 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

package guestlist

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	"acorn_go/pkg/donors"
	"acorn_go/pkg/spreadsheet"
	"os"
	"strconv"
	"testing"
)

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

const guestlistFile = "/home/bozo/golang/acorn_go/data/guestlist.xlsx"
const guestlistTab = "guestlist"

const donorListFile = "/home/bozo/golang/acorn_go/data/donors.xlsx"
const donorListTab = "Sheet1"

// ----------------------------------------------------------------------------
// Test Main
// ----------------------------------------------------------------------------

func TestMain(m *testing.M) {
	exitVal := m.Run()
	os.Exit(exitVal)
}

// ----------------------------------------------------------------------------
// Test functions
// ----------------------------------------------------------------------------

// Test_NewGuestlist checks the reading of the guestlist.xlsm spreadsheet.
func Test_NewGuestlist(t *testing.T) {
	var sprdsht spreadsheet.Spreadsheet
	var donorList donors.DonorList
	var err error = nil
	var guestlist Guestlist

	//
	// Open the spreadsheet
	//
	sprdsht, err = spreadsheet.ProcessData(guestlistFile, guestlistTab)
	if err != nil {
		t.Error(err.Error())
	}
	//
	// Fetch donor list
	//
	donorList, err = generateDonorList()
	if err != nil {
		t.Error(err.Error())
	}
	//
	// Build the donor list
	//
	guestlist, err = NewGuestlist(&sprdsht, &donorList)
	if err != nil {
		t.Error(err.Error())
	}
	//
	// Test function
	//
	var testFunction = func(t *testing.T) {
		var email = "juliettaapple0@gmail.com"
		var name = "Ms. Julie Apple"

		var count = (&guestlist).Count()
		if count < 30 {
			t.Error("guest count less than 30: " + strconv.Itoa(count))
		}
		if !(&guestlist).Contains(email) {
			t.Error("guestlist does not contain guest with email: " + email)
		}
		var guest = (&guestlist).Get(email)
		if guest.Name() != name {
			t.Error("get returned wrong donor: " + guest.Name())
		}
	}

	t.Run("Test_New", testFunction)
}

// ----------------------------------------------------------------------------
// Support functions
// ----------------------------------------------------------------------------

// generateDonorList creates the collection of donors
func generateDonorList() (donors.DonorList, error) {
	var sprdsht spreadsheet.Spreadsheet
	var donorList donors.DonorList
	var err error
	//
	// Obtain spreadsheet data
	//
	sprdsht, err = spreadsheet.ProcessData(donorListFile, donorListTab)
	//
	// Generate donor list
	//
	if err == nil {
		donorList, err = donors.NewDonorAddressList(&sprdsht)
	}

	return donorList, err
}
