// ----------------------------------------------------------------------------
//
// Guest
//
// Author: William Shaffer
// Version: 29-Jun-2024
//
// Copyright (c) 2024 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

// The guestlist package reads the guestlist.xlmx spreadsheet exported from
// paperless posts.  It builds a map with the guest information including whether
// they have responded to the email invitation.

package guestlist

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	a "acorn_go/pkg/address"
	d "acorn_go/pkg/donors"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type Guest struct {
	d.Donor
	status   Status
	isMember bool
}

// ----------------------------------------------------------------------------
// Factory Function
// ----------------------------------------------------------------------------

// New returns a guest based on the specified inputs.
func New(nm string, adr a.Address, eml string, sts string) Guest {
	guest := Guest{
		Donor:    d.New(nm, nm, adr, eml, 0, false),
		status:   NewStatus(sts),
		isMember: false,
	}
	return guest
}

// ----------------------------------------------------------------------------
// Methods
// ----------------------------------------------------------------------------

// Status returns the status of the guest
func (guest Guest) Status() Status {
	return guest.status
}

// NoResponse is true if the guest has not responded to the invitation
func (guest Guest) NoResponse() bool {
	var status = guest.Status()
	var result = (status != Attending) && (status != Regrets) && (status != Unsubscribed)
	return result
}

// SetMember sets the isMember field to true.
func (guest *Guest) SetMember() {
	guest.isMember = true
}

// IsMember returns true if the guest is a member as opposed to management.
func (guest Guest) IsMember() bool {
	return guest.isMember
}
