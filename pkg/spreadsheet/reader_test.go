// ----------------------------------------------------------------------------
//
// Spreadsheet reader
//
// Author: William Shaffer
// Version: 15-Apr-2024
//
// Copyright (c) William Shaffer
//
// ----------------------------------------------------------------------------

// The spreadsheet package reads and processes Excel spreadsheets.
package spreadsheet

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	"fmt"
	"os"
	"strconv"
	"testing"
)

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

// Test_ReadSpreadsheet checks that the specified spreadsheet can be read.
func Test_ReadSpreadsheet(t *testing.T) {
	var spreadsheet, err = ProcessData()
	if err != nil {
		t.Error("error reading spreadsheet: " + err.Error())
	}
	var size = spreadsheet.Size()
	fmt.Println("Spreadsheet size = ", size)
	if size == 0 {
		t.Error("spreadsheet was read without rows")
	}
}

// Test_ColumnSearch checks that the spreadsheet can identify the proper
// column when given a textual heading.
func Test_Column(t *testing.T) {
	var err error
	var spreadsheet Spreadsheet
	var column int

	spreadsheet, err = ProcessData()
	if err != nil {
		t.Error("error reading spreadsheet: " + err.Error())
	}
	//
	// Valid heading
	//
	column, err = spreadsheet.column("Payee")
	if err != nil {
		t.Error(err.Error())
	} else if column != 1 {
		t.Error("column should be 1, but is: " + strconv.Itoa(column))
	}
	//
	// Invalid heading
	//
	_, err = spreadsheet.column("XXX")
	if err == nil {
		t.Error("column did not identify an invalid heading")
	} else {
		fmt.Println(err)
	}
}

// Test_Cell checks that a cell from the spreadsheet can be resolved.
func Test_Cell(t *testing.T) {
	var err error
	var spreadsheet Spreadsheet
	var cell string

	spreadsheet, err = ProcessData()
	if err != nil {
		t.Error("error reading spreadsheet: " + err.Error())
	}
	cell, err = spreadsheet.Cell(1, "Type")
	if err != nil {
		t.Error(err.Error())
	} else if cell != "Payment" {
		t.Error("incorrect value of cell: " + cell)
	} else {
		fmt.Println(cell)
	}
}
