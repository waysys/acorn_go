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

var sortOrderAccountType = map[AccountType]int{
	Associate:       0,
	Dependent:       1,
	IndividualGrant: 2,
	UnknownAccount:  3,
}

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

// DetermineAccountType returns the account type for a given tag
func DetermineAccountType(tag Tag) AccountType {
	if accountType, ok := tagTable[tag]; ok {
		return accountType
	}
	return UnknownAccount
}

// CompareAccountType returns the sort indicator
func CompareAccountType(acctType1 AccountType, acctType2 AccountType) int {
	var order1 = sortOrderAccountType[acctType1]
	var order2 = sortOrderAccountType[acctType2]

	var result = 0
	if order1 < order2 {
		result = -1
	} else if order1 > order2 {
		result = 1
	}

	return result
}
