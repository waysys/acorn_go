// ----------------------------------------------------------------------------
//
// Column Handler Tests
//
// This file contains test of the Column Handler.
//
// Author: William Shaffer
//
// Copyright (c) 2025 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

package spreadsheet

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	a "acorn_go/pkg/accounting"
	"testing"
)

// ----------------------------------------------------------------------------
// Tests
// ----------------------------------------------------------------------------

func TestNextLetterBasic(t *testing.T) {
	if NextLetter("a") != "b" {
		t.Fatalf("expected b, got %s", NextLetter("a"))
	}
	if NextLetter("y") != "z" {
		t.Fatalf("expected z, got %s", NextLetter("y"))
	}
	if NextLetter("z") != "a" {
		t.Fatalf("expected a, got %s", NextLetter("z"))
	}
	if NextLetter("A") != "B" {
		t.Fatalf("expected B, got %s", NextLetter("A"))
	}
	if NextLetter("Z") != "A" {
		t.Fatalf("expected A, got %s", NextLetter("Z"))
	}
}

func TestNextSecondLetter(t *testing.T) {
	if NextSecondLetter("a") != "c" {
		t.Fatalf("expected c, got %s", NextSecondLetter("a"))
	}
	if NextSecondLetter("y") != "a" { // y -> z -> a
		t.Fatalf("expected a, got %s", NextSecondLetter("y"))
	}
}

func TestIsValidColumn(t *testing.T) {
	if !isValidColumn("a") {
		t.Fatalf("a should be valid")
	}
	if isValidColumn("") {
		t.Fatalf("empty string should be invalid")
	}
	if isValidColumn("ab") {
		t.Fatalf("ab should be invalid")
	}
}

func TestColumnHandlerBasics(t *testing.T) {
	ch := NewColumnHandler("A")
	// Count columns should have length NumFiscalYears
	if len(ch.countColumns) != a.NumFiscalYears {
		t.Fatalf("countColumns length expected %d, got %d", a.NumFiscalYears, len(ch.countColumns))
	}
	if len(ch.paymentColumns) != a.NumFiscalYears {
		t.Fatalf("paymentColumns length expected %d, got %d", a.NumFiscalYears, len(ch.paymentColumns))
	}
	// Validate CountColumn and PaymentColumn for each fiscal year
	for i := 0; i < a.NumFiscalYears; i++ {
		fy := a.FYIndicator(i)
		_ = ch.CountColumn(fy)
		_ = ch.PaymentColumn(fy)
	}
	// TotalColumn should be non-empty
	if ch.TotalColumn() == "" {
		t.Fatalf("total column should not be empty")
	}
}

// helper to assert that f panics
func assertPanics(t *testing.T, f func(), name string) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic for %s, but none occurred", name)
		}
	}()
	f()
}

func TestNextLetterPanicsOnInvalid(t *testing.T) {
	// multi-character string should cause assertion panic
	assertPanics(t, func() { NextLetter("ab") }, "NextLetter with multi-char")
	// non-letter should cause assertion panic
	assertPanics(t, func() { NextLetter("1") }, "NextLetter with digit")
	// Non-ASCII rune (e.g., multibyte) also should panic because not an ASCII letter
	assertPanics(t, func() { NextLetter("λ") }, "NextLetter with non-ASCII")
}

func TestIsValidColumnMultibyte(t *testing.T) {
	// Multibyte single-rune should be considered length 1 by RuneCountInString
	if !isValidColumn("λ") {
		t.Fatalf("λ should be considered a valid single-character column by isValidColumn")
	}
}

func TestColumnHandlerInvalidFYIndicatorPanics(t *testing.T) {
	ch := NewColumnHandler("A")
	// OutOfRange is not a valid FYIndicator; methods should assert and panic
	assertPanics(t, func() { ch.CountColumn(a.OutOfRange) }, "CountColumn with OutOfRange")
	assertPanics(t, func() { ch.PaymentColumn(a.OutOfRange) }, "PaymentColumn with OutOfRange")
}
