// ----------------------------------------------------------------------------
//
// Grant List
//
// Author: William Shaffer
// Version: 30-May-2024
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
	q "acorn_go/pkg/quickbooks"
	"sort"
	"strconv"

	dec "github.com/shopspring/decimal"
	"github.com/waysys/assert/assert"
	s "github.com/waysys/set/pkg/stringset"
	d "github.com/waysys/waydate/pkg/date"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type GrantList struct {
	trans []*Transaction
	count int
}

// ----------------------------------------------------------------------------
// Factory Function
// ----------------------------------------------------------------------------

// NewGrantList returns an initialized grant list
func NewGrantList() GrantList {
	var grantList = GrantList{
		trans: make([]*Transaction, 1000),
		count: 0,
	}

	return grantList
}

// ----------------------------------------------------------------------------
// Factory Function for Scholarships
// ----------------------------------------------------------------------------

// AssembleGrantList populates the grant list with bill and AP transaction data.
func AssembleGrantList(billList *q.BillList, transList *q.TransList) (GrantList, error) {
	var err error = nil
	//
	// Create grant list
	//
	var grantList = NewGrantList()
	//
	// Add grants and transfers
	//
	processBills(billList, &grantList)
	//
	// Add write-offs and payments
	//
	processOtherTransactions(transList, &grantList)

	return grantList, err
}

// processBills cycles through bills and populates the grants and transfers
func processBills(billList *q.BillList, grantList *GrantList) {
	var numBills = billList.Size()
	var bill *q.EducationBill
	var date d.Date
	var amount dec.Decimal
	var recipient *q.Recipient
	var edInst *q.Vendor
	var transType TransType
	var account string
	//
	// Cycle through list of bills
	//
	for index := range numBills {
		//
		// Extract data
		//
		bill = billList.Get(index)
		date = bill.TransactionDate()
		amount = bill.Amount()
		recipient = bill.Recipient()
		edInst = bill.Vendor()
		transType = convertBillType(bill.BType())
		account = bill.Account()
		//
		// Create transaction
		//
		var transaction = NewTransaction(
			date,
			transType,
			recipient,
			edInst,
			amount,
			account,
		)
		//
		// Add transaction
		//
		grantList.Add(&transaction)
	}
}

// processOtherTransactions cylces through the TransList and adds them to
// the grant list.
func processOtherTransactions(transList *q.TransList, grantList *GrantList) {
	var numTrans = transList.Size()
	var apTrans *q.APTransaction
	//
	// Cycle through the transaction list
	//
	for index := range numTrans {
		apTrans = transList.Get(index)
		if selectTransaction(apTrans) {
			var transaction = processTransaction(apTrans)
			grantList.Add(&transaction)
		}
	}
}

// selectTransaction returns true if the transaction should be included in the
// grant list.  Bills and transfers are already included, This function
// returns true for transactions that are not bills or transfers.
func selectTransaction(apTrans *q.APTransaction) bool {
	var result = apTrans.IsVendorCredit() || apTrans.IsPayment() || apTrans.IsDeposit()
	return result
}

// processTransaction extracts data from the AP transaction and populates
// the grant transaction
func processTransaction(apTrans *q.APTransaction) Transaction {
	//
	// Extract the data
	//
	var date = apTrans.TransactionDate()
	var transType = convertTransType(apTrans)
	var recipient = apTrans.Recipient()
	var vendor = apTrans.Vendor()
	var amount = dec.Decimal(apTrans.Amount())
	var account = apTrans.Account()
	//
	// Create transaction
	//
	var transaction = NewTransaction(
		date,
		transType,
		recipient,
		vendor,
		amount,
		account,
	)
	return transaction
}

// convertTransType computes the Grant transaction from the AP transaction
func convertTransType(apTrans *q.APTransaction) TransType {
	var transType = Grant
	switch {
	case apTrans.IsBill():
		transType = Grant
	case apTrans.IsVendorCredit():
		transType = WriteOff
	case apTrans.IsPayment():
		transType = GrantPayment
	case apTrans.IsDeposit():
		transType = Refund
	default:
		assert.Assert(false, "unrecognized grant transaction type")
	}
	return transType
}

// convertTranType converts the bill type to the grant transacton type
func convertBillType(billType q.BillType) TransType {
	var tranType TransType
	switch billType {
	case q.Grant:
		tranType = Grant
	case q.Transfer:
		tranType = Transfer
	case q.Individual:
		tranType = Grant
	case q.Dependent:
		tranType = Grant
	default:
		assert.Assert(false, "Unknown bill type: "+strconv.Itoa(int(billType)))
	}
	return tranType
}

// ----------------------------------------------------------------------------
// Factory Methods for Individual Grants
// ----------------------------------------------------------------------------

// AssembleIndividualGrantList creates a grant list for individual grants.
func AssembleIndividualGrantList(transList *q.TransList) (GrantList, error) {
	var err error = nil
	var apTrans *q.APTransaction
	//
	// Create grant list
	//
	var grantList = NewGrantList()
	//
	// Cycle through the accounts payable transaction list
	//
	var numTrans = transList.Size()
	for index := range numTrans {
		apTrans = transList.Get(index)
		if apTrans.IsIndividualGrant() {
			var transaction = processIndividualGrant(apTrans)
			grantList.Add(&transaction)
		}
	}
	return grantList, err
}

// processIndividualGrant returns a grant transaction based on the
// accounts payable transaction
func processIndividualGrant(apTrans *q.APTransaction) Transaction {
	//
	// Extract the data
	//
	var date = apTrans.TransactionDate()
	var transType = convertTransType(apTrans)
	var vendor = apTrans.Vendor()
	var vendorName = vendor.Name()
	vendorName = q.ConvertName(vendorName)
	var recipient = q.NewRecipient(vendorName)
	var amount = dec.Decimal(apTrans.Amount())
	var account = apTrans.Account()
	//
	// Create transaction
	//
	var transaction = NewTransaction(
		date,
		transType,
		&recipient,
		vendor,
		amount,
		account,
	)
	return transaction
}

// ----------------------------------------------------------------------------
// Methods
// ----------------------------------------------------------------------------

// Add adds a grant transaction to the grant list.
func (grantList *GrantList) Add(transaction *Transaction) {
	grantList.trans[grantList.count] = transaction
	grantList.count++
}

// Size returns the number of transaction
func (grantList *GrantList) Size() int {
	return grantList.count
}

// Get returns the transaction at the position indicated by the index
func (grantList *GrantList) Get(index int) *Transaction {
	assert.Assert(index >= 0 && index < grantList.count,
		"index is out of range in grant list: "+strconv.Itoa(index))
	return grantList.trans[index]
}

// TotalTransAmount returns the total amount of transactions for a fiscal year
// and transaction type
func (grantList *GrantList) TotalTransAmount(fiscalYear a.FYIndicator, transType TransType) dec.Decimal {
	var total dec.Decimal = dec.Zero
	var numTrans = grantList.Size()
	var transaction Transaction

	var match = func(tran Transaction) bool {
		var transactionDate = tran.TransactionDate()
		var tt = tran.TransType()
		var result = false
		if tt == transType {
			var fyIndicator = a.FiscalYearIndicator(transactionDate)
			result = fyIndicator == fiscalYear
		}
		return result
	}

	for index := range numTrans {
		transaction = *grantList.Get(index)
		if match(transaction) {
			total = total.Add(transaction.Amount())
		}
	}
	return total
}

// GrandTotalTransactions returns the total amount of transactions of a particular type
// of transaction for all fiscal years.
func (grantList *GrantList) GrandTotalTransactions(transType TransType) dec.Decimal {
	var total = dec.Zero
	for _, fy := range a.FYIndicators {
		var amount = grantList.TotalTransAmount(fy, transType)
		total = total.Add(amount)
	}
	return total
}

// TotalNetWriteOff returns the net writeoff which is the gross writeoff minus the transfers
func (grantList *GrantList) TotalNetWriteOff(fiscalYear a.FYIndicator) dec.Decimal {
	var totalWriteOffs = grantList.TotalTransAmount(fiscalYear, WriteOff)
	var totalTransfers = grantList.TotalTransAmount(fiscalYear, Transfer)
	var totalNetWriteOff = totalWriteOffs.Sub(totalTransfers)
	return totalNetWriteOff
}

// GrandTotalNetWriteoff returns the total net writeoffs.
func (grantList *GrantList) GrandTotalNetWriteoff() dec.Decimal {
	var total = dec.Zero
	for _, fy := range a.FYIndicators {
		var amount = grantList.TotalNetWriteOff(fy)
		total = total.Add(amount)
	}
	return total
}

// NetBlance returns the net balance for the specified fiscal year.  The net balance is:
//
//	grants - payments - writeoff + transfers + refunds
func (grantList *GrantList) NetBalance(fiscalYear a.FYIndicator) dec.Decimal {
	var net = grantList.TotalTransAmount(fiscalYear, Grant)
	net = net.Sub(grantList.TotalTransAmount(fiscalYear, GrantPayment))
	net = net.Sub(grantList.TotalTransAmount(fiscalYear, WriteOff))
	net = net.Add(grantList.TotalTransAmount(fiscalYear, Transfer))
	net = net.Add(grantList.TotalTransAmount(fiscalYear, Refund))
	return net
}

// TotalNetBalance returns the sum of net balances for all fiscal years.
func (grantList *GrantList) TotalNetBalance() dec.Decimal {
	var total = dec.Zero
	for _, fy := range a.FYIndicators {
		total = total.Add(grantList.NetBalance(fy))
	}
	return total
}

// SortGrantList returns the grant list sorted in this order:
//
//	Recipient
//	Educational Institution
//	Transaction Date
func (grantList *GrantList) Sort() {
	var less = func(i int, j int) bool {
		var trans1 = grantList.Get(i)
		var trans2 = grantList.Get(j)
		var result = compare(trans1, trans2)
		return result
	}
	sort.Slice(grantList.trans[:grantList.count], less)
}

// Names returns an alphabetized name of recipients.
func (grantList *GrantList) Names() []string {
	var names = s.New()
	var index = 0
	var numTrans = grantList.Size()

	for index = range numTrans {
		var transaction = grantList.Get(index)
		var name = transaction.Recipient()
		names.Add(name)
	}

	var nameList = names.ToSlice()
	sort.Strings(nameList)
	return nameList
}
