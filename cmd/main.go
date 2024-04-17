// ----------------------------------------------------------------------------
//
// Donor analysis program
//
// Author: William Shaffer
// Version: 12-Apr-2024
//
// Copyright (c) William Shaffer
//
// ----------------------------------------------------------------------------

// This package processes donation information based on the payment transactons
// from a spreadsheet obtain from QuickBooks.  The package will then output
// the results in another spreadsheet.
package main

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	"fmt"
	"os"

	"acorn_go/pkg/donor_info"
	"acorn_go/pkg/spreadsheet"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// Functions
// ----------------------------------------------------------------------------

// main supervises the processing of the donation data.
func main() {
	var sprdsht spreadsheet.Spreadsheet
	var donorList donor_info.DonorList
	var err error

	fmt.Println("-----------------------------------------------------------")
	fmt.Println("Acorn Scholarship Fund Donor Analysis")
	fmt.Println("-----------------------------------------------------------")
	//
	// Obtain spreadsheet data
	//
	sprdsht, err = spreadsheet.ProcessData()
	if err != nil {
		fmt.Println("Error processing spreadsheet: " + err.Error())
		os.Exit(1)
	}
	//
	// Generate donor list
	//
	donorList, err = donor_info.NewDonorList(&sprdsht)
	if err != nil {
		fmt.Println("Error generating donor list: " + err.Error())
	}
	if donorList != nil {
		fmt.Println("the end")
	}
	os.Exit(0)
}
