// ----------------------------------------------------------------------------
//
// Donor List
//
// Author: William Shaffer
//
// Copyright (c) 2024, 2025 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

package donations

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	a "acorn_go/pkg/accounting"
	dn "acorn_go/pkg/donors"
	"acorn_go/pkg/spreadsheet"
	"fmt"
	"sort"
	"strings"

	d "github.com/waysys/waydate/pkg/date"

	"github.com/waysys/assert/assert"

	dec "github.com/shopspring/decimal"
)

// This file contains functions that manage a map of donors with their donation
// information.

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// DonationList maps the donors name to the donor information structure
type DonationList map[string]*dn.Donor

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

const (
	columnNameDonor       = "Payee"
	columnTransactionType = "Type"
	columnDate            = "Date"
	columnPayment         = "Payment"
)

const (
	payment       = "Payment"
	excludedDonor = "Nadine L. Tolman Trust"
)

// ----------------------------------------------------------------------------
// Factory Methods
// ----------------------------------------------------------------------------

// NewDonationList creates a donor list from the information in a spreadsheet.
func NewDonationList(sprdsht *spreadsheet.Spreadsheet) (DonationList, error) {
	var donationList = make(DonationList)
	var numRows = sprdsht.Size()
	var err error
	var okToSelect bool
	//
	// Loop through all the rows in the spreadsheet
	//
	for row := 1; row < numRows; row++ {
		okToSelect, err = selectRow(sprdsht, row)
		if err != nil {
			return donationList, err
		}
		if okToSelect {
			err = processPayment(donationList, sprdsht, row)
		}
	}
	return donationList, err
}

// selectRow decides whether a row should be included in the analysis.
func selectRow(sprdsht *spreadsheet.Spreadsheet, row int) (bool, error) {
	var err error = nil
	var tranType string
	var nameDonor string
	var result bool
	//
	// Fetch transaction type
	//
	tranType, err = sprdsht.Cell(row, columnTransactionType)
	if err != nil {
		return false, err
	}
	tranType = strings.TrimSpace(tranType)
	//
	// Fetch donor name
	//
	nameDonor, err = sprdsht.Cell(row, columnNameDonor)
	if err != nil {
		return false, err
	}
	nameDonor = strings.TrimSpace(nameDonor)
	//
	// Include row if value == Payment and nameDonor != Nadine L. Tolman Trust
	//
	if (tranType == payment) && (nameDonor != excludedDonor) {
		result = true
	} else {
		result = false
	}
	return result, err
}

// processPayment creates a Donor structure
func processPayment(donationList DonationList, sprdsht *spreadsheet.Spreadsheet, row int) error {
	var value string
	var err error = nil
	var amountDonation dec.Decimal
	var dateDonation d.Date
	//
	// Obtain donor name
	//
	value, err = sprdsht.Cell(row, columnNameDonor)
	if err != nil {
		return err
	}
	var nameDonor = value
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
	//
	// Obtain amount of donation
	//
	amountDonation, err = sprdsht.CellDecimal(row, columnPayment)
	if err != nil {
		return err
	}
	//
	// Create an entry in the donor list if there is not already one
	// for this donor.
	//
	if !donationList.Contains(nameDonor) {
		var donor = dn.NewDonorWithDonation(nameDonor)
		donationList[nameDonor] = &donor
	}
	//
	// Update the donor information
	//
	err = donationList.AddDonation(nameDonor, amountDonation, dateDonation)

	return err
}

// ----------------------------------------------------------------------------
// Methods
// ----------------------------------------------------------------------------

// AddDonation adds a donation to the list.  If the donor is not already
// in the list, a non-nil error is returned.  If the donor is in the list,
// the donation is added to the donors values, based on the donation date.
func (donationList DonationList) AddDonation(
	nameDonor string,
	amountDonation dec.Decimal,
	dateDonation d.Date) error {
	var err error = nil
	var donorPtr *dn.Donor
	//
	// Validate inputs
	//
	if nameDonor == "" {
		return fmt.Errorf("name of donor must not be empty")
	}
	if amountDonation.LessThan(dec.Zero) {
		return fmt.Errorf("amount of donation must not be negative: %s", amountDonation.String())
	}
	//
	// Retrieve donor structure for the named donor
	//
	if !donationList.Contains(nameDonor) {
		return fmt.Errorf("donor name not found in donation list: %s", nameDonor)
	}
	donorPtr = donationList.Get(nameDonor)
	//
	// Update donation amounts based on donation date.
	//
	var fy = a.FiscalYearIndicator(dateDonation)
	donorPtr.AddDonation(amountDonation, fy)
	var year = a.YIndicator(dateDonation)
	donorPtr.AddCalDonation(amountDonation, year)
	donorPtr.AddMonthDonation(amountDonation, dateDonation)
	return err
}

// hasDonor returns true if the donor's name is in the list.  Otherwise, it
// returns false.
func (donationList DonationList) Contains(nameDonor string) bool {
	_, found := donationList[nameDonor]
	return found
}

// Get returns a pointer to the donor structure for the named donor.
// If the donor is not found, an assertion failure occurs.
func (donationList DonationList) Get(nameDonor string) *dn.Donor {
	donorPtr, found := donationList[nameDonor]
	assert.Assert(found, "donor name not found in donation list: "+nameDonor)
	return donorPtr
}

// DonorCount returns the number of donors in the list.
func (donationList DonationList) DonorCount() int {
	return len(donationList)
}

// DonorKeys returns a alphabetically sorted slice of donor list keys
func (donationList DonationList) DonorKeys() []string {
	keys := make([]string, 0, len(donationList))
	for k := range donationList {
		keys = append(keys, k)
	}

	//sort the slice of keys alphabetically
	sort.Strings(keys)
	return keys
}
