// ----------------------------------------------------------------------------
//
// Donor series program
//
// Author: William Shaffer
// Version: 20-May-2024
//
// Copyright (c) 2024 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

// The donors package manages the names, addresses, and email addresses of
// donors
package donors

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	a "acorn_go/pkg/address"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type Donor struct {
	key             string
	name            string
	address         a.Address
	email           string
	numberHousehold int
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// Factory Functions
// ----------------------------------------------------------------------------

// New creates a new donor
func New(ky string, nm string, adr a.Address, eml string, count int) Donor {
	donor := Donor{
		key:             ky,
		name:            nm,
		address:         adr,
		email:           eml,
		numberHousehold: count,
	}
	return donor
}

// ----------------------------------------------------------------------------
// Properties
// ----------------------------------------------------------------------------

// Key returns the abbreviated donor name used to identify the donor.
func (donor Donor) Key() string {
	return donor.key
}

// Name returns the full names of the donors
func (donor Donor) Name() string {
	return donor.name
}

// Street returns the street address of the donor
func (donor Donor) Street() string {
	return donor.address.Street
}

// City returns the city of the donor
func (donor Donor) City() string {
	return donor.address.City
}

// State returns the state of the donor
func (donor Donor) State() string {
	return donor.address.State
}

// Zip returns the zip code of the donor
func (donor Donor) Zip() string {
	return donor.address.Zip
}

// Email returns the email address of the donor
func (donor Donor) Email() string {
	return donor.email
}

// NumberInHousehold return the number of people in the
// donor household
func (donor Donor) NumberInHousehold() int {
	return donor.numberHousehold
}
