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
