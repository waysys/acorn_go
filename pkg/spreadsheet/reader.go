// ----------------------------------------------------------------------------
//
// Spreadsheet reader
//
// Author: William Shaffer
// Version: 15-Apr-2024
//
// Copyright (c) 2024 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

// The spreadsheet package reads and processes Excel spreadsheets.
package spreadsheet

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	dec "github.com/shopspring/decimal"
	"github.com/xuri/excelize/v2"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type Spreadsheet struct {
	headings []string
	rows     [][]string
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// Functions
// ----------------------------------------------------------------------------

// readData extracts the donation data from Excel file and
// returns an array of string arrays.
func readData(excelFillName string, tab string) ([][]string, error) {
	var err error
	var file *excelize.File
	var rows [][]string
	//
	// Open file
	//
	file, err = excelize.OpenFile(excelFillName)
	if err != nil {
		return rows, err
	}
	//
	// Function to close file
	//
	defer func() {
		// Close the spreadsheet.
		if err := file.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	//
	// Read rows
	//
	rows, err = file.GetRows(tab)
	if err != nil {
		return rows, err
	}
	return rows, nil
}

// ProcessData reads the donation Excel file and returns the column headings
// and a slice of the data in a Spreadsheet structure.
func ProcessData(fileName string, tab string) (Spreadsheet, error) {
	var rows [][]string
	var err error
	var spreadsheet Spreadsheet
	//
	// Retieve data form .xlsx file
	//
	rows, err = readData(fileName, tab)
	if err != nil {
		return spreadsheet, err
	}
	if len(rows) == 0 {
		err = errors.New("spreadsheet is empty")
	}
	//
	// Form heading.  The first row in the spreadsheet must contain the headings.
	//
	spreadsheet.headings = rows[0]
	spreadsheet.rows = rows
	return spreadsheet, err
}

// ----------------------------------------------------------------------------
// Methods
// ----------------------------------------------------------------------------

// Size returns the number of rows in the spreadsheet, including the header
// row.
func (spreadsheet *Spreadsheet) Size() int {
	return len(spreadsheet.rows)
}

// Column returns an integer indicating the column position of a string
// in the header.  If the string is not found in the header, an error
// is returned.
func (spreadsheet *Spreadsheet) column(heading string) (int, error) {
	var err error
	var column = 0

	if heading == "" {
		err = errors.New("heading must not be empty")
		return column, err
	}

	for column = 0; column < len(spreadsheet.headings); column++ {
		if spreadsheet.headings[column] == heading {
			return column, nil
		}
	}
	err = errors.New("heading not found in headings: " + heading)
	return column, err
}

// Cell returns the value in a cell of the spreadsheet.
func (spreadsheet *Spreadsheet) Cell(row int, heading string) (string, error) {
	var column int
	var err error
	var cell = ""

	column, err = spreadsheet.column(heading)
	if err != nil {
		return cell, err
	}
	if row < 1 || row >= spreadsheet.Size() {
		err = errors.New("invalid row for spreadsheet: " + strconv.Itoa(row))
		return cell, err
	}
	cell = spreadsheet.rows[row][column]
	cell = strings.TrimSpace(cell)
	return cell, nil
}

// CellDecimal returns the value in the cell of the spreadsheet.
func (spreadsheet *Spreadsheet) CellDecimal(row int, heading string) (dec.Decimal, error) {
	var value string
	var err error
	var amount dec.Decimal = dec.Zero

	value, err = spreadsheet.Cell(row, heading)
	if err == nil {
		amount, err = dec.NewFromString(value)
	}

	return amount, err
}
