// ----------------------------------------------------------------------------
//
// # Annual letter program
//
// This program produces a spreadsheet with the names and addresses of
// donors who have donated in a specific calendar year.
//
// Author: William Shaffer
// Version: 30-Dec-2024
//
// Copyright (c) 2024 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

package main

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	a "acorn_go/pkg/accounting"
	dna "acorn_go/pkg/donations"
	dns "acorn_go/pkg/donors"
	s "acorn_go/pkg/spreadsheet"
	sp "acorn_go/pkg/support"
	"fmt"
	"os"
	"strconv"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

const donationsFile = "/home/bozo/golang/acorn_go/data/donations.xlsx"
const tabDonations = "Worksheet"
const donorFile = "/home/bozo/golang/acorn_go/data/donors.xlsx"
const tabDonors = "Sheet1"

const outputFile = "/home/bozo/Downloads/annualletter.xlsx"
const outputTab = "donors"
const outputTab2 = "Month Donors"

const reportYear = a.Y2025

// ----------------------------------------------------------------------------
// Functions
// ----------------------------------------------------------------------------

// main supervises the processing of the donation data.
func main() {
	var donationList dna.DonationList
	var donorList dns.DonorList
	var err error
	var output s.SpreadsheetFile
	//
	// Read input data
	//
	printHeader()
	donorList = generateAddressList()
	donationList = generateDonationList()
	//
	// Create output spreadsheet
	//
	output, err = s.New(outputFile, outputTab)
	sp.Check(err, "Error opening output file: ")
	//
	// Output annual donors to spreadsheet
	//
	outputAnnualDonors(&donorList, &donationList, output)
	//
	// Output prior month donors
	//
	output, err = output.AddSheet(outputTab2)
	sp.Check(err, "Error adding Month Donor sheet: ")
	outputMonthlysDonors(&donorList, &donationList, output)
	err = output.Save()
	sp.Check(err, "Error saving output file")
	err = output.Close()
	sp.Check(err, "Error closing output file")
	os.Exit(0)
}

// generateAddressList creates the donor list with addresses
func generateAddressList() dns.DonorList {
	var sprdsht s.Spreadsheet
	var donorList dns.DonorList
	var err error
	//
	// Obtain spreadsheet data
	//
	sprdsht, err = s.ProcessData(donorFile, tabDonors)
	sp.Check(err, "Error processing spreadsheet: ")
	//
	// Generate donor list
	//
	donorList, err = dns.NewDonorAddressList(&sprdsht)
	sp.Check(err, "Error generating donor list: ")
	return donorList
}

// generateDonationList creates the donation list
func generateDonationList() dna.DonationList {
	var sprdsht s.Spreadsheet
	var donationList dna.DonationList
	var err error
	//
	// Obtain spreadsheet data
	//
	sprdsht, err = s.ProcessData(donationsFile, tabDonations)
	sp.Check(err, "Error processing spreadsheet: ")
	//
	// Obtain donation list
	//
	donationList, err = dna.NewDonationList(&sprdsht)
	sp.Check(err, "Error generating donation list: ")
	return donationList
}

// ----------------------------------------------------------------------------
// Print Functions
// ----------------------------------------------------------------------------

// printHeader places the header information at the top of the page
func printHeader() {
	fmt.Println("-----------------------------------------------------------")
	fmt.Println("Acorn Scholarship Fund  Annual Letter List")
	fmt.Println("-----------------------------------------------------------")
}

// ----------------------------------------------------------------------------
// Output Functions
// ----------------------------------------------------------------------------

// outputDonors creates a spreadsheets with the donor names and addresses
// for donors that have donated in the specified calendar year.
// outputAddresses outputs the addresses to a spreadsheet
// This is the full list of donors
func outputAnnualDonors(
	donorList *dns.DonorList,
	donationList *dna.DonationList,
	output s.SpreadsheetFile) {
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
	row++
	//
	// Process donors
	//
	var personCount = 0
	var keys = donationList.DonorKeys()
	for _, key := range keys {
		var donor = donationList.Get(key)
		if selectDonor(donor) {
			donor = donorList.Get(key)
			s.WriteCell(&output, "A", row, donor.Name())
			s.WriteCell(&output, "B", row, donor.Email())
			s.WriteCell(&output, "C", row, donor.Street())
			s.WriteCell(&output, "D", row, donor.City())
			s.WriteCell(&output, "E", row, donor.State())
			s.WriteCell(&output, "F", row, donor.Zip())
			row++
			personCount++
		}
	}
	//
	// Output to person count
	//
	fmt.Println("Number of annual donors: " + strconv.Itoa(personCount))
}

// selectDonor returns true if the donor is to be output
func selectDonor(donor *dns.Donor) bool {
	var result = false
	if donor.Key() != "Tolman, Nadine" {
		result = donor.IsCalDonor(reportYear)
	}
	return result
}

// outputDonors creates a spreadsheets with the donor names and addresses
// for donors that have donated in the specified calendar year.
// outputAddresses outputs the addresses to a spreadsheet
// This is the full list of donors
func outputMonthlysDonors(
	donorList *dns.DonorList,
	donationList *dna.DonationList,
	output s.SpreadsheetFile) {
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
	s.WriteCell(&output, "G", row, "Key")
	row++
	//
	// Process donors
	//
	var personCount = 0
	var keys = donationList.DonorKeys()
	for _, key := range keys {
		var donor = donationList.Get(key)
		if selectDonor(donor) {
			if donor.IsCurrentDonor() {
				donor = donorList.Get(key)
				s.WriteCell(&output, "A", row, donor.Name())
				s.WriteCell(&output, "B", row, donor.Email())
				s.WriteCell(&output, "C", row, donor.Street())
				s.WriteCell(&output, "D", row, donor.City())
				s.WriteCell(&output, "E", row, donor.State())
				s.WriteCell(&output, "F", row, donor.Zip())
				s.WriteCell(&output, "G", row, key)
				row++
				personCount++
			}
		}
	}
	//
	// Output to person count
	//
	fmt.Println("Number of monthly donors: " + strconv.Itoa(personCount))
	fmt.Println("Current month is: " + a.CurrentMonth())
}
