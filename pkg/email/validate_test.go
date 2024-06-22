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
