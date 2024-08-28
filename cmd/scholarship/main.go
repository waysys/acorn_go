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

	"github.com/waysys/assert/assert"
	d "github.com/waysys/waydate/pkg/date"

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
	output, err = output.AddSheet("Recipients")
	check(err, "Error: ")
	outputRecipientList(&output, &grantList)
	//
	// Produce recipient summary tab
	//
	output, err = output.AddSheet("RecipSum")
	check(err, "Error: ")
	outputRecipientSummary(&output, &grantList)
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
	var row = 1
	var numRows = grantList.Size()
	var lastRecipient string = ""
	var balance = dec.Zero
	var totalGrants = dec.Zero
	var totalPayments = dec.Zero
	var totalWriteOffs = dec.Zero
	var totalTransfers = dec.Zero
	var totalBalance = dec.Zero
	grantList.Sort()

	//
	// Insert title and headings
	//
	writeCell(output, "A", row, "Recipient Actions")
	row += 2
	writeCell(output, "A", row, "Transaction Date")
	writeCell(output, "B", row, "Recipient")
	writeCell(output, "C", row, "Educational Institution")
	writeCell(output, "D", row, "Grant")
	writeCell(output, "E", row, "Payment")
	writeCell(output, "F", row, "Write-Off")
	writeCell(output, "G", row, "Transfer")
	writeCell(output, "H", row, "Net Balance")
	row++
	//
	// Loop through the transactions in the grant list
	//
	for index := 0; index < numRows; index++ {
		//
		// Extract data from transaction
		//
		var tran = grantList.Get(index)
		var transactionDate = tran.TransactionDate()
		var recipient = tran.Recipient()
		var edInst = tran.EducationalInstitution()
		var amount = tran.Amount()
		var transType = tran.TransType()
		//
		// Check if recipient has changed
		//
		if recipient != lastRecipient {
			writeCellDecimal(output, "H", row, balance)
			totalBalance = totalBalance.Add(balance)
			row++
			lastRecipient = recipient
			balance = dec.Zero
		}
		//
		// Write data to spreadsheet
		//
		writeCellDate(output, "A", row, transactionDate)
		writeCell(output, "B", row, recipient)
		writeCell(output, "C", row, edInst)
		if transType == g.Grant {
			totalGrants = totalGrants.Add(amount)
			writeCellDecimal(output, "D", row, amount)
			balance = balance.Add(amount)
		} else if transType == g.GrantPayment {
			totalPayments = totalPayments.Add(amount)
			writeCellDecimal(output, "E", row, amount)
			balance = balance.Sub(amount)
		} else if transType == g.WriteOff {
			totalWriteOffs = totalWriteOffs.Add(amount)
			writeCellDecimal(output, "F", row, amount)
			balance = balance.Sub(amount)
		} else if transType == g.Transfer {
			totalTransfers = totalTransfers.Add(amount)
			writeCellDecimal(output, "G", row, amount)
			balance = balance.Add(amount)
		} else {
			assert.Assert(false, "invalid value for transType: "+strconv.Itoa(int(transType)))
		}
		row++
	}
	totalBalance = totalBalance.Add(balance)
	writeCellDecimal(output, "H", row, balance)
	//
	// Write totals to spreadsheet
	//
	row += 2
	writeCell(output, "C", row, "Totals")
	writeCellDecimal(output, "D", row, totalGrants)
	writeCellDecimal(output, "E", row, totalPayments)
	writeCellDecimal(output, "F", row, totalWriteOffs)
	writeCellDecimal(output, "G", row, totalTransfers)
	writeCellDecimal(output, "H", row, totalBalance)
}

// outputRecipientSummary creates the RecipSum tab with the recipient summary.
func outputRecipientSummary(output *spreadsheet.SpreadsheetFile, grantList *g.GrantList) {
	//
	// Create the recipient summary list
	//
	var row = 1
	var list, err = g.AssembleRecipientSumList(grantList)
	check(err, "Error: ")
	var totalPayments = []dec.Decimal{dec.Zero, dec.Zero}
	var totalCount = []int{0, 0}
	//
	// Create Headings
	//
	writeCell(output, "A", row, "Recipient Summary")
	row += 2
	writeCell(output, "A", row, "Recipient Name")
	writeCell(output, "B", row, "FY2023 Count")
	writeCell(output, "C", row, "FY2023 Payments")
	writeCell(output, "D", row, "FY2024 Count")
	writeCell(output, "E", row, "FY2024 Payments")
	row++
	//
	// Loop through the recipient summaries
	//
	var names = list.Names()
	for _, name := range names {
		var recipientSum, err = list.Get(name)
		check(err, "Error: ")
		writeCell(output, "A", row, name)
		//
		// Compute FY2023 values
		//
		outputSumData(output, a.FY2023, row, "B", "C", recipientSum)
		if recipientSum.IsRecipient(a.FY2023) {
			totalCount[a.FY2023] += 1
		}
		var amount = recipientSum.PaymentTotal(a.FY2023)
		totalPayments[a.FY2023] = totalPayments[a.FY2023].Add(amount)
		//
		// Compute FY2024 values
		//
		outputSumData(output, a.FY2024, row, "D", "E", recipientSum)
		if recipientSum.IsRecipient(a.FY2024) {
			totalCount[a.FY2024] += 1
		}
		amount = recipientSum.PaymentTotal(a.FY2024)
		totalPayments[a.FY2024] = totalPayments[a.FY2024].Add(amount)
		row++
	}
	//
	// Create totals
	//
	writeCell(output, "A", row, "Totals")
	writeCellInt(output, "B", row, totalCount[a.FY2023])
	writeCellDecimal(output, "C", row, totalPayments[a.FY2023])
	writeCellInt(output, "D", row, totalCount[a.FY2024])
	writeCellDecimal(output, "E", row, totalPayments[a.FY2024])
}

// outputSumData inserts the recipient summary data into spreadsheet.
func outputSumData(
	output *spreadsheet.SpreadsheetFile,
	fiscalYear a.FYIndicator,
	row int,
	columnCount string,
	columnAmount string,
	recipientSum *g.RecipientSum) {
	var count = 0
	if recipientSum.IsRecipient(fiscalYear) {
		count = 1
	}
	var amount = recipientSum.PaymentTotal(fiscalYear)
	writeCellInt(output, columnCount, row, count)
	writeCellDecimal(output, columnAmount, row, amount)
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
	var err = outputPtr.SetCellDecimal(cell, value, spreadsheet.FormatMoney)
	check(err, "Error writing cell "+cell+": ")
}

// writeDate outputs a date value to the specified cell.
func writeCellDate(
	outputPtr *spreadsheet.SpreadsheetFile,
	column string,
	row int,
	date d.Date) {
	var cell = cellName(column, row)
	var err = outputPtr.SetCellDate(cell, date)
	check(err, "Error writing cell "+cell+":")
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
