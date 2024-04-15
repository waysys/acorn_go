// ----------------------------------------------------------------------------
//
// Donor List
//
// Author: William Shaffer
// Version: 14-Apr-2024
//
// Copyright (c) William Shaffer
//
// ----------------------------------------------------------------------------

package donor_info

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	d "acorn_go/pkg/date"
	"errors"

	dec "github.com/shopspring/decimal"
)

// This file contains functions that manage a map of donors with their donation
// information.

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// DonorList maps the donors name to the donor information structure
type DonorList map[string]*Donor

// ----------------------------------------------------------------------------
// Operations
// ----------------------------------------------------------------------------

// AddDonation adds a donation to the list.  If the donor is not already
// in the list, a new donation structure is created.  If the donor is in the list,
// the donation is added to the donors values, based on the donation date.
func (donorList *DonorList) AddDonation(nameDonor string, amountDonation dec.Decimal, dateDonation d.Date) error {
	var err error = nil
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
	return err
}
