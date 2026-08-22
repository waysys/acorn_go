// ----------------------------------------------------------------------------
//
// Scholarship Entry
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
	dec "github.com/shopspring/decimal"
	d "github.com/waysys/waydate/pkg/date"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type Scholarship struct {
	termPeriod       TermPeriod
	institutionType  InstitutionType
	accountType      AccountType
	enrollmentStatus EnrollmentStatus
}

type Entry struct {
	amount      dec.Decimal
	scholarship Scholarship
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// Factory Function
// ----------------------------------------------------------------------------

// NewEntry creates a new entry using input from the bill spreadsheet
func NewEntry(
	billDate d.Date,
	university string,
	tag Tag,
	amount dec.Decimal) (Entry, error) {

	var err error
	var entry Entry
	var institutionType InstitutionType
	var accountType AccountType
	var enrollmentStatus EnrollmentStatus
	//
	// Calculate inputs
	//
	var termPeriod = DetermineTermPeriod(billDate)

	accountType = DetermineAccountType(tag)

	if (accountType == Associate) || (accountType == Dependent) {
		institutionType, err = DetermineInstitutionType(university)
		if err != nil {
			return entry, err
		}
	} else {
		institutionType = UnknownInstitution
	}

	enrollmentStatus, err = DetermineEnrollmentStatus(institutionType, amount)
	if err != nil {
		return entry, err
	}
	//
	// Define the Scholarship structure
	//
	var scholarship = Scholarship{
		termPeriod:       termPeriod,
		institutionType:  institutionType,
		accountType:      accountType,
		enrollmentStatus: enrollmentStatus,
	}
	//
	// Define the Entry structure
	//
	entry = Entry{
		amount:      amount,
		scholarship: scholarship,
	}

	return entry, nil
}

// ----------------------------------------------------------------------------
// Methods - Scholarship
// ----------------------------------------------------------------------------

// TermPeriod returns the time period
func (sch Scholarship) TermPeriod() TermPeriod {
	return sch.termPeriod
}

// InstitutionType returns the institution type
func (sch Scholarship) InstitutionType() InstitutionType {
	return sch.institutionType
}

// AccountType returns the account type
func (sch Scholarship) AccountType() AccountType {
	return sch.accountType
}

// EnrollmentStatus returns the enrollment status
func (sch Scholarship) EnrollmentStatus() EnrollmentStatus {
	return sch.enrollmentStatus
}

// Compare this scholarship to another
func (sch Scholarship) Compare(anotherSch Scholarship) int {
	var result = CompareInstitutionType(sch.InstitutionType(), anotherSch.InstitutionType())
	if result == 0 {
		result = CompareTermPeriod(sch.TermPeriod(), anotherSch.TermPeriod())
	}
	if result == 0 {
		result = CompareEnrollmentStatus(sch.EnrollmentStatus(), anotherSch.EnrollmentStatus())
	}
	if result == 0 {
		result = CompareAccountType(sch.AccountType(), anotherSch.AccountType())
	}
	return result
}

// ----------------------------------------------------------------------------
// Methods - Entry
// ----------------------------------------------------------------------------

// Amount returns the amount
func (entry Entry) Amount() dec.Decimal {
	return entry.amount
}

// Scholarship returns the scholarship
func (entry Entry) Scholarship() Scholarship {
	return entry.scholarship
}
