// ----------------------------------------------------------------------------
//
// Donor series program
//
// Author: William Shaffer
// Version: 20-May-2024
//
// Copyright (c) 2024 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

package donors

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	"acorn_go/pkg/spreadsheet"
	"os"
	"strconv"
	"testing"
)

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

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

// Test_NewDonorList checks the reading of the donor spreadsheet and the
// creation of the donor list.
func Test_NewDonorList(t *testing.T) {
	var sprdsht spreadsheet.Spreadsheet
	var err error = nil
	var donorList DonorList
	//
	// Open the spreadsheet
	//
	sprdsht, err = spreadsheet.ProcessData(donorListFile, donorListTab)
	if err != nil {
		t.Error(err.Error())
	}
	//
	// Build the donor list
	//
	donorList, err = NewDonorAddressList(&sprdsht)
	if err != nil {
		t.Error(err.Error())
	}
	//
	// Test function
	//
	var testFunction = func(t *testing.T) {
		var apple = "Apple, Julietta"

		var count = (&donorList).Count()
		if count != 177 {
			t.Error("donor count is not 169: " + strconv.Itoa(count))
		}
		if !(&donorList).Contains(apple) {
			t.Error("donor list does not contain donor: " + apple)
		}
		var donor = donorList.Get(apple)
		if donor.Key() != apple {
			t.Error("get returned wrong donor: " + donor.Key())
		}
	}

	t.Run("Test_New", testFunction)
}
