// ----------------------------------------------------------------------------
//
// Email Validation Program
//
// Author: William Shaffer
// Version: 22-Jun-2024
//
// Copyright (c) 2024 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

package main

// This program validates emails from the donors.xlsx spreadsheet.

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	"fmt"
	"os"

	"acorn_go/pkg/donors"
	v "acorn_go/pkg/email"
	"acorn_go/pkg/spreadsheet"
)

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

const addressListFile = "/home/bozo/golang/acorn_go/data/donors.xlsx"
const addressListTab = "Sheet1"

// ----------------------------------------------------------------------------
// Functions
// ----------------------------------------------------------------------------

func main() {
	var addressList donors.DonorList

	printHeader()
	//
	// Fetch email addresses
	//
	addressList = generateAddresses()
	//
	// Validate each emil
	//
	var keys = addressList.Keys()
	for _, key := range keys {
		var address = addressList.Get(key)
		if address.HasEmail() {
			var email = address.Email()
			var result, err = v.Validate(email)
			check(err, "Error validating email belonging to: "+address.Name())
			if !result {
				fmt.Println("invalid email for " + key + ": " + email)
			}
		} else {
			fmt.Println("no email for: " + address.Name())
		}

	}
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
	addressList, err = donors.NewDonorAddressList(&sprdsht)
	check(err, "Error generating address list: ")
	return addressList
}

// ----------------------------------------------------------------------------
// Print Functions
// ----------------------------------------------------------------------------

// printHeader places the header information at the top of the page
func printHeader() {
	fmt.Println("-----------------------------------------------------------")
	fmt.Println("Acorn Scholarship Fund Email Validation")
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
