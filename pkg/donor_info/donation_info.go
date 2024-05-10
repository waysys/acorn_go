// ----------------------------------------------------------------------------
//
// Donoation information
//
// Author: William Shaffer
// Version: 12-Apr-2024
//
// Copyright (c) 2024 William Shaffer All rights Reserved
//
// ----------------------------------------------------------------------------

package donor_info

// This file defines a structure containing the number of donations and the
// amont of donations in a month

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type DonationInfo struct {
	count  int
	amount float64
}

// ----------------------------------------------------------------------------
// Factory Functions
// ----------------------------------------------------------------------------

func NewDonationInfo() DonationInfo {
	var donationInfo = DonationInfo{
		count:  0,
		amount: 0.0,
	}
	return donationInfo
}

// ----------------------------------------------------------------------------
// Methods
// ----------------------------------------------------------------------------

// Count returns the number of donations
func (info DonationInfo) Count() int {
	return info.count
}

// Amount returns the amount of donations
func (info DonationInfo) Amount() float64 {
	return info.amount
}

// AddCount adds a value to the number of donations
func (info *DonationInfo) AddCount(value int) {
	info.count += value
}

// AddAmount adds a value to the amount of donations
func (info *DonationInfo) AddAmount(amount float64) {
	info.amount += amount
}
