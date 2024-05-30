// ----------------------------------------------------------------------------
//
// Scholarship Analysis Program
//
// Author: William Shaffer
// Version: 10-May-2024
//
// Copyright (c) 2024 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

// This program produces a spreadsheet with the anaysis of scholarship payments.

package main

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	g "acorn_go/pkg/grants"
	q "acorn_go/pkg/quickbooks"
	"acorn_go/pkg/spreadsheet"
	"fmt"
	"os"
	"strconv"
)

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

const outputFile = "/home/bozo/Downloads/scholarship.xlsx"

// ----------------------------------------------------------------------------
// Main Function
// ----------------------------------------------------------------------------

// main supervises the processing of scholarship data.
func main() {
	var err error = nil
	var apTranlist q.TransList
	var billList q.BillList
	var grantList g.GrantList
	var output spreadsheet.SpreadsheetFile

	printHeader()
	//
	// Read the accounts payable transaction
	//
	apTranlist, err = q.ReadAPTransactions()
	//
	// Read the bills
	//
	if err == nil {
		billList, err = q.ReadBills(&apTranlist)
	}
	//
	// Generate grant transactions
	//
	if err == nil {
		grantList, err = g.AssembleGrantList(&billList, &apTranlist)
	}
	//
	// Output results
	//
	if err == nil {
		output, err = spreadsheet.New(outputFile, "Summary")
	}
	if err == nil {
		outputGrantSummary(&output, &grantList)
	}

}

// ----------------------------------------------------------------------------
// Spreadsheet Functions
// ----------------------------------------------------------------------------

// outputGrantSummary populates the Summary tab with the amounts by transaction
// type
func outputGrantSummary(output *spreadsheet.SpreadsheetFile, grantList *g.GrantList) {

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
