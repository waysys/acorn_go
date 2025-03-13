// ----------------------------------------------------------------------------
//
// Recipient Summary List
//
// Author: William Shaffer
// Version: 27-Aug-2024
//
// Copyright (c) 2024 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

package grants

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	"errors"
	"maps"
	"slices"
	"sort"

	"github.com/waysys/assert/assert"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type RecipientSumList struct {
	sums map[string]*RecipientSum
}

// ----------------------------------------------------------------------------
// Factory Function
// ----------------------------------------------------------------------------

// NewRecipientSumList returns a recipient summary list with an empty list.
func NewRecipientSumList() RecipientSumList {
	var list = RecipientSumList{
		sums: make(map[string]*RecipientSum),
	}
	return list
}

// AssembleRecipientSumList returns a recipient summary list with the list
// filled from the data in the grant list.
func AssembleRecipientSumList(grantList *GrantList) (RecipientSumList, error) {
	var err error = nil
	var list = NewRecipientSumList()
	var numTrans = grantList.Size()
	for index := 0; index < numTrans; index++ {
		var transaction = grantList.Get(index)
		if transaction.IsPayment() {
			err = processPaymentForSum(&list, transaction)
			if err != nil {
				break
			}
		}
		if transaction.IsRefund() {
			err = processRefundForSum(&list, transaction)
			if err != nil {
				break
			}
		}
	}
	return list, err
}

// processPaymentForSum creates a recipient summary if one does not already
// exists.  It then adds the payment to the total payments for the fiscal year.
func processPaymentForSum(list *RecipientSumList, transaction *Transaction) error {
	//
	// Create a new recipient summary if one does not alreaady exist
	//
	var err error = nil
	var name = transaction.Recipient()
	var recipientSum = list.FetchSummary(name)
	//
	// Add payment amount to the total for the fiscal year
	//
	assert.Assert(list.Contains(name), "recipient summary is not present: "+name)
	var fiscalYear = transaction.FiscalYear()
	var amount = transaction.Amount()
	recipientSum.AddPayment(fiscalYear, amount)
	return err
}

// processRefundForSum adds the refund to the total refunds for the fiscal year.
func processRefundForSum(list *RecipientSumList, transaction *Transaction) error {
	//
	// Fetch summary
	//
	var err error = nil
	var name = transaction.Recipient()
	assert.Assert(list.Contains(name), "recipient summary is not present: "+name)
	var recipientSum = list.FetchSummary(name)
	var fiscalYear = transaction.FiscalYear()
	//
	// Add refund to total for the fiscal year
	//
	var amount = transaction.Amount()
	recipientSum.AddRefund(fiscalYear, amount)
	return err
}

// ----------------------------------------------------------------------------
// Methods
// ----------------------------------------------------------------------------

// Contains returns true if the recipient summary list contains a recipient
// summary for the recipient with the specified name.
func (sumList *RecipientSumList) Contains(name string) bool {
	var _, exists = sumList.sums[name]
	return exists
}

// Add adds a recipient summary to the list
func (sumList *RecipientSumList) Add(recipientSum *RecipientSum) {
	var name = recipientSum.RecipientName()
	if !sumList.Contains(name) {
		sumList.sums[name] = recipientSum
	}
	assert.Assert(sumList.Contains(name), "RecipientSum for recipient was not added: "+name)
}

// Get returns the recipient summary with the specified name.
func (sumList *RecipientSumList) Get(name string) (*RecipientSum, error) {
	var err error = nil
	var recipientSum, exists = sumList.sums[name]
	if !exists {
		err = errors.New("Recipient summary for recipient does not exist: " + name)
	}
	return recipientSum, err
}

// Names returns an alphabetized list of recipient names
func (sumList *RecipientSumList) Names() []string {
	var names = slices.Collect(maps.Keys(sumList.sums))
	sort.Strings(names)
	assert.Assert(sort.StringsAreSorted(names), "names are not sorted")
	return names
}

// Return an instance of RecipientSum.  If there is already a recipient summary in the
// list with the specified name, return that summary.  Otherwise create a new
// summary and add it to the list.
func (sumList *RecipientSumList) FetchSummary(name string) *RecipientSum {
	var recipientSum *RecipientSum

	if !sumList.Contains(name) {
		var summary = NewRecipientSum(name)
		recipientSum = &summary
		sumList.Add(recipientSum)
	} else {
		recipientSum, _ = sumList.Get(name)
	}
	return recipientSum
}
