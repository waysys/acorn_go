// ----------------------------------------------------------------------------
//
// Email Validation Package
//
// Author: William Shaffer
// Version: 21-Jun-2024
//
// Copyright (c) 2024 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

// This package validates email addresses using the Zero Bounce API
package email

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	"errors"
	"os"

	"github.com/waysys/assert/assert"
	z "github.com/zerobounce/zerobouncego"
)

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

const key = "ZERO_BOUNCE_API_KEY"

// ----------------------------------------------------------------------------
// Functions
// ----------------------------------------------------------------------------

// Validate returns true if the specified email is an existing email.
// This service depends on an API key from Zero Bounce being
// established as an evironmental variable for the user.
func Validate(email string) (bool, error) {
	var err error = nil
	var response *z.ValidateResponse
	var result = false
	//
	// Preconditions
	//
	assert.Precondition(hasKey())
	assert.Precondition(hasCredits())
	//
	// Query Zero Balance
	//
	response, err = z.Validate(email, "")
	if err == nil {
		result = checkStatus(response.Status)
	}
	return result, err
}

// hasKey returns an error if an environmental variable ZERO_BOUNCE_API_KEY
// has not been defined for the user.
func hasKey() error {
	var err error = nil
	var code = os.Getenv(key)
	if code == "" {
		err = errors.New("ZERO_BOUNCE_API_KEY is not defined")
	}
	return err
}

// hasCredits returns an error if the Zero Bounce account does not have
// a positive credit balance.
func hasCredits() error {
	var err error = nil
	var credits *z.CreditsResponse

	credits, err = z.GetCredits()
	if err == nil {
		if credits.Credits() <= 0 {
			err = errors.New("account does not have positive credit balance")
		}
	}
	return err
}

// checkStatus returns true is the status of the response indicates
// that the email can be used.
func checkStatus(status string) bool {
	var result = false
	switch status {
	case z.S_VALID:
		result = true
	case z.S_INVALID:
		result = false
	case z.S_CATCH_ALL:
		result = true
	case z.S_SPAMTRAP:
		result = false
	case z.S_ABUSE:
		result = true
	case z.S_DO_NOT_MAIL:
		result = false
	case z.S_UNKNOWN:
		result = true
	}
	return result
}
