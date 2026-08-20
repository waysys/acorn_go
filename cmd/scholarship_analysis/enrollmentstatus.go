// ----------------------------------------------------------------------------
//
// EnrollmentStatus determines whether the applicant is going full time
// or part time.
//
//
// Author: William Shaffer
//
// Copyright (c) 2026 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

package main

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	"fmt"

	dec "github.com/shopspring/decimal"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type EnrollmentStatus string

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

const (
	UnknownEnrollment EnrollmentStatus = ""
	FullTime          EnrollmentStatus = "FullTime"
	PartTime          EnrollmentStatus = "PartTime"
)

// ----------------------------------------------------------------------------
// Data
// ----------------------------------------------------------------------------

var maxTwoYearPartTimeAmount = dec.NewFromInt(600)
var maxFourYearPartTimeAmount = dec.NewFromInt(2000)

// ----------------------------------------------------------------------------
// Function
// ----------------------------------------------------------------------------

// DetermineEnrollmentStatus determines the enrollment status based on
// the institution type and amount of scholarship payment. It returns
// an error if institutionType is not one of the recognized
// InstitutionType values.
func DetermineEnrollmentStatus(
	institutionType InstitutionType,
	amount dec.Decimal) (EnrollmentStatus, error) {
	var maxPartTimeAmount dec.Decimal

	switch institutionType {
	case FourYear:
		maxPartTimeAmount = maxFourYearPartTimeAmount
	case TwoYear:
		maxPartTimeAmount = maxTwoYearPartTimeAmount
	default:
		return UnknownEnrollment, fmt.Errorf("unrecognized institution type: %s", institutionType)
	}

	// amount of payment must be greater than the max amount to
	// qualify as full time
	if amount.GreaterThan(maxPartTimeAmount) {
		return FullTime, nil
	}
	return PartTime, nil
}
