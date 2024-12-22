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
	a "acorn_go/pkg/accounting"
	"errors"

	q "acorn_go/pkg/quickbooks"

	r "github.com/waysys/waydate/pkg/daterange"

	d "github.com/waysys/waydate/pkg/date"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// An award group is a collection of scholarships awarded on a specified date.
// The award date is the date the board approves the scholarships.
// An award group has a date range in which bills may exist.
type AwardGroup struct {
	awardDate  d.Date
	awardRange r.DateRange
	groupName  string
	bills      []*q.EducationBill
}

// ----------------------------------------------------------------------------
// Factory Function
// ----------------------------------------------------------------------------

// NewAwardGroup creates a new award group.
func NewAwardGroup(name string, date d.Date, rng r.DateRange) (AwardGroup, error) {
	var err error = nil
	var group AwardGroup = AwardGroup{}
	//
	// Preconditions
	//
	if name == "" {
		err = errors.New("award group name must not be an empty string")
		return group, err
	}
	if !rng.InRange(date) {
		err = errors.New("the award range must contain the award date: " + date.String())
		return group, err
	}
	//
	// Initialization
	//
	group = AwardGroup{
		awardDate:  date,
		awardRange: rng,
		groupName:  name,
	}
	return group, err
}

// ----------------------------------------------------------------------------
// Properties
// ----------------------------------------------------------------------------

// GroupName returns the name of the group.
func (group *AwardGroup) GroupName() string {
	return group.groupName
}

// AwardDate returns the date of the scholarship awards.
func (group *AwardGroup) AwardDate() d.Date {
	return group.awardDate
}

// InAwardGroup returns true if the specified date is in the
// award range.
func (group *AwardGroup) InAwardGroup(date d.Date) bool {
	return group.awardRange.InRange(date)
}

// FiscalYear returns the fiscal year indicator based on
// the award date.
func (group *AwardGroup) FiscalYear() a.FYIndicator {
	return a.FiscalYearIndicator(group.AwardDate())
}

// AddBill adds an educaiton bill to the award group
func (group *AwardGroup) AddBill(bill *q.EducationBill) {
	group.bills = append(group.bills, bill)
}

// Count returns the number of education bills associated
// with this award group.
func (group *AwardGroup) Count() int {
	return len(group.bills)
}
