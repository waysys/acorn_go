// ----------------------------------------------------------------------------
//
// Accounting package test
//
// Author: William Shaffer
// Version: 13-Apr-2024
//
// Copyright (c) William Shaffer
//
// ----------------------------------------------------------------------------

// Package accounting contains items related to accounting.
package accounting

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	"os"
	"testing"

	d "github.com/waysys/waydate/pkg/date"
)

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

// Test_FYIndicator checks the FiscalYearIndicator function
func Test_FYIndicator(t *testing.T) {

	var date2023_1, _ = d.New(9, 1, 2022)
	var date2023_2, _ = d.New(8, 31, 2023)
	var date2024_1, _ = d.New(9, 1, 2023)
	var date2024_2, _ = d.New(8, 31, 2024)
	var date2025_1, _ = d.New(9, 1, 2024)
	var date2025_2, _ = d.New(8, 31, 2025)
	var dateOutOfRange, _ = d.New(9, 1, 2025)

	var testFunction = func(t *testing.T) {
		var fy = FiscalYearIndicator(date2023_1)
		if fy != FY2023 {
			t.Error("Fiscal year should be FY2023, not: " + fy.String())
		}
		fy = FiscalYearIndicator(date2023_2)
		if fy != FY2023 {
			t.Error("Fiscal year should be FY2023, not: " + fy.String())
		}
		fy = FiscalYearIndicator(date2024_1)
		if fy != FY2024 {
			t.Error("Fiscal year should be FY2024, not: " + fy.String())
		}
		fy = FiscalYearIndicator(date2024_2)
		if fy != FY2024 {
			t.Error("Fiscal year should be FY2024, not: " + fy.String())
		}
		fy = FiscalYearIndicator(date2025_1)
		if fy != FY2025 {
			t.Error("Fiscal year should be FY2025, not: " + fy.String())
		}
		fy = FiscalYearIndicator(date2025_2)
		if fy != FY2025 {
			t.Error("Fiscal year should be FY2025, not: " + fy.String())
		}
		fy = FiscalYearIndicator(dateOutOfRange)
		if fy != OutOfRange {
			t.Error("Fiscal year should be OutOfRange, not: " + fy.String())
		}
	}

	t.Run("Test_FYIndicator", testFunction)
}
