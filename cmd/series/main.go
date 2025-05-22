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

	a "acorn_go/pkg/accounting"
	ds "acorn_go/pkg/donationseries"
	s "acorn_go/pkg/spreadsheet"
	sp "acorn_go/pkg/support"

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
	inputFile      = "/home/bozo/golang/acorn_go/data/donations.xlsx"
	tab            = "Worksheet"
	outputFileName = "/home/bozo/Downloads/donations_series.xlsx"
	sheetName      = "Donations"
)

const (
	beginYear  = 2022
	beginMonth = 9
	endYear    = 2025
	endMonth   = 8
)

// ----------------------------------------------------------------------------
// Functions
// ----------------------------------------------------------------------------

// main supervises the execution of this program.  It produces a spreadsheet
// with the donation series by year and month
func main() {
	var sprdsht s.Spreadsheet
	var donationSeries ds.DonationSeries
	var err error

	printHeader()
	//
	// Obtain spreadsheet data
	//
	sprdsht, err = s.ProcessData(inputFile, tab)
	sp.Check(err, "Error processing spreadsheet")
	//
	// Generate donation series
	//
	donationSeries, err = ds.NewDonationSeries(&sprdsht)
	sp.Check(err, "Error generating donation series")
	//
	// Output donation series to spreadsheet
	//
	err = outputDonationSeries(&donationSeries)
	sp.Check(err, "Error writing output")
	os.Exit(0)
}

// printHeader places the header information at the top of the page
func printHeader() {
	fmt.Println("-----------------------------------------------------------")
	fmt.Println("Acorn Scholarship Fund Donation Series")
	fmt.Println("-----------------------------------------------------------")
}

// outputDonationSeries produces a spreadsheet with the donation time series
func outputDonationSeries(dsPtr *(ds.DonationSeries)) error {
	var err error = nil
	var output s.SpreadsheetFile
	var keys []d.YearMonth
	var totalDonations = []float64{0.0, 0.0, 0.0}
	var row int
	var fiscalYear a.FYIndicator
	var yearMonth d.YearMonth
	var avg float64 = 0.0
	var start, _ = d.NewYearMonth(beginYear, beginMonth)
	var thru, _ = d.NewYearMonth(endYear, endMonth)

	if dsPtr == nil {
		err = errors.New("pointer to donation series is nil")
		return err
	}

	output, err = s.New(outputFileName, sheetName)
	if err != nil {
		return err
	}
	defer output.Close()
	//
	// Enter values in into spreadsheet
	//
	keys, err = d.Keys(start, thru)
	if err != nil {
		return err
	}
	//
	// Insert title
	//
	s.WriteCell(&output, "A", 1, "Year")
	s.WriteCell(&output, "B", 1, "Month")
	s.WriteCell(&output, "C", 1, "Donation Amount")
	s.WriteCell(&output, "D", 1, "Number of Donations")
	s.WriteCell(&output, "E", 1, "Average Amount of Donation")
	s.WriteCell(&output, "F", 1, "Fiscal Year")
	s.WriteCell(&output, "G", 1, "Total Donations for Fiscal Year")

	for row, yearMonth = range keys {
		//
		// Insert month and year
		//
		var yms = yearMonth.MonthString()
		s.WriteCellInt(&output, "A", row+2, yms.Year)
		s.WriteCell(&output, "B", row+2, yms.Month)
		//
		// Insert donation for that month and year
		//
		amount := math.Round((*dsPtr).GetAmount(yearMonth))
		s.WriteCellFloat(&output, "C", row+2, amount)
		//
		// Insert the number of donors for that month and year
		//
		count := (*dsPtr).GetCount(yearMonth)
		s.WriteCellInt(&output, "D", row+2, count)
		//
		// Insert the average donation
		//
		if count == 0 {
			avg = 0.00
		} else {
			avg = math.Round(amount / float64(count))
		}
		s.WriteCellFloat(&output, "E", row+2, avg)
		//
		// Calculate total donations for FY2023 through FY2025
		//
		calcTotalDonations(yearMonth, amount, &totalDonations)
		fiscalYear, err = a.FiscalYearFromYearMonth(yearMonth)
		s.WriteCell(&output, "F", row+2, fiscalYear.String())
		s.WriteCellFloat(&output, "G", row+2, totalDonations[fiscalYear])
	}
	//
	// Ootput totals
	//
	row += 4
	outputTotals(output, "FY2023 Donations", totalDonations[a.FY2023], row)
	outputTotals(output, "FY2024 Donations", totalDonations[a.FY2024], row+1)
	outputTotals(output, "FY2025 Donations", totalDonations[a.FY2025], row+2)
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
	totalDonations *[]float64) {
	var err error = nil
	var indicator a.FYIndicator

	indicator, err = a.FiscalYearFromYearMonth(yearMonth)
	assert.Assert(err == nil, "Unable to handle year month: "+yearMonth.String())
	assert.Assert(indicator != a.OutOfRange, "Fiscal year indicator is out of range")
	(*totalDonations)[indicator] += amount
}

// outputTotals places the donation totals for each fiscal year in rows below the
// time series
func outputTotals(output s.SpreadsheetFile, title string, totalDonation float64, row int) {
	s.WriteCell(&output, "A", row, title)
	s.WriteCellFloat(&output, "B", row, totalDonation)
}
