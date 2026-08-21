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
	sp "acorn_go/pkg/spreadsheet"
	s "acorn_go/pkg/support"
	"fmt"
	"maps"
	"slices"
	"strconv"

	dec "github.com/shopspring/decimal"
	d "github.com/waysys/waydate/pkg/date"
)

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

const outputFile = "/home/bozo/Downloads/scholarship_analysis.xlsx"
const outputTab1 = "Scholarship"
const outputTab2 = "Individual Grant"

const (
	billFile = "/home/bozo/golang/acorn_go/data/bills.xlsx"
	billTab  = "Sheet1"
)

const (
	columnTransDate  = "Date"
	columnVendorName = "Vendor"
	columnBillType   = "Bills"
	columnAmount     = "Amount"
)

// ----------------------------------------------------------------------------
// Functions
// ----------------------------------------------------------------------------

// main supervises the processing of scholarship analysis.
func main() {
	var err error
	var sprdsht sp.Spreadsheet
	var billList []Entry
	var output sp.SpreadsheetFile

	printHeader()
	//
	// Read the bills
	//
	sprdsht, err = readBills()
	s.Check(err, "Error: ")
	//
	// Build bill list
	//
	billList = processBills(&sprdsht)
	//
	// Analyze the bill list
	//
	var billCount = processEntryCount(billList)
	fmt.Println("Number of scholarship categories for count: " +
		strconv.Itoa(len(billCount)))
	var billAmount = processEntryAmount(billList)
	fmt.Println("Number of scholarship categories for amount: " +
		strconv.Itoa(len(billAmount)))
	//
	// Create output spreadsheet
	//
	output, err = sp.New(outputFile, outputTab1)
	s.Check(err, "Error opening output file: ")
	var finish = func() {
		var err1 = output.Save()
		s.Check(err1, "Error saving output file: ")
		err1 = output.Close()
		s.Check(err1, "Error closing output file: ")
	}
	defer finish()
	//
	// Produce the scholarship analysis
	//
	outputScholarships(billCount, billAmount, &output)
	//
	// Produce the individual grant analysis
	//
	output, err = output.AddSheet(outputTab2)
	s.Check(err, "Error adding individual grant tab: ")
	outputIndividualGrants(billCount, billAmount, &output)
	printFooter()
}

// ----------------------------------------------------------------------------
// Spreadsheet Functions
// ----------------------------------------------------------------------------

// readBills reads the bill spreadsheet and returns the data
func readBills() (sp.Spreadsheet, error) {
	var sprdsht sp.Spreadsheet
	var err error
	//
	// Obtain spreadsheet data
	//
	sprdsht, err = sp.ProcessData(billFile, billTab)
	return sprdsht, err
}

// processBills reads the spreadsheet and creates an array of entries.
func processBills(sprdsht *sp.Spreadsheet) []Entry {
	var numRows = sprdsht.Size()
	var billList []Entry
	//
	// Loop through the spreadsheet
	//
	for row := 1; row < numRows; row++ {
		var entry, err = processBill(sprdsht, row)
		if err != nil {
			fmt.Printf("spreadsheet row %d: %s \n", row+1, err)
		} else if entry.Scholarship().TermPeriod() != Outside {
			billList = append(billList, entry)
		}
	}
	return billList
}

// processBill extracts a row from the spreadsheet, converts it to an
// entry, and returns the entry
func processBill(sprdsht *sp.Spreadsheet, row int) (Entry, error) {
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

// processEntryCount counts the number of bills with the same
// scholarship categories.
func processEntryCount(billList []Entry) map[Scholarship]int {
	var billCount = make(map[Scholarship]int)

	for _, entry := range billList {
		var scholarship = entry.Scholarship()
		billCount[scholarship]++
	}
	return billCount
}

// processEntryAmount adds up the amounts of scholarship
func processEntryAmount(billList []Entry) map[Scholarship]dec.Decimal {
	var billAmount = make(map[Scholarship]dec.Decimal)

	for _, entry := range billList {
		var scholarship = entry.Scholarship()
		var amount = entry.Amount()
		billAmount[scholarship] = billAmount[scholarship].Add(amount)
	}

	return billAmount
}

// ----------------------------------------------------------------------------
// Output Functions
// ----------------------------------------------------------------------------

// outputScholarships creates a spreadsheet with the number of scholarships
// in each scholarship category
func outputScholarships(
	billCount map[Scholarship]int,
	billAmount map[Scholarship]dec.Decimal,
	output *sp.SpreadsheetFile,
) {
	//
	// Insert Heading
	//
	var row = 1
	sp.WriteCell(output, "A", row, "Term Period")
	sp.WriteCell(output, "B", row, "Institution Type")
	sp.WriteCell(output, "C", row, "Account Type")
	sp.WriteCell(output, "D", row, "Enrollment Status")
	sp.WriteCell(output, "E", row, "Count")
	sp.WriteCell(output, "F", row, "Scholarship Amount")
	sp.WriteCell(output, "G", row, "Average Scholarship")
	row++
	//
	// Loop through the bill count
	//
	var totalCount = 0
	var totalAmount = dec.Zero
	scholarships := sortScholarships(billCount)
	for _, scholarship := range scholarships {
		if scholarship.AccountType() != IndividualGrant {
			//
			// Process Counts
			//
			var count = billCount[scholarship]
			sp.WriteCell(output, "A", row, string(scholarship.TermPeriod()))
			sp.WriteCell(output, "B", row, string(scholarship.InstitutionType()))
			sp.WriteCell(output, "C", row, string(scholarship.AccountType()))
			sp.WriteCell(output, "D", row, string(scholarship.EnrollmentStatus()))
			sp.WriteCellInt(output, "E", row, count)
			totalCount += count
			//
			// Process Scholarship Amounts
			//
			var amount = billAmount[scholarship]
			totalAmount = totalAmount.Add(amount)
			sp.WriteCellDecimal(output, "F", row, amount.Round(0))
			var avg = average(amount, count)
			sp.WriteCellDecimal(output, "G", row, avg)
			row++
		}
	}
	row++
	sp.WriteCell(output, "D", row, "Totals")
	sp.WriteCellInt(output, "E", row, totalCount)
	sp.WriteCellDecimal(output, "F", row, totalAmount.Round(0))
	var totalAverage = average(totalAmount, totalCount)
	sp.WriteCellDecimal(output, "G", row, totalAverage)
}

// average computes the average scholarship given the amount and the number
// of scholarships
func average(amount dec.Decimal, count int) dec.Decimal {
	var average = dec.Zero
	if count > 0 {
		var countDec = dec.NewFromInt(int64(count))
		average = amount.Div(countDec)
		average = average.Round(0)
	}
	return average
}

// sortScholarships fetches the keys and sorts them in a meaningful order
func sortScholarships(billCount map[Scholarship]int) []Scholarship {
	var compare = func(a, b Scholarship) int {
		return a.Compare(b)
	}
	var scholarships = slices.Collect(maps.Keys(billCount))
	slices.SortFunc(scholarships, compare)
	return scholarships
}

// outputIndividualGrant populates the second tab of the spreadsheet
// with information about individual grants.
func outputIndividualGrants(
	billCount map[Scholarship]int,
	billAmount map[Scholarship]dec.Decimal,
	output *sp.SpreadsheetFile,
) {
	//
	// Insert Heading
	//
	var row = 1
	sp.WriteCell(output, "A", row, "Term Period")
	sp.WriteCell(output, "B", row, "Count")
	sp.WriteCell(output, "C", row, "Total Grant Amount")
	sp.WriteCell(output, "D", row, "Average Grant Amount")
	row++
	//
	// Loop through the bill count
	//
	var totalCount = 0
	var totalAmount = dec.Zero
	scholarships := sortScholarships(billCount)
	for _, scholarship := range scholarships {
		if scholarship.AccountType() == IndividualGrant {
			//
			// Process Counts
			//
			var count = billCount[scholarship]
			sp.WriteCell(output, "A", row, string(scholarship.TermPeriod()))
			sp.WriteCellInt(output, "B", row, count)
			totalCount += count
			//
			// Process Scholarship Amounts
			//
			var amount = billAmount[scholarship]
			totalAmount = totalAmount.Add(amount)
			sp.WriteCellDecimal(output, "C", row, amount.Round(0))
			var avg = average(amount, count)
			sp.WriteCellDecimal(output, "D", row, avg)
			row++
		}
	}
	row++
	sp.WriteCell(output, "A", row, "Totals")
	sp.WriteCellInt(output, "B", row, totalCount)
	sp.WriteCellDecimal(output, "C", row, totalAmount.Round(0))
	var totalAverage = average(totalAmount, totalCount)
	sp.WriteCellDecimal(output, "D", row, totalAverage)
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
	fmt.Println("-----------------------------------------------------------")
	fmt.Println("Program is finished.")
	fmt.Println("-----------------------------------------------------------")
}
