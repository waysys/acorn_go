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
	a "acorn_go/pkg/accounting"
	g "acorn_go/pkg/grants"
	q "acorn_go/pkg/quickbooks"
	"acorn_go/pkg/spreadsheet"
	"fmt"
	"os"
	"strconv"

	dec "github.com/shopspring/decimal"
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
	check(err, "Error: ")

	defer output.Close()
	//
	// Produce summary tab
	//
	outputGrantSummary(&output, &grantList)
	//
	// Produce recipient tab
	//
	// output, err = output.AddSheet("Recipients")
	// outputRecipientList(&output, &grantList)
	//
	// Save results
	//
	output.Save()
	printFooter()
}

// ----------------------------------------------------------------------------
// Spreadsheet Functions
// ----------------------------------------------------------------------------

// outputGrantSummary populates the Summary tab with the amounts by transaction
// type
func outputGrantSummary(output *spreadsheet.SpreadsheetFile, grantList *g.GrantList) {
	//
	// Insert title
	//
	var row = 1
	writeCell(output, "A", row, "Summary of Scholarship Payments")
	row += 2
	//
	// Insert headings
	//
	writeCell(output, "A", row, "Fiscal Year")
	writeCell(output, "B", row, "Total Grants")
	writeCell(output, "C", row, "Total Payments")
	writeCell(output, "D", row, "Total Net Write-Offs")
	row++
	//
	// Amounts
	//
	writeCell(output, "A", row, "FY2023")
	var grantTotal2023 = grantList.TotalTransAmount(a.FY2023, g.Grant)
	writeCellDecimal(output, "B", row, grantTotal2023)
	var paymentTotal2023 = grantList.TotalTransAmount(a.FY2023, g.GrantPayment)
	writeCellDecimal(output, "C", row, paymentTotal2023)
	var writeOffTotal2023 = grantList.TotalNetWriteOff(a.FY2023)
	writeCellDecimal(output, "D", row, writeOffTotal2023)
	row++
	//
	writeCell(output, "A", row, "FY2024")
	var grantTotal2024 = grantList.TotalTransAmount(a.FY2024, g.Grant)
	writeCellDecimal(output, "B", row, grantTotal2024)
	var paymentTotal2024 = grantList.TotalTransAmount(a.FY2024, g.GrantPayment)
	writeCellDecimal(output, "C", row, paymentTotal2024)
	var writeOffTotal2024 = grantList.TotalNetWriteOff(a.FY2024)
	writeCellDecimal(output, "D", row, writeOffTotal2024)
	row += 2
	//
	writeCell(output, "A", row, "Total All Years")
	var grantTotal = grantTotal2023.Add(grantTotal2024)
	writeCellDecimal(output, "B", row, grantTotal)
	var paymentTotal = paymentTotal2023.Add(paymentTotal2024)
	writeCellDecimal(output, "C", row, paymentTotal)
	var writeOffTotal = writeOffTotal2023.Add(writeOffTotal2024)
	writeCellDecimal(output, "D", row, writeOffTotal)
	row += 2
	//
	var remainingGrants = grantTotal.Sub(paymentTotal).Sub(writeOffTotal)
	writeCell(output, "A", row, "Remaining Grants")
	writeCellDecimal(output, "B", row, remainingGrants)
}

// outputRecipientList produces a list of recipients organized by fiscal year and
// award group
func outputRecipientList(output *spreadsheet.SpreadsheetFile, grantList *g.GrantList) {

}

// ----------------------------------------------------------------------------
// Print Functions
// ----------------------------------------------------------------------------

// printHeader places the header information at the top of the page
func printHeader() {
	fmt.Println("-----------------------------------------------------------")
	fmt.Println("Acorn Scholarship Fund Scholarship Analysis")
	fmt.Println("-----------------------------------------------------------")
}

// printFooter prints a message indicating completion of the process
func printFooter() {
	fmt.Println("Program is finished.")
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

// writeCellFloat outputs a float64 value to the specified cell
func writeCellDecimal(
	outputPtr *spreadsheet.SpreadsheetFile,
	column string,
	row int,
	value dec.Decimal) {

	var cell = cellName(column, row)
	var amount = value.String()
	var err = outputPtr.SetCell(cell, amount)
	check(err, "Error writing cell "+cell+": ")
}
