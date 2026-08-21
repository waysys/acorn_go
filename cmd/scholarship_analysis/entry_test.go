// ----------------------------------------------------------------------------
//
// # entry_test
//
// Tests for NewEntry and the Entry and Scholarship accessors.
//
// Author: William Shaffer
//
// Copyright (c) 2026 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

package main

import (
	"testing"

	dec "github.com/shopspring/decimal"
	d "github.com/waysys/waydate/pkg/date"
)

// ----------------------------------------------------------------------------
// Tests
// ----------------------------------------------------------------------------

func TestNewEntry(t *testing.T) {
	tests := []struct {
		name            string
		month           d.Month
		day             d.Day
		year            d.Year
		university      string
		tag             Tag
		amount          dec.Decimal
		expectedTerm    TermPeriod
		expectedInst    InstitutionType
		expectedAccount AccountType
		expectedEnroll  EnrollmentStatus
	}{
		{"four year full time dependent", 1, 15, 2026,
			"North Carolina State University", "Dependent", dec.NewFromInt(5000),
			Spring, FourYear, Dependent, FullTime},
		{"four year part time at threshold", 5, 1, 2026,
			"East Carolina University", "Grant", dec.NewFromInt(2000),
			Summer, FourYear, Associate, PartTime},
		{"two year full time", 8, 1, 2026,
			"Wake Technical Community College", "Dependent", dec.NewFromInt(601),
			Fall, TwoYear, Dependent, FullTime},
		{"two year part time at threshold", 9, 1, 2025,
			"Pitt Community College", "", dec.NewFromInt(600),
			PriorTerm, TwoYear, Associate, PartTime},
		{"individual grant skips institution lookup", 2, 1, 2026,
			"Some Employer", "Individual", dec.NewFromInt(3000),
			Spring, UnknownInstitution, IndividualGrant, UnknownEnrollment},
		{"unrecognized tag", 2, 1, 2026,
			"Shaw University", "Scholarship", dec.NewFromInt(3000),
			Spring, UnknownInstitution, UnknownAccount, UnknownEnrollment},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			billDate := mustDate(test.month, test.day, test.year)
			entry, err := NewEntry(billDate, test.university, test.tag, test.amount)
			if err != nil {
				t.Fatalf("NewEntry returned unexpected error: %v", err)
			}
			if !entry.Amount().Equal(test.amount) {
				t.Errorf("Amount() = %s, expected %s", entry.Amount(), test.amount)
			}
			scholarship := entry.Scholarship()
			if scholarship.TermPeriod() != test.expectedTerm {
				t.Errorf("TermPeriod() = %q, expected %q",
					scholarship.TermPeriod(), test.expectedTerm)
			}
			if scholarship.InstitutionType() != test.expectedInst {
				t.Errorf("InstitutionType() = %q, expected %q",
					scholarship.InstitutionType(), test.expectedInst)
			}
			if scholarship.AccountType() != test.expectedAccount {
				t.Errorf("AccountType() = %q, expected %q",
					scholarship.AccountType(), test.expectedAccount)
			}
			if scholarship.EnrollmentStatus() != test.expectedEnroll {
				t.Errorf("EnrollmentStatus() = %q, expected %q",
					scholarship.EnrollmentStatus(), test.expectedEnroll)
			}
		})
	}
}

func TestNewEntryUnknownUniversity(t *testing.T) {
	billDate := mustDate(1, 15, 2026)
	entry, err := NewEntry(billDate, "Not A Real University", "Grant", dec.NewFromInt(1000))
	if err == nil {
		t.Error("NewEntry with unknown university expected error, got nil")
	}
	if entry != (Entry{}) {
		t.Errorf("NewEntry with unknown university returned non-zero entry: %+v", entry)
	}
}

func TestNewEntryOutsideTermPeriod(t *testing.T) {
	tests := []struct {
		name  string
		month d.Month
		day   d.Day
		year  d.Year
	}{
		{"before range", 8, 31, 2025},
		{"after range", 9, 1, 2026},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			billDate := mustDate(test.month, test.day, test.year)
			entry, err := NewEntry(billDate, "Shaw University", "Grant", dec.NewFromInt(1000))
			if err != nil {
				t.Fatalf("NewEntry returned unexpected error: %v", err)
			}
			if entry.Scholarship().TermPeriod() != Outside {
				t.Errorf("TermPeriod() = %q, expected %q",
					entry.Scholarship().TermPeriod(), Outside)
			}
		})
	}
}

func TestEntryAccessors(t *testing.T) {
	scholarship := Scholarship{
		termPeriod:       Summer,
		institutionType:  FourYear,
		accountType:      Dependent,
		enrollmentStatus: FullTime,
	}
	entry := Entry{
		amount:      dec.NewFromFloat(1234.56),
		scholarship: scholarship,
	}

	if !entry.Amount().Equal(dec.NewFromFloat(1234.56)) {
		t.Errorf("Amount() = %s, expected 1234.56", entry.Amount())
	}
	if entry.Scholarship() != scholarship {
		t.Errorf("Scholarship() = %+v, expected %+v", entry.Scholarship(), scholarship)
	}
}

func TestScholarshipAccessors(t *testing.T) {
	scholarship := Scholarship{
		termPeriod:       Fall,
		institutionType:  TwoYear,
		accountType:      IndividualGrant,
		enrollmentStatus: PartTime,
	}

	if scholarship.TermPeriod() != Fall {
		t.Errorf("TermPeriod() = %q, expected %q", scholarship.TermPeriod(), Fall)
	}
	if scholarship.InstitutionType() != TwoYear {
		t.Errorf("InstitutionType() = %q, expected %q", scholarship.InstitutionType(), TwoYear)
	}
	if scholarship.AccountType() != IndividualGrant {
		t.Errorf("AccountType() = %q, expected %q", scholarship.AccountType(), IndividualGrant)
	}
	if scholarship.EnrollmentStatus() != PartTime {
		t.Errorf("EnrollmentStatus() = %q, expected %q", scholarship.EnrollmentStatus(), PartTime)
	}
}

func TestScholarshipCompare(t *testing.T) {
	var base = Scholarship{
		termPeriod:       Spring,
		institutionType:  FourYear,
		accountType:      Dependent,
		enrollmentStatus: FullTime,
	}
	var twoYear = base
	twoYear.institutionType = TwoYear
	var fall = base
	fall.termPeriod = Fall
	var associate = base
	associate.accountType = Associate
	var partTime = base
	partTime.enrollmentStatus = PartTime

	tests := []struct {
		name     string
		a        Scholarship
		b        Scholarship
		expected int
	}{
		{"identical scholarships", base, base, 0},
		{"institution type decides first", twoYear, base, -1},
		{"institution type reversed", base, twoYear, 1},
		{"term period breaks institution tie", base, fall, -1},
		{"account type breaks term tie", associate, base, -1},
		{"enrollment status breaks account tie", base, partTime, -1},
		{"enrollment status reversed", partTime, base, 1},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := test.a.Compare(test.b)
			if actual != test.expected {
				t.Errorf("Compare() = %d, expected %d", actual, test.expected)
			}
		})
	}
}
