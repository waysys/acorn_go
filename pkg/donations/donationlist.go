// ----------------------------------------------------------------------------
//
// Donor List
//
// Author: William Shaffer
// Version: 14-Apr-2024
//
// Copyright (c) 2024 William Shaffer All Rights Reserved
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
	"errors"
	"sort"
	"strings"

	d "github.com/waysys/waydate/pkg/date"

	"github.com/waysys/assert/assert"

	dec "github.com/shopspring/decimal"
)

// This file contains functions that manage a map of donors with their donation
// information.

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
// Types
// ----------------------------------------------------------------------------

// DonationList maps the donors name to the donor information structure
type DonationList map[string]*dn.Donor

// ----------------------------------------------------------------------------
// Factory Methods
// ----------------------------------------------------------------------------

// NewDonorationList creates a donor list from the information in a spreadsheet.
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
			err = processPayment(&donationList, sprdsht, row)
		}
		if err != nil {
			return donationList, err
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
func processPayment(donationListPtr *DonationList, sprdsht *spreadsheet.Spreadsheet, row int) error {
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
	if !donationListPtr.Contains(nameDonor) {
		var donor = dn.NewDonorWithDonation(nameDonor)
		(*donationListPtr)[nameDonor] = &donor
	}
	//
	// Update the donor information
	//
	err = donationListPtr.AddDonation(nameDonor, amountDonation, dateDonation)
	return err
}

// ----------------------------------------------------------------------------
// Methods
// ----------------------------------------------------------------------------

// AddDonation adds a donation to the list.  If the donor is not already
// in the list, a new donation structure is created.  If the donor is in the list,
// the donation is added to the donors values, based on the donation date.
func (donationListPtr *DonationList) AddDonation(nameDonor string, amountDonation dec.Decimal, dateDonation d.Date) error {
	var err error = nil
	var donorPtr *dn.Donor
	//
	// Validate inputs
	//
	if nameDonor == "" {
		err = errors.New("name of donor must not be empty")
		return err
	}
	if amountDonation.LessThan(dec.Zero) {
		err = errors.New("amount of donation must not be negative: " + amountDonation.String())
		return err
	}
	//
	// Retrieve donor structure for the named donor
	//
	if !donationListPtr.Contains(nameDonor) {
		err = errors.New("donor name not found in donation list: " + nameDonor)
		return err
	}
	donorPtr = donationListPtr.Get(nameDonor)
	//
	// Update donation amounts based on donation date.
	//
	var fy = a.FiscalYearIndicator(dateDonation)
	donorPtr.AddDonation(amountDonation, fy)
	return err
}

// hasDonor returns true if the donor's name is in the list.  Otherwise, it
// returns false.
func (donationListPtr *DonationList) Contains(nameDonor string) bool {
	var _, found = (*donationListPtr)[nameDonor]
	return found
}

// Get returns a pointer to the donor structure for the named donor.
func (donationListPtr *DonationList) Get(nameDonor string) *dn.Donor {
	var donorPtr, found = (*donationListPtr)[nameDonor]
	assert.Assert(found, "donor name not found in donation list: "+nameDonor)
	return donorPtr
}

// DonorCount returns the number of donors in the list.
func (donationListPtr *DonationList) DonorCount() int {
	return len(*donationListPtr)
}

// DonorKeys returns a alphabetically sorted slice of donor list keys
func (donationListPtr *DonationList) DonorKeys() []string {
	keys := make([]string, 0, len(*donationListPtr))
	for k := range *donationListPtr {
		keys = append(keys, k)
	}

	//sort the slice of keys alphabetically
	sort.Strings(keys)
	return keys
}
