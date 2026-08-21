// ----------------------------------------------------------------------------
//
// # main_test
//
// Tests for the pure functions in main: processEntryCount,
// processEntryAmount, average, and sortScholarships.
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
// Test Helpers
// ----------------------------------------------------------------------------

// makeEntry builds an entry directly from a scholarship and amount,
// bypassing the spreadsheet and NewEntry validation.
func makeEntry(scholarship Scholarship, amount int64) Entry {
	return Entry{
		amount:      dec.NewFromInt(amount),
		scholarship: scholarship,
	}
}

var (
	schFourYearSpringFull = Scholarship{
		termPeriod:       Spring,
		institutionType:  FourYear,
		accountType:      Dependent,
		enrollmentStatus: FullTime,
	}
	schFourYearFallPart = Scholarship{
		termPeriod:       Fall,
		institutionType:  FourYear,
		accountType:      Associate,
		enrollmentStatus: PartTime,
	}
	schTwoYearFallPart = Scholarship{
		termPeriod:       Fall,
		institutionType:  TwoYear,
		accountType:      Dependent,
		enrollmentStatus: PartTime,
	}
	schTwoYearPriorFull = Scholarship{
		termPeriod:       PriorTerm,
		institutionType:  TwoYear,
		accountType:      Associate,
		enrollmentStatus: FullTime,
	}
)

// ----------------------------------------------------------------------------
// Tests
// ----------------------------------------------------------------------------

func TestProcessEntryCount(t *testing.T) {
	billList := []Entry{
		makeEntry(schFourYearSpringFull, 5000),
		makeEntry(schFourYearSpringFull, 3000),
		makeEntry(schTwoYearFallPart, 600),
	}

	billCount := processEntryCount(billList)

	if len(billCount) != 2 {
		t.Errorf("len(billCount) = %d, expected 2", len(billCount))
	}
	if billCount[schFourYearSpringFull] != 2 {
		t.Errorf("billCount[four year spring] = %d, expected 2",
			billCount[schFourYearSpringFull])
	}
	if billCount[schTwoYearFallPart] != 1 {
		t.Errorf("billCount[two year fall] = %d, expected 1",
			billCount[schTwoYearFallPart])
	}
}

func TestProcessEntryCountEmpty(t *testing.T) {
	billCount := processEntryCount(nil)
	if len(billCount) != 0 {
		t.Errorf("len(billCount) = %d, expected 0", len(billCount))
	}
}

func TestProcessEntryAmount(t *testing.T) {
	billList := []Entry{
		makeEntry(schFourYearSpringFull, 5000),
		makeEntry(schFourYearSpringFull, 3000),
		makeEntry(schTwoYearFallPart, 600),
	}

	billAmount := processEntryAmount(billList)

	if len(billAmount) != 2 {
		t.Errorf("len(billAmount) = %d, expected 2", len(billAmount))
	}
	if !billAmount[schFourYearSpringFull].Equal(dec.NewFromInt(8000)) {
		t.Errorf("billAmount[four year spring] = %s, expected 8000",
			billAmount[schFourYearSpringFull])
	}
	if !billAmount[schTwoYearFallPart].Equal(dec.NewFromInt(600)) {
		t.Errorf("billAmount[two year fall] = %s, expected 600",
			billAmount[schTwoYearFallPart])
	}
}

func TestAverage(t *testing.T) {
	tests := []struct {
		name     string
		amount   dec.Decimal
		count    int
		expected dec.Decimal
	}{
		{"even division", dec.NewFromInt(8000), 2, dec.NewFromInt(4000)},
		{"fractional result rounds down", dec.NewFromInt(1000), 3, dec.NewFromInt(333)},
		{"fractional result rounds up", dec.NewFromInt(1001), 3, dec.NewFromInt(334)},
		{"half rounds away from zero", dec.NewFromInt(1005), 2, dec.NewFromInt(503)},
		{"single scholarship", dec.NewFromInt(1500), 1, dec.NewFromInt(1500)},
		{"zero count yields zero", dec.NewFromInt(5000), 0, dec.Zero},
		{"zero amount", dec.Zero, 4, dec.Zero},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := average(test.amount, test.count)
			if !actual.Equal(test.expected) {
				t.Errorf("average(%s, %d) = %s, expected %s",
					test.amount, test.count, actual, test.expected)
			}
		})
	}
}

func TestSortScholarships(t *testing.T) {
	billCount := map[Scholarship]int{
		schFourYearFallPart:   1,
		schTwoYearFallPart:    2,
		schFourYearSpringFull: 3,
		schTwoYearPriorFull:   1,
	}

	scholarships := sortScholarships(billCount)

	// Institution type is the primary sort key (Two Year before Four Year),
	// term period is the tiebreaker (Prior Term before Fall, Spring before Fall).
	expected := []Scholarship{
		schTwoYearPriorFull,
		schTwoYearFallPart,
		schFourYearSpringFull,
		schFourYearFallPart,
	}

	if len(scholarships) != len(expected) {
		t.Fatalf("len(scholarships) = %d, expected %d", len(scholarships), len(expected))
	}
	for index, scholarship := range scholarships {
		if scholarship != expected[index] {
			t.Errorf("scholarships[%d] = %+v, expected %+v",
				index, scholarship, expected[index])
		}
	}
}

func TestSortScholarshipsEmpty(t *testing.T) {
	scholarships := sortScholarships(map[Scholarship]int{})
	if len(scholarships) != 0 {
		t.Errorf("len(scholarships) = %d, expected 0", len(scholarships))
	}
}
