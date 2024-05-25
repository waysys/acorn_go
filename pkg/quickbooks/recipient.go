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
	"acorn_go/pkg/assert"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// Recipient is an Associate who receives a grant.
type Recipient struct {
	name string
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// Factory Functions
// ----------------------------------------------------------------------------

// NewRecipient returns a recipient structure with the specified name.
func NewRecipient(name string) Recipient {
	assert.Assert(name != "", "recipient name must not be an empty string")
	var recipient = Recipient{
		name: name,
	}
	return recipient
}

// ----------------------------------------------------------------------------
// Properties
// ----------------------------------------------------------------------------

// Name returns the name of the recipient.
func (recipient *Recipient) Name() string {
	return recipient.name
}
