// ----------------------------------------------------------------------------
//
// Scholarship Analysis Program
//
// Author: William Shaffer
//
// Copyright (c) 2024, 2025 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

// This program produces a spreadsheet with the anaysis of scholarship payments.
// Three reports are generated:
// -- Summary of Scholarship Payments
// -- Recipient Actions
// -- Recipient Summary
// -- Name Tags

package main

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	a "acorn_go/pkg/accounting"
	g "acorn_go/pkg/grants"
	q "acorn_go/pkg/quickbooks"
	sp "acorn_go/pkg/spreadsheet"
	s "acorn_go/pkg/support"
	"fmt"

	dec "github.com/shopspring/decimal"
	"github.com/waysys/assert/assert"
)

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

const outputFile = "/home/bozo/Downloads/scholarship.xlsx"

// output tabs
const (
	summary          = "Summary"
	recipients       = "Recipients"
	recipientSummary = "RecipSum"
	nameTags         = "Name Tags"
)

// startColumn is the first column where the fiscal year column starts
const startColumn = "B"

// ----------------------------------------------------------------------------
// Main Function
// ----------------------------------------------------------------------------

// main supervises the processing of scholarship data.
func main() {
	var err error = nil
	var apTranlist q.TransList
	var billList q.BillList
	var grantList g.GrantList
	var output sp.SpreadsheetFile

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
		output, err = sp.New(outputFile, summary)
	}
	s.Check(err, "Error: ")
	//
	// Defer function to close spreadsheet
	//
	var finish = func() {
		err = output.Save()
		s.Check(err, "Error saving output file: ")
		err = output.Close()
		s.Check(err, "Error closing output file: ")
	}
	defer finish()
	//
	// Produce summary tab
	//
	outputGrantSummary(&output, &grantList)
	//
	// Produce recipient tab
	//
	output, err = output.AddSheet(recipients)
	s.Check(err, "Error: ")
	outputRecipientList(&output, &grantList)
	//
	// Produce recipient summary tab
	//
	output, err = output.AddSheet(recipientSummary)
	s.Check(err, "Error: ")
	outputRecipientSummary(&output, &grantList)
	//
	// Produce the name tag list
	//
	output, err = output.AddSheet(nameTags)
	s.Check(err, "Error: ")
	outputNameTagList(&output, &grantList)
	//
	// Save results
	//
	output.Save()
	s.Check(err, "Error: ")
	printFooter()
}

// ----------------------------------------------------------------------------
// Spreadsheet Functions
// ----------------------------------------------------------------------------

// outputGrantSummary populates the Summary tab with the amounts by transaction
// type
func outputGrantSummary(output *sp.SpreadsheetFile, grantList *g.GrantList) {
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
	sp.WriteCell(output, "F", row, "Net Balance")
	row++
	//
	// Amounts
	//
	for _, fy := range a.FYIndicators {
		outputGrantSummaryLine(output, grantList, row, fy)
		row++
	}
	row++
	outputTotalLine(output, grantList, row)
}

// outputGrantSummaryLine inserts transactions amounts for the specified fiscal year
// into the spreadsheet.
func outputGrantSummaryLine(
	output *sp.SpreadsheetFile,
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
	var netBalance = grantList.NetBalance(fy)
	sp.WriteCellDecimal(output, "F", row, netBalance)
}

// outputTotalLine inserts the totals for the types of transactions.
func outputTotalLine(
	output *sp.SpreadsheetFile,
	grantList *g.GrantList,
	row int) {
	sp.WriteCell(output, "A", row, "Totals")
	sp.WriteCellDecimal(output, "B", row, grantList.GrandTotalTransactions(g.Grant))
	sp.WriteCellDecimal(output, "C", row, grantList.GrandTotalTransactions(g.GrantPayment))
	sp.WriteCellDecimal(output, "D", row, grantList.GrandTotalNetWriteoff())
	sp.WriteCellDecimal(output, "E", row, grantList.GrandTotalTransactions(g.Refund))
	sp.WriteCellDecimal(output, "F", row, grantList.TotalNetBalance())
}

// outputRecipientList produces a list of recipients organized by fiscal year and
// award group
func outputRecipientList(output *sp.SpreadsheetFile, grantList *g.GrantList) {
	var row = 1
	var numRows = grantList.Size()
	var lastRecipient string = ""
	var balance = dec.Zero
	var totals = []dec.Decimal{dec.Zero, dec.Zero, dec.Zero, dec.Zero, dec.Zero}
	var columns = []string{"E", "F", "G", "H", "I"}
	var totalBalance = dec.Zero

	grantList.Sort()

	//
	// Insert title and headings
	//
	sp.WriteCell(output, "A", row, "Recipient Actions")
	row += 2
	sp.WriteCell(output, "A", row, "Transaction Date")
	sp.WriteCell(output, "B", row, "Dependent?")
	sp.WriteCell(output, "C", row, "Recipient")
	sp.WriteCell(output, "D", row, "Educational Institution")
	sp.WriteCell(output, columns[g.Grant], row, "Grant")
	sp.WriteCell(output, columns[g.GrantPayment], row, "Payment")
	sp.WriteCell(output, columns[g.WriteOff], row, "Write-Off")
	sp.WriteCell(output, columns[g.Transfer], row, "Transfer")
	sp.WriteCell(output, columns[g.Refund], row, "Refunds")
	sp.WriteCell(output, "J", row, "Net Balance")
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
		var dependent = tran.Dependent()
		//
		// Check if recipient has changed
		//
		if recipient != lastRecipient {
			sp.WriteCellDecimal(output, "J", row, balance)
			totalBalance = totalBalance.Add(balance)
			row++
			lastRecipient = recipient
			balance = dec.Zero
		}
		//
		// Write data to spreadsheet
		//
		sp.WriteCellDate(output, "A", row, transactionDate)
		sp.WriteCell(output, "B", row, dependent)
		sp.WriteCell(output, "C", row, recipient)
		sp.WriteCell(output, "D", row, edInst)
		sp.WriteCellDecimal(output, columns[transType], row, amount)
		totals[transType] = totals[transType].Add(amount)
		balance = computeBalance(balance, amount, transType)
		row++
	}
	totalBalance = totalBalance.Add(balance)
	sp.WriteCellDecimal(output, "J", row, balance)
	//
	// Write totals to spreadsheet
	//
	row += 2
	sp.WriteCell(output, "D", row, "Totals")
	for _, transType := range g.TransTypes {
		sp.WriteCellDecimal(output, columns[transType], row, totals[transType])
	}
	sp.WriteCellDecimal(output, "J", row, totalBalance)
}

// computeBalance calculates the remaning balance based on the type of transaction.
func computeBalance(balance dec.Decimal, amount dec.Decimal, transType g.TransType) dec.Decimal {
	switch transType {
	case g.Grant:
		balance = balance.Add(amount)
	case g.GrantPayment:
		balance = balance.Sub(amount)
	case g.WriteOff:
		balance = balance.Sub(amount)
	case g.Transfer:
		balance = balance.Add(amount)
	case g.Refund:
		balance = balance.Add(amount)
	default:
		assert.Assert(false, "Invalid transaction type: "+transType.String())
	}
	return balance
}

// outputRecipientSummary creates the RecipSum tab with the recipient summary.
func outputRecipientSummary(output *sp.SpreadsheetFile, grantList *g.GrantList) {
	//
	// Create the recipient summary list
	//
	var row = 1
	var totalPayments []dec.Decimal
	var totalCount []int
	var amount = dec.Zero
	var count = 0
	//
	// Initialize arrays for each fiscal year
	//
	for index := 0; index < a.NumFiscalYears; index++ {
		totalCount = append(totalCount, 0)
		totalPayments = append(totalPayments, dec.Zero)
	}
	//
	// Create Headings
	//
	sp.CreateHeadings(output, row, startColumn, "Recipient Summary")
	row += 3
	//
	// Loop through the recipient summaries
	//
	var list, err = g.AssembleRecipientSumList(grantList)
	s.Check(err, "Error: ")
	var names = list.Names()
	var grandCount = 0
	var ch = sp.NewColumnHandler(startColumn)
	//
	// Cycle through recipients
	//
	for _, name := range names {
		var recipientSum, err = list.Get(name)
		s.Check(err, "Error: ")
		sp.WriteCell(output, "A", row, name)
		grandCount++
		//
		// Cycle through fiscal years
		//
		for _, fy := range a.FYIndicators {
			if recipientSum.IsRecipient(fy) {
				count = 1
			} else {
				count = 0
			}
			amount = recipientSum.NetPaymentTotal(fy)
			totalCount[fy] += count
			totalPayments[fy] = totalPayments[fy].Add(amount)
			//
			// Output the values
			//
			sp.OutputSumData(output, fy, row, &ch, count, amount)
		}
		row++
	}
	//
	// Insert total counts, total payments, and grand count
	//
	row++
	sp.OutputTotals(output, row, &ch, totalCount, totalPayments, grandCount)
}

// outputNameTagList produces a list of recipients that have been given awards,
// but not necessarily received payments.
func outputNameTagList(output *sp.SpreadsheetFile, grantList *g.GrantList) {
	var row = 1
	var names = grantList.Names()
	sp.WriteCell(output, "A", row, "Recipient")
	row++
	//
	// Loop through grant list names
	//
	for _, name := range names {
		name = q.RevertName(name)
		sp.WriteCell(output, "A", row, name)
		row++
	}
}

// ----------------------------------------------------------------------------
// Print Functions
// ----------------------------------------------------------------------------

// printHeader prints a message indicating that the program has started.
func printHeader() {
	fmt.Println("-----------------------------------------------------------")
	fmt.Println("Acorn Scholarship Fund Scholarship Analysis")
	fmt.Println("-----------------------------------------------------------")
}

// printFooter prints a message indicating completion of the program
func printFooter() {
	fmt.Println("Program is finished.")
	fmt.Println("-----------------------------------------------------------")
}
