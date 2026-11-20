package bsdate

import (
	"fmt"
	"time"
)

// Date is a date in the Bikram Sambat calendar.
//
// A Date created with New, MustNew, FromAD or Today is always valid. The zero
// value is not a valid date; its methods return zero values (see IsZero).
// Date is a small value type and is safe to copy and compare with ==.
type Date struct {
	year, month, day int
}

// New returns the BS date year-month-day. Month is 1 (Baisakh) to 12
// (Chaitra). It returns an error wrapping ErrBSDateOutOfRange if the year is
// outside MinYear–MaxYear or the month or day does not exist.
func New(year, month, day int) (Date, error) {
	if err := validate(year, month, day); err != nil {
		return Date{}, err
	}
	return Date{year: year, month: month, day: day}, nil
}

// MustNew is like New but panics if the date is invalid. It is meant for
// dates known to be valid, such as constants in code and tests.
func MustNew(year, month, day int) Date {
	d, err := New(year, month, day)
	if err != nil {
		panic(err)
	}
	return d
}

// FromAD converts an AD (Gregorian) time to a BS date.
//
// Only the calendar date of t in its own location is used; the time of day is
// ignored. It returns an error wrapping ErrADDateOutOfRange if the date is
// outside 1943-04-14 to 2034-04-13.
func FromAD(t time.Time) (Date, error) {
	y, m, d := t.Date()
	day := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	last := epoch.AddDate(0, 0, totalDays()-1)
	if day.Before(epoch) || day.After(last) {
		return Date{}, fmt.Errorf("%w: %s is outside %s to %s", ErrADDateOutOfRange,
			day.Format(time.DateOnly), epoch.Format(time.DateOnly), last.Format(time.DateOnly))
	}
	return fromOffset(int(day.Sub(epoch).Hours()) / 24), nil
}

// Today returns the current date in BS, using the local time zone.
// It panics if the current date is outside the supported range.
func Today() Date {
	d, err := FromAD(time.Now())
	if err != nil {
		panic(err)
	}
	return d
}

// DaysInMonth returns the number of days in the given BS month (1–12).
func DaysInMonth(year, month int) (int, error) {
	if err := validate(year, month, 1); err != nil {
		return 0, err
	}
	return monthDays[year][month-1], nil
}

// Year returns the BS year.
func (d Date) Year() int { return d.year }

// Month returns the BS month, 1 (Baisakh) to 12 (Chaitra).
func (d Date) Month() int { return d.month }

// Day returns the day of the month.
func (d Date) Day() int { return d.day }

// IsZero reports whether d is the zero value, which is not a valid date.
func (d Date) IsZero() bool { return d == Date{} }

// ToAD returns the equivalent AD (Gregorian) date at midnight UTC.
// It returns the zero time.Time if d is the zero Date.
func (d Date) ToAD() time.Time {
	if d.IsZero() {
		return time.Time{}
	}
	return epoch.AddDate(0, 0, d.offset())
}

// Weekday returns the day of the week.
func (d Date) Weekday() time.Weekday { return d.ToAD().Weekday() }

// MonthName returns the month name in Devanagari, e.g. "भाद्र".
func (d Date) MonthName() string {
	if d.IsZero() {
		return ""
	}
	return monthNames[d.month-1]
}

// MonthNameRomanized returns the month name in Latin script, e.g. "Bhadra".
func (d Date) MonthNameRomanized() string {
	if d.IsZero() {
		return ""
	}
	return monthNamesRomanized[d.month-1]
}

// DayName returns the weekday name in Devanagari, e.g. "बुधबार".
func (d Date) DayName() string {
	if d.IsZero() {
		return ""
	}
	return weekdayNames[d.Weekday()]
}

// DayNameRomanized returns the Nepali weekday name in Latin script,
// e.g. "Budhbar".
func (d Date) DayNameRomanized() string {
	if d.IsZero() {
		return ""
	}
	return weekdayNamesRomanized[d.Weekday()]
}

// DayNameEnglish returns the English weekday name, e.g. "Wednesday".
func (d Date) DayNameEnglish() string {
	if d.IsZero() {
		return ""
	}
	return d.Weekday().String()
}

// Nepali formats the date in Nepali script, e.g. "भाद्र २५, २०८२".
func (d Date) Nepali() string {
	if d.IsZero() {
		return ""
	}
	return d.MonthName() + " " + NepaliDigits(d.day) + ", " + NepaliDigits(d.year)
}

// String formats the date as YYYY-MM-DD, e.g. "2082-05-25".
func (d Date) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.year, d.month, d.day)
}

// offset returns the number of days from epoch to d.
func (d Date) offset() int {
	n := yearStart[d.year-MinYear]
	for m := 0; m < d.month-1; m++ {
		n += monthDays[d.year][m]
	}
	return n + d.day - 1
}

// fromOffset converts a day offset from epoch (0 <= n < totalDays()) to a
// Date. A day belongs to the next year/month once n reaches that period's
// length, hence >= (the JS original used >, which was off by one).
func fromOffset(n int) Date {
	y := MinYear
	for n >= yearStart[y-MinYear+1] {
		y++
	}
	n -= yearStart[y-MinYear]
	m := 0
	for n >= monthDays[y][m] {
		n -= monthDays[y][m]
		m++
	}
	return Date{year: y, month: m + 1, day: n + 1}
}

func validate(year, month, day int) error {
	if year < MinYear || year > MaxYear {
		return fmt.Errorf("%w: year %d is outside %d-%d", ErrBSDateOutOfRange, year, MinYear, MaxYear)
	}
	if month < 1 || month > 12 {
		return fmt.Errorf("%w: month %d is outside 1-12", ErrBSDateOutOfRange, month)
	}
	if n := monthDays[year][month-1]; day < 1 || day > n {
		return fmt.Errorf("%w: day %d is outside 1-%d for %d-%02d", ErrBSDateOutOfRange, day, n, year, month)
	}
	return nil
}
