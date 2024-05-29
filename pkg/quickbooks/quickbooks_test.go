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
	"os"
	"strconv"
	"testing"

	d "github.com/waysys/waydate/pkg/date"

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
	var transactionType = Bill
	var amount Money = Money(dec.NewFromInt(100))
	var account = "7050 Grants, contracts, & direct assistance:Awards & grants - individuals"
	var transaction = NewAPTransaction(transactionDate, vendorName, recipientName, transactionType, amount, account)
	var transPtr = &transaction

	var testFunction = func(t *testing.T) {
		if transPtr.transactionDate != transactionDate {
			t.Error("transaction date did not matching origin")
		}
		if transPtr.Vendor().Name() != vendorName {
			t.Error("vendor name is not correct: " + transPtr.Vendor().Name())
		}
		if transPtr.Recipient().Name() != recipientName {
			t.Error("recipient name is not correct: " + transPtr.Recipient().Name())
		}
		if transPtr.Amount() != amount {
			t.Error("amount of transaction is incorrect: " + dec.Decimal(transPtr.Amount()).String())
		}
		if !transPtr.GTZero() {
			t.Error("amount of transaction is not considered greater than zero")
		}
		if !transPtr.IsBill() {
			t.Error("transaction is not considered a valid bill")
		}
	}

	t.Run("Test_NewAPTransaction", testFunction)
}

// Test_ReadAPTransaction tests the ReadAPTransaction function.
func Test_ReadAPTransaction(t *testing.T) {
	var transList, err = ReadAPTransactions()
	if err != nil {
		t.Error(err.Error())
	}
	var size = (&transList).Size()
	if size < 50 {
		t.Error("transList is too small: " + strconv.Itoa(size))
	}
}
