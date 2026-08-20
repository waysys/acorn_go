// ----------------------------------------------------------------------------
//
// # termperiod_test
//
// Tests for DetermineTermPeriod.
//
// Author: William Shaffer
//
// Copyright (c) 2026 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

package main

import (
	"testing"

	d "github.com/waysys/waydate/pkg/date"
)

// ----------------------------------------------------------------------------
// Tests
// ----------------------------------------------------------------------------

func TestDetermineTermPeriod(t *testing.T) {
	tests := []struct {
		name     string
		month    d.Month
		day      d.Day
		year     d.Year
		expected TermPeriod
	}{
		// Prior Term: 9/1/2025 - 10/31/2025
		{"prior term first day", 9, 1, 2025, PriorTerm},
		{"prior term interior", 10, 15, 2025, PriorTerm},
		{"prior term last day", 10, 31, 2025, PriorTerm},

		// Spring: 11/1/2025 - 3/31/2026
		{"spring first day", 11, 1, 2025, Spring},
		{"spring interior same year", 12, 15, 2025, Spring},
		{"spring interior next year", 2, 15, 2026, Spring},
		{"spring last day of february", 2, 28, 2026, Spring},
		{"spring last day", 3, 31, 2026, Spring},

		// Summer: 4/1/2026 - 6/30/2026
		{"summer first day", 4, 1, 2026, Summer},
		{"summer interior", 5, 15, 2026, Summer},
		{"summer last day", 6, 30, 2026, Summer},

		// Fall: 7/1/2026 - 8/31/2026
		{"fall first day", 7, 1, 2026, Fall},
		{"fall interior", 8, 15, 2026, Fall},
		{"fall last day", 8, 31, 2026, Fall},

		// Outside: before 9/1/2025 or after 8/31/2026
		{"outside day before prior term", 8, 31, 2025, Outside},
		{"outside far past", 1, 1, 2020, Outside},
		{"outside day after fall", 9, 1, 2026, Outside},
		{"outside far future", 12, 31, 2030, Outside},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			billDate := mustDate(test.month, test.day, test.year)
			actual := DetermineTermPeriod(billDate)
			if actual != test.expected {
				t.Errorf("DetermineTermPeriod(%d/%d/%d) = %q, expected %q",
					test.month, test.day, test.year, actual, test.expected)
			}
		})
	}
}
