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

// The donationseries creates a time series of donations by month.
package donationseries

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	"acorn_go/pkg/spreadsheet"
	"strings"

	d "github.com/waysys/waydate/pkg/date"

	dec "github.com/shopspring/decimal"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type DonationSeries map[d.YearMonth]*DonationInfo

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

const (
	columnNameDonor       = "Payee"
	columnTransactionType = "Type"
	columnDate            = "Date"
	columnPayment         = "Payment"
	columnPayee           = "Payee"
)

const (
	unusualDonor = "Nadine L. Tolman Trust"
)

// ----------------------------------------------------------------------------
// Factory Functions
// ----------------------------------------------------------------------------

// NewDonationSeries creates a donation series from the donation spreadsheet
func NewDonationSeries(sprdsht *spreadsheet.Spreadsheet) (DonationSeries, error) {
	var donationSeries = make(DonationSeries)
	var err error = nil
	var numRows = sprdsht.Size()
	var value string
	var donor string
	//
	// Loop through all rows in the spreadsheet except row 0 which has
	// the column headings
	//
	for row := 1; row < numRows; row++ {
		value, err = sprdsht.Cell(row, columnTransactionType)
		if err != nil {
			return donationSeries, err
		}

		donor, err = sprdsht.Cell(row, columnPayee)
		if err != nil {
			return donationSeries, err
		}

		if selectRow(value, donor) {
			err = processSeries(&donationSeries, sprdsht, row)
			if err != nil {
				return donationSeries, err
			}
		}
	}
	return donationSeries, err
}

// selectRow determines if a row will be processed
func selectRow(value string, donor string) bool {
	var result = value == columnPayment
	result = result && donor != unusualDonor
	return result
}

// processSeries processes row in the spreadsheet and adds the data
// to the donation series.
func processSeries(dsPtr *DonationSeries, sprdsht *spreadsheet.Spreadsheet, row int) error {
	var value string
	var err error = nil
	var yearMonth d.YearMonth
	var amountDonation dec.Decimal
	var dateDonation d.Date
	var donationInfo DonationInfo
	var donationInfoPtr *DonationInfo
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
		donationInfo = NewDonationInfo()
		(*dsPtr)[yearMonth] = &donationInfo
	}
	//
	// Add amount to donation series
	//
	donationInfoPtr = (*dsPtr)[yearMonth]
	(*donationInfoPtr).AddAmount(amountDonation.InexactFloat64())
	(*donationInfoPtr).AddCount(1)

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
// an entry for the year and month, a zero amount is returned.
func (dsPtr *DonationSeries) GetAmount(yearMonth d.YearMonth) float64 {
	var amount float64 = 0.0
	var diPtr, found = (*dsPtr)[yearMonth]
	if !found {
		amount = 0.00
	} else {
		amount = (*diPtr).Amount()
	}
	return amount
}

// GetCount returns the number of donations for the year and month.  If the
// donation series does not contain an entry for the year and month, zero
// is returned.
func (dsPtr *DonationSeries) GetCount(yearMonth d.YearMonth) int {
	var count = 0
	var diPtr, found = (*dsPtr)[yearMonth]
	if !found {
		count = 0
	} else {
		count = (*diPtr).Count()
	}
	return count
}
