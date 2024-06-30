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

package address

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type Address struct {
	Street string
	Unit   string
	City   string
	State  string
	Zip    string
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// Factory Functions
// ----------------------------------------------------------------------------

// New creates a new address structure based on the inputs
func New(street string, unit string, city string, state string, zip string) Address {
	address := Address{
		Street: street,
		Unit:   unit,
		City:   city,
		State:  state,
		Zip:    zip,
	}
	return address
}

// Blank address returns a blank address
func BlankAddress() Address {
	address := Address{}
	return address
}
