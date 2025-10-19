// ----------------------------------------------------------------------------
//
// Accounting package test
//
// Author: William Shaffer
//
// Copyright (c) 2024, 2025 William Shaffer All Rights Reserved
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
	var date2026_1, _ = d.New(9, 1, 2025)
	var date2026_2, _ = d.New(8, 31, 2026)
	var dateOutOfRange, _ = d.New(9, 1, 2026)

	var testFunction = func(t *testing.T) {
		var fy = FiscalYearIndicator(date2026_1)
		if fy != FY2026 {
			t.Error("Fiscal year should be FY2026, not: " + fy.String())
		}
		fy = FiscalYearIndicator(date2026_2)
		if fy != FY2026 {
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

		fy = FiscalYearIndicator(date2023_1)
		if fy != FY2023 {
			t.Error("Fiscal year should be FY2023, not: " + fy.String())
		}

		fy = FiscalYearIndicator(date2023_2)
		if fy != FY2023 {
			t.Error("Fiscal year should be FY2023, not: " + fy.String())
		}
		fy = FiscalYearIndicator(dateOutOfRange)
		if fy != OutOfRange {
			t.Error("Fiscal year should be OutOfRange, not: " + fy.String())
		}
	}

	t.Run("Test_FYIndicator", testFunction)
}

// Test_CalendarYear checks YIndicator and IsCurrentMonth
func Test_CalendarYear(t *testing.T) {
	var d2022, _ = d.New(6, 1, 2022)
	var d2023, _ = d.New(6, 1, 2023)
	var d2024, _ = d.New(6, 1, 2024)
	var d2025, _ = d.New(6, 1, 2025)
	var d2026, _ = d.New(6, 1, 2026)

	if YIndicator(d2022) != Y2022 {
		t.Error("YIndicator for 2022 should be Y2022, not: " + YIndicator(d2022).String())
	}
	if YIndicator(d2023) != Y2023 {
		t.Error("YIndicator for 2023 should be Y2023, not: " + YIndicator(d2023).String())
	}
	if YIndicator(d2024) != Y2024 {
		t.Error("YIndicator for 2024 should be Y2024, not: " + YIndicator(d2024).String())
	}
	if YIndicator(d2025) != Y2025 {
		t.Error("YIndicator for 2025 should be Y2025, not: " + YIndicator(d2025).String())
	}
	if YIndicator(d2026) != Y2026 {
		t.Error("YIndicator for 2026 should be Y2026, not: " + YIndicator(d2026).String())
	}

	// current month is 2025-09-01 .. 2025-09-30 in calendar_year.go
	var inside, _ = d.New(9, 15, 2025)
	var outside, _ = d.New(8, 1, 2025)

	if !IsCurrentMonth(inside) {
		t.Error("IsCurrentMonth should be true for 2025-09-15")
	}
	if IsCurrentMonth(outside) {
		t.Error("IsCurrentMonth should be false for 2025-08-01")
	}
}
