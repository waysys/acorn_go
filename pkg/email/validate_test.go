// ----------------------------------------------------------------------------
//
// Email Validation Tests
//
// Author: William Shaffer
// Version: 21-Jun-2024
//
// Copyright (c) 2024 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

package email

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	"os"
	"testing"
)

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// Test Main
// ----------------------------------------------------------------------------

func TestMain(m *testing.M) {
	exitVal := m.Run()
	os.Exit(exitVal)
}

// ----------------------------------------------------------------------------
// Test functions
// ----------------------------------------------------------------------------

// Test_Valid checks that a valid email will pass.
func Test_Valid(t *testing.T) {
	var email = "valid@example.com"
	var valid, err = Validate(email)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if !valid {
		t.Error("email was not found to be valid: " + email)
	}
}

func Test_MultiValid(t *testing.T) {
	type aTest struct {
		email string
		value bool
	}
	var data = []aTest{
		{"invalid@example.com", false},
		{"donotmail@example.com", false},
		{"unknown@example.com", true},
		{"failed_syntax_check@example.com", false},
		{"mailbox_quota_exceeded@example.com", false},
		{"possible_typo@example.com", false},
		{"role_based_catch_all@example.com", false},
	}

	var tt aTest
	var testFunction = func(t *testing.T) {
		var err error = nil
		var email = tt.email
		var valid = false

		valid, err = Validate(email)
		if err != nil {
			t.Error(err.Error())
		} else if valid != tt.value {
			if !valid {
				t.Error("Validation did not recognize valid email: " + email)
			} else {
				t.Error("Validation did not recognize invalid email: " + email)
			}
		}

	}
	for _, d := range data {
		tt = d
		t.Run(d.email, testFunction)
	}
}
