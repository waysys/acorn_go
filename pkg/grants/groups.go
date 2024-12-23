// ----------------------------------------------------------------------------
//
// Award Group
//
// Author: William Shaffer
// Version: 12-May-2024
//
// Copyright (c) 2024 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

package grants

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	"errors"

	r "github.com/waysys/waydate/pkg/daterange"

	d "github.com/waysys/waydate/pkg/date"
)

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// Summer 2023 Award Group
var summer2023AwardDate, _ = d.New(7, 23, 2023)
var summer2023AwardRangeBegin = summer2023AwardDate
var summer2023AwardRangeEnd, _ = d.New(7, 31, 2023)
var summer2023AwardRange, _ = r.New(summer2023AwardRangeBegin, summer2023AwardRangeEnd)
var summer2023AwardGroup, _ = NewAwardGroup("Summer 2023", summer2023AwardDate, summer2023AwardRange)

// Fall 2023 Award Group
var fall2023AwardDate, _ = d.New(11, 9, 2023)
var fall2023AwardRangeBegin = fall2023AwardDate
var fall2023AwardRangeEnd, _ = d.New(5, 31, 2024)
var fall2023AwardRange, _ = r.New(fall2023AwardRangeBegin, fall2023AwardRangeEnd)
var fall2023AwardGroup, _ = NewAwardGroup("Fall 2023", fall2023AwardDate, fall2023AwardRange)

// Summer 2024 Award Group
var summer2024AwardDate, _ = d.New(7, 11, 2024)
var summer2024AwardRangeBegin = summer2024AwardDate
var summer2024AwardRangeEnd, _ = d.New(9, 30, 2024)
var summer2024AwardRange, _ = r.New(summer2024AwardRangeBegin, summer2024AwardRangeEnd)
var summer2024AwardGroup, _ = NewAwardGroup("Summer 2024", summer2024AwardDate, summer2024AwardRange)

// Fall 2024 Award Group
var fall2024AwardDate, _ = d.New(11, 11, 2024)
var fall2024AwardRangeBegin = fall2024AwardDate
var fall2024AwardRangeEnd, _ = d.New(12, 31, 2024)
var fall2024AwardRange, _ = r.New(fall2024AwardRangeBegin, fall2024AwardRangeEnd)
var fall2024AwardGroup, _ = NewAwardGroup("Fall 2024", fall2024AwardDate, fall2024AwardRange)

var groups = []*AwardGroup{
	&summer2023AwardGroup,
	&fall2023AwardGroup,
	&summer2024AwardGroup,
	&fall2024AwardGroup,
}

// ----------------------------------------------------------------------------
// Functions
// ----------------------------------------------------------------------------

// FindAwardGroup returns an award group based on the transaction date supplied.
func FindAwardGroup(transactionDate d.Date) (*AwardGroup, error) {
	var err error = errors.New("no award group found for transaction date: " + transactionDate.String())
	var awardGroup *AwardGroup

	for _, group := range groups {
		if group.awardRange.InRange(transactionDate) {
			awardGroup = group
			err = nil
			break
		}
	}
	return awardGroup, err
}

// Groups returns the slice of award groups.
func Groups() []*AwardGroup {
	return groups
}
