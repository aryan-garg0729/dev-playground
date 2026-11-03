// Package game turns graded submissions into XP, levels, streaks and league
// standing. It reads a Result; it never changes one (grading is the only
// authority on pass/fail).
package game

// Every tunable number for v1 lives here (SPEC "Gamification engine").
const (
	// FirstTryBonus and NoHintBonus are fractions of the exercise's base XP.
	FirstTryBonus = 0.50
	NoHintBonus   = 0.25

	// Leaving level n needs LevelBase * n^LevelExponent cumulative XP.
	LevelBase     = 100.0
	LevelExponent = 1.5

	// FreezeEvery streak days earns one Streak Freeze, up to MaxFreezes held.
	FreezeEvery = 7
	MaxFreezes  = 2

	// Leagues: groups of up to LeagueGroupSize per tier and week. At the weekly
	// close the top LeaguePromote move up a tier and the bottom LeagueDemote move
	// down (nobody is both, so small groups promote first).
	LeagueGroupSize = 30
	LeaguePromote   = 5
	LeagueDemote    = 5

	// BossBonus is extra XP for defeating a boss (its first pass).
	BossBonus = 100

	DefaultDailyGoal = 2
	MaxDailyGoal     = 10
)

// Tiers in ascending order.
var Tiers = []string{"bronze", "silver", "gold", "diamond"}

// ReasonPass is the ledger reason for a lesson's first pass (ref_id = exercise
// id); the ledger's unique key allows one per user and exercise.
const ReasonPass = "pass"
