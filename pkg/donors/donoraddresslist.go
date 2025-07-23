// ----------------------------------------------------------------------------
//
// Donor address program
//
// Author: William Shaffer
// Version: 20-May-2024
//
// Copyright (c) 2024 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

package donors

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	ac "acorn_go/pkg/accounting"
	a "acorn_go/pkg/address"
	"acorn_go/pkg/spreadsheet"
	"errors"
	"sort"
	"strconv"

	dec "github.com/shopspring/decimal"
	"github.com/waysys/assert/assert"
	d "github.com/waysys/waydate/pkg/date"
)

// This file contains functions that manage a map of donors with their names,
// addresses, and emails.  Deceased donors are excluded.

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// Donor column names
const (
	columnKey             = "Donor full name"
	columnEmail           = "Email"
	columnStreet          = "Bill street"
	columnCity            = "Bill city"
	columnState           = "Bill state"
	columnZip             = "Bill zip"
	columnName            = "Invite Name"
	columnNumberHousehold = "NumberHousehold"
	columnDeceased        = "Deceased"
)

// Donation column names
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

// DonorList maps the donors name to the donor information structure
type DonorList map[string]*Donor

// ----------------------------------------------------------------------------
// Factory Functions
// ----------------------------------------------------------------------------

// NewDonorList creates a donor list from the information in a spreadsheet.
func NewDonorAddressList(sprdsht *spreadsheet.Spreadsheet) (DonorList, error) {
	var donorList = make(DonorList)
	var numRows = sprdsht.Size()
	var err error
	var value string
	//
	// Loop through all the rows in the spreadsheet
	//
	for row := 1; row < numRows; row++ {
		value, err = sprdsht.Cell(row, columnDeceased)
		if err != nil {
			return donorList, err
		}
		//
		// Deceased donors are not included
		//
		if value != "Yes" {
			err = processDonor(&donorList, sprdsht, row)
			if err != nil {
				return donorList, err
			}
		}
	}
	return donorList, err
}

// processDonor processes a single donor
func processDonor(donorList *DonorList, sprdsht *spreadsheet.Spreadsheet, row int) error {
	var key string
	var name string
	var email string
	var count int
	var address a.Address
	var value string
	var donor Donor
	var err error = nil
	//
	// Extract values from spreadsheet
	//
	key, err = sprdsht.Cell(row, columnKey)
	if err == nil {
		name, err = sprdsht.Cell(row, columnName)
	}
	if err == nil {
		address, err = processAddress(sprdsht, row)
	}
	if err == nil {
		email, err = sprdsht.Cell(row, columnEmail)
	}
	if err == nil {
		value, err = sprdsht.Cell(row, columnNumberHousehold)
	}
	if err == nil {
		count, err = strconv.Atoi(value)
		if err != nil {
			count = 1
			err = nil
		}
	}
	//
	// Create donor
	//
	if err == nil {
		donor = New(key, name, address, email, count)
		donorList.Add(&donor)
	}
	return err
}

// processAddress processes returns an address structure based on the
// data in the spreadsheet.
func processAddress(sprdsht *spreadsheet.Spreadsheet, row int) (a.Address, error) {
	var addr a.Address = a.Address{}
	var err error = nil
	//
	// Extract address from spreadsheet
	//
	addr.Street, err = sprdsht.Cell(row, columnStreet)
	if err == nil {
		addr.City, err = sprdsht.Cell(row, columnCity)
	}
	if err == nil {
		addr.State, err = sprdsht.Cell(row, columnState)
	}
	if err == nil {
		addr.Zip, err = sprdsht.Cell(row, columnZip)
	}
	return addr, err
}

// AddDonations adds donation information to donor list
func AddDonations(sprdsht *spreadsheet.Spreadsheet, donarList *DonorList) error {
	var numRows = sprdsht.Size()
	var err error

	for row := 1; row < numRows; row++ {
		err = processDonation(sprdsht, donarList, row)
	}
	return err
}

// processDonation processes a row in the donation spreadsheet
func processDonation(sprdsht *spreadsheet.Spreadsheet, donorList *DonorList, row int) error {
	var err error = nil
	var value = ""
	var dateDonation d.Date
	var amountDonation dec.Decimal

	//
	// Obtain donor name
	//
	value, err = sprdsht.Cell(row, columnNameDonor)
	if err != nil {
		return err
	}
	var nameDonor = value
	//
	// Obtain transaction type
	//
	value, err = sprdsht.Cell(row, columnTransactionType)
	if err != nil {
		return err
	}
	var transType = value
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
	// Add contribution to the donor for the right calendar year
	//
	if donorList.Contains(nameDonor) {
		var donor = donorList.Get(nameDonor)
		if selectDonation(transType, nameDonor) {
			donor.AddCalDonation(amountDonation, ac.YIndicator(dateDonation))
		}
	}
	return err
}

// selectDonation returns true if the donation should be added to the
// donor.
func selectDonation(transType string, nameDonor string) bool {
	var result = transType == payment
	result = result && (nameDonor != excludedDonor)

	return result
}

// ----------------------------------------------------------------------------
// Methods
// ----------------------------------------------------------------------------

// Contains returns true if the donor list contains a key equal to specified value.
func (donorList *DonorList) Contains(key string) bool {
	var _, ok = (*donorList)[key]
	return ok
}

// Add inserts a donor into the donor list
func (donorList *DonorList) Add(donor *Donor) {
	assert.Assert(!donorList.Contains(donor.Key()), "donor is already in donor list: "+donor.Key())
	(*donorList)[donor.Key()] = donor
}

// Get returns a pointer to a donor with the specified key.
// The donor must be in the list.
func (donorList *DonorList) Get(key string) *Donor {
	assert.Assert(donorList.Contains(key), "donor is not in donor list: "+key)
	return (*donorList)[key]
}

// Count returns the number of entries in the donor list.
func (donorList *DonorList) Count() int {
	return len(*donorList)
}

// DonorKeys returns a alphabetically sorted slice of donor list keys
func (donorList *DonorList) Keys() []string {
	keys := make([]string, 0, len(*donorList))
	for k := range *donorList {
		keys = append(keys, k)
	}

	//sort the slice of keys alphabetically
	sort.Strings(keys)
	return keys
}

// GetByEmail returns the donor with the specified email.
// If no donor has the specified email, an error is returned.
func (donorList *DonorList) GetByEmail(email string) (*Donor, error) {
	var donor *Donor = nil
	var err error = errors.New("unable to find donor with this email: " + email)
	var keys = donorList.Keys()

	for _, key := range keys {
		donor = donorList.Get(key)
		if donor.Email() == email {
			err = nil
			break
		}
	}
	return donor, err
}
