// ----------------------------------------------------------------------------
//
// Project Grants
//
// Author: William Shaffer
// Version: 10-May-2024
//
// Copyright (c) 2024 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

// This program produces a spreadsheet grants.xlsx that contains the number
// of scholarshp awards by award groups.

package main

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	g "acorn_go/pkg/grants"
	sp "acorn_go/pkg/spreadsheet"
	s "acorn_go/pkg/support"
	"errors"
	"fmt"

	d "github.com/waysys/waydate/pkg/date"
)

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

const outputFile = "/home/bozo/Downloads/grants.xlsx"

const (
	billFile = "/home/bozo/golang/acorn_go/data/bills.xlsx"
	billTab  = "Sheet1"
)

// ----------------------------------------------------------------------------
// Main Function
// ----------------------------------------------------------------------------

// main supervises the calculation of grants per term
func main() {
	var err error = nil
	var sprdsht sp.Spreadsheet

	printHeader()
	//
	//  Read the bill spreadsheet
	//
	sprdsht, err = sp.ProcessData(billFile, billTab)
	s.Check(err, "Error: ")
	var numRows = sprdsht.Size()
	//
	// Loop through spreadsheet to count bills
	//
	for row := 1; row < numRows; row++ {
		err = processBill(&sprdsht, row)
		s.Check(err, "Error: ")
	}
	//
	// Output the results
	//
	outputGrantCount()
	printFooter()
}

// ----------------------------------------------------------------------------
// Processing Function
// ----------------------------------------------------------------------------

// processBill processes a row in the bills.xslx spreadsheet and adds the
// bill to the appropriate award group.
func processBill(sprdsht *sp.Spreadsheet, row int) error {
	var err error
	var billDate d.Date
	var vendorName string
	var billType string
	var awardGroup *g.AwardGroup

	const (
		columnTransDate     = "Bill date"
		columnVendorName    = "Vendor"
		columnRecipientName = "Memo"
		columnBillType      = "Bills"
	)
	//
	// Read data from spreadsheet
	//
	billDate, err = sprdsht.CellDate(row, columnTransDate)
	if err == nil {
		vendorName, err = sprdsht.Cell(row, columnVendorName)
	}
	if err == nil {
		billType, err = sprdsht.Cell(row, columnBillType)
	}
	//
	// Determine award group for this bill
	//
	if err == nil {
		if isEducationalBill(vendorName, billType) {
			awardGroup, err = g.FindAwardGroup(billDate)
			if err == nil {
				awardGroup.AddBill()
			}
		}
	}

	return err
}

// isEducationalBill returns true if the row is a valid education bill
func isEducationalBill(recipientName string, billType string) bool {
	var result = false

	switch billType {
	case "Grant":
		result = true
	case "Transfer":
		result = false
	case "":
		result = isEducational(recipientName)
	default:
		var err = errors.New("Unknown bill type: " + billType)
		s.Check(err, "")
	}

	return result
}

// isEducational returns true if the recipient is an educational institution.
func isEducational(recipientName string) bool {
	var result = false

	switch recipientName {
	case "Mr. Albert Martinelli":
		result = false
	case "Cypress of Raleigh Club":
		result = false
	case "Mr. Larry Potter":
		result = false
	default:
		result = true
	}
	return result
}

// ----------------------------------------------------------------------------
// Output Functions
// ----------------------------------------------------------------------------

// outputGrantCount creates a spreadsheet with the counts of the grants
func outputGrantCount() {
	var output sp.SpreadsheetFile
	var err error

	output, err = sp.New(outputFile, "Summary")
	s.Check(err, "Error: ")
	defer output.Close()
	//
	// Write title
	//
	var row = 1
	sp.WriteCell(&output, "A", row, "Summary of Grants")
	row += 2
	//
	// Write heading
	//
	sp.WriteCell(&output, "A", row, "Award Group")
	sp.WriteCell(&output, "B", row, "Number of Awards")
	row++
	//
	// Loop through slice of award groups
	//
	var groups = g.Groups()
	for _, awardGroup := range groups {
		sp.WriteCell(&output, "A", row, awardGroup.GroupName())
		sp.WriteCellInt(&output, "B", row, awardGroup.Count())
		row++
	}
	output.Save()
}

// ----------------------------------------------------------------------------
// Print Functions
// ----------------------------------------------------------------------------

// printHeader places the header information at the top of the page
func printHeader() {
	fmt.Println("-----------------------------------------------------------")
	fmt.Println("Acorn Scholarship Fund Grant Projection")
	fmt.Println("-----------------------------------------------------------")
}

// printFooter prints a message indicating completion of the process
func printFooter() {
	fmt.Println("Program is finished.")
	fmt.Println("-----------------------------------------------------------")
}
