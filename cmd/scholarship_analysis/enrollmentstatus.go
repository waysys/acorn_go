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

var sortOrderEnrollmentStatus = map[EnrollmentStatus]int{
	FullTime:          0,
	PartTime:          1,
	UnknownEnrollment: 2,
}

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
	case UnknownInstitution:
		return UnknownEnrollment, nil
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

// CompareEnrollmentStatus using the associated sort order
func CompareEnrollmentStatus(
	enrollmentStatus1 EnrollmentStatus,
	enrollmentStatus2 EnrollmentStatus,
) int {
	var order1 = sortOrderEnrollmentStatus[enrollmentStatus1]
	var order2 = sortOrderEnrollmentStatus[enrollmentStatus2]

	var result = 0
	if order1 < order2 {
		result = -1
	} else if order1 > order2 {
		result = 1
	}

	return result
}
