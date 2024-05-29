// ----------------------------------------------------------------------------
//
// Donor series program
//
// Author: William Shaffer
// Version: 12-Apr-2024
//
// Copyright (c) 2024 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

// This program creates a spreasheet with a timeline of
// amount of donations.

package main

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"

	a "acorn_go/pkg/accounting"
	"acorn_go/pkg/donorinfo"
	"acorn_go/pkg/spreadsheet"

	d "github.com/waysys/waydate/pkg/date"

	"github.com/waysys/assert/assert"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

const (
	inputFile      = "/home/bozo/golang/acorn_go/data/register.xlsx"
	tab            = "Worksheet"
	outputFileName = "/home/bozo/Downloads/donations.xlsx"
	sheetName      = "Donations"
)

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

// ----------------------------------------------------------------------------
// Functions
// ----------------------------------------------------------------------------

// main supervises the execution of this program.  It produces a spreadsheet
// with the donation series by year and month
func main() {
	var sprdsht spreadsheet.Spreadsheet
	var donationSeries donorinfo.DonationSeries
	var err error

	printHeader()
	//
	// Obtain spreadsheet data
	//
	sprdsht, err = spreadsheet.ProcessData(inputFile, tab)
	check(err, "Error processing spreadsheet")
	//
	// Generate donation series
	//
	donationSeries, err = donorinfo.NewDonationSeries(&sprdsht)
	check(err, "Error generating donation series")
	//
	// Output donation series to spreadsheet
	//
	err = outputDonationSeries(&donationSeries)
	check(err, "Error writing output")
	os.Exit(0)
}

// printHeader places the header information at the top of the page
func printHeader() {
	fmt.Println("-----------------------------------------------------------")
	fmt.Println("Acorn Scholarship Fund Donation Series")
	fmt.Println("-----------------------------------------------------------")
}

// outputDonationSeries produces a spreadsheet with the donation time series
func outputDonationSeries(dsPtr *(donorinfo.DonationSeries)) error {
	var err error = nil
	var output spreadsheet.SpreadsheetFile
	var keys []d.YearMonth
	var value string
	var totalDonationsFY2023 float64 = 0.0
	var totalDonationsFY2024 float64 = 0.0
	var row int
	var yearMonth d.YearMonth
	var avg float64 = 0.0

	if dsPtr == nil {
		err = errors.New("pointer to donation series is nil")
		return err
	}

	output, err = spreadsheet.New(outputFileName, sheetName)
	if err != nil {
		return err
	}
	defer output.Close()
	//
	// Enter values in into spreadsheet
	//
	keys, err = d.Keys()
	if err != nil {
		return err
	}
	//
	// Insert title
	//
	writeCell(&output, "A", 1, "Month/Year")
	writeCell(&output, "B", 1, "Donation Amount")
	writeCell(&output, "C", 1, "Number of Donations")
	writeCell(&output, "D", 1, "Average Amount of Donation")

	for row, yearMonth = range keys {
		//
		// Insert month and year
		//
		value = yearMonth.String()
		writeCell(&output, "A", row+2, value)
		//
		// Insert donation for that month and year
		//
		amount := math.Round((*dsPtr).GetAmount(yearMonth))
		writeCellFloat(&output, "B", row+2, amount)
		//
		// Insert the number of donors for that month and year
		//
		count := (*dsPtr).GetCount(yearMonth)
		writeCellInt(&output, "C", row+2, count)
		//
		// Insert the average donation
		//
		if count == 0 {
			avg = 0.00
		} else {
			avg = math.Round(amount / float64(count))
		}
		writeCellFloat(&output, "D", row+2, avg)
		//
		// Calculate total donations for FY2023 and FY2024
		//
		totalDonationsFY2023, totalDonationsFY2024 =
			calcTotalDonations(yearMonth, amount, totalDonationsFY2023, totalDonationsFY2024)
	}
	//
	// Ootput totals
	//
	row += 4
	outputTotals(output, "FY2023 Donations", totalDonationsFY2023, row)
	outputTotals(output, "FY2024 Donations", totalDonationsFY2024, row+1)
	//
	// Save the spreadsheet
	//
	err = output.Save()
	return err
}

// caldTotalDonations determines the fiscal year of the donation from the year month and
// adds the amount to the appropriate total donation.
func calcTotalDonations(
	yearMonth d.YearMonth,
	amount float64,
	totDonFY2023 float64,
	totDonFY2024 float64) (float64, float64) {
	var err error = nil
	var indicator a.FYIndicator

	indicator, err = a.FiscalYearFromYearMonth(yearMonth)
	assert.Assert(err == nil, "Unable to handle year month: "+yearMonth.String())
	switch indicator {
	case a.FY2023:
		totDonFY2023 += amount
	case a.FY2024:
		totDonFY2024 += amount
	default:
		assert.Assert(false, "unmatched fiscal year: "+yearMonth.String())
	}
	return totDonFY2023, totDonFY2024
}

// outputTotals places the donation totals for each fiscal year in rows below the
// time series
func outputTotals(output spreadsheet.SpreadsheetFile, title string, totalDonation float64, row int) {
	writeCell(&output, "A", row, title)
	writeCellFloat(&output, "B", row, totalDonation)
}
