// ----------------------------------------------------------------------------
//
// # paymentcount_test
//
// Tests for createPaymentMap, the PaymentCount add method, and the
// PaymentCount accessors.
//
// Author: William Shaffer
//
// Copyright (c) 2026 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

package main

import (
	"testing"
)

// ----------------------------------------------------------------------------
// Test Helpers
// ----------------------------------------------------------------------------

var (
	schFourYearSpringPart = Scholarship{
		termPeriod:       Spring,
		institutionType:  FourYear,
		accountType:      Dependent,
		enrollmentStatus: PartTime,
	}
	schIndividualGrantSummer = Scholarship{
		termPeriod:       Summer,
		institutionType:  UnknownInstitution,
		accountType:      IndividualGrant,
		enrollmentStatus: UnknownEnrollment,
	}
	schUnknownAccountFall = Scholarship{
		termPeriod:       Fall,
		institutionType:  UnknownInstitution,
		accountType:      UnknownAccount,
		enrollmentStatus: UnknownEnrollment,
	}
)

// ----------------------------------------------------------------------------
// Tests
// ----------------------------------------------------------------------------

func TestCreatePaymentMap(t *testing.T) {
	billCount := map[Scholarship]int{
		schFourYearSpringFull:    2,
		schFourYearSpringPart:    3,
		schFourYearFallPart:      1,
		schTwoYearFallPart:       4,
		schTwoYearPriorFull:      5,
		schIndividualGrantSummer: 6,
	}

	paymentMap := createPaymentMap(billCount)

	if len(paymentMap) != 4 {
		t.Errorf("len(paymentMap) = %d, expected 4", len(paymentMap))
	}

	spring := paymentMap[Spring]
	if spring.DependentScholarshipCount() != 5 {
		t.Errorf("Spring DependentScholarshipCount() = %d, expected 5",
			spring.DependentScholarshipCount())
	}
	if spring.AssociateScholarshipCount() != 0 {
		t.Errorf("Spring AssociateScholarshipCount() = %d, expected 0",
			spring.AssociateScholarshipCount())
	}

	fall := paymentMap[Fall]
	if fall.AssociateScholarshipCount() != 1 {
		t.Errorf("Fall AssociateScholarshipCount() = %d, expected 1",
			fall.AssociateScholarshipCount())
	}
	if fall.DependentScholarshipCount() != 4 {
		t.Errorf("Fall DependentScholarshipCount() = %d, expected 4",
			fall.DependentScholarshipCount())
	}

	priorTerm := paymentMap[PriorTerm]
	if priorTerm.AssociateScholarshipCount() != 5 {
		t.Errorf("PriorTerm AssociateScholarshipCount() = %d, expected 5",
			priorTerm.AssociateScholarshipCount())
	}

	summer := paymentMap[Summer]
	if summer.IndividualGrantCount() != 6 {
		t.Errorf("Summer IndividualGrantCount() = %d, expected 6",
			summer.IndividualGrantCount())
	}
}

func TestCreatePaymentMapEmpty(t *testing.T) {
	paymentMap := createPaymentMap(map[Scholarship]int{})
	if paymentMap == nil {
		t.Error("createPaymentMap returned nil map")
	}
	if len(paymentMap) != 0 {
		t.Errorf("len(paymentMap) = %d, expected 0", len(paymentMap))
	}
}

func TestCreatePaymentMapUnknownAccount(t *testing.T) {
	billCount := map[Scholarship]int{
		schUnknownAccountFall: 7,
	}

	paymentMap := createPaymentMap(billCount)

	fall := paymentMap[Fall]
	if fall.AssociateScholarshipCount() != 0 ||
		fall.DependentScholarshipCount() != 0 ||
		fall.IndividualGrantCount() != 0 {
		t.Errorf("Fall counts = %+v, expected all zero for unknown account type", fall)
	}
}

func TestPaymentCountAdd(t *testing.T) {
	tests := []struct {
		name              string
		accountType       AccountType
		count             int
		expectedAssociate int
		expectedDependent int
		expectedGrant     int
	}{
		{"associate", Associate, 3, 3, 0, 0},
		{"dependent", Dependent, 4, 0, 4, 0},
		{"individual grant", IndividualGrant, 5, 0, 0, 5},
		{"unknown account adds nothing", UnknownAccount, 6, 0, 0, 0},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			paymentCount := PaymentCount{}.add(test.accountType, test.count)
			if paymentCount.AssociateScholarshipCount() != test.expectedAssociate {
				t.Errorf("AssociateScholarshipCount() = %d, expected %d",
					paymentCount.AssociateScholarshipCount(), test.expectedAssociate)
			}
			if paymentCount.DependentScholarshipCount() != test.expectedDependent {
				t.Errorf("DependentScholarshipCount() = %d, expected %d",
					paymentCount.DependentScholarshipCount(), test.expectedDependent)
			}
			if paymentCount.IndividualGrantCount() != test.expectedGrant {
				t.Errorf("IndividualGrantCount() = %d, expected %d",
					paymentCount.IndividualGrantCount(), test.expectedGrant)
			}
		})
	}
}

func TestPaymentCountAddAccumulates(t *testing.T) {
	paymentCount := PaymentCount{}.
		add(Associate, 2).
		add(Associate, 3).
		add(Dependent, 1)

	if paymentCount.AssociateScholarshipCount() != 5 {
		t.Errorf("AssociateScholarshipCount() = %d, expected 5",
			paymentCount.AssociateScholarshipCount())
	}
	if paymentCount.DependentScholarshipCount() != 1 {
		t.Errorf("DependentScholarshipCount() = %d, expected 1",
			paymentCount.DependentScholarshipCount())
	}
	if paymentCount.IndividualGrantCount() != 0 {
		t.Errorf("IndividualGrantCount() = %d, expected 0",
			paymentCount.IndividualGrantCount())
	}
}

func TestPaymentCountAccessors(t *testing.T) {
	paymentCount := PaymentCount{
		associateScholarshipCount: 1,
		dependentScholarshipCount: 2,
		individualGrantCount:      3,
	}

	if paymentCount.AssociateScholarshipCount() != 1 {
		t.Errorf("AssociateScholarshipCount() = %d, expected 1",
			paymentCount.AssociateScholarshipCount())
	}
	if paymentCount.DependentScholarshipCount() != 2 {
		t.Errorf("DependentScholarshipCount() = %d, expected 2",
			paymentCount.DependentScholarshipCount())
	}
	if paymentCount.IndividualGrantCount() != 3 {
		t.Errorf("IndividualGrantCount() = %d, expected 3",
			paymentCount.IndividualGrantCount())
	}
}
