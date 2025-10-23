// ----------------------------------------------------------------------------
//
// # Individual
//
// This program produces individual grant spreadsheet which includes the amount
// for each recipient for the fiscal year and the number of recipients by
// fiscal year.
//
// Author: William Shaffer
//
//	Copyright (c) 2025 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------
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
)

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

const outputFile = "/home/bozo/Downloads/individual.xlsx"

// ----------------------------------------------------------------------------
// Main Function
// ----------------------------------------------------------------------------

// main supervises the processing of individual data.
func main() {
	var err error = nil
	var apTranlist q.TransList
	var grantList g.GrantList
	var output sp.SpreadsheetFile

	printHeader()
	//
	// Read the accounts payable transaction
	//
	apTranlist, err = q.ReadAPTransactions()
	s.Check(err, "Error: ")
	//
	// Generate grant transactions
	//
	grantList, err = g.AssembleIndividualGrantList(&apTranlist)
	s.Check(err, "Error: ")
	//
	// Output results
	//
	output, err = sp.New(outputFile, "Individual Grants")
	s.Check(err, "Error: ")
	var finish = func() {
		err = output.Save()
		s.Check(err, "Error saving output file: ")
		err = output.Close()
		s.Check(err, "Error closing output file: ")
	}
	defer finish()
	outputRecipientSummary(&output, &grantList)
	output.Save()
	printFooter()
}

// ----------------------------------------------------------------------------
// Main Function
// ----------------------------------------------------------------------------

// outputRecipientSummary creates the RecipSum tab with the recipient summary.
func outputRecipientSummary(output *sp.SpreadsheetFile, grantList *g.GrantList) {
	//
	// Create the recipient summary list
	//
	var row = 1
	var list, err = g.AssembleRecipientSumList(grantList)
	s.Check(err, "Error: ")
	var totalPayments = []dec.Decimal{dec.Zero, dec.Zero, dec.Zero}
	var totalCount = []int{0, 0, 0}
	var countColumns = []string{"B", "D", "F"}
	var paymentColumns = []string{"C", "E", "G"}
	var amount = dec.Zero
	//
	// Create Headings
	//
	sp.WriteCell(output, "A", row, "Individual Grant Recipient Summary")
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
		// Cycle through fiscal years
		//
		for _, fy := range a.FYIndicators {
			outputSumData(output, fy, row, countColumns[fy], paymentColumns[fy], recipientSum)
			if recipientSum.IsIndividualGrant(fy) {
				totalCount[fy] += 1
			}
			amount = recipientSum.GrantTotal(fy)
			totalPayments[fy] = totalPayments[fy].Add(amount)
		}
		row++
	}
	//
	// Create total payments
	//
	row++
	sp.WriteCell(output, "A", row, "Total Payments")
	for _, fy := range a.FYIndicators {
		sp.WriteCellInt(output, countColumns[fy], row, totalCount[fy])
		sp.WriteCellDecimal(output, paymentColumns[fy], row, totalPayments[fy])
	}
}

// outputSumData inserts the recipient summary data into spreadsheet.
func outputSumData(
	output *sp.SpreadsheetFile,
	fiscalYear a.FYIndicator,
	row int,
	columnCount string,
	columnAmount string,
	recipientSum *g.RecipientSum) {
	var count = 0
	if recipientSum.IsIndividualGrant(fiscalYear) {
		count = 1
	}
	var amount = recipientSum.GrantTotal(fiscalYear)
	sp.WriteCellInt(output, columnCount, row, count)
	sp.WriteCellDecimal(output, columnAmount, row, amount)
}

// ----------------------------------------------------------------------------
// Print Functions
// ----------------------------------------------------------------------------

// printHeader places the header information at the top of the page
func printHeader() {
	fmt.Println("-----------------------------------------------------------")
	fmt.Println("Acorn Scholarship Fund Individual Grant Analysis")
	fmt.Println("-----------------------------------------------------------")
}

// printFooter prints a message indicating completion of the process
func printFooter() {
	fmt.Println("Program is finished.")
	fmt.Println("-----------------------------------------------------------")
}
