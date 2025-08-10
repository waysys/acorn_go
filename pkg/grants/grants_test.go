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

package grants

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	a "acorn_go/pkg/accounting"
	"os"
	"strconv"
	"testing"

	r "github.com/waysys/waydate/pkg/daterange"

	q "acorn_go/pkg/quickbooks"

	d "github.com/waysys/waydate/pkg/date"

	dec "github.com/shopspring/decimal"
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

// Test_NewTransaction tests the creation of a new transaction.
func Test_NewTransaction(t *testing.T) {
	var transType = Grant
	const recipientName = "Jone, Jack"
	var recipient q.Recipient = q.NewRecipient(recipientName)
	const edInstName = "Wake Tech"
	var edInst q.Vendor = q.NewVendor(edInstName)
	var amount = dec.NewFromInt(1000)
	var account = "7040"
	var transaction = NewTransaction(awardDate, transType, &recipient, &edInst, amount, account)

	var testFunction = func(t *testing.T) {
		if transaction.Recipient() != recipientName {
			t.Error("recipient not returned correctly: " +
				transaction.Recipient())
		}
		if transaction.TransType() != Grant {
			t.Error("grant was not returned correctly: " +
				transaction.TransType().String())
		}
		if transaction.TransactionDate() != awardDate {
			t.Error("transaction date was not returned correctly: " +
				transaction.transactionDate.String())
		}
		if transaction.EducationalInstitution() != edInstName {
			t.Error("education institution was not returned correctly: " +
				transaction.EducationalInstitution())
		}
	}

	t.Run("Test_NewTransaction", testFunction)
}

// Test_GrantList tests the creation and use of grant lists
func Test_GrantList(t *testing.T) {
	var apTranLIst q.TransList
	var err error
	var billList q.BillList
	var grantList GrantList

	apTranLIst, err = q.ReadAPTransactions()
	if err == nil {
		billList, err = q.ReadBills(&apTranLIst)
	}
	if err == nil {
		grantList, err = AssembleGrantList(&billList, &apTranLIst)
	}

	var testFunction = func(t *testing.T) {
		if err != nil {
			t.Error(err.Error())
		}
		if grantList.Size() < 10 {
			t.Error("grant list is too small: " + strconv.Itoa(grantList.Size()))
		}
	}

	t.Run("Test_GrantList", testFunction)
}

// Test_PaymentCalculation
func Test_PaymentCalculation(t *testing.T) {
	//
	// Create a grant payment
	//
	var transType = GrantPayment
	const recipientName = "Jone, Jack"
	var recipient q.Recipient = q.NewRecipient(recipientName)
	const edInstName = "Wake Tech"
	var edInst q.Vendor = q.NewVendor(edInstName)
	var amount = dec.NewFromInt(1000)
	var transaction = NewTransaction(awardDate, transType, &recipient, &edInst, amount, "")
	var grantList = NewGrantList()
	grantList.Add(&transaction)
	var total = grantList.TotalTransAmount(a.FY2023, transType)

	var testFunction = func(t *testing.T) {
		var anotherAmount = dec.NewFromInt(1000)
		if !anotherAmount.Equal(total) {
			t.Error("total not equal 1000: " + total.String())
		}
	}

	t.Run("Test_GrantList", testFunction)
}
