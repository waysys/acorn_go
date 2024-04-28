// ----------------------------------------------------------------------------
//
// Donor analysis program
//
// Author: William Shaffer
// Version: 12-Apr-2024
//
// Copyright (c) William Shaffer
//
// ----------------------------------------------------------------------------

// This program creates a spreasheet with a timeline of the number and
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

	"acorn_go/pkg/assert"
	d "acorn_go/pkg/date"
	"acorn_go/pkg/donor_info"
	"acorn_go/pkg/spreadsheet"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

const (
	outputFileName = "/home/bozo/Downloads/donations.xlsx"
	sheetName      = "Donations"
)

// ----------------------------------------------------------------------------
// Functions
// ----------------------------------------------------------------------------

// main supervises the execution of this program.  It produces a spreadsheet
// with the donation series by year and month
func main() {
	var sprdsht spreadsheet.Spreadsheet
	var donationSeries donor_info.DonationSeries
	var err error

	printHeader()
	//
	// Obtain spreadsheet data
	//
	sprdsht, err = spreadsheet.ProcessData()
	if err != nil {
		fmt.Println("Error processing spreadsheet: " + err.Error())
		os.Exit(1)
	}
	//
	// Generate donor list
	//
	donationSeries, err = donor_info.NewDonationSeries(&sprdsht)
	if err != nil {
		fmt.Println("Error generating donor list: " + err.Error())
	}
	//
	// Output donation series to spreadsheet
	//
	err = outputDonationSeries(&donationSeries)
	if err != nil {
		fmt.Println("Error writing output: " + err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}

// printHeader places the header information at the top of the page
func printHeader() {
	fmt.Println("-----------------------------------------------------------")
	fmt.Println("Acorn Scholarship Fund Donation Series")
	fmt.Println("-----------------------------------------------------------")
}

// outputDonationSeries produces a spreadsheet with the donation time series
func outputDonationSeries(dsPtr *(donor_info.DonationSeries)) error {
	var err error = nil
	var output spreadsheet.SpreadsheetFile
	var keys []d.YearMonth
	var cell string
	var value string
	var totalDonationsFY2023 float64 = 0.0
	var totalDonationsFY2024 float64 = 0.0
	var row int
	var yearMonth d.YearMonth

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
	for row, yearMonth = range keys {
		//
		// Insert month and year
		//
		cell = cellName("A", row+2)
		value = yearMonth.String()
		output.SetCell(cell, value)
		//
		// Insert donation for that month and year
		//
		cell = cellName("B", row+2)
		amount := math.Round((*dsPtr).GetAmount(yearMonth))
		output.SetCellFloat(cell, amount)
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

// cellName generates a string representing a cell in the spreadsheet.
func cellName(column string, row int) string {
	var cellName = column + strconv.Itoa(row)
	return cellName
}

// caldTotalDonations determines the fiscal year of the donation from the year month and
// adds the amount to the appropriate total donation.
func calcTotalDonations(
	yearMonth d.YearMonth,
	amount float64,
	totDonFY2023 float64,
	totDonFY2024 float64) (float64, float64) {
	var err error = nil
	var indicator donor_info.FYIndicator

	indicator, err = donor_info.FiscalYearFromYearMonth(yearMonth)
	assert.Assert(err == nil, "Unable to handle year month: "+yearMonth.String())
	switch indicator {
	case donor_info.FY2023:
		totDonFY2023 += amount
	case donor_info.FY2024:
		totDonFY2024 += amount
	default:
		assert.Assert(false, "unmatched fiscal year: "+yearMonth.String())
	}
	return totDonFY2023, totDonFY2024
}

// outputTotals places the donation totals for each fiscal year in rows below the
// time series
func outputTotals(output spreadsheet.SpreadsheetFile, title string, totalDonation float64, row int) {
	var cell string

	cell = cellName("A", row)
	output.SetCell(cell, title)
	cell = cellName("B", row)
	output.SetCellFloat(cell, totalDonation)
}
