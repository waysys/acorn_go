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
	"strconv"

	dec "github.com/shopspring/decimal"
	d "github.com/waysys/waydate/pkg/date"
	"github.com/xuri/excelize/v2"
	excel "github.com/xuri/excelize/v2"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type SpreadsheetFile struct {
	filename  string
	sheetname string
	filePtr   *excel.File
}

type FormatIndex int

// ----------------------------------------------------------------------------
// Constant
// ----------------------------------------------------------------------------

const (
	FormatPerCent FormatIndex = 9
	FormatDate    FormatIndex = 15
	FormatMoney   FormatIndex = 2
	FormatInt     FormatIndex = 1
)

// ----------------------------------------------------------------------------
// Factory Functions
// ----------------------------------------------------------------------------

// Create an Excel spreadsheet with the specified file name and a sheet
// with a specified name.
func New(filename string, sheetname string) (SpreadsheetFile, error) {
	var spFile SpreadsheetFile
	var err error = nil
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
	f := excel.NewFile()
	spFile.filePtr = f
	//
	// Create a new sheet.
	//
	_, err = f.NewSheet(sheetname)
	if err != nil {
		return spFile, err
	}
	f.DeleteSheet("Sheet1")
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

// AddSheet creates a new sheet in the spreadsheet file and returns
// a new spreadsheet file structure pointing to the same file, but
// with the new sheet name.
func (spFilePtr *SpreadsheetFile) AddSheet(sheetname string) (SpreadsheetFile, error) {
	var err error = nil
	var spFile SpreadsheetFile
	//
	// Precondition
	//
	if sheetname == "" {
		err = errors.New("sheetname must not be an empty string")
		return spFile, err
	}
	//
	// Make copy of spreadsheet file structure
	//
	spFile.filename = spFilePtr.filename
	spFile.filePtr = spFilePtr.filePtr
	spFile.sheetname = sheetname
	//
	// Create new sheet
	//
	_, err = spFilePtr.filePtr.NewSheet(sheetname)
	return spFile, err
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

// SetCellFloat sets the value of a cell to a floating point number
func (spFilePtr *SpreadsheetFile) SetCellFloat(cell string, value float64) error {
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
	err = file.SetCellFloat(sheetname, cell, value, 0, 64)
	return err
}

// SetCellInt sets the value of a cell to an integer
func (spFilePtr *SpreadsheetFile) SetCellInt(cell string, value int) error {
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
	err = file.SetCellInt(sheetname, cell, int64(value))
	if err == nil {
		err = spFilePtr.SetNumFmt(cell, FormatInt)
	}
	return err
}

// SetCellInt sets the value of a cell to a decimal number
func (spFilePtr *SpreadsheetFile) SetCellDecimal(cell string, amount dec.Decimal, index FormatIndex) error {
	var err error = nil
	var sheetname = spFilePtr.sheetname
	var file = spFilePtr.filePtr
	var value string = ""
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
	value = amount.String()
	err = spFilePtr.SetNumFmt(cell, index)
	if err == nil {
		err = file.SetCellValue(sheetname, cell, value)
	}
	return err
}

// SetCellDate sets the value of a cell to the specified date
func (spFilePtr *SpreadsheetFile) SetCellDate(cell string, date d.Date) error {
	var err error = nil
	var sheetname = spFilePtr.sheetname
	var file = spFilePtr.filePtr
	var value string = ""
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
	value = date.String()
	err = spFilePtr.SetNumFmt(cell, FormatDate)
	if err == nil {
		err = file.SetCellValue(sheetname, cell, value)
	}
	return err
}

// SetNumFmt sets the number format on a cell.
// For format codes, see https://xuri.me/excelize/en/style.html#number_format
func (spFilePtr *SpreadsheetFile) SetNumFmt(cell string, index FormatIndex) error {
	var style = excelize.Style{
		NumFmt: int(index),
	}
	var err error = nil
	var st int
	//
	// Precondition
	//
	if index < 0 || index > 49 {
		err = errors.New("Invalid value for format index: " + strconv.Itoa(int(index)))
		return err
	}
	st, err = spFilePtr.filePtr.NewStyle(&style)
	if err != nil {
		return err
	}
	err = spFilePtr.filePtr.SetCellStyle(spFilePtr.sheetname, cell, cell, st)
	return err
}
