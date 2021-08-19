package sas

import (
	// Standard Library Imports
	"errors"
)

var (
	ErrInvalidVersion          = errors.New("error parsing signed version")
	ErrInvalidStartDateFormat  = errors.New("invalid date format provided for signed start, must be ISO 8601 formatted date string")
	ErrInvalidExpiryDateFormat = errors.New("invalid date format provided for signed expiry, must be ISO 8601 formatted date string")
)
