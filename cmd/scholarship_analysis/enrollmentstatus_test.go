// ----------------------------------------------------------------------------
//
// # enrollmentstatus_test
//
// Tests for DetermineEnrollmentStatus.
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
)

// ----------------------------------------------------------------------------
// Tests
// ----------------------------------------------------------------------------

func TestDetermineEnrollmentStatus(t *testing.T) {
	tests := []struct {
		name            string
		institutionType InstitutionType
		amount          dec.Decimal
		expected        EnrollmentStatus
	}{
		// Two Year threshold: $600
		{"two year below threshold", TwoYear, dec.NewFromInt(599), PartTime},
		{"two year at threshold", TwoYear, dec.NewFromInt(600), PartTime},
		{"two year above threshold", TwoYear, dec.NewFromInt(601), FullTime},
		{"two year fractional below threshold", TwoYear, dec.NewFromFloat(599.99), PartTime},
		{"two year fractional above threshold", TwoYear, dec.NewFromFloat(600.01), FullTime},

		// Four Year threshold: $2000
		{"four year below threshold", FourYear, dec.NewFromInt(1999), PartTime},
		{"four year at threshold", FourYear, dec.NewFromInt(2000), PartTime},
		{"four year above threshold", FourYear, dec.NewFromInt(2001), FullTime},
		{"four year fractional above threshold", FourYear, dec.NewFromFloat(2000.01), FullTime},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := DetermineEnrollmentStatus(test.institutionType, test.amount)
			if err != nil {
				t.Fatalf("DetermineEnrollmentStatus(%q, %s) returned unexpected error: %v",
					test.institutionType, test.amount, err)
			}
			if actual != test.expected {
				t.Errorf("DetermineEnrollmentStatus(%q, %s) = %q, expected %q",
					test.institutionType, test.amount, actual, test.expected)
			}
		})
	}
}

func TestDetermineEnrollmentStatusUnrecognizedInstitution(t *testing.T) {
	tests := []struct {
		name            string
		institutionType InstitutionType
		expectError     bool
	}{
		{"unknown institution", UnknownInstitution, false},
		{"unrecognized institution", "FiveYear", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := DetermineEnrollmentStatus(test.institutionType, dec.NewFromInt(1000))
			if test.expectError && err == nil {
				t.Errorf("DetermineEnrollmentStatus(%q, 1000) expected error, got nil",
					test.institutionType)
			}
			if !test.expectError && err != nil {
				t.Errorf("DetermineEnrollmentStatus(%q, 1000) returned unexpected error: %v",
					test.institutionType, err)
			}
			if actual != UnknownEnrollment {
				t.Errorf("DetermineEnrollmentStatus(%q, 1000) = %q, expected %q",
					test.institutionType, actual, UnknownEnrollment)
			}
		})
	}
}
