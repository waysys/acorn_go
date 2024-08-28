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

	dec "github.com/shopspring/decimal"
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
		if transaction.transType == GrantPayment && transaction.Amount().GreaterThanOrEqual(dec.Zero) {
			err = processTransactionForSum(&list, transaction)
			if err != nil {
				break
			}
		}
	}
	return list, err
}

// processTransactions creates a recipient summary if one does not already
// exists.  It then adds the payment to the total payments for the fiscal year.
func processTransactionForSum(list *RecipientSumList, transaction *Transaction) error {
	//
	// Create a new recipient summary if one does not alreaady exist
	//
	var err error = nil
	var recipientSum *RecipientSum
	var name = transaction.Recipient()
	if !list.Contains(name) {
		var summary = NewRecipientSum(name)
		recipientSum = &summary
		list.Add(recipientSum)
	} else {
		recipientSum, err = list.Get(name)
	}
	//
	// Add payment amount to the total for the fiscal year
	//
	if err == nil {
		assert.Assert(list.Contains(name), "recipient summary is not present: "+name)
		var fiscalYear = transaction.FiscalYear()
		var amount = transaction.Amount()
		recipientSum.AddPayment(fiscalYear, amount)
	}
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
