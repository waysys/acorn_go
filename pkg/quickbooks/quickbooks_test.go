// ----------------------------------------------------------------------------
//
// AP Transaction
//
// Author: William Shaffer
// Version: 25-May-2024
//
// Copyright (c) 2024 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

package quickbooks

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	d "acorn_go/pkg/date"
	"os"
	"testing"

	dec "github.com/shopspring/decimal"
)

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

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

// Test_NewVendor checks the creation of vendors.
func Test_NewVendor(t *testing.T) {
	var vendorList = (&APVendorList)
	var vendorName = "Wake Technical Community College"
	var unknownVendor = "XXX"
	var vendorPtr *Vendor = nil

	vendorPtr = vendorList.Add(vendorName)

	var testFunction = func(t *testing.T) {
		if vendorPtr.Name() != vendorName {
			t.Error("vendor.Name() returned wrong value: " + vendorPtr.Name())
		}
		if !vendorList.Contains(vendorName) {
			t.Error("vendor list does not contain: " + vendorName)
		}
		var anotherVendor = vendorList.Get(vendorName)
		if anotherVendor.Name() != vendorName {
			t.Error("vendor list Get did not return a valid structures: " + anotherVendor.Name())
		}
		if vendorList.Contains(unknownVendor) {
			t.Error("vendor list Contains found unknown vendor: " + unknownVendor)
		}

	}

	t.Run("Test_New", testFunction)
}

// Test_NewRecipient checks the creation of recipients.
func Test_NewRecipient(t *testing.T) {
	var recipientList = (&APRecipientList)
	var recipientName = "Jones, Jack"
	var unknownRecipientName = "XXX"
	var recipientPtr *Recipient = nil

	recipientPtr = recipientList.Add(recipientName)

	var testFunction = func(t *testing.T) {
		if recipientPtr.Name() != recipientName {
			t.Error("recipient name is wrong: " + recipientPtr.Name())
		}
		if !recipientList.Contains(recipientName) {
			t.Error("recipient list could not find: " + recipientName)
		}
		var anotherRecipient = recipientList.Get(recipientName)
		if recipientPtr != anotherRecipient {
			t.Error("recipient pointer does not equal another recipient pointer")
		}
		if recipientList.Contains(unknownRecipientName) {
			t.Error("recipient list found unknown recipient: " + unknownRecipientName)
		}
	}

	t.Run("Test_New", testFunction)
}

// Test_APTransaction checks the creation of AP transactions.
func Test_NewAPTransaction(t *testing.T) {
	var transactionDate, _ = d.New(1, 2, 2024)
	var vendorName = "Wake Technical Community College"
	var recipientName = "Jones, Jake"
	var transactionType = "Bill"
	var amount Money = Money(dec.NewFromInt(100))
	var transaction = NewAPTransaction(transactionDate, vendorName, recipientName, transactionType, amount)

	var testFunction = func(t *testing.T) {
		if transaction.transactionDate != transactionDate {
			t.Error("transaction date did not matching origin")
		}
	}

	t.Run("Test_New", testFunction)
}