// ----------------------------------------------------------------------------
//
// Non-Repeat Data
//
// Author: William Shaffer
// Version: 07-May-2024
//
// Copyright (c) 2024 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

package donorinfo

// This file builds a slice of donation data.

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	"acorn_go/pkg/assert"
	d "acorn_go/pkg/date"
	"acorn_go/pkg/spreadsheet"
	"sort"
	"strings"

	dec "github.com/shopspring/decimal"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type namelist map[string]bool

// ----------------------------------------------------------------------------
// Factory Functions
// ----------------------------------------------------------------------------

func ComputeNonRepeatDonors(donorListPtr *DonorList, sprdsht *spreadsheet.Spreadsheet) []NonRepeat {
	var nonRepeats []NonRepeat
	var numRows = sprdsht.Size()
	var err error
	var value string
	var dateDonation d.Date
	var amountDonation dec.Decimal
	//
	// Generate a map of names of non-repeat donors
	//
	var donorNames = buildNonRepeatNames(donorListPtr)
	//
	// Generate a slice of non-repeat donations by looping through
	// the spreadsheet rows.
	//
	for row := 1; row < numRows; row++ {
		//
		// Obtain donor name
		//
		var donorName = getValue(sprdsht, row, columnNameDonor)
		//
		// If donor is a non-repeat donor, add the information
		//
		if donorNames[donorName] {
			//
			// get date of donation
			//
			value = getValue(sprdsht, row, columnDate)
			value = strings.TrimSpace(value)
			dateDonation, err = d.NewFromString(value)
			assert.Assert(err == nil, "Error converting string to date - "+value)
			//
			// get amount of donations
			//
			value = getValue(sprdsht, row, columnPayment)
			value = strings.ReplaceAll(value, ",", "")
			value = strings.TrimSpace(value)
			if value == "" {
				amountDonation = dec.Zero
			} else {
				amountDonation, err = dec.NewFromString(value)
				assert.Assert(err == nil, "Error converting amount to decimal number: ")
			}
			//
			// Create nonrepeat structure
			//
			if amountDonation.GreaterThan(ZERO) {
				var nonrepeat = NewNonRepeat(donorName, dateDonation, amountDonation)
				nonRepeats = append(nonRepeats, nonrepeat)
			}
		}

	}
	sortResult(&nonRepeats)
	return nonRepeats
}

// ----------------------------------------------------------------------------
// Support Functions
// ----------------------------------------------------------------------------

// buildNonRepeatNames generates a map with the names of non-repeat
// donors as the key
func buildNonRepeatNames(donorListPtr *DonorList) namelist {
	var names = make(namelist)
	for name, donorPtr := range *donorListPtr {
		if donorPtr.IsFY23DonorOnly() {
			names[name] = true
		}
	}
	return names
}

// getValue returns a value from the spreadsheet
func getValue(sprdsht *spreadsheet.Spreadsheet, row int, heading string) string {
	var value string
	var err error

	value, err = sprdsht.Cell(row, heading)
	assert.Assert(err == nil, "Error accessing spreadsheet cell: "+heading+": ")
	return value
}

// sortResult sorts the non-repeat donors
func sortResult(nonRepeatPtr *[]NonRepeat) {
	var less = func(i, j int) bool {
		return (*nonRepeatPtr)[i].nameDonor < (*nonRepeatPtr)[j].nameDonor
	}
	sort.Slice(*nonRepeatPtr, less)
}
