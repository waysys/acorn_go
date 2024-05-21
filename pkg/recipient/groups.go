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
