// ----------------------------------------------------------------------------
//
// Guestlist Status
//
// Author: William Shaffer
// Version: 29-Jun-2024
//
// Copyright (c) 2024 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

package guestlist

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type Status string

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

const (
	Opened       Status = "Opened"
	Attending    Status = "Attending"
	Sent         Status = "Sent"
	Regrets      Status = "Regrets"
	Unsubscribed Status = "Unsubscribed"
	Unknown             = "Unknown"
)

// ----------------------------------------------------------------------------
// Functions
// ----------------------------------------------------------------------------

// NewStatus accepts a string returns a Status.  If the string is not a valid
// status, Unknown is returned.
func NewStatus(value string) Status {
	var result Status
	switch value {
	case "Opened":
		result = Opened
	case "Attending":
		result = Attending
	case "Sent":
		result = Sent
	case "Regrets":
		result = Regrets
	case "Unsubscribed":
		result = Unsubscribed
	default:
		result = Unknown
	}
	return result
}
