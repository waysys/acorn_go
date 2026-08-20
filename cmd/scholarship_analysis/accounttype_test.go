// ----------------------------------------------------------------------------
//
// # accounttype_test
//
// Tests for DetermineAccountType.
//
// Author: William Shaffer
//
// Copyright (c) 2026 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

package main

import "testing"

// ----------------------------------------------------------------------------
// Tests
// ----------------------------------------------------------------------------

func TestDetermineAccountType(t *testing.T) {
	tests := []struct {
		name     string
		tag      Tag
		expected AccountType
	}{
		// Known tags
		{"dependent tag", "Dependent", Dependent},
		{"grant tag", "Grant", Associate},
		{"individual tag", "Individual", IndividualGrant},
		{"blank tag defaults to associate", "", Associate},

		// Unknown tags fall back to UnknownAccount
		{"unrecognized tag", "Scholarship", UnknownAccount},
		{"lowercase tag", "grant", UnknownAccount},
		{"tag with trailing space", "Grant ", UnknownAccount},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := DetermineAccountType(test.tag)
			if actual != test.expected {
				t.Errorf("DetermineAccountType(%q) = %q, expected %q",
					test.tag, actual, test.expected)
			}
		})
	}
}
