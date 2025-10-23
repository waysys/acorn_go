// ----------------------------------------------------------------------------
//
// Column Handler
//
// This code handles the writing of cells to a spreadsheet based on the
// fiscal years available.  It is used when the output spreadsheet has columns
// of data for each fiscal year.
//
// Author: William Shaffer
//
// Copyright (c) 2025 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

package spreadsheet

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	a "acorn_go/pkg/accounting"

	utf "unicode/utf8"

	"github.com/waysys/assert/assert"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type ColumnHandler struct {
	countColumns   []string
	paymentColumns []string
	totalColumn    string
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// Factory Functions
// ----------------------------------------------------------------------------

// NewColumnHandler returns an instance of the column handler structure
// initialized.
func NewColumnHandler(startColumn string) ColumnHandler {
	var countColumns []string
	var paymentColumns []string
	var letterCount = startColumn
	var letterPayment string

	for index := 0; index < a.NumFiscalYears; index++ {
		countColumns = append(countColumns, letterCount)
		letterPayment = NextLetter(letterCount)
		paymentColumns = append(paymentColumns, letterPayment)
		letterCount = NextSecondLetter(letterCount)
	}
	var totalColumn = letterCount
	var handler = ColumnHandler{
		countColumns:   countColumns,
		paymentColumns: paymentColumns,
		totalColumn:    totalColumn,
	}
	return handler
}

// ----------------------------------------------------------------------------
// Methods
// ----------------------------------------------------------------------------

// Return the column for the recipient count
func (ch *ColumnHandler) CountColumn(fy a.FYIndicator) string {
	assert.Assert(a.IsFYIndicator(fy), "Invalid fiscal year indicator: "+fy.String())
	return ch.countColumns[int(fy)]
}

// Return the column for the recipient scholarship
func (ch *ColumnHandler) PaymentColumn(fy a.FYIndicator) string {
	assert.Assert(a.IsFYIndicator(fy), "Invalid fiscal year indicator: "+fy.String())
	return ch.paymentColumns[int(fy)]
}

// Return the column for the total scholarships
func (ch *ColumnHandler) TotalColumn() string {
	return ch.totalColumn
}

// Return the column name for the count given the fiscal year
func (ch *ColumnHandler) CountColumnLabel(fy a.FYIndicator) string {
	var fyString = fy.String()
	return fyString + " Count"
}

// Return the column name for the payment given the fiscal year
func (ch *ColumnHandler) PaymentColumnLabel(fy a.FYIndicator) string {
	var fyString = fy.String()
	return fyString + " Payment"
}

// ----------------------------------------------------------------------------
// Support Functions
// ----------------------------------------------------------------------------

// NextLetter returns the next ASCII letter in the alphabet.
// - 'a'..'y' -> next lowercase letter
// - 'z' -> 'a'
// - 'A'..'Y' -> next uppercase letter
// - 'Z' -> 'A'
// For any rune that is not an ASCII letter the input is returned unchanged.
func NextLetter(str string) string {
	assert.Assert(isValidColumn(str),
		"str is not a single character: "+str)
	var r = convertToRune(str)
	var result rune
	switch {
	case r >= 'a' && r <= 'y':
		result = r + 1
	case r == 'z':
		result = 'a'
	case r >= 'A' && r <= 'Y':
		result = r + 1
	case r == 'Z':
		result = 'A'
	default:
		assert.Assert(false, "str does not contain an ASCII letter: "+str)
	}
	return string(result)
}

// convertToRune converts a string of one character to the corresponding rune.
func convertToRune(str string) rune {
	r := []rune(str)[0]
	return r
}

// NextSecondLetter accepts a string with a single letter and
// returns a string with the letter that is 2 positions in alphabetical order
func NextSecondLetter(str string) string {
	var result = NextLetter(str)
	result = NextLetter(result)
	return result
}

// isValidColumn returns true if a string is a valid column header: it is a string
// with one character.
func isValidColumn(str string) bool {
	return utf.RuneCountInString(str) == 1
}
