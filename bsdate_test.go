package bsdate

import (
	"errors"
	"testing"
	"time"
)

func ad(y int, m time.Month, d int) time.Time {
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

// Known BS <-> AD pairs, including range edges and month/year boundaries.
var knownDates = []struct {
	bs Date
	ad time.Time
}{
	{MustNew(2000, 1, 1), ad(1943, time.April, 14)},   // epoch
	{MustNew(2000, 1, 30), ad(1943, time.May, 13)},    // last day of Baisakh 2000
	{MustNew(2000, 2, 1), ad(1943, time.May, 14)},     // JS returned 2000-01-31
	{MustNew(2000, 12, 31), ad(1944, time.April, 12)}, // last day of 2000
	{MustNew(2001, 1, 1), ad(1944, time.April, 13)},   // first day of 2001
	{MustNew(2082, 5, 18), ad(2025, time.September, 3)},
	{MustNew(2082, 5, 25), ad(2025, time.September, 10)},
	{MustNew(2090, 12, 30), ad(2034, time.April, 13)}, // last supported day
}

func TestToAD(t *testing.T) {
	for _, tt := range knownDates {
		if got := tt.bs.ToAD(); !got.Equal(tt.ad) {
			t.Errorf("%v.ToAD() = %s, want %s", tt.bs, got.Format(time.DateOnly), tt.ad.Format(time.DateOnly))
		}
	}
}

func TestFromAD(t *testing.T) {
	for _, tt := range knownDates {
		got, err := FromAD(tt.ad)
		if err != nil {
			t.Errorf("FromAD(%s) error: %v", tt.ad.Format(time.DateOnly), err)
			continue
		}
		if got != tt.bs {
			t.Errorf("FromAD(%s) = %v, want %v", tt.ad.Format(time.DateOnly), got, tt.bs)
		}
	}
}

func TestFromADIgnoresTimeOfDayAndUsesOwnLocation(t *testing.T) {
	kathmandu := time.FixedZone("NPT", 5*3600+45*60)
	tests := []struct {
		name string
		in   time.Time
		want Date
	}{
		// JS returned 2000-02-02 for noon because of Math.ceil.
		{"noon UTC", time.Date(1943, time.May, 14, 12, 0, 0, 0, time.UTC), MustNew(2000, 2, 1)},
		{"end of day UTC", time.Date(1943, time.May, 14, 23, 59, 59, 0, time.UTC), MustNew(2000, 2, 1)},
		// 00:30 in Kathmandu is still the previous day in UTC; the local date wins.
		{"early morning Kathmandu", time.Date(2025, time.September, 10, 0, 30, 0, 0, kathmandu), MustNew(2082, 5, 25)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromAD(tt.in)
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromADOutOfRange(t *testing.T) {
	for _, in := range []time.Time{
		ad(1943, time.April, 13),
		ad(2034, time.April, 14),
		ad(1, time.January, 1),
		ad(9999, time.December, 31),
	} {
		if _, err := FromAD(in); !errors.Is(err, ErrADDateOutOfRange) {
			t.Errorf("FromAD(%s) error = %v, want ErrADDateOutOfRange", in.Format(time.DateOnly), err)
		}
	}
}

func TestNewInvalid(t *testing.T) {
	tests := []struct{ y, m, d int }{
		{1999, 5, 25},
		{2091, 5, 25},
		{2000, 0, 1},
		{2000, 13, 25},
		{2000, 1, 0},
		{2000, 1, 33},
		{2000, 1, 31}, // Baisakh 2000 has 30 days
		{0, 0, 0},
	}
	for _, tt := range tests {
		if _, err := New(tt.y, tt.m, tt.d); !errors.Is(err, ErrBSDateOutOfRange) {
			t.Errorf("New(%d, %d, %d) error = %v, want ErrBSDateOutOfRange", tt.y, tt.m, tt.d, err)
		}
	}
}

func TestMustNewPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("MustNew(1999, 1, 1) did not panic")
		}
	}()
	MustNew(1999, 1, 1)
}

// TestRoundTrip converts every supported day AD -> BS -> AD and checks that
// consecutive AD days map to consecutive BS days.
func TestRoundTrip(t *testing.T) {
	prev := Date{}
	count := 0
	for day := epoch; !day.After(ad(2034, time.April, 13)); day = day.AddDate(0, 0, 1) {
		bs, err := FromAD(day)
		if err != nil {
			t.Fatalf("FromAD(%s): %v", day.Format(time.DateOnly), err)
		}
		if _, err := New(bs.Year(), bs.Month(), bs.Day()); err != nil {
			t.Fatalf("FromAD(%s) produced invalid date %v: %v", day.Format(time.DateOnly), bs, err)
		}
		if back := bs.ToAD(); !back.Equal(day) {
			t.Fatalf("round trip %s -> %v -> %s", day.Format(time.DateOnly), bs, back.Format(time.DateOnly))
		}
		if !prev.IsZero() && !isNextDay(prev, bs) {
			t.Fatalf("%v is not the day after %v", bs, prev)
		}
		prev = bs
		count++
	}
	if count != totalDays() {
		t.Errorf("visited %d days, want %d", count, totalDays())
	}
}

func isNextDay(a, b Date) bool {
	switch {
	case a.day < monthDays[a.year][a.month-1]:
		return b == Date{a.year, a.month, a.day + 1}
	case a.month < 12:
		return b == Date{a.year, a.month + 1, 1}
	default:
		return b == Date{a.year + 1, 1, 1}
	}
}

func TestNames(t *testing.T) {
	d := MustNew(2082, 5, 18) // Wednesday, 2025-09-03
	tests := []struct {
		name, got, want string
	}{
		{"Nepali", d.Nepali(), "भाद्र १८, २०८२"},
		{"MonthName", d.MonthName(), "भाद्र"},
		{"MonthNameRomanized", d.MonthNameRomanized(), "Bhadra"},
		{"DayName", d.DayName(), "बुधबार"},
		{"DayNameRomanized", d.DayNameRomanized(), "Budhbar"},
		{"DayNameEnglish", d.DayNameEnglish(), "Wednesday"},
		{"String", d.String(), "2082-05-18"},
	}
	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("%s() = %q, want %q", tt.name, tt.got, tt.want)
		}
	}
	if d.Weekday() != time.Wednesday {
		t.Errorf("Weekday() = %v, want Wednesday", d.Weekday())
	}
}

func TestZeroValue(t *testing.T) {
	var d Date
	if !d.IsZero() {
		t.Error("IsZero() = false for zero Date")
	}
	if !d.ToAD().IsZero() {
		t.Errorf("ToAD() = %v, want zero time", d.ToAD())
	}
	for name, s := range map[string]string{
		"Nepali": d.Nepali(), "MonthName": d.MonthName(), "MonthNameRomanized": d.MonthNameRomanized(),
		"DayName": d.DayName(), "DayNameRomanized": d.DayNameRomanized(), "DayNameEnglish": d.DayNameEnglish(),
	} {
		if s != "" {
			t.Errorf("%s() = %q for zero Date, want empty", name, s)
		}
	}
}

func TestDaysInMonth(t *testing.T) {
	if n, err := DaysInMonth(2082, 5); err != nil || n != 31 {
		t.Errorf("DaysInMonth(2082, 5) = %d, %v; want 31, nil", n, err)
	}
	if _, err := DaysInMonth(2082, 13); !errors.Is(err, ErrBSDateOutOfRange) {
		t.Errorf("DaysInMonth(2082, 13) error = %v, want ErrBSDateOutOfRange", err)
	}
}

func TestToday(t *testing.T) {
	now := time.Now()
	today := Today()
	y, m, d := now.Date()
	if got := today.ToAD(); !got.Equal(ad(y, m, d)) {
		t.Errorf("Today().ToAD() = %s, want %04d-%02d-%02d", got.Format(time.DateOnly), y, m, d)
	}
}

func TestNepaliDigits(t *testing.T) {
	tests := map[int]string{0: "०", 7: "७", 123: "१२३", 2082: "२०८२", -45: "-४५"}
	for in, want := range tests {
		if got := NepaliDigits(in); got != want {
			t.Errorf("NepaliDigits(%d) = %q, want %q", in, got, want)
		}
	}
}
