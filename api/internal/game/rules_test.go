package game

import (
	"slices"
	"testing"
	"time"
)

func date(y int, m time.Month, d int) time.Time { return time.Date(y, m, d, 0, 0, 0, 0, time.UTC) }

func TestAwardFor(t *testing.T) {
	tests := []struct {
		name           string
		base           int
		firstTry, hint bool
		want           Breakdown
		wantTotal      int
	}{
		{"first try, no hint", 20, true, false, Breakdown{20, 10, 5}, 35},
		{"retry, no hint", 20, false, false, Breakdown{20, 0, 5}, 25},
		{"first try with hint", 20, true, true, Breakdown{20, 10, 0}, 30},
		{"retry with hint", 20, false, true, Breakdown{20, 0, 0}, 20},
		{"rounds half up", 10, true, false, Breakdown{10, 5, 3}, 18},
		{"zero xp", 0, true, false, Breakdown{}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AwardFor(tt.base, tt.firstTry, tt.hint)
			if got != tt.want || got.Total() != tt.wantTotal {
				t.Fatalf("AwardFor = %+v (total %d), want %+v (total %d)", got, got.Total(), tt.want, tt.wantTotal)
			}
		})
	}
}

func TestLevels(t *testing.T) {
	// Leaving level n needs round(100 * n^1.5) cumulative XP.
	starts := map[int]int{1: 0, 2: 100, 3: 283, 4: 520, 7: 1470, 8: 1852}
	for n, xp := range starts {
		if got := LevelStart(n); got != xp {
			t.Errorf("LevelStart(%d) = %d, want %d", n, got, xp)
		}
	}
	tests := []struct{ xp, level int }{{0, 1}, {99, 1}, {100, 2}, {282, 2}, {283, 3}, {1469, 6}, {1470, 7}, {1851, 7}, {1852, 8}}
	for _, tt := range tests {
		if got := LevelFor(tt.xp); got != tt.level {
			t.Errorf("LevelFor(%d) = %d, want %d", tt.xp, got, tt.level)
		}
	}
}

func TestLocalDateAcrossTimezones(t *testing.T) {
	kolkata, _ := time.LoadLocation("Asia/Kolkata")
	la, _ := time.LoadLocation("America/Los_Angeles")
	auckland, _ := time.LoadLocation("Pacific/Auckland")
	// 2026-11-03 20:00 UTC is already the 4th in Kolkata and Auckland, still the 3rd in LA.
	at := time.Date(2026, 11, 3, 20, 0, 0, 0, time.UTC)
	tests := []struct {
		loc  *time.Location
		want time.Time
	}{
		{time.UTC, date(2026, 11, 3)},
		{kolkata, date(2026, 11, 4)},
		{la, date(2026, 11, 3)},
		{auckland, date(2026, 11, 4)},
	}
	for _, tt := range tests {
		if got := LocalDate(at, tt.loc); !got.Equal(tt.want) {
			t.Errorf("LocalDate(%s) = %s, want %s", tt.loc, got.Format(time.DateOnly), tt.want.Format(time.DateOnly))
		}
	}
}

func TestApplyPass(t *testing.T) {
	d := func(day int) time.Time { return date(2026, 11, day) }
	tests := []struct {
		name  string
		in    Streak
		today time.Time
		want  StreakChange
	}{
		{"first ever pass", Streak{}, d(3),
			StreakChange{Streak: Streak{Days: 1, Last: d(3)}, Extended: true}},
		{"second pass same day", Streak{Days: 4, Last: d(3)}, d(3),
			StreakChange{Streak: Streak{Days: 4, Last: d(3)}}},
		{"next day extends", Streak{Days: 4, Last: d(3)}, d(4),
			StreakChange{Streak: Streak{Days: 5, Last: d(4)}, Extended: true}},
		{"one missed day, no freeze: reset", Streak{Days: 4, Last: d(3)}, d(5),
			StreakChange{Streak: Streak{Days: 1, Last: d(5)}, Extended: true, Reset: true}},
		{"one missed day uses a freeze", Streak{Days: 4, Freezes: 1, Last: d(3)}, d(5),
			StreakChange{Streak: Streak{Days: 5, Freezes: 0, Last: d(5)}, Extended: true, FreezesUsed: 1}},
		{"two missed days use two freezes", Streak{Days: 9, Freezes: 2, Last: d(3)}, d(6),
			StreakChange{Streak: Streak{Days: 10, Freezes: 0, Last: d(6)}, Extended: true, FreezesUsed: 2}},
		{"three missed days, two freezes: reset, freezes kept", Streak{Days: 9, Freezes: 2, Last: d(3)}, d(7),
			StreakChange{Streak: Streak{Days: 1, Freezes: 2, Last: d(7)}, Extended: true, Reset: true}},
		{"7th day earns a freeze", Streak{Days: 6, Last: d(3)}, d(4),
			StreakChange{Streak: Streak{Days: 7, Freezes: 1, Last: d(4)}, Extended: true, FreezeEarned: true}},
		{"14th day earns a second", Streak{Days: 13, Freezes: 1, Last: d(3)}, d(4),
			StreakChange{Streak: Streak{Days: 14, Freezes: 2, Last: d(4)}, Extended: true, FreezeEarned: true}},
		{"no more than two held", Streak{Days: 20, Freezes: 2, Last: d(3)}, d(4),
			StreakChange{Streak: Streak{Days: 21, Freezes: 2, Last: d(4)}, Extended: true}},
		{"freeze spent then earned on day 7", Streak{Days: 6, Freezes: 1, Last: d(3)}, d(5),
			StreakChange{Streak: Streak{Days: 7, Freezes: 1, Last: d(5)}, Extended: true, FreezesUsed: 1, FreezeEarned: true}},
		{"date before last (timezone moved west) changes nothing", Streak{Days: 3, Last: d(4)}, d(3),
			StreakChange{Streak: Streak{Days: 3, Last: d(4)}}},
		{"across a month boundary", Streak{Days: 2, Last: date(2026, 10, 31)}, d(1),
			StreakChange{Streak: Streak{Days: 3, Last: d(1)}, Extended: true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ApplyPass(tt.in, tt.today)
			if got != tt.want {
				t.Fatalf("ApplyPass(%+v, %s)\n got  %+v\n want %+v", tt.in, tt.today.Format(time.DateOnly), got, tt.want)
			}
		})
	}
}

// A learner in Kolkata (UTC+5:30) passing late at night and early the next
// morning in local time keeps a two-day streak, though both are the same UTC day.
func TestStreakFollowsUserTimezone(t *testing.T) {
	kolkata, _ := time.LoadLocation("Asia/Kolkata")
	la, _ := time.LoadLocation("America/Los_Angeles")
	tests := []struct {
		name     string
		loc      *time.Location
		a, b     time.Time
		wantDays int
	}{
		{"Kolkata: 23:30 then 00:30 local is two days", kolkata,
			time.Date(2026, 11, 3, 18, 0, 0, 0, time.UTC), time.Date(2026, 11, 3, 19, 0, 0, 0, time.UTC), 2},
		{"UTC: same instants are one day", time.UTC,
			time.Date(2026, 11, 3, 18, 0, 0, 0, time.UTC), time.Date(2026, 11, 3, 19, 0, 0, 0, time.UTC), 1},
		{"LA: 20:00 then next 21:00 local, a different UTC split", la,
			time.Date(2026, 11, 4, 4, 0, 0, 0, time.UTC), time.Date(2026, 11, 5, 5, 0, 0, 0, time.UTC), 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := ApplyPass(Streak{}, LocalDate(tt.a, tt.loc)).Streak
			s = ApplyPass(s, LocalDate(tt.b, tt.loc)).Streak
			if s.Days != tt.wantDays {
				t.Fatalf("days = %d, want %d", s.Days, tt.wantDays)
			}
		})
	}
}

func TestCurrent(t *testing.T) {
	d := func(day int) time.Time { return date(2026, 11, day) }
	tests := []struct {
		name  string
		s     Streak
		today time.Time
		want  int
	}{
		{"none", Streak{}, d(3), 0},
		{"active today", Streak{Days: 5, Last: d(3)}, d(3), 5},
		{"active yesterday, today still open", Streak{Days: 5, Last: d(2)}, d(3), 5},
		{"missed one, no freeze", Streak{Days: 5, Last: d(1)}, d(3), 0},
		{"missed one, freeze covers it", Streak{Days: 5, Freezes: 1, Last: d(1)}, d(3), 5},
		{"missed three, two freezes", Streak{Days: 5, Freezes: 2, Last: d(1)}, d(5), 0},
	}
	for _, tt := range tests {
		if got := Current(tt.s, tt.today); got != tt.want {
			t.Errorf("%s: Current = %d, want %d", tt.name, got, tt.want)
		}
	}
}

func TestWeekStart(t *testing.T) {
	tests := []struct{ at, want time.Time }{
		{time.Date(2026, 11, 2, 0, 0, 0, 0, time.UTC), date(2026, 11, 2)},   // Monday 00:00
		{time.Date(2026, 11, 8, 23, 59, 0, 0, time.UTC), date(2026, 11, 2)}, // Sunday night
		{time.Date(2026, 11, 9, 0, 0, 0, 0, time.UTC), date(2026, 11, 9)},
		{time.Date(2026, 11, 4, 12, 0, 0, 0, time.FixedZone("x", 14*3600)), date(2026, 11, 2)},
	}
	for _, tt := range tests {
		if got := WeekStart(tt.at); !got.Equal(tt.want) {
			t.Errorf("WeekStart(%s) = %s, want %s", tt.at, got.Format(time.DateOnly), tt.want.Format(time.DateOnly))
		}
	}
}

func TestMoves(t *testing.T) {
	count := func(m []string, r string) int {
		n := 0
		for _, x := range m {
			if x == r {
				n++
			}
		}
		return n
	}
	tests := []struct {
		tier              string
		n                 int
		promoted, demoted int
	}{
		{"silver", 30, 5, 5},
		{"silver", 12, 5, 5},
		{"silver", 8, 5, 3},
		{"silver", 3, 3, 0},
		{"bronze", 30, 5, 0},  // nobody below bronze
		{"diamond", 30, 0, 5}, // nobody above diamond
		{"diamond", 4, 0, 4},
	}
	for _, tt := range tests {
		m := Moves(tt.tier, tt.n)
		if count(m, Promoted) != tt.promoted || count(m, Demoted) != tt.demoted || len(m) != tt.n {
			t.Errorf("Moves(%s, %d) = %v; want %d promoted, %d demoted", tt.tier, tt.n, m, tt.promoted, tt.demoted)
		}
		// Promotions are the top ranks, demotions the bottom ones.
		if tt.promoted > 0 && m[0] != Promoted {
			t.Errorf("Moves(%s, %d): rank 1 is %s", tt.tier, tt.n, m[0])
		}
		if tt.demoted > 0 && m[tt.n-1] != Demoted {
			t.Errorf("Moves(%s, %d): last rank is %s", tt.tier, tt.n, m[tt.n-1])
		}
	}
	if got := []string{NextTier("bronze", Promoted), NextTier("gold", Demoted), NextTier("diamond", Promoted), NextTier("bronze", Demoted)}; !slices.Equal(got, []string{"silver", "silver", "diamond", "bronze"}) {
		t.Errorf("NextTier = %v", got)
	}
}
