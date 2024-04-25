// ----------------------------------------------------------------------------
//
// Spreadsheet writer
//
// Author: William Shaffer
// Version: 25-Apr-2024
//
// Copyright (c) William Shaffer
//
// ----------------------------------------------------------------------------

// The spreadsheet package writes an Excel spreadsheet.
package spreadsheet

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	"errors"

	"github.com/xuri/excelize/v2"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type SpreadsheetFile struct {
	filename  string
	sheetname string
	filePtr   *excelize.File
}

// ----------------------------------------------------------------------------
// Factory Functions
// ----------------------------------------------------------------------------

// Create an Excel spreadsheet with the specified file name and a sheet
// with a specified name.
func New(filename string, sheetname string) (SpreadsheetFile, error) {
	var spFile SpreadsheetFile
	var err error = nil
	var index int
	//
	// Preconditions
	//
	if filename == "" {
		err = errors.New("spreadsheet filename must not be an empty string")
		return spFile, err
	}
	spFile.filename = filename
	if sheetname == "" {
		err = errors.New("sheetname must not be an empty string")
		return spFile, err
	}
	spFile.sheetname = sheetname
	//
	// create the file
	//
	f := excelize.NewFile()
	spFile.filePtr = f
	//
	// Create a new sheet.
	//
	index, err = f.NewSheet(sheetname)
	if err != nil {
		return spFile, err
	}
	//
	// Set active sheet of the workbook.
	//
	f.SetActiveSheet(index)
	return spFile, nil
}

// ----------------------------------------------------------------------------
// Methods
// ----------------------------------------------------------------------------

// Save saves the Excel file with the name specified in the NewFile function.
func (spFilePtr *SpreadsheetFile) Save() error {
	var err error = nil
	//
	// Preconditions
	//
	if spFilePtr == nil {
		err = errors.New("pointer to spreadsheet file is nil")
		return err
	}
	//
	// Save file
	//
	var filename = (*spFilePtr).filename
	err = (*spFilePtr).filePtr.SaveAs(filename)
	return err
}

// Close closes the spreadsheet file.
func (spFilePtr *SpreadsheetFile) Close() error {
	var err error = nil
	//
	// Preconditions
	//
	if spFilePtr == nil {
		err = errors.New("pointer to spreadsheet file is nil")
		return err
	}
	//
	// Close the file
	//
	err = ((*spFilePtr).filePtr).Close()
	return err
}

// SetCell sets the value of a cell in the specified spreadsheet
func (spFilePtr *SpreadsheetFile) SetCell(cell string, value string) error {
	var err error = nil
	var sheetname = (*spFilePtr).sheetname
	var file = (*spFilePtr).filePtr
	//
	// Preconditions
	//
	if cell == "" {
		err = errors.New("cell name must not be empty")
		return err
	}
	//
	// Set value
	//
	err = file.SetCellValue(sheetname, cell, value)
	return err
}
