// ----------------------------------------------------------------------------
//
// Guestlist
//
// Author: William Shaffer
// Version: 29-Jun-2024
//
// Copyright (c) 2024 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

package guestlist

// The Guestlist object is a map of guests with a key based on the email
// of the guest.  Guestlist reads the guestlist.xslx spreadsheet to build
// the map.

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	"acorn_go/pkg/address"
	"acorn_go/pkg/donors"
	"acorn_go/pkg/spreadsheet"

	"github.com/waysys/assert/assert"
)

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

const (
	columnName   = "Full Name"
	columnEmail  = "Email/Phone Number"
	columnStatus = "Status"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type Guestlist map[string]*Guest

// ----------------------------------------------------------------------------
// Factory Functions
// ----------------------------------------------------------------------------

// NewGuestList returns a guestlist populated from the guestlist.xslx spreadsheet.
func NewGuestlist(sprdsht *spreadsheet.Spreadsheet, donorList *donors.DonorList) (Guestlist, error) {
	var guestlist = make(Guestlist)
	var numRows = sprdsht.Size()
	var err error = nil

	//
	// Loop through all the rows in the spreadsheet
	//
	for row := 1; row < numRows; row++ {
		var guest, err1 = processRow(donorList, sprdsht, row)
		if err1 != nil {
			err = err1
			break
		}
		if guest.Email() != "" {
			guestlist.Add(&guest)
		} else {
			continue
		}

	}
	return guestlist, err
}

// processRow extracts the data from one row of the spreadsheets and builds
// a guest object to place in the guestlist
func processRow(
	donorList *donors.DonorList,
	sprdsht *spreadsheet.Spreadsheet,
	row int) (Guest, error) {
	var err error = nil
	var donor *donors.Donor = nil
	var name string
	var email string
	var status string
	var guest Guest
	var blankAddress = address.BlankAddress()
	//
	// Extract data from spreadsheet
	//
	name, err = sprdsht.Cell(row, columnName)
	if err == nil {
		email, err = sprdsht.Cell(row, columnEmail)
	}
	if err == nil {
		status, err = sprdsht.Cell(row, columnStatus)
	}
	if err == nil {
		donor, err = donorList.GetByEmail(email)
		if err != nil {
			var dn = donors.New(name, name, blankAddress, email, 0)
			donor = &dn
			err = nil
		}
	}
	if err == nil {
		guest = New(name, donor.Address(), email, status)
	}
	return guest, err
}

// ----------------------------------------------------------------------------
// Methods
// ----------------------------------------------------------------------------

// Contains returns true if the guestlist contains a guest with the
// specified email.
func (guestlist *Guestlist) Contains(email string) bool {
	var _, ok = (*guestlist)[email]
	return ok
}

// Add adds a guest to the guestlist
func (guestlist *Guestlist) Add(guest *Guest) {
	var email = guest.Email()
	assert.Assert(!guestlist.Contains(email),
		"Guest with this email is already in list: "+email)
	(*guestlist)[email] = guest
}

// Get returns the guest with the specified email.
func (guestlist *Guestlist) Get(email string) *Guest {
	assert.Assert(guestlist.Contains(email),
		"Guest with this email is not in list: "+email)
	var guest, _ = (*guestlist)[email]
	return guest
}

// Count returns the number of guests in the guestlist.
func (guestlist *Guestlist) Count() int {
	return len(*guestlist)
}

// Keys returns a slice containing the keys.  The order
// of the keys is indeterminate.
func (guestlist *Guestlist) Keys() []string {
	keys := make([]string, 0, len(*guestlist))
	for key := range *guestlist {
		keys = append(keys, key)
	}
	return keys
}
