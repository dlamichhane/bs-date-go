# bs-date-go

[![Go Reference](https://pkg.go.dev/badge/github.com/dlamichhane/bs-date-go.svg)](https://pkg.go.dev/github.com/dlamichhane/bs-date-go)
[![CI](https://github.com/dlamichhane/bs-date-go/actions/workflows/ci.yml/badge.svg)](https://github.com/dlamichhane/bs-date-go/actions/workflows/ci.yml)

A Go library for converting between Bikram Sambat (BS/Nepali) and Anno Domini (AD/Gregorian) dates. It has no dependencies.

This is the Go version of the [bs-date](https://github.com/dlamichhane/bs-date) JavaScript library.

## Installation

```bash
go get github.com/dlamichhane/bs-date-go
```

```go
import bsdate "github.com/dlamichhane/bs-date-go"
```

The package name is `bsdate`, because Go identifiers can't contain hyphens. The alias in the import above is optional, but it makes the name obvious to readers.

## Usage

```go
package main

import (
	"errors"
	"fmt"
	"time"

	bsdate "github.com/dlamichhane/bs-date-go"
)

func main() {
	// BS -> AD
	d, err := bsdate.New(2082, 5, 25)
	if err != nil {
		panic(err)
	}
	fmt.Println(d.ToAD().Format(time.DateOnly)) // 2025-09-10

	// AD -> BS
	bs, err := bsdate.FromAD(time.Date(2025, time.September, 10, 0, 0, 0, 0, time.UTC))
	if err != nil {
		panic(err)
	}
	fmt.Println(bs)                       // 2082-05-25
	fmt.Println(bs.Nepali())              // भाद्र २५, २०८२
	fmt.Println(bs.MonthNameRomanized())  // Bhadra
	fmt.Println(bs.DayName())             // बुधबार
	fmt.Println(bs.DayNameRomanized())    // Budhbar
	fmt.Println(bs.DayNameEnglish())      // Wednesday

	// Today's date in BS (local time zone)
	fmt.Println(bsdate.Today())

	// Errors can be checked with errors.Is
	_, err = bsdate.New(2100, 1, 1)
	fmt.Println(errors.Is(err, bsdate.ErrBSDateOutOfRange)) // true
}
```

## API

| Function / method | Description |
|---|---|
| `New(year, month, day) (Date, error)` | Create a validated BS date (month 1 = Baisakh … 12 = Chaitra) |
| `MustNew(year, month, day) Date` | Like `New`, but panics on an invalid date |
| `FromAD(time.Time) (Date, error)` | Convert an AD date to BS. Uses the calendar date of `t` in its own time zone and ignores the time of day |
| `Today() Date` | Today's date in BS |
| `DaysInMonth(year, month) (int, error)` | Number of days in a BS month |
| `NepaliDigits(int) string` | `2082` → `"२०८२"` |
| `d.ToAD() time.Time` | Convert to AD (midnight UTC) |
| `d.Year()`, `d.Month()`, `d.Day()` | Date parts |
| `d.Weekday() time.Weekday` | Day of the week |
| `d.MonthName()` / `d.MonthNameRomanized()` | `"भाद्र"` / `"Bhadra"` |
| `d.DayName()` / `d.DayNameRomanized()` / `d.DayNameEnglish()` | `"बुधबार"` / `"Budhbar"` / `"Wednesday"` |
| `d.Nepali()` | `"भाद्र २५, २०८२"` |
| `d.String()` | `"2082-05-25"` |

Errors wrap `ErrBSDateOutOfRange` (invalid or unsupported BS date) or `ErrADDateOutOfRange` (AD date outside the supported range).

## Supported range

| Calendar | From | To |
|---|---|---|
| BS | 2000-01-01 | 2090-12-30 |
| AD | 1943-04-14 | 2034-04-13 |

## Development

```bash
go test -race -cover ./...
go vet ./...
gofmt -l .
```

## License

[MIT](LICENSE)
