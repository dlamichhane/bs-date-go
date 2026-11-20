// Package bsdate converts dates between the Bikram Sambat (BS, the official
// calendar of Nepal) and Anno Domini (AD, Gregorian) calendars.
//
// Supported range: BS 2000-01-01 to 2090-12-30, which is AD 1943-04-14 to
// 2034-04-13.
//
// Convert BS to AD:
//
//	d, err := bsdate.New(2082, 5, 25)
//	if err != nil {
//		// handle invalid date
//	}
//	t := d.ToAD() // 2025-09-10 00:00:00 UTC
//
// Convert AD to BS:
//
//	d, err := bsdate.FromAD(time.Date(2025, 9, 10, 0, 0, 0, 0, time.UTC))
//	fmt.Println(d)          // 2082-05-25
//	fmt.Println(d.Nepali()) // भाद्र २५, २०८२
//
// Errors wrap ErrBSDateOutOfRange or ErrADDateOutOfRange, so they can be
// checked with errors.Is.
package bsdate
