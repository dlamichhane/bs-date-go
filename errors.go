package bsdate

import "errors"

var (
	// ErrBSDateOutOfRange is returned when a Bikram Sambat date is invalid or
	// outside the supported range (MinYear–MaxYear). Check with errors.Is.
	ErrBSDateOutOfRange = errors.New("bsdate: BS date out of range")

	// ErrADDateOutOfRange is returned when an AD (Gregorian) date falls outside
	// the range that can be converted to BS (1943-04-14 to 2034-04-13).
	// Check with errors.Is.
	ErrADDateOutOfRange = errors.New("bsdate: AD date out of range")
)
