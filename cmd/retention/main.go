// ----------------------------------------------------------------------------
//
// Non-repeat donor analysis
//
// Author: William Shaffer
// Version: 7-May-2024
//
// Copyright (c) 2024 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

// This program generates a list of FY2023 donors who did not donate in
// FY2024 with the dates and amounts they donated.

package main

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	d "acorn_go/pkg/date"
	"acorn_go/pkg/donorinfo"
	"acorn_go/pkg/spreadsheet"
	"fmt"
	"os"
	"strconv"

	dec "github.com/shopspring/decimal"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

const (
	inputFile      = "/home/bozo/golang/acorn_go/data/donations.xlsx"
	tab            = "Worksheet"
	outputFileName = "/home/bozo/Downloads/nonrepeat.xlsx"
	sheetName      = "Non-Repeat Donors"
)

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

// writeCellDate outputs a date value to the specified cell
func writeCellDate(
	outputPtr *spreadsheet.SpreadsheetFile,
	column string,
	row int,
	date d.Date) {

	var cell = cellName(column, row)
	var err = outputPtr.SetCell(cell, date.String())
	check(err, "Error writing cell "+cell+": ")
}

// writeCellDecimal outputs a float64 value to the specified cell
func writeCellDecimal(
	outputPtr *spreadsheet.SpreadsheetFile,
	column string,
	row int,
	num dec.Decimal) {

	var cell = cellName(column, row)
	var value = num.String()
	var err = outputPtr.SetCell(cell, value)
	check(err, "Error writing cell "+cell+": ")
}

// ----------------------------------------------------------------------------
// Functions
// ----------------------------------------------------------------------------

// main supervises the execution of this program.  It produces a spreadsheet
// with the list of non-repeat donors
func main() {
	var sprdsht spreadsheet.Spreadsheet
	var donorList donorinfo.DonorList
	var nonrepeats []donorinfo.NonRepeat
	var err error

	printHeader()
	//
	// Obtain spreadsheet data
	//
	sprdsht, err = spreadsheet.ProcessData(inputFile, tab)
	check(err, "Error processing spreadsheet")
	//
	// Generate donor list
	//
	donorList, err = donorinfo.NewDonorList(&sprdsht)
	check(err, "Error generating donor list: ")
	//
	// Generate donor list
	//
	nonrepeats = donorinfo.ComputeNonRepeatDonors(&donorList, &sprdsht)
	check(err, "Error generating donor list")
	//
	// Output donation series to spreadsheet
	//
	err = outputRetention(nonrepeats)
	check(err, "Error writing output")
	os.Exit(0)
}

// printHeader places the header information at the top of the page
func printHeader() {
	fmt.Println("-----------------------------------------------------------")
	fmt.Println("Acorn Scholarship Fund Retention Analysis")
	fmt.Println("-----------------------------------------------------------")
}

// outputRetention produces a spreadsheet with the non-repeat donors names
// and donation dates and amounts listed.
func outputRetention(nonrepeats []donorinfo.NonRepeat) error {
	var err error
	var output spreadsheet.SpreadsheetFile
	//
	// Create output spreadsheet
	//
	output, err = spreadsheet.New(outputFileName, sheetName)
	check(err, "Error opening output file: ")
	//
	// Insert Heading
	//
	var row = 1
	writeCell(&output, "A", row, "List of non-repeat donors")
	row += 2
	writeCell(&output, "A", row, "Donor Name")
	writeCell(&output, "B", row, "Date of Donation")
	writeCell(&output, "C", row, "Amount of Donation")
	row++
	//
	// Insert donor information
	//
	for i := 0; i < len(nonrepeats); i++ {
		var name = nonrepeats[i].NameDonor()
		var date = nonrepeats[i].DateDonation()
		var amount = nonrepeats[i].AmountDonation()
		writeCell(&output, "A", row, name)
		writeCellDate(&output, "B", row, date)
		writeCellDecimal(&output, "C", row, amount)
		row++
	}
	//
	// Finish
	//
	err = output.Save()
	check(err, "Error saving output file")
	err = output.Close()
	check(err, "Error closing output file")
	return err
}
