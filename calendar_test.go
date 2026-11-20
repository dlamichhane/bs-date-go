package bsdate

import "testing"

func TestCalendarData(t *testing.T) {
	if len(monthDays) != MaxYear-MinYear+1 {
		t.Errorf("monthDays has %d years, want %d", len(monthDays), MaxYear-MinYear+1)
	}
	for y := MinYear; y <= MaxYear; y++ {
		total := 0
		for m, n := range monthDays[y] {
			if n < 29 || n > 32 {
				t.Errorf("year %d month %d has %d days, want 29-32", y, m+1, n)
			}
			total += n
		}
		if total < 365 || total > 366 {
			t.Errorf("year %d has %d days, want 365-366", y, total)
		}
	}
}
