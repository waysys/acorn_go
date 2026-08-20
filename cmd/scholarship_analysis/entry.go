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
	"errors"

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
	count       int
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
	count int,
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
	if termPeriod == Outside {
		err = errors.New("Bill date is outside of date range of interest")
		return entry, err
	}
	institutionType, err = DetermineInstitutionType(university)
	if err != nil {
		return entry, err
	}
	accountType = DetermineAccountType(tag)
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
		count:       count,
		amount:      amount,
		scholarship: scholarship,
	}

	return entry, err
}

// ----------------------------------------------------------------------------
// Methods - Scholarship
// ----------------------------------------------------------------------------

// Return the time period
func (sch Scholarship) TermPeriod() TermPeriod {
	return sch.termPeriod
}

// Return the institution type
func (sch Scholarship) InstitutionType() InstitutionType {
	return sch.institutionType
}

// Return the account type
func (sch Scholarship) AccountType() AccountType {
	return sch.accountType
}

// Return enrollment status
func (sch Scholarship) EnrollmentStatus() EnrollmentStatus {
	return sch.enrollmentStatus
}

// ----------------------------------------------------------------------------
// Methods - Entry
// ----------------------------------------------------------------------------

// Return the entry count
func (entry Entry) Count() int {
	return entry.count
}

// Return the amount
func (entry Entry) Amount() dec.Decimal {
	return entry.amount
}

// Return the scholarship
func (entry Entry) Scholarship() Scholarship {
	return entry.scholarship
}
