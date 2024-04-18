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

	printHeader()
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
	//
	// Calculate donor counts
	//
	var donorCount = donor_info.ComputeDonorCount(&donorList)
	//
	// Output results
	//
	printDonorCount(donorCount)
	printNonRepeatDonors(&donorList)
	os.Exit(0)
}

// printHeader places the header information at the top of the page
func printHeader() {
	fmt.Println("-----------------------------------------------------------")
	fmt.Println("Acorn Scholarship Fund Donor Analysis")
	fmt.Println("-----------------------------------------------------------")
}

// printDonorCount will print out the conor account informaiton
func printDonorCount(donorCount donor_info.DonorCount) {
	fmt.Printf("FY2023 donor count:     %d\n", donorCount.TotalDonorsFY2023)
	fmt.Printf("FY2024 donor count:     %d\n", donorCount.TotalDonorsFY2024)
	fmt.Printf("Total number of donors: %d\n", donorCount.TotalDonors)
	fmt.Println("")
	fmt.Printf("Donors donating in FY2023 only: %d\n", donorCount.DonorsFY2023Only)
	fmt.Printf("Donors donating in FY2024 only: %d\n", donorCount.DonorsFY2024Only)
	fmt.Printf("Donors donating in both years:  %d\n", donorCount.DonorsFY2023AndFY2024)
}

// printNonRepeatDonors prints a list of donors who donated in FY2023 but not FY2024.
func printNonRepeatDonors(donorListPtr *donor_info.DonorList) {
	var names = donor_info.NonRepeatDonors(donorListPtr)
	var count = 0
	fmt.Printf("\nDonors who donated in FY2023 but not FY2024\n\n")
	for _, name := range names {
		count++
		fmt.Println(name)
	}
	fmt.Printf("\nNumber of non-repeat donors: %d", count)
}
