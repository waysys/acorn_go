// ----------------------------------------------------------------------------
//
// Donor Count Tests
//
// Author: William Shaffer
//
// Copyright (c) 2025 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

package donations

import (
	"testing"

	a "acorn_go/pkg/accounting"
	dn "acorn_go/pkg/donors"

	dec "github.com/shopspring/decimal"
)

// Test_NewDonorCount verifies NewDonorCount initialization
func Test_NewDonorCount(t *testing.T) {
	dc := NewDonorCount(a.FY2025)
	if dc.FY() != a.FY2025 {
		t.Fatalf("expected FY %s, got %s", a.FY2025.String(), dc.FY().String())
	}
	if dc.TotalDonorCount() != 0 {
		t.Fatalf("expected total 0, got %d", dc.TotalDonorCount())
	}
	if len(dc.donorCount) != a.NumFiscalYears {
		t.Fatalf("donorCount length expected %d, got %d", a.NumFiscalYears, len(dc.donorCount))
	}
}

// Test_Add_Count_Total tests Add, Count and TotalDonorCount
func Test_Add_Count_Total(t *testing.T) {
	dc := NewDonorCount(a.FY2026)
	dc.Add(CurrentYear, 2)
	dc.Add(PriorYear, 1)
	if dc.Count(CurrentYear) != 2 {
		t.Fatalf("expected CurrentYear count 2, got %d", dc.Count(CurrentYear))
	}
	if dc.Count(PriorYear) != 1 {
		t.Fatalf("expected PriorYear count 1, got %d", dc.Count(PriorYear))
	}
	if dc.TotalDonorCount() != 3 {
		t.Fatalf("expected total 3, got %d", dc.TotalDonorCount())
	}
}

// Test_Add_Negative_Panic ensures Add panics on negative count
func Test_Add_Negative_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic on negative Add")
		}
	}()
	dc := NewDonorCount(a.FY2026)
	dc.Add(CurrentYear, -1)
}

// Test_ApplyDonorCount checks Prior/PriorPrior/Current logic
func Test_ApplyDonorCount(t *testing.T) {
	analysisFy := a.FY2026

	// donorA donated in FY2025 -> should count as PriorYear
	donorA := dn.NewDonorWithDonation("A")
	donorA.AddDonation(dec.NewFromInt(100), a.FY2025)

	// donorB donated in FY2024 only -> PriorPriorYear
	donorB := dn.NewDonorWithDonation("B")
	donorB.AddDonation(dec.NewFromInt(50), a.FY2024)

	// donorC donated in FY2026 only -> CurrentYear
	donorC := dn.NewDonorWithDonation("C")
	donorC.AddDonation(dec.NewFromInt(75), a.FY2026)

	dc := NewDonorCount(analysisFy)
	// use pointer receiver
	dc.ApplyDonorCount(&donorA, analysisFy)
	dc.ApplyDonorCount(&donorB, analysisFy)
	dc.ApplyDonorCount(&donorC, analysisFy)

	if dc.Count(PriorYear) != 1 {
		t.Fatalf("expected PriorYear 1, got %d", dc.Count(PriorYear))
	}
	if dc.Count(PriorPriorYear) != 1 {
		t.Fatalf("expected PriorPriorYear 1, got %d", dc.Count(PriorPriorYear))
	}
	if dc.Count(CurrentYear) != 1 {
		t.Fatalf("expected CurrentYear 1, got %d", dc.Count(CurrentYear))
	}
}
