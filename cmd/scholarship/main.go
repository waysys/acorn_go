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
	sp "acorn_go/pkg/spreadsheet"
	s "acorn_go/pkg/support"
	"fmt"
	"strconv"

	"github.com/waysys/assert/assert"

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
	s.Check(err, "Error: ")

	defer output.Close()
	//
	// Produce summary tab
	//
	outputGrantSummary(&output, &grantList)
	//
	// Produce recipient tab
	//
	output, err = output.AddSheet("Recipients")
	s.Check(err, "Error: ")
	outputRecipientList(&output, &grantList)
	//
	// Produce recipient summary tab
	//
	output, err = output.AddSheet("RecipSum")
	s.Check(err, "Error: ")
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
	sp.WriteCell(output, "A", row, "Summary of Scholarship Payments")
	row += 2
	//
	// Insert headings
	//
	sp.WriteCell(output, "A", row, "Fiscal Year")
	sp.WriteCell(output, "B", row, "Total Grants")
	sp.WriteCell(output, "C", row, "Total Payments")
	sp.WriteCell(output, "D", row, "Total Net Write-Offs")
	sp.WriteCell(output, "E", row, "Total Refunds")
	row++
	//
	// Amounts
	//
	for _, fy := range a.FYIndicators {
		outputGrantSummaryLine(output, grantList, row, fy)
		row++
	}
	row++
}

// outputGrantSummaryLine inserts transactions amounts for the specified fiscal year
// into the spreadsheet.
func outputGrantSummaryLine(
	output *spreadsheet.SpreadsheetFile,
	grantList *g.GrantList,
	row int,
	fy a.FYIndicator) {
	sp.WriteCell(output, "A", row, fy.String())
	var grantTotal = grantList.TotalTransAmount(fy, g.Grant)
	sp.WriteCellDecimal(output, "B", row, grantTotal)
	var paymentTotal = grantList.TotalTransAmount(fy, g.GrantPayment)
	sp.WriteCellDecimal(output, "C", row, paymentTotal)
	var writeOffTotal = grantList.TotalNetWriteOff(fy)
	sp.WriteCellDecimal(output, "D", row, writeOffTotal)
	var refundTotal = grantList.TotalTransAmount(fy, g.Refund)
	sp.WriteCellDecimal(output, "E", row, refundTotal)
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
	var totalRefunds = dec.Zero
	var totalBalance = dec.Zero
	grantList.Sort()

	//
	// Insert title and headings
	//
	sp.WriteCell(output, "A", row, "Recipient Actions")
	row += 2
	sp.WriteCell(output, "A", row, "Transaction Date")
	sp.WriteCell(output, "B", row, "Recipient")
	sp.WriteCell(output, "C", row, "Educational Institution")
	sp.WriteCell(output, "D", row, "Grant")
	sp.WriteCell(output, "E", row, "Payment")
	sp.WriteCell(output, "F", row, "Write-Off")
	sp.WriteCell(output, "G", row, "Transfer")
	sp.WriteCell(output, "H", row, "Refunds")
	sp.WriteCell(output, "I", row, "Net Balance")
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
			sp.WriteCellDecimal(output, "I", row, balance)
			totalBalance = totalBalance.Add(balance)
			row++
			lastRecipient = recipient
			balance = dec.Zero
		}
		//
		// Write data to spreadsheet
		//
		sp.WriteCellDate(output, "A", row, transactionDate)
		sp.WriteCell(output, "B", row, recipient)
		sp.WriteCell(output, "C", row, edInst)
		switch transType {
		case g.Grant:
			totalGrants = totalGrants.Add(amount)
			sp.WriteCellDecimal(output, "D", row, amount)
			balance = balance.Add(amount)
		case g.GrantPayment:
			totalPayments = totalPayments.Add(amount)
			sp.WriteCellDecimal(output, "E", row, amount)
			balance = balance.Sub(amount)
		case g.WriteOff:
			totalWriteOffs = totalWriteOffs.Add(amount)
			sp.WriteCellDecimal(output, "F", row, amount)
			balance = balance.Sub(amount)
		case g.Transfer:
			totalTransfers = totalTransfers.Add(amount)
			sp.WriteCellDecimal(output, "G", row, amount)
			balance = balance.Add(amount)
		case g.Refund:
			totalRefunds = totalRefunds.Add(amount)
			sp.WriteCellDecimal(output, "H", row, amount)
			balance = balance.Add(amount)
		default:
			assert.Assert(false, "invalid value for transType: "+strconv.Itoa(int(transType)))
		}
		row++
	}
	totalBalance = totalBalance.Add(balance)
	sp.WriteCellDecimal(output, "H", row, balance)
	//
	// Write totals to spreadsheet
	//
	row += 2
	sp.WriteCell(output, "C", row, "Totals")
	sp.WriteCellDecimal(output, "D", row, totalGrants)
	sp.WriteCellDecimal(output, "E", row, totalPayments)
	sp.WriteCellDecimal(output, "F", row, totalWriteOffs)
	sp.WriteCellDecimal(output, "G", row, totalTransfers)
	sp.WriteCellDecimal(output, "H", row, totalRefunds)
	sp.WriteCellDecimal(output, "I", row, totalBalance)
}

// outputRecipientSummary creates the RecipSum tab with the recipient summary.
func outputRecipientSummary(output *spreadsheet.SpreadsheetFile, grantList *g.GrantList) {
	//
	// Create the recipient summary list
	//
	var row = 1
	var list, err = g.AssembleRecipientSumList(grantList)
	s.Check(err, "Error: ")
	var totalPayments = []dec.Decimal{dec.Zero, dec.Zero, dec.Zero}
	var totalCount = []int{0, 0, 0}
	//
	// Create Headings
	//
	sp.WriteCell(output, "A", row, "Recipient Summary")
	row += 2
	sp.WriteCell(output, "A", row, "Recipient Name")
	sp.WriteCell(output, "B", row, "FY2023 Count")
	sp.WriteCell(output, "C", row, "FY2023 Payments")
	sp.WriteCell(output, "D", row, "FY2024 Count")
	sp.WriteCell(output, "E", row, "FY2024 Payments")
	sp.WriteCell(output, "F", row, "FY2025 Count")
	sp.WriteCell(output, "G", row, "FY2025 Payments")
	row++
	//
	// Loop through the recipient summaries
	//
	var names = list.Names()
	for _, name := range names {
		var recipientSum, err = list.Get(name)
		s.Check(err, "Error: ")
		sp.WriteCell(output, "A", row, name)
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
		//
		// Compute FY2025 values
		//
		outputSumData(output, a.FY2025, row, "F", "G", recipientSum)
		if recipientSum.IsRecipient(a.FY2025) {
			totalCount[a.FY2025] += 1
		}
		amount = recipientSum.PaymentTotal(a.FY2025)
		totalPayments[a.FY2025] = totalPayments[a.FY2025].Add(amount)
		row++
	}
	//
	// Create totals
	//
	row++
	sp.WriteCell(output, "A", row, "Totals")
	sp.WriteCellInt(output, "B", row, totalCount[a.FY2023])
	sp.WriteCellDecimal(output, "C", row, totalPayments[a.FY2023])
	sp.WriteCellInt(output, "D", row, totalCount[a.FY2024])
	sp.WriteCellDecimal(output, "E", row, totalPayments[a.FY2024])
	sp.WriteCellInt(output, "F", row, totalCount[a.FY2025])
	sp.WriteCellDecimal(output, "G", row, totalPayments[a.FY2025])
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
	sp.WriteCellInt(output, columnCount, row, count)
	sp.WriteCellDecimal(output, columnAmount, row, amount)
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
