// ----------------------------------------------------------------------------
//
// Email Validation Program
//
// Author: William Shaffer
// Version: 20-Jun-2024
//
// Copyright (c) 2024 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

package main

// This program validates emails

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	"fmt"

	"github.com/zerobounce/zerobouncego"
)

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

const email = "wshaffer@waysysweb.com"

// ----------------------------------------------------------------------------
// Functions
// ----------------------------------------------------------------------------

func main() {

	// For Querying a single E-Mail and IP
	// IP can also be an empty string
	response, error_ := zerobouncego.Validate(email, "")

	if error_ != nil {
		fmt.Println("error occurred: ", error_.Error())
	} else {
		// Now you can check status
		if response.Status == zerobouncego.S_VALID {
			fmt.Println("This email is valid")
		}

		if response.Status == zerobouncego.S_INVALID {
			fmt.Println("This email is invalid")
			if response.SubStatus == zerobouncego.SS_POSSIBLE_TYPO {
				fmt.Println("This email might have a typo")
			}
		}

	}
}
