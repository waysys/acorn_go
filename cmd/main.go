// ----------------------------------------------------------------------------
//
// Donor analysis program
//
// Author: William Shaffer
// Version: 12-Apr-2024
//
// Copyright (c) William Shaffer
//
// ----------------------------------------------------------------------------

// This package processes donation information based on the payment transactons
// from a spreadsheet obtain from QuickBooks.  The package will then output
// the results in another spreadsheet.
package main

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	"fmt"
	"log"

	"github.com/xuri/excelize/v2"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

const inputFile = "/home/bozo/golang/acorn_go/data/register.xlsx"
const tab = "Worksheet"

// ----------------------------------------------------------------------------
// Functions
// ----------------------------------------------------------------------------

// main supervises the processing of the donation data.
func main() {

	data := readData(inputFile, tab)

	for _, row := range data {
		for _, col := range row {
			fmt.Print(col, "\t")
		}
		fmt.Println()
	}
}

// readData extracts the donation data from Excel file and
// returns an array of string arrays.
func readData(excelFillName string, tab string) [][]string {
	//
	// Open file
	//
	file, err := excelize.OpenFile(excelFillName)
	if err != nil {
		log.Fatal(err)
	}
	//
	// Read rows
	//
	rows, err := file.GetRows(tab)
	if err != nil {
		log.Fatal(err)
	}
	return rows
}
