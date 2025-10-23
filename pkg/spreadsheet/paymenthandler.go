// ----------------------------------------------------------------------------
//
// Payment Handler
//
// This code handles the writing of cells to a spreadsheet based on the
// fiscal years available.  It is used when the output spreadsheet has columns
// of data for each fiscal year.
//
// Author: William Shaffer
//
// Copyright (c) 2025 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

package spreadsheet

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	a "acorn_go/pkg/accounting"

	dec "github.com/shopspring/decimal"
)

// ----------------------------------------------------------------------------
// Support Functions
// ----------------------------------------------------------------------------

// Create headings for recipient summary
func CreateHeadings(output *SpreadsheetFile,
	row int,
	startColumn string,
	heading string) {
	var column string
	var label string
	var ch = NewColumnHandler(startColumn)
	//
	// Place title
	//
	WriteCell(output, "A", row, heading)
	row += 2
	//
	// Place headings
	//
	WriteCell(output, "A", row, "Recipient Name")
	for fy := range a.FYIndicators {
		column = ch.CountColumn((a.FYIndicator(fy)))
		label = ch.CountColumnLabel(a.FYIndicator(fy))
		WriteCell(output, column, row, label)
		column = ch.PaymentColumn((a.FYIndicator(fy)))
		label = ch.PaymentColumnLabel((a.FYIndicator(fy)))
		WriteCell(output, column, row, label)
	}
	column = ch.TotalColumn()
	WriteCell(output, column, row, "Total Count")
}

// outputSumData inserts the recipient summary data into spreadsheet.
func OutputSumData(
	output *SpreadsheetFile,
	fy a.FYIndicator,
	row int,
	ch *ColumnHandler,
	count int,
	amount dec.Decimal) {
	var column string

	column = ch.CountColumn(fy)
	WriteCellInt(output, column, row, count)
	column = ch.PaymentColumn(fy)
	WriteCellDecimal(output, column, row, amount)
	column = ch.TotalColumn()
	WriteCellInt(output, column, row, 1)
}

// outputTotals inserts the totals for all recipients on a row
func OutputTotals(
	output *SpreadsheetFile,
	row int,
	ch *ColumnHandler,
	counts []int,
	amounts []dec.Decimal,
	grandCount int) {
	var column string

	WriteCell(output, "A", row, "Total Payments")
	for _, fy := range a.FYIndicators {
		column = ch.CountColumn(fy)
		WriteCellInt(output, column, row, counts[fy])
		column = ch.PaymentColumn(fy)
		WriteCellDecimal(output, column, row, amounts[fy])
	}
	column = ch.TotalColumn()
	WriteCellInt(output, column, row, grandCount)
}
