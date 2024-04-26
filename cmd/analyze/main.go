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
		os.Exit(1)
	}
	//
	// Calculate donor counts
	//
	var donorCount = donor_info.ComputeDonorCount(&donorList)
	//
	// Output results from donor counts
	//
	printDonorCount(donorCount)
	printNonRepeatDonors(&donorList)
	//
	// Calculate donations
	//
	var donations = donor_info.ComputeDonations(&donorList)
	//
	// Output results from donations and repeat donors
	//
	printAmountAnalysis(donations)
	printRepeatAnalysis(donations)
	//
	// Calculate major donor statistics
	//
	var majorDonor = donor_info.ComputeMajorDonors(&donorList)
	printMajorDonorAnalysis(majorDonor)

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

// printAmountAnalysis prints the amounts of donations by fiscal year, by repeat donors
// versus non-repeat donors, and the comparisons of average donations between the
// current fiscal year and the prior fiscal year.
func printAmountAnalysis(donations donor_info.Donations) {
	fmt.Printf("\n\n\nDonations\n\n")
	fmt.Printf("FY2023 donations: $%v\n", donations.DonationsFY2023)
	fmt.Printf("FY2024 donations: $%v\n", donations.DonationsFY2024)
	fmt.Printf("Total donations:  $%v\n", donations.DonationsTotal)
}

// printRepeatDonations prints the values comparing FY2023 and FY2024 donations
// from repeat donors.
func printRepeatAnalysis(donations donor_info.Donations) {
	fmt.Printf("\n\nRepeat Donation Analysis\n\n")
	fmt.Printf("Number of repeat donors: %d\n", donations.CountRepeatDonors)
	fmt.Printf("Average FY2023 donations for repeat donors: $%5.0f\n", donations.AvgDonationFY2023)
	fmt.Printf("Average FY2024 donations for repeat donors: $%5.0f\n", donations.AvgDonationFY2024)
	fmt.Printf("Percent change in average donations: %3.0f percent\n", donations.DonationChange)
}

// printMajorDonorAnalysis prints the values for major donors
func printMajorDonorAnalysis(majorDonor donor_info.MajorDonor) {
	fmt.Printf("\n\nMajor Donor Analysis\n\n")
	fmt.Printf("Number of major donors in FY2023: %d\n", majorDonor.MajorDonorsFY2023)
	fmt.Printf("Number of major donors in FY2024: %d\n", majorDonor.MajorDonorsFY2024)
	fmt.Printf("Major donations in FY2023: $%v\n", majorDonor.DonationsMajorFY2023.IntPart())
	fmt.Printf("Major donations in FY2024: $%v\n", majorDonor.DonationsMajorFY2024.IntPart())
	fmt.Printf("Average major donations in FY2023: $%5.0f\n", majorDonor.AvgMajorDonationFY2023)
	fmt.Printf("Average major donations in FY2024: $%5.0f\n", majorDonor.AvgMajorDonationFY2024)
	fmt.Printf("Percent change in average major donations: %3.0f percent\n", majorDonor.DonationChange)
	fmt.Printf("Percent of total donations by major donors FY2023: %3.0f peercent\n",
		majorDonor.PercentTotalDonationsFY2023)
	fmt.Printf("Percent of total donations by major donors FY2024: %3.0f percent\n",
		majorDonor.PercentTotalDonationsFY2024)
}
