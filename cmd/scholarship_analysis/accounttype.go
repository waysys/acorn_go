// ----------------------------------------------------------------------------
//
// Account defines whether the scholarship is for:
//   o Associate
//   o Dependent
//   o IndividualGrant
//
//
// Author: William Shaffer
//
// Copyright (c) 2026 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

package main

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type AccountType string

type Tag string

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

const (
	Associate       AccountType = "Associate"
	Dependent       AccountType = "Dependent"
	IndividualGrant AccountType = "IndividualGrant"
	UnknownAccount  AccountType = "Unknown"
)

const (
	dependent  Tag = "Dependent"
	grant      Tag = "Grant"
	individual Tag = "Individual"
	blank      Tag = ""
)

var tagTable = map[Tag]AccountType{
	dependent:  Dependent,
	grant:      Associate,
	individual: IndividualGrant,
	blank:      Associate,
}

// ----------------------------------------------------------------------------
// Function
// ----------------------------------------------------------------------------

func DetermineAccountType(tag Tag) AccountType {
	if accountType, ok := tagTable[tag]; ok {
		return accountType
	}
	return UnknownAccount
}
