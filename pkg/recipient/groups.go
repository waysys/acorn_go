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

package recipient

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	d "acorn_go/pkg/date"
	r "acorn_go/pkg/daterange"
	"errors"
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
var fall2023AwardRangeEnd, _ = d.New(1, 31, 2024)
var fall2023AwardRange, _ = r.New(fall2023AwardRangeBegin, fall2023AwardRangeEnd)
var fall2023AwardGroup, _ = NewAwardGroup("Fall 2023", fall2023AwardDate, fall2023AwardRange)

// May 2024 Award Group
var may2024AwardDate, _ = d.New(5, 22, 2024)
var may2024AwardRangeBegin = may2024AwardDate
var may2024AwardRangeEnd = may2024AwardDate
var may2024AwardRange, _ = r.New(may2024AwardRangeBegin, may2024AwardRangeEnd)
var may2023AwardGroup, _ = NewAwardGroup("May 2023", may2024AwardDate, may2024AwardRange)

var groups = []AwardGroup{
	summer2023AwardGroup,
	fall2023AwardGroup,
	may2023AwardGroup,
}

// ----------------------------------------------------------------------------
// Functions
// ----------------------------------------------------------------------------

// FinaAwardGroup returns an award group based on the transaction date supplied.
func FinaAwardGroup(transactionDate d.Date) (AwardGroup, error) {
	var err error = errors.New("no award group found for transaction date: " + transactionDate.String())
	var awardGroup = AwardGroup{}

	for _, group := range groups {
		if group.awardRange.InRange(transactionDate) {
			awardGroup = group
			err = nil
			break
		}
	}
	return awardGroup, err
}
