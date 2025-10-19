// ----------------------------------------------------------------------------
//
// Donor List Test
//
// Author: William Shaffer
//
// Copyright (c) 2024, 2025 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

package donations

import (
	"strings"
	"testing"

	dn "acorn_go/pkg/donors"

	dec "github.com/shopspring/decimal"
	d "github.com/waysys/waydate/pkg/date"
)

// Test_AddDonation_MissingDonor_Error ensures AddDonation returns an error
// when the donor is not present in the list.
func Test_AddDonation_MissingDonor_Error(t *testing.T) {
	dl := make(DonationList)
	date, _ := d.New(9, 1, 2025)
	err := dl.AddDonation("Alice", dec.NewFromInt(10), date)
	if err == nil || !strings.Contains(err.Error(), "donor name not found") {
		t.Fatalf("expected donor-not-found error, got: %v", err)
	}
}

// Test_AddDonation_NegativeAmount_Error ensures AddDonation rejects negative amounts
func Test_AddDonation_NegativeAmount_Error(t *testing.T) {
	dl := make(DonationList)
	donor := dn.NewDonorWithDonation("Bob")
	dl["Bob"] = &donor
	date, _ := d.New(9, 1, 2025)
	err := dl.AddDonation("Bob", dec.NewFromInt(-5), date)
	if err == nil || !strings.Contains(err.Error(), "must not be negative") {
		t.Fatalf("expected negative-amount error, got: %v", err)
	}
}

// Test_DonorKeys_Sorted checks that DonorKeys returns keys sorted alphabetically
func Test_DonorKeys_Sorted(t *testing.T) {
	dl := make(DonationList)
	a := dn.NewDonorWithDonation("alice")
	b := dn.NewDonorWithDonation("bob")
	c := dn.NewDonorWithDonation("charlie")
	dl["bob"] = &b
	dl["alice"] = &a
	dl["charlie"] = &c

	keys := dl.DonorKeys()
	if len(keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(keys))
	}
	if keys[0] != "alice" || keys[1] != "bob" || keys[2] != "charlie" {
		t.Fatalf("keys not sorted: %v", keys)
	}
}

// Test_Get_PanicOnMissing ensures Get panics when the donor is not present.
func Test_Get_PanicOnMissing(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic from Get on missing donor")
		}
	}()
	dl := make(DonationList)
	_ = dl.Get("does-not-exist")
}
