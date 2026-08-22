// ----------------------------------------------------------------------------
//
// paymentcount.go calculates the number of associate scholarships, dependent
// scholarships, and individual grants by term period.
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

type PaymentCount struct {
	associateScholarshipCount int
	dependentScholarshipCount int
	individualGrantCount      int
}

// ----------------------------------------------------------------------------
// Functions
// ----------------------------------------------------------------------------

// createPaymentMap creates the map of payment counts using the billCount
func createPaymentMap(billCount map[Scholarship]int) map[TermPeriod]PaymentCount {
	var paymentMap = make(map[TermPeriod]PaymentCount)

	for scholarship, count := range billCount {
		var termPeriod = scholarship.TermPeriod()
		var accountType = scholarship.AccountType()
		var paymentCount = paymentMap[termPeriod]
		paymentCount = paymentCount.add(accountType, count)
		paymentMap[termPeriod] = paymentCount
	}
	return paymentMap
}

// ----------------------------------------------------------------------------
// Methods
// ----------------------------------------------------------------------------

// add adds a count to the appropriate attribute of the payment count
func (paymentCount PaymentCount) add(
	accountType AccountType,
	count int,
) PaymentCount {
	switch accountType {
	case Associate:
		paymentCount.associateScholarshipCount += count
	case Dependent:
		paymentCount.dependentScholarshipCount += count
	case IndividualGrant:
		paymentCount.individualGrantCount += count
	}
	return paymentCount
}

// AssociateScholarshipCount returns the number of associate scholarships
func (paymentCount PaymentCount) AssociateScholarshipCount() int {
	return paymentCount.associateScholarshipCount
}

// DependentScholarshipCount returns the number of dependent scholarships
func (paymentCount PaymentCount) DependentScholarshipCount() int {
	return paymentCount.dependentScholarshipCount
}

// IndividualGrantCount returns the number of individual grants
func (paymentCount PaymentCount) IndividualGrantCount() int {
	return paymentCount.individualGrantCount
}
