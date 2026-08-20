// ----------------------------------------------------------------------------
//
// # Scholarship Analysis
//
// This program calculates metrics about the scholarships and grants
// in the specified year.
//
// Author: William Shaffer
//
// Copyright (c) 2026 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

package main

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	"acorn_go/pkg/spreadsheet"
	s "acorn_go/pkg/support"
	"fmt"

	dec "github.com/shopspring/decimal"
	d "github.com/waysys/waydate/pkg/date"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

const SCHOLARSHIP_PROGRAM = "100 Scholarship Program"
const INDIVIDUAL_GRANT_PROGRAM = "110 Individual Grant Program"

const DEPENDENT = "Dependent"
const ASSOCIATE = "Grant"

const outputFile = "/home/bozo/Downloads/scholarship_analysis.xlsx"

const (
	billFile = "/home/bozo/golang/acorn_go/data/bills.xlsx"
	billTab  = "Sheet1"
)

const (
	columnTransDate  = "Date"
	columnVendorName = "Vendor"
	columnBillType   = "Bills"
	columnItemType   = "Item class"
	columnAmount     = "Amount"
)

// ----------------------------------------------------------------------------
// Functions
// ----------------------------------------------------------------------------

// main supervises the processing of scholarship analysis.
func main() {
	var err error = nil
	var sprdsht spreadsheet.Spreadsheet
	var billList []Entry

	printHeader()
	//
	// Read the bills
	//
	sprdsht, err = readBills()
	s.Check(err, "Error: ")
	//
	// Build bill list
	//
	billList, err = processBills(&sprdsht)
	s.Check(err, "Error: ")
	//
	// Analyze the bill list
	//
	processEntries(billList)
}

// ----------------------------------------------------------------------------
// Spreadsheet Functions
// ----------------------------------------------------------------------------

// readBills reads the bill spreadsheet and returns the data
func readBills() (spreadsheet.Spreadsheet, error) {
	var sprdsht spreadsheet.Spreadsheet
	var err error
	//
	// Obtain spreadsheet data
	//
	sprdsht, err = spreadsheet.ProcessData(billFile, billTab)
	return sprdsht, err
}

// processBills reads the spreadsheet and creates an array of entries.
func processBills(sprdsht *spreadsheet.Spreadsheet) ([]Entry, error) {
	var numRows = sprdsht.Size()
	var billList []Entry
	//
	// Loop through the spreadsheet
	//
	for row := 1; row < numRows; row++ {
		var entry, err = processBill(sprdsht, row)
		if err != nil {
			return billList, err
		}
		if entry.Scholarship().termPeriod != Outside {
			billList = append(billList, entry)
		}
	}
	return billList, nil
}

// processBill extracts a row from the spreadsheet, converts it to an
// entry, and returns the entry
func processBill(sprdsht *spreadsheet.Spreadsheet, row int) (Entry, error) {
	var billDate d.Date
	var university string
	var tag Tag
	var amount dec.Decimal
	var value string
	var err error
	var entry Entry
	//
	// Read data from spreadsheet
	//
	billDate, err = sprdsht.CellDate(row, columnTransDate)
	if err == nil {
		university, err = sprdsht.Cell(row, columnVendorName)
	}
	if err == nil {
		value, err = sprdsht.Cell(row, columnBillType)
	}
	if err == nil {
		tag = Tag(value)
	}
	if err == nil {
		value, err = sprdsht.Cell(row, columnAmount)
	}
	if err == nil {
		amount, err = dec.NewFromString(value)
	}
	//
	// Convert data to an entry
	//
	if err == nil {
		entry, err = NewEntry(
			row,
			billDate,
			university,
			tag,
			amount,
		)
	}
	return entry, err
}

// ----------------------------------------------------------------------------
// Process List of Bills
// ----------------------------------------------------------------------------

// Generate analysis of scholarships
func processEntries(billList []Entry) {
	var total = dec.Zero
	for _, entry := range billList {
		var amount = entry.Amount()
		total = total.Add(amount)
	}
	fmt.Print("Total scholarships and grants: " + total.String())
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
