// ----------------------------------------------------------------------------
//
// # universities
//
// This file contain a map of university names to the
// two year or four year indicator.
//
// Author: William Shaffer
//
// Copyright (c) 2026 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

package main

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import "fmt"

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type InstitutionType string

// ----------------------------------------------------------------------------
// Data
// ----------------------------------------------------------------------------

const (
	UnknownInstitution InstitutionType = ""
	FourYear           InstitutionType = "Four Year"
	TwoYear            InstitutionType = "Two Year"
)

var universities = map[string]InstitutionType{
	"Accelerated Academy":                      FourYear,
	"American Intercontinental University":     FourYear,
	"Appalachian State University":             FourYear,
	"Arizona College of Nursing":               FourYear,
	"Arizona State University":                 FourYear,
	"Barton College":                           FourYear,
	"Campbell University":                      FourYear,
	"Care One Health Training Institute":       FourYear,
	"Catawba Valley Community College":         TwoYear,
	"Chamberlain University":                   FourYear,
	"Chevres, Athiana":                         FourYear,
	"Colorado Christian University":            FourYear,
	"Davidson-Davie Community College":         TwoYear,
	"Durham Technical Community College":       TwoYear,
	"East Carolina University":                 FourYear,
	"East Coast Polytechnic Institute":         FourYear,
	"Edgecombe Community College":              TwoYear,
	"Fayetteville State University":            FourYear,
	"Guilford College":                         FourYear,
	"Hernandez, Jose":                          FourYear,
	"Johnson C. Smith University":              FourYear,
	"Methodist University":                     FourYear,
	"Nash Community College":                   TwoYear,
	"Norfolk State University":                 FourYear,
	"North Carolina State University":          FourYear,
	"Paredes-Arellano, Brandon":                FourYear,
	"Pitt Community College":                   TwoYear,
	"Shades of Purple":                         TwoYear,
	"Shaw University":                          FourYear,
	"Stepful, Inc.":                            FourYear,
	"University of North Carolina":             FourYear,
	"University of North Carolina - Charlotte": FourYear,
	"University of North Carolina - Pembroke":  FourYear,
	"University of North Carolina Greensboro":  FourYear,
	"University of North Carolina Wilmington":  FourYear,
	"Vance-Granville Community College":        TwoYear,
	"Wake Technical Community College":         TwoYear,
	"William Peace University":                 FourYear,
	"Wilson Community College":                 TwoYear,
	"Winston-Salem State University":           FourYear,
}

// ----------------------------------------------------------------------------
// Functions
// ----------------------------------------------------------------------------

func DetermineInstitutionType(university string) (InstitutionType, error) {
	institutionType, ok := universities[university]
	if ok {
		return institutionType, nil
	}
	return UnknownInstitution, fmt.Errorf("no entry for this university: %s", university)
}
