// ----------------------------------------------------------------------------
//
// Recipient Test
//
// Author: William Shaffer
// Version: 12-May-2024
//
// Copyright (c) 2024 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

package recipient

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	a "acorn_go/pkg/accounting"
	d "acorn_go/pkg/date"
	r "acorn_go/pkg/daterange"
	"os"
	"testing"
)

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

var awardDate, _ = d.New(7, 28, 2023)
var endDate, _ = d.New(8, 31, 2023)
var awardRange, _ = r.New(awardDate, endDate)

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

// Test_NewAwardGroup checks the factory function for creating a new award group.
func Test_NewAwardGroup(t *testing.T) {
	var err error = nil
	var group AwardGroup
	var name = "Spring 2023"

	var testFunction = func(t *testing.T) {
		group, err = NewAwardGroup(name, awardDate, awardRange)
		if err != nil {
			t.Error("unable to create new award group")
		}
		if group.AwardDate() != awardDate {
			t.Error("AwardDate() does not equal award date")
		}
		if group.GroupName() != name {
			t.Error("GroupName() does not match name: " + group.GroupName())
		}
		if !group.InAwardGroup(awardDate) {
			t.Error("award date is not considered in award range")
		}
		if group.FiscalYear() != a.FY2023 {
			t.Error("fiscal year not determined correctly")
		}
	}

	t.Run("Test_NewAwardGroup", testFunction)

}
