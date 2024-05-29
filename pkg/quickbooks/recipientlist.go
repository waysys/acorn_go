// ----------------------------------------------------------------------------
//
// AP Transaction
//
// Author: William Shaffer
// Version: 24-May-2024
//
// Copyright (c) 2024 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

package quickbooks

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	"github.com/waysys/assert/assert"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type RecipientList map[string]*Recipient

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

var APRecipientList = make(RecipientList)

// ----------------------------------------------------------------------------
// Methods
// ----------------------------------------------------------------------------

// Add adds a new recipient to the list if it is not already in the list.
func (recipientList *RecipientList) Add(name string) *Recipient {
	var recipient *Recipient = nil
	if recipientList.Contains(name) {
		recipient = recipientList.Get(name)
	} else {
		var r = NewRecipient(name)
		recipient = &r
		(*recipientList)[name] = recipient
	}
	return recipient
}

// Contains returns true if the recipient list holds a recipient with the
// specified name.
func (recipientList *RecipientList) Contains(name string) bool {
	var _, ok = (*recipientList)[name]
	return ok
}

// Get returns a recipient with the specified name.  The recipient with
// that name must exist in the list.
func (recipientList *RecipientList) Get(name string) *Recipient {
	assert.Assert(recipientList.Contains(name), "recipient must exist in list: "+name)
	return (*recipientList)[name]
}
