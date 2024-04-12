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

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/xuri/excelize/v2"
)

func main() {
	inputFile := "/home/bozo/golang/acorn_go/data/register.xlsx"
	tab := "Worksheet"
	//
	// Check current directory
	//
	mydir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(mydir)

	file, err := excelize.OpenFile(inputFile)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := file.GetRows(tab)
	if err != nil {
		log.Fatal(err)
	}

	for _, row := range rows {
		for _, col := range row {
			fmt.Print(col, "\t")
		}
		fmt.Println()
	}
}
