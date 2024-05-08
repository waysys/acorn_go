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
	"strconv"

	"acorn_go/pkg/donor_info"
	"acorn_go/pkg/spreadsheet"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

const outputFile = "/home/bozo/Downloads/analysis.xlsx"

var donorTypes = []string{
	"FY2023 Only Donors",
	"FY2024 Only Donors",
	"FY2023 and FY2024 Donors",
}

// ----------------------------------------------------------------------------
// Functions
// ----------------------------------------------------------------------------

// main supervises the processing of the donation data.
func main() {
	var sprdsht spreadsheet.Spreadsheet
	var donorList donor_info.DonorList
	var err error
	var output spreadsheet.SpreadsheetFile

	printHeader()
	//
	// Obtain spreadsheet data
	//
	sprdsht, err = spreadsheet.ProcessData()
	check(err, "Error processing spreadsheet: ")
	//
	// Generate donor list
	//
	donorList, err = donor_info.NewDonorList(&sprdsht)
	check(err, "Error generating donor list: ")
	//
	// Calculate donor counts
	//
	var donorCount = donor_info.ComputeDonorCount(&donorList)
	//
	// Create output spreadsheet
	//
	output, err = spreadsheet.New(outputFile, "Donor Count")
	check(err, "Error opening output file: ")
	//
	// Output results from donor counts
	//
	printDonorCount(donorCount)
	outputDonorCount(donorCount, &output)
	//
	// Calculate donations
	//
	var donations = donor_info.ComputeDonations(&donorList)
	//
	// Output results from donations and repeat donors
	//
	printAmountAnalysis(donations)
	printRepeatAnalysis(donations)
	output, err = output.AddSheet("Donation Analysis")
	check(err, "Error adding donor analysis sheet: ")
	outputDonations(donations, &output)
	//
	// Calculate major donor statistics
	//
	var majorDonor = donor_info.ComputeMajorDonors(&donorList)
	printMajorDonorAnalysis(majorDonor)
	output, err = output.AddSheet("Major Donor")
	check(err, "Error adding major donor sheet")
	outputMajorDonor(majorDonor, &output)
	//
	// Finish up
	//
	err = output.Save()
	check(err, "Error saving output file: ")
	err = output.Close()
	check(err, "Error closing output file: ")
	os.Exit(0)
}

// ----------------------------------------------------------------------------
// Print Functions
// ----------------------------------------------------------------------------

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

// printAmountAnalysis prints the amounts of donations by fiscal year, by repeat donors
// versus non-repeat donors, and the comparisons of average donations between the
// current fiscal year and the prior fiscal year.
func printAmountAnalysis(donations donor_info.Donations) {
	fmt.Printf("\n\n\nDonations\n\n")
	fmt.Printf("FY2023 donations: $%v\n", donations.FYDonation(donor_info.FY2023))
	fmt.Printf("FY2024 donations: $%v\n", donations.FYDonation(donor_info.FY2024))
	fmt.Printf("Total donations:  $%v\n", donations.TotalDonation())
}

// printRepeatDonations prints the values comparing FY2023 and FY2024 donations
// from repeat donors.
func printRepeatAnalysis(donations donor_info.Donations) {
	fmt.Printf("\n\nRepeat Donation Analysis\n\n")
	fmt.Printf("Number of repeat donors: %d\n", donations.DonorCount(donor_info.DonorFY2023AndFY2024))
	fmt.Printf("Average FY2023 donations for repeat donors: $%5.0f\n",
		donations.AvgDonation(donor_info.DonorFY2023AndFY2024, donor_info.FY2023))
	fmt.Printf("Average FY2024 donations for repeat donors: $%5.0f\n",
		donations.AvgDonation(donor_info.DonorFY2023AndFY2024, donor_info.FY2024))
	fmt.Printf("Percent change in average donations: %3.0f percent\n",
		donations.DonationChange(donor_info.DonorFY2023AndFY2024))
}

// printMajorDonorAnalysis prints the values for major donors
func printMajorDonorAnalysis(majorDonor donor_info.MajorDonor) {
	fmt.Printf("\n\nMajor Donor Analysis\n\n")
	fmt.Printf("Number of major donors in FY2023: %d\n", majorDonor.MajorDonorCount(donor_info.FY2023))
	fmt.Printf("Number of major donors in FY2024: %d\n", majorDonor.MajorDonorCount(donor_info.FY2024))
	fmt.Printf("Major donations in FY2023: $%6.0f\n", majorDonor.DonationsMajor(donor_info.FY2023))
	fmt.Printf("Major donations in FY2024: $%6.0f\n", majorDonor.DonationsMajor(donor_info.FY2024))
	fmt.Printf("Average major donations in FY2023: $%5.0f\n", majorDonor.AvgDonation(donor_info.FY2023))
	fmt.Printf("Average major donations in FY2024: $%5.0f\n", majorDonor.AvgDonation(donor_info.FY2024))
	fmt.Printf("Percent change in average major donations: %3.0f percent\n", majorDonor.PercentChange())
	fmt.Printf("Percent of total donations by major donors FY2023: %3.0f percent\n",
		majorDonor.PercentDonation(donor_info.FY2023))
	fmt.Printf("Percent of total donations by major donors FY2024: %3.0f percent\n",
		majorDonor.PercentDonation(donor_info.FY2024))
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

// writeCellFloat outputs a float64 value to the specified cell
func writeCellFloat(
	outputPtr *spreadsheet.SpreadsheetFile,
	column string,
	row int,
	value float64) {

	var cell = cellName(column, row)
	var err = outputPtr.SetCellFloat(cell, value)
	check(err, "Error writing cell "+cell+": ")
}

// ----------------------------------------------------------------------------
// Spreadsheet Functions
// ----------------------------------------------------------------------------

// outputDonorCount inserts the donor count data into a sheet in the
// spreadsheet file
func outputDonorCount(
	donorCount donor_info.DonorCount,
	outputPtr *spreadsheet.SpreadsheetFile) {
	var row int = 1
	//
	// Place Title
	//
	writeCell(outputPtr, "A", row, "Donor Count Analysis")
	//
	// Place Headings
	//
	row += 2
	writeCell(outputPtr, "A", row, "Donor Type")
	writeCell(outputPtr, "B", row, "FY2023 Donors")
	writeCell(outputPtr, "C", row, "FY2024 Donors")
	//
	// Place data rows
	//
	row++
	writeCell(outputPtr, "A", row, donorTypes[0])
	writeCellInt(outputPtr, "B", row, donorCount.DonorsFY2023Only)
	writeCellInt(outputPtr, "C", row, 0)
	row++
	writeCell(outputPtr, "A", row, donorTypes[1])
	writeCellInt(outputPtr, "B", row, 0)
	writeCellInt(outputPtr, "C", row, donorCount.DonorsFY2024Only)
	row++
	writeCell(outputPtr, "A", row, donorTypes[2])
	writeCellInt(outputPtr, "B", row, donorCount.DonorsFY2023AndFY2024)
	writeCellInt(outputPtr, "C", row, donorCount.DonorsFY2023AndFY2024)
	row++
	writeCell(outputPtr, "A", row, "Total Donors")
	writeCellInt(outputPtr, "B", row, donorCount.TotalDonorsFY2023)
	writeCellInt(outputPtr, "C", row, donorCount.TotalDonorsFY2024)
	row += 2
	writeCell(outputPtr, "A", row, "Total donors")
	writeCellInt(outputPtr, "B", row, donorCount.TotalDonors)
	row++
	writeCell(outputPtr, "A", row, "Retention Rate")
	writeCellFloat(outputPtr, "B", row, donorCount.RetentionRate())
	row++
	writeCell(outputPtr, "A", row, "Acquisition Rate")
	writeCellFloat(outputPtr, "B", row, donorCount.AcquisitionRate())
}

// outputDonations inserts the donation data into the spreadsheet.
func outputDonations(
	donations donor_info.Donations,
	outputPtr *spreadsheet.SpreadsheetFile) {

	var row = 1
	//
	// Place Title
	//
	writeCell(outputPtr, "A", row, "Donation analysis")
	//
	// Place Headings
	//
	row += 2
	writeCell(outputPtr, "A", row, "Donor Type")
	writeCell(outputPtr, "B", row, "FY2023 Donations")
	writeCell(outputPtr, "C", row, "FY2024 Donations")
	writeCell(outputPtr, "E", row, "FY2023 Average Donations")
	writeCell(outputPtr, "F", row, "FY2023 Average Donations")
	//
	// Ouput Data
	//
	var donorType donor_info.DonorType
	for donorType = donor_info.DonorFY2023Only; donorType <= donor_info.DonorFY2023AndFY2024; donorType++ {
		row++
		writeCell(outputPtr, "A", row, donorTypes[donorType])
		writeCellFloat(outputPtr, "B", row, donations.Donation(donorType, donor_info.FY2023))
		writeCellFloat(outputPtr, "C", row, donations.Donation(donorType, donor_info.FY2024))
		writeCellFloat(outputPtr, "E", row, donations.AvgDonation(donorType, donor_info.FY2023))
		writeCellFloat(outputPtr, "F", row, donations.AvgDonation(donorType, donor_info.FY2024))
	}
	row++
	writeCell(outputPtr, "A", row, "Total Donations")
	writeCellFloat(outputPtr, "B", row, donations.FYDonation(donor_info.FY2023))
	writeCellFloat(outputPtr, "C", row, donations.FYDonation(donor_info.FY2024))
	writeCellFloat(outputPtr, "E", row, donations.FYAvgDonation(donor_info.FY2023))
	writeCellFloat(outputPtr, "F", row, donations.FYAvgDonation(donor_info.FY2024))
	row += 2
	writeCell(outputPtr, "A", row, "Total Donations for Both Years")
	writeCellFloat(outputPtr, "B", row, donations.TotalDonation())
}

// outputMajorDonor inserts the major donor data into the spreadsheet.
func outputMajorDonor(
	md donor_info.MajorDonor,
	outputPtr *spreadsheet.SpreadsheetFile) {
	var row = 1
	//
	// Place Title
	//
	writeCell(outputPtr, "A", row, "Major Donor Analysis")
	//
	// Place Headings
	//
	row += 2
	writeCell(outputPtr, "A", row, "")
	writeCell(outputPtr, "B", row, "FY2023")
	writeCell(outputPtr, "C", row, "FY2024")
	//
	// Ouput Data
	//
	row += 2
	writeCell(outputPtr, "A", row, "Major donors")
	writeCellInt(outputPtr, "B", row, md.MajorDonorCount(donor_info.FY2023))
	writeCellInt(outputPtr, "C", row, int(md.MajorDonorCount(donor_info.FY2024)))
	row++
	writeCell(outputPtr, "A", row, "Major donor donations")
	writeCellFloat(outputPtr, "B", row, md.DonationsMajor(donor_info.FY2023))
	writeCellFloat(outputPtr, "C", row, md.DonationsMajor(donor_info.FY2024))
	row++
	writeCell(outputPtr, "A", row, "Average major donor donation")
	writeCellFloat(outputPtr, "B", row, md.AvgDonation(donor_info.FY2023))
	writeCellFloat(outputPtr, "C", row, md.AvgDonation(donor_info.FY2024))
	row++
	writeCell(outputPtr, "A", row, "Percent of total donations by major donors")
	writeCellFloat(outputPtr, "B", row, md.PercentDonation(donor_info.FY2023))
	writeCellFloat(outputPtr, "C", row, md.PercentDonation(donor_info.FY2024))
	row += 2
	writeCell(outputPtr, "A", row, "Percent change in average donations")
	writeCellFloat(outputPtr, "B", row, md.PercentChange())
}
