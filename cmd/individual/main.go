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
	var totalPayments []dec.Decimal
	var totalCount []int
	var amount = dec.Zero
	var grandCount = 0
	// The column where the first count column is placed
	var startColumn = "B"
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
	sp.CreateHeadings(output, row, startColumn, "Individual Grant Recipient Summary")
	row += 3
	//
	// Loop through the recipient summaries
	//
	var names = list.Names()
	var ch = sp.NewColumnHandler(startColumn)
	for _, name := range names {
		var recipientSum, err = list.Get(name)
		s.Check(err, "Error: ")
		sp.WriteCell(output, "A", row, name)
		//
		// Cycle through fiscal years
		//
		for _, fy := range a.FYIndicators {
			if recipientSum.IsIndividualGrant(fy) {
				count = 1
				totalCount[fy] += 1
			} else {
				count = 0
			}
			amount = recipientSum.GrantTotal(fy)
			totalPayments[fy] = totalPayments[fy].Add(amount)
			//
			// Output results
			//
			sp.OutputSumData(output, fy, row, &ch, count, amount)
		}
		row++
	}
	//
	// Create total payments
	//
	row++
	// grandCount is the total number of recipients (one row per name).
	grandCount = len(names)
	sp.OutputTotals(output, row, &ch, totalCount, totalPayments, grandCount)
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
