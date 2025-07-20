// ----------------------------------------------------------------------------
//
// AP Transaction
//
// Author: William Shaffer
// Version: 27-May-2024
//
// Copyright (c) 2024 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

package quickbooks

// BillList manages the list of bills and associated AP transactions.

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	"acorn_go/pkg/spreadsheet"
	"strconv"

	"github.com/waysys/assert/assert"
	d "github.com/waysys/waydate/pkg/date"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type BillList struct {
	bills []*EducationBill
	count int
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

const (
	billFile = "/home/bozo/golang/acorn_go/data/bills.xlsx"
	billTab  = "Sheet1"
)

// ----------------------------------------------------------------------------
// Factory Function
// ----------------------------------------------------------------------------

// ReadBills produces a bill list by reading the bills.xlsx spreadsheet and
// the transaction list.
func ReadBills(transList *TransList) (BillList, error) {
	var sprdsht spreadsheet.Spreadsheet
	var err error
	var billList = BillList{
		bills: make([]*EducationBill, 1000),
		count: 0,
	}
	//
	// Obtain spreadsheet data
	//
	sprdsht, err = spreadsheet.ProcessData(billFile, billTab)
	if err == nil {
		err = processBills(&sprdsht, &billList, transList)
	}
	return billList, err
}

// processBills reads the spreadsheet and creates bills.
func processBills(sprdsht *spreadsheet.Spreadsheet,
	billList *BillList,
	transList *TransList) error {
	var numRows = sprdsht.Size()
	//
	// Loop through the spreadsheet
	//
	for row := 1; row < numRows; row++ {
		var bill, err, okToUse = processBill(sprdsht, row, transList)
		if err != nil {
			return err
		}
		if okToUse {
			billList.Add(&bill)
		}
	}
	return nil
}

// processBill process a single row of the bill spreadsheet
func processBill(
	sprdsht *spreadsheet.Spreadsheet,
	row int,
	transList *TransList) (EducationBill, error, bool) {
	var err error = nil
	var bill = EducationBill{}
	var billDate d.Date
	var vendorName string
	var billType string
	var recipientName string
	var trans *APTransaction
	var okToUse bool

	const (
		columnTransDate     = "Date"
		columnVendorName    = "Vendor"
		columnRecipientName = "Memo"
		columnBillType      = "Bills"
		columnItemType      = "Item class"
	)
	//
	// Read data from spreadsheet
	//
	billDate, err = sprdsht.CellDate(row, columnTransDate)
	if err == nil {
		vendorName, err = sprdsht.Cell(row, columnVendorName)
	}
	if err == nil {
		recipientName, err = sprdsht.Cell(row, columnRecipientName)
	}
	if err == nil {
		recipientName = ConvertName(recipientName)
		billType, err = sprdsht.Cell(row, columnBillType)
	}
	//
	// Find associated AP transaction
	//
	if err == nil {
		trans = transList.Find(Bill, vendorName, recipientName, billDate)
	} else {
		trans = nil
	}

	//
	// Create bill
	//
	if trans != nil {
		bill = NewEducationBill(trans, billType)
		okToUse = true
	} else {
		bill = EducationBill{}
		okToUse = false
	}
	return bill, err, okToUse
}

// ----------------------------------------------------------------------------
// Methods
// ----------------------------------------------------------------------------

// Add adds a bill to the bill list.
func (billList *BillList) Add(bill *EducationBill) {
	billList.bills[billList.count] = bill
	billList.count++
}

// Size returns the number of bills in the list
func (billList *BillList) Size() int {
	return billList.count
}

// Get returns the bill at the specified index.
func (billList *BillList) Get(index int) *EducationBill {
	assert.Assert(0 <= index && index < billList.count, "invalid index: "+
		strconv.Itoa(index))

	return billList.bills[index]
}
