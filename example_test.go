package bsdate_test

import (
	"errors"
	"fmt"
	"time"

	bsdate "github.com/dlamichhane/bs-date-go"
)

func ExampleNew() {
	d, err := bsdate.New(2082, 5, 25)
	if err != nil {
		panic(err)
	}
	fmt.Println(d.ToAD().Format(time.DateOnly))
	// Output: 2025-09-10
}

func ExampleNew_invalid() {
	_, err := bsdate.New(2082, 13, 1)
	fmt.Println(errors.Is(err, bsdate.ErrBSDateOutOfRange))
	fmt.Println(err)
	// Output:
	// true
	// bsdate: BS date out of range: month 13 is outside 1-12
}

func ExampleFromAD() {
	d, err := bsdate.FromAD(time.Date(2025, time.September, 10, 0, 0, 0, 0, time.UTC))
	if err != nil {
		panic(err)
	}
	fmt.Println(d)
	fmt.Println(d.Year(), d.Month(), d.Day())
	// Output:
	// 2082-05-25
	// 2082 5 25
}

func ExampleDate_Nepali() {
	d := bsdate.MustNew(2082, 5, 18)
	fmt.Println(d.Nepali())
	fmt.Println(d.MonthNameRomanized(), d.Day())
	fmt.Println(d.DayName(), d.DayNameRomanized(), d.DayNameEnglish())
	// Output:
	// भाद्र १८, २०८२
	// Bhadra 18
	// बुधबार Budhbar Wednesday
}

func ExampleNepaliDigits() {
	fmt.Println(bsdate.NepaliDigits(2082))
	// Output: २०८२
}
