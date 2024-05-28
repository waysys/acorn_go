// ----------------------------------------------------------------------------
//
// AP Transaction
//
// Author: William Shaffer
// Version: 27-May-2024
//
// Copyright (c) 2027 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

package quickbooks

// BillList manages the list of bills and associated AP transactions.

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	d "acorn_go/pkg/date"
	"acorn_go/pkg/spreadsheet"
	"errors"
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
	billTab  = "Worksheet"
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
		bills: make([]*EducationBill, 200),
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
		var bill, err = processBill(sprdsht, row, transList)
		if err != nil {
			return err
		}
		billList.Add(&bill)
	}
	return nil
}

// processBill process a single row of the bill spreadsheet
func processBill(
	sprdsht *spreadsheet.Spreadsheet,
	row int,
	transList *TransList) (EducationBill, error) {
	var err error = nil
	var bill = EducationBill{}
	var billDate d.Date
	var vendorName string
	var billType string
	var recipientName string
	const (
		columnTransDate     = "Bill date"
		columnVendorName    = "Vendor"
		columnRecipientName = "Memo"
		columnBillType      = "Bills"
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
		billType, err = sprdsht.Cell(row, columnBillType)
	}
	//
	// Find associated AP transaction
	//
	var trans = transList.Find(Bill, vendorName, recipientName, billDate)
	if trans == nil {
		var message = "transaction not found: \n" +
			"vendor: " + vendorName + "\n" +
			"recipient: " + recipientName + "\n" +
			"date: " + billDate.String()
		err = errors.New(message)
	}
	//
	// Create bill
	//
	bill = NewEducationBill(trans, billType)

	return bill, err
}

// ----------------------------------------------------------------------------
// Methods
// ----------------------------------------------------------------------------

// Add adds a bill to the bill list.
func (billList *BillList) Add(bill *EducationBill) {
	billList.bills[billList.count] = bill
	billList.count++
}
