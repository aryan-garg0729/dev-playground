package game

import (
	"math"
	"slices"
	"time"
)

// Breakdown is the XP for one first pass.
type Breakdown struct {
	Base     int `json:"base"`
	FirstTry int `json:"first_try"`
	NoHint   int `json:"no_hint"`
}

// Total is the sum of the parts.
func (b Breakdown) Total() int { return b.Base + b.FirstTry + b.NoHint }

// AwardFor is the XP for a first pass of an exercise worth base XP.
func AwardFor(base int, firstTry, usedHint bool) Breakdown {
	b := Breakdown{Base: max(base, 0)}
	if firstTry {
		b.FirstTry = int(math.Round(float64(b.Base) * FirstTryBonus))
	}
	if !usedHint {
		b.NoHint = int(math.Round(float64(b.Base) * NoHintBonus))
	}
	return b
}

// LevelStart is the cumulative XP at which level n begins (level 1 starts at 0).
func LevelStart(n int) int {
	if n <= 1 {
		return 0
	}
	return int(math.Round(LevelBase * math.Pow(float64(n-1), LevelExponent)))
}

// LevelFor is the level reached with xp cumulative XP.
func LevelFor(xp int) int {
	n := 1
	for LevelStart(n+1) <= xp {
		n++
	}
	return n
}

// Streak is a user's streak state. Last is the local calendar date of the last
// day that counted, as midnight UTC (zero when there is none).
type Streak struct {
	Days    int
	Freezes int
	Last    time.Time
}

// StreakChange is the effect of a pass on a streak.
type StreakChange struct {
	Streak
	Extended     bool // today now counts (it did not before)
	FreezesUsed  int  // freezes spent on missed days
	FreezeEarned bool
	Reset        bool // missed days could not be covered; the streak restarted
}

// LocalDate is the calendar date of t in loc, as midnight UTC.
func LocalDate(t time.Time, loc *time.Location) time.Time {
	y, m, d := t.In(loc).Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

func daysBetween(a, b time.Time) int { return int(b.Sub(a).Hours() / 24) }

// ApplyPass updates s for a pass on local date today.
//
// The first pass of a day extends the streak. Each missed day in between uses
// one Streak Freeze; if there are not enough, the streak restarts at 1 and the
// freezes are kept. Every FreezeEvery days of streak earns a freeze (max
// MaxFreezes). A date at or before Last (e.g. after a timezone change) changes
// nothing.
func ApplyPass(s Streak, today time.Time) StreakChange {
	c := StreakChange{Streak: s}
	if !s.Last.IsZero() && !today.After(s.Last) {
		return c
	}
	c.Extended = true
	missed := 0
	if !s.Last.IsZero() {
		missed = daysBetween(s.Last, today) - 1
	}
	switch {
	case s.Last.IsZero() || s.Days == 0:
		c.Days = 1
	case missed > s.Freezes:
		c.Days, c.Reset = 1, true
	default:
		c.Freezes -= missed
		c.FreezesUsed = missed
		c.Days++
	}
	c.Last = today
	if c.Days%FreezeEvery == 0 && c.Freezes < MaxFreezes {
		c.Freezes++
		c.FreezeEarned = true
	}
	return c
}

// Current is the streak to show on local date today: still alive when the
// missed days since Last could be covered by freezes, else 0.
func Current(s Streak, today time.Time) int {
	if s.Last.IsZero() || s.Days == 0 {
		return 0
	}
	if missed := daysBetween(s.Last, today) - 1; missed > s.Freezes {
		return 0
	}
	return s.Days
}

// WeekStart is the Monday 00:00 UTC that starts t's league week.
func WeekStart(t time.Time) time.Time {
	t = t.UTC()
	d := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	offset := (int(d.Weekday()) + 6) % 7 // Monday = 0
	return d.AddDate(0, 0, -offset)
}

// Result of a week in a league group.
const (
	Promoted = "promoted"
	Demoted  = "demoted"
	Stayed   = "stayed"
)

// Moves gives each rank (index in ranked order) its result in tier: the top
// LeaguePromote move up, then the bottom LeagueDemote of the rest move down.
// There is no promotion from the top tier and no demotion from the bottom one.
func Moves(tier string, n int) []string {
	out := make([]string, n)
	i := slices.Index(Tiers, tier)
	promote := 0
	if i < len(Tiers)-1 {
		promote = min(LeaguePromote, n)
	}
	demote := 0
	if i > 0 {
		demote = min(LeagueDemote, n-promote)
	}
	for r := range out {
		switch {
		case r < promote:
			out[r] = Promoted
		case r >= n-demote:
			out[r] = Demoted
		default:
			out[r] = Stayed
		}
	}
	return out
}

// NextTier is tier after result.
func NextTier(tier, result string) string {
	i := slices.Index(Tiers, tier)
	if i < 0 {
		return Tiers[0]
	}
	switch result {
	case Promoted:
		i = min(i+1, len(Tiers)-1)
	case Demoted:
		i = max(i-1, 0)
	}
	return Tiers[i]
}
