// ----------------------------------------------------------------------------
//
// Fiscal Year
//
// Author: William Shaffer
// Version: 25-Apr-2024
//
// Copyright (c) William Shaffer
//
// ----------------------------------------------------------------------------

package donor_info

// This file composes a time series of total donations by year and month.

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	d "acorn_go/pkg/date"
	"acorn_go/pkg/spreadsheet"
	"strings"

	dec "github.com/shopspring/decimal"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type DonationSeries map[d.YearMonth]float64

// ----------------------------------------------------------------------------
// Factory Functions
// ----------------------------------------------------------------------------

// NewDonationSeries creates a donation series from the spreadsheet
func NewDonationSeries(sprdsht *spreadsheet.Spreadsheet) (DonationSeries, error) {
	var donationSeries = make(DonationSeries)
	var err error = nil
	var numRows = sprdsht.Size()
	var value string
	//
	// Loop through all rows in the spreadsheet except row 0 which has
	// the column headings
	//
	for row := 1; row < numRows; row++ {
		value, err = sprdsht.Cell(row, columnTransactionType)
		if err != nil {
			return donationSeries, err
		}
		value = strings.TrimSpace(value)
		if value == payment {
			err = processSeries(&donationSeries, sprdsht, row)
			if err != nil {
				return donationSeries, err
			}
		}
	}
	return donationSeries, err
}

// processSeries processes row in the spreadsheet and adds the data
// to the donation series.
func processSeries(dsPtr *DonationSeries, sprdsht *spreadsheet.Spreadsheet, row int) error {
	var value string
	var err error = nil
	var yearMonth d.YearMonth
	var amountDonation dec.Decimal
	var dateDonation d.Date
	//
	// Obtain date of donation
	//
	value, err = sprdsht.Cell(row, columnDate)
	if err != nil {
		return err
	}
	dateDonation, err = d.NewFromString(value)
	if err != nil {
		return err
	}
	yearMonth, err = d.NewYearMonthFromDate(dateDonation)
	if err != nil {
		return err
	}
	//
	// Obtain amount of donation
	//
	value, err = sprdsht.Cell(row, columnPayment)
	if err != nil {
		return err
	}
	value = strings.ReplaceAll(value, ",", "")
	value = strings.TrimSpace(value)
	if value == "" {
		amountDonation = dec.Zero
	} else {
		amountDonation, err = dec.NewFromString(value)
	}
	if err != nil {
		return err
	}
	//
	// Create an entry in the donation series for this year and month
	// if one does not already exist.
	//
	if !dsPtr.hasYearMonth(yearMonth) {
		(*dsPtr)[yearMonth] = 0.0
	}
	//
	// Add amount to donation series
	//
	(*dsPtr)[yearMonth] += amountDonation.InexactFloat64()
	return err
}

// ----------------------------------------------------------------------------
// Methods
// ----------------------------------------------------------------------------

// hasYearMonth returns true if the donation series has an entry for the month
// and year.
func (dsPtr *DonationSeries) hasYearMonth(yearMonth d.YearMonth) bool {
	var _, found = (*dsPtr)[yearMonth]
	return found
}

// GetAmount returns the amount of donations for a year and month rounded to
// the nearest whole dollar amount.  If the donation series does not contain
// and entry for the year and month, a zero amount is returned.
func (dsPtr *DonationSeries) GetAmount(yearMonth d.YearMonth) float64 {
	var amount, found = (*dsPtr)[yearMonth]
	if !found {
		amount = 0.00
	}
	return amount
}
