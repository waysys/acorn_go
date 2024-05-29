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

// This file contains code to read the accounts_payable.xlsx spreadsheet and
// build a slice containing all the transactions.

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	d "acorn_go/pkg/date"
	"acorn_go/pkg/spreadsheet"
	"strconv"

	"github.com/waysys/assert/assert"

	dec "github.com/shopspring/decimal"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type TransList struct {
	trans []*APTransaction
	count int
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

const (
	inputFile = "/home/bozo/golang/acorn_go/data/accounts_payable.xlsx"
	tab       = "Worksheet"
)

var ZERO = Money(dec.Zero)

const (
	columnTransactionDate = "Date"
	columnVendor          = "Payee"
	columnRecipient       = "Memo"
	columnBilled          = "Billed"
	columnPaid            = "Paid"
	columnTransactionType = "Type"
	columnAccount         = "Account"
)

// ----------------------------------------------------------------------------
// Factory Function
// ----------------------------------------------------------------------------

// ReadAPtransactions reads the accounts_payable.xlsx spreadsheet and
// generates the AP transaction list
func ReadAPTransactions() (TransList, error) {
	var sprdsht spreadsheet.Spreadsheet
	var err error
	var transList = TransList{
		trans: make([]*APTransaction, 200),
		count: 0,
	}
	//
	// Obtain spreadsheet data
	//
	sprdsht, err = spreadsheet.ProcessData(inputFile, tab)
	if err == nil {
		err = processTransactions(&sprdsht, &transList)
	}
	return transList, err
}

// processTransaction processes the data in the spreadsheet to populate
// the transaction list.
func processTransactions(
	sprdshtPtr *spreadsheet.Spreadsheet,
	transList *TransList) error {
	var numRows = sprdshtPtr.Size()

	for row := 1; row < numRows; row++ {
		var transaction, err = processTransaction(sprdshtPtr, row)
		if err != nil {
			return err
		}
		if selectTransaction(&transaction) {
			transList.Add(&transaction)
		}
	}
	return nil
}

// selectTransaction returns true if the transaction should be added to
// the transaction list
func selectTransaction(transaction *APTransaction) bool {
	var result = transaction.IsBill() || transaction.IsPayment() || transaction.IsVendorCredit()
	return result
}

// processTransaction processes a single row in the spreadshet
func processTransaction(sprdshtPtr *spreadsheet.Spreadsheet, row int) (APTransaction, error) {
	var err error = nil
	var transaction = APTransaction{}
	var account string
	var transactionDate d.Date
	var vendorName string
	var recipientName string
	var amount Money
	var transactonType QuickbooksTransactionType
	//
	// Read values from spreadsheet
	//
	account, err = sprdshtPtr.Cell(row, columnAccount)
	if err == nil {
		transactionDate, err = sprdshtPtr.CellDate(row, columnTransactionDate)
	}
	if err == nil {
		vendorName, err = sprdshtPtr.Cell(row, columnVendor)
	}
	if err == nil {
		recipientName, err = sprdshtPtr.Cell(row, columnRecipient)
	}
	if err == nil {
		transactonType, err = retrieveType(sprdshtPtr, row)
	}
	if err == nil {
		amount, err = retrieveAmount(sprdshtPtr, row, transactonType)
	}
	//
	// Create transaction
	//
	if err == nil {
		transaction = NewAPTransaction(transactionDate,
			vendorName,
			recipientName,
			transactonType,
			amount,
			account)
	}
	return transaction, err
}

// retrieveType retrieves the transaction type and converts it to
// a QuickbooksTransactionType.
func retrieveType(sprdshtPtr *spreadsheet.Spreadsheet, row int) (QuickbooksTransactionType, error) {
	var transactionType QuickbooksTransactionType = Unknown
	var err error = nil
	var value string

	value, err = sprdshtPtr.Cell(row, columnTransactionType)
	if err == nil {
		transactionType = NewQuickbooksTransactionType(value)
	}

	return transactionType, err
}

// retrieveAmount returns the amount of the transaction based on the
// transaction type.
func retrieveAmount(
	sprdshtPtr *spreadsheet.Spreadsheet,
	row int,
	transactionType QuickbooksTransactionType) (Money, error) {
	var err error = nil
	var amount dec.Decimal = dec.Zero

	switch transactionType {
	case Bill:
		amount, err = sprdshtPtr.CellDecimal(row, columnBilled)
	case BillPayment:
		amount, err = sprdshtPtr.CellDecimal(row, columnPaid)
	case VendorCredit:
		amount, err = sprdshtPtr.CellDecimal(row, columnPaid)
	case Unknown:
		amount = dec.Zero
	}

	return Money(amount), err
}

// ----------------------------------------------------------------------------
// Methods
// ----------------------------------------------------------------------------

// Add adds a transaction to the list of transactions.
func (transList *TransList) Add(trans *APTransaction) {
	transList.trans[transList.count] = trans
	transList.count++
}

// Size returns the number of transactions
func (transList *TransList) Size() int {
	return transList.count
}

// Get returns a transaction at the specified index
func (transList *TransList) Get(index int) *APTransaction {
	assert.Assert(0 <= index && index < transList.count,
		"index is out of range: "+strconv.Itoa(index))
	return transList.trans[index]
}

// Return the transaction that matches the supplied attributes
func (transList *TransList) Find(
	transType QuickbooksTransactionType,
	vendorName string,
	recipientName string,
	transDate d.Date,
) *APTransaction {
	var trans *APTransaction = nil
	var num = transList.count

	var match = func(tran *APTransaction) bool {
		var result = true
		result = result && tran.TransactionType() == transType
		result = result && tran.Vendor().name == vendorName
		result = result && tran.Vendor().name == recipientName
		result = result && tran.TransactionDate() == transDate
		return result
	}

	for index := 0; index < num; num++ {
		var item = transList.Get(index)
		if match(item) {
			trans = item
			break
		}
	}

	return trans
}
