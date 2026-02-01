// ----------------------------------------------------------------------------
//
// Donor series program
//
// Author: William Shaffer
// Version: 20-May-2024
//
// Copyright (c) 2024 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

// The donors package manages the names, addresses, and email addresses of
// donors.  It also manages the donations made by the donors.
package donors

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	ac "acorn_go/pkg/accounting"
	a "acorn_go/pkg/address"
	s "strings"

	dec "github.com/shopspring/decimal"
	"github.com/waysys/assert/assert"
	d "github.com/waysys/waydate/pkg/date"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type Donor struct {
	key                   string
	name                  string
	address               a.Address
	email                 string
	numberHousehold       int
	donations             []dec.Decimal
	donationsCalendarYear []dec.Decimal
	donationsCurrentMonth dec.Decimal
	deceased              bool
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

var ZERO = dec.Zero
var MajorDonorLimit = dec.NewFromInt(1000)

// ----------------------------------------------------------------------------
// Factory Functions
// ----------------------------------------------------------------------------

// New creates a new donor
func New(ky string, nm string, adr a.Address, eml string, count int, decsd bool) Donor {
	var dn1 []dec.Decimal
	var dn2 []dec.Decimal
	var numYears = ac.MaxYear - ac.MinYear + 1
	for index := 0; index < numYears; index++ {
		dn1 = append(dn1, ZERO)
		dn2 = append(dn2, ZERO)
	}

	donor := Donor{
		key:                   ky,
		name:                  nm,
		address:               adr,
		email:                 eml,
		numberHousehold:       count,
		donations:             dn1,
		donationsCalendarYear: dn2,
		donationsCurrentMonth: ZERO,
		deceased:              decsd,
	}
	return donor
}

// NewDonorWithDonation creates a new donor.
func NewDonorWithDonation(name string) Donor {
	var blankAddress = a.BlankAddress()
	var donor = New(name, name, blankAddress, "", 0, false)
	return donor
}

// ----------------------------------------------------------------------------
// Properties
// ----------------------------------------------------------------------------

// Key returns the abbreviated donor name used to identify the donor.
func (donor Donor) Key() string {
	return donor.key
}

// Name returns the full names of the donors
func (donor Donor) Name() string {
	return donor.name
}

// Street returns the street address of the donor
func (donor Donor) Street() string {
	var street = donor.address.Street
	street = s.Replace(street, "DriveUnit", "Drive, Unit", 1)
	street = s.Replace(street, "\n", ", ", 1)
	street = s.Replace(street, "DrUnit", "Drive, Unit", 1)
	return street
}

// City returns the city of the donor
func (donor Donor) City() string {
	return donor.address.City
}

// State returns the state of the donor
func (donor Donor) State() string {
	return donor.address.State
}

// Zip returns the zip code of the donor
func (donor Donor) Zip() string {
	return donor.address.Zip
}

// Address returns the address structure associated with this
// donor
func (donor Donor) Address() a.Address {
	return donor.address
}

// Email returns the email address of the donor
func (donor Donor) Email() string {
	return donor.email
}

// HasEmail returns true if the donor has an email address.
func (donor Donor) HasEmail() bool {
	var result = true
	var email = donor.Email()
	if s.TrimSpace(email) == "--" {
		result = false
	}
	return result
}

// NumberInHousehold return the number of people in the
// donor household
func (donor Donor) NumberInHousehold() int {
	return donor.numberHousehold
}

// Deceased returns true if the donor is deceased
func (donor Donor) Deceased() bool {
	return donor.deceased
}

// ----------------------------------------------------------------------------
// Donation Properties - Fiscal Year
// ----------------------------------------------------------------------------

// Donation returns the donation for the specified fiscal year
func (donor Donor) Donation(fy ac.FYIndicator) dec.Decimal {
	assert.Assert(ac.IsFYIndicator(fy), "Invalid FYIndicator: "+fy.String())
	var amount = donor.donations[fy]
	return amount
}

// TotalDonation returns the total donations for this donor.
func (donor Donor) TotalDonation() dec.Decimal {
	var amount = ZERO
	for _, fy := range ac.FYIndicators {
		var donation = donor.Donation(fy)
		amount = amount.Add(donation)
	}
	return amount
}

// IsDonor returns true if the donor donated an amount greater than zero
// for the specified fiscal year.
func (donor Donor) IsDonor(fy ac.FYIndicator) bool {
	var donation = donor.Donation(fy)
	var result = donation.GreaterThan(ZERO)
	return result
}

// IsMajorDonor returns true if the donor donated the major donor limit
func (donor Donor) IsMajorDonor(fy ac.FYIndicator) bool {
	var donation = donor.Donation(fy)
	var result = donation.GreaterThanOrEqual(MajorDonorLimit)
	return result
}

// IsMajorDonorOverall returns true if the donor has been a major donor
// in any of the fiscal years.
func (donor Donor) IsMajorDonorOverall() bool {
	var result = false
	for _, fy := range ac.FYIndicators {
		if donor.IsMajorDonor(fy) {
			result = true
			break
		}
	}
	return result
}

// IsNonRepeatDonor returns true if the donor donated in the prior year to the specified
// fiscal year, but did not donate in the current fiscal year.
func (donor Donor) IsNonRepeatDonor(fy ac.FYIndicator) bool {
	var result = false
	var priorFy = fy.Prior()
	result = (priorFy != ac.OutOfRange) && !donor.IsDonor(fy) && donor.IsDonor(priorFy)
	return result
}

// ----------------------------------------------------------------------------
// Donation Properties - Setters - Fiscal Year
// ----------------------------------------------------------------------------

// AddDonation adds the amount to the donations for the specified fical year.
func (donor Donor) AddDonation(amount dec.Decimal, fy ac.FYIndicator) {
	assert.Assert(ac.IsFYIndicator(fy), "Invalid FYIndicator: "+fy.String())
	donor.donations[fy] = donor.donations[fy].Add(amount)
}

// ----------------------------------------------------------------------------
// Donation Properties - Calendar Year
// ----------------------------------------------------------------------------

// CalDonation returns the donations for a donor for the specified calendar year.
func (donor Donor) CalDonation(year ac.YearIndicator) dec.Decimal {
	assert.Assert(ac.IsYearIndicator(year), "Invalid year indicator: "+year.String())
	var amount = donor.donationsCalendarYear[year]
	return amount
}

// IsDonor returns true if the donor donated in the calendar year.
func (donor Donor) IsCalDonor(year ac.YearIndicator) bool {
	assert.Assert(ac.IsYearIndicator(year), "Invalid year indicator: "+year.String())
	var donation = donor.CalDonation(year)
	var result = donation.GreaterThan(ZERO)
	return result
}

// ----------------------------------------------------------------------------
// Donation Method - Setters - Calendar Year
// ----------------------------------------------------------------------------

// AddCalDonation adds an amount to the calendar donations for the specified year.
func (donor *Donor) AddCalDonation(amount dec.Decimal, year ac.YearIndicator) {
	assert.Assert(ac.IsYearIndicator(year), "Invalid year indicator: "+year.String())
	donor.donationsCalendarYear[year] = donor.donationsCalendarYear[year].Add(amount)
}

// ----------------------------------------------------------------------------
// Donation Properties - Current Month
// ----------------------------------------------------------------------------

// isCurrentMonthDonor returns true if the donor has contributed
// in the current month.
func (donor Donor) IsCurrentDonor() bool {
	var result = donor.CurrentMonthDonation().GreaterThan(ZERO)
	return result
}

// CurrentMonthDonation returns the total amount of current month donations
// for this donor.
func (donor Donor) CurrentMonthDonation() dec.Decimal {
	return donor.donationsCurrentMonth
}

// ----------------------------------------------------------------------------
// Donation Setters - Current Month
// ----------------------------------------------------------------------------

// AddMonthDonation adds the amount of the donation if the payment date
// is in the current month date range.
func (donor *Donor) AddMonthDonation(amount dec.Decimal, date d.Date) {
	if ac.IsCurrentMonth(date) {
		donor.donationsCurrentMonth = donor.donationsCurrentMonth.Add(amount)
	}
}
