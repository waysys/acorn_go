// ----------------------------------------------------------------------------
//
// Convert Name
//
// Author: William Shaffer
// Version: 15-May-2025
//
// Copyright (c) 2025 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

package quickbooks

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	"strings"
)

// ----------------------------------------------------------------------------
// Functions
// ----------------------------------------------------------------------------

// ConvertName changes a name in the form of "first-name last-name" to
// the form "last-name, first_name"
func ConvertName(name string) string {
	var conv = ""
	var mult = strings.Split(name, " ")
	var length = len(mult)

	switch length {
	case 0:
		conv = "Error!"
	case 1:
		conv = mult[0]
	default:
		conv = mult[length-1] + ", " + mult[0]
	}
	return conv
}

// RevertName changes a name in the form of "last-name, first-name" to
// "first-name last-name"
func RevertName(name string) string {
	var revert = ""
	var mult = strings.Split(name, ",")
	var length = len(mult)

	switch length {
	case 0:
		revert = "Error!"
	case 1:
		revert = mult[0]
	default:
		revert = mult[length-1] + " " + mult[0]
	}
	return revert
}
