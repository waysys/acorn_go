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

package donorinfo

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	a "acorn_go/pkg/accounting"
	d "acorn_go/pkg/date"
	"acorn_go/pkg/spreadsheet"
	"errors"
	"sort"
	"strings"

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
	payment = "Payment"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// DonorList maps the donors name to the donor information structure
type DonorList map[string]*Donor

// ----------------------------------------------------------------------------
// Factory Methods
// ----------------------------------------------------------------------------

// NewDonorList creates a donor list from the information in a spreadsheet.
func NewDonorList(sprdsht *spreadsheet.Spreadsheet) (DonorList, error) {
	var donorList = make(DonorList)
	var numRows = sprdsht.Size()
	var err error
	var value string
	//
	// Loop through all the rows in the spreadsheet
	//
	for row := 1; row < numRows; row++ {
		value, err = sprdsht.Cell(row, columnTransactionType)
		if err != nil {
			return donorList, err
		}
		value = strings.TrimSpace(value)
		if value == payment {
			err = processPayment(&donorList, sprdsht, row)
			if err != nil {
				return donorList, err
			}
		}
	}
	return donorList, err
}

// processPayment creates a Donor structure
func processPayment(donorListPtr *DonorList, sprdsht *spreadsheet.Spreadsheet, row int) error {
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
	if !donorListPtr.Contains(nameDonor) {
		var donor = New(nameDonor)
		(*donorListPtr)[nameDonor] = &donor
	}
	//
	// Update the donor information
	//
	err = (*donorListPtr).AddDonation(nameDonor, amountDonation, dateDonation)
	return err
}

// ----------------------------------------------------------------------------
// Methods
// ----------------------------------------------------------------------------

// AddDonation adds a donation to the list.  If the donor is not already
// in the list, a new donation structure is created.  If the donor is in the list,
// the donation is added to the donors values, based on the donation date.
func (donorListPtr *DonorList) AddDonation(nameDonor string, amountDonation dec.Decimal, dateDonation d.Date) error {
	var err error = nil
	var donorPtr *Donor
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
	if !donorListPtr.Contains(nameDonor) {
		err = errors.New("donor name not found in donation list: " + nameDonor)
		return err
	}
	donorPtr = donorListPtr.Get(nameDonor)
	//
	// Update donation amounts based on donation date.
	//
	var fiscalYearIndicator = a.FiscalYearIndicator(dateDonation)
	switch fiscalYearIndicator {
	case a.FY2023:
		donorPtr.AddFY23(amountDonation)
	case a.FY2024:
		donorPtr.AddFY24(amountDonation)
	default:
		err = errors.New("date of donation is not in either fiscal year: " + dateDonation.String())
	}
	return err
}

// hasDonor returns true if the donor's name is in the list.  Otherwise, it
// returns false.
func (donorListPtr *DonorList) Contains(nameDonor string) bool {
	var _, found = (*donorListPtr)[nameDonor]
	return found
}

// Get returns a pointer to the donor structure for the named donor.
func (donorListPtr *DonorList) Get(nameDonor string) *Donor {
	var donorPtr, found = (*donorListPtr)[nameDonor]
	assert.Assert(found, "donor name not found in donation list: "+nameDonor)
	return donorPtr
}

// DonorCount returns the number of donors in the list.
func (donorListPtr *DonorList) DonorCount() int {
	return len(*donorListPtr)
}

// DonorKeys returns a alphabetically sorted slice of donor list keys
func (donorListPtr *DonorList) DonorKeys() []string {
	keys := make([]string, 0, len(*donorListPtr))
	for k := range *donorListPtr {
		keys = append(keys, k)
	}

	//sort the slice of keys alphabetically
	sort.Strings(keys)
	return keys
}
