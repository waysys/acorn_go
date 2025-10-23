// ----------------------------------------------------------------------------
//
// Donor analysis program
//
// Author: William Shaffer
// Version: 12-Apr-2024
//
// Copyright (c) 2024 William Shaffer All Rights Reserved
//
// ----------------------------------------------------------------------------

// This package provides functions that help handle errors in main programs.
package support

// ----------------------------------------------------------------------------
// Imports
// ----------------------------------------------------------------------------

import (
	"log/slog"
	"os"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// Error Handling Functions
// ----------------------------------------------------------------------------

// check tests an error to see if it is null.  If not, it prints an
// error message and exits the program.
func Check(err error, message string) {
	if err != nil {
		logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
		logger.Error(
			message+err.Error(),
			"error", err,
		)
		os.Exit(1)
	}
}
