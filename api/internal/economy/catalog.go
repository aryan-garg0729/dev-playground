// Package economy is the game economy: Bytes, power-ups, quests and chests,
// medals and boss rewards. It rewards learning but never changes grading:
// nothing here can mark a check as passed or reveal a solution, and grading
// never imports this package (the game engine calls it through an interface).
package economy

import "fmt"

// Byte rewards (SPEC "Bytes").
const (
	ChestMinBytes   = 40
	ChestMaxBytes   = 80
	ChapterBytes    = 150
	BossBytes       = 300
	LeagueTop5Bytes = 200
	WeeklyQuestByte = 150

	BetMin = 100
	BetMax = 500

	FocusBoostMinutes = 15
	RepairWindowDays  = 2  // a streak lost in the last 48 hours
	RepairEveryDays   = 30 // once a month
	MaxFreezes        = 2  // matches game.MaxFreezes and the user_stats check
)

// Power-up ids.
const (
	StreakFreeze    = "streak_freeze"
	StreakRepair    = "streak_repair"
	HolidayMode     = "holiday_mode"
	FocusBoost      = "focus_boost"
	ExtraHint       = "extra_hint"
	Peek            = "peek"
	SecondWind      = "second_wind"
	LeagueShield    = "league_shield"
	DoubleOrNothing = "double_or_nothing"
)

// PowerUp is a shop item.
type PowerUp struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Note        string `json:"note"`
	Cost        int    `json:"cost"` // 0: not bought (bet stake, Pro)
	// Unavailable explains why it cannot be bought yet ("pro", "coming_soon").
	Unavailable string `json:"unavailable,omitempty"`
}

// PowerUps in shop order.
var PowerUps = []PowerUp{
	{ID: StreakFreeze, Name: "Streak Freeze", Description: "Covers one missed day so your streak survives.", Note: "Hold up to 2 · 1 free every 7-day streak", Cost: 200},
	{ID: FocusBoost, Name: "Focus Boost", Description: "Double XP for the next 15 minutes of lessons.", Note: "1 per day", Cost: 150},
	{ID: SecondWind, Name: "Second Wind", Description: "Keeps your first-try bonus after one failed submit.", Note: "1 per chapter", Cost: 250},
	{ID: ExtraHint, Name: "Extra Hint", Description: "One more tutor hint after your daily limit.", Note: "Removes the no-hint bonus, like any hint", Cost: 100, Unavailable: "coming_soon"},
	{ID: Peek, Name: "Peek", Description: "Shows the input of one hidden check you failed. Never the expected output.", Note: "Removes the first-try bonus", Cost: 120},
	{ID: LeagueShield, Name: "League Shield", Description: "You won't be demoted at the end of this week.", Note: "Not two weeks in a row", Cost: 400},
	{ID: StreakRepair, Name: "Streak Repair", Description: "Brings back a streak you lost in the last 48 hours.", Note: "Once a month", Cost: 1500},
	{ID: HolidayMode, Name: "Holiday Mode", Description: "Pause your streak for up to 7 days while you travel.", Note: "Turn on at least 3 days ahead", Unavailable: "pro"},
	{ID: DoubleOrNothing, Name: "Double or Nothing", Description: "Bet Bytes that you'll hit your daily goal 7 days in a row.", Note: "One bet at a time"},
}

// PowerUpByID finds a power-up.
func PowerUpByID(id string) (PowerUp, bool) {
	for _, p := range PowerUps {
		if p.ID == id {
			return p, true
		}
	}
	return PowerUp{}, false
}

// ChestOdds are the daily and boss chests' power-up odds (weights, shown to users).
var ChestOdds = []struct {
	PowerUp string `json:"powerup"`
	Weight  int    `json:"weight"`
}{
	{FocusBoost, 40}, {Peek, 30}, {SecondWind, 20}, {LeagueShield, 10},
}

// Tier is one level of a medal.
type Tier struct {
	Tier  string `json:"tier"` // bronze | silver | gold | diamond
	Gte   int    `json:"gte"`
	Bytes int    `json:"bytes"`
	Title string `json:"title"`
}

// Medal is a medal and its tiers. Stat names a counter or a user stat.
type Medal struct {
	ID          string `json:"id"`
	Family      string `json:"family"`
	Name        string `json:"name"`
	Unit        string `json:"unit"` // "days", "lessons", ...
	Stat        string `json:"stat"`
	Secret      bool   `json:"secret"`
	ComingSoon  bool   `json:"coming_soon"`
	Description string `json:"description"`
	Tiers       []Tier `json:"tiers"`
}

// Families in display order.
var Families = []struct{ ID, Name string }{
	{"consistency", "Consistency"}, {"mastery", "Mastery"}, {"craft", "Craft"},
	{"league", "League"}, {"memory", "Memory"}, {"secret", "Secret"},
}

var tierBytes = map[string]int{"bronze": 50, "silver": 150, "gold": 300, "diamond": 500}
var roman = []string{"I", "II", "III", "IV"}

// tiers builds bronze/silver/gold tiers (or the given names) with titles.
func tiers(name string, gte ...int) []Tier {
	names := []string{"bronze", "silver", "gold", "diamond"}
	out := make([]Tier, len(gte))
	for i, g := range gte {
		out[i] = Tier{Tier: names[i], Gte: g, Bytes: tierBytes[names[i]], Title: fmt.Sprintf("%s %s", name, roman[i])}
	}
	return out
}

func one(name, tier string, gte int) []Tier {
	return []Tier{{Tier: tier, Gte: gte, Bytes: tierBytes[tier], Title: name}}
}

// Medals in display order. Stats: user stats (longest_streak, total_xp) or
// counters bumped by the rewards hook, the league close and the rollover.
var Medals = []Medal{
	{ID: "streak", Family: "consistency", Name: "Streak", Stat: "longest_streak", Unit: "days", Description: "Keep a daily streak", Tiers: tiers("Streak", 7, 30, 100)},
	{ID: "goal_getter", Family: "consistency", Name: "Goal Getter", Stat: "goal_days", Unit: "days on goal", Description: "Hit your daily goal", Tiers: tiers("Goal Getter", 7, 30, 100)},
	{ID: "night_owl", Family: "consistency", Name: "Night Owl", Stat: "night", Unit: "late lessons", Description: "Pass lessons after 10 pm", Tiers: tiers("Night Owl", 3, 15, 50)},
	{ID: "early_bird", Family: "consistency", Name: "Early Bird", Stat: "early", Unit: "lessons before 8 am", Description: "Pass lessons before 8 am", Tiers: tiers("Early Bird", 3, 15, 50)},

	{ID: "boss_slayer", Family: "mastery", Name: "Boss Slayer", Stat: "bosses", Unit: "bosses", Description: "Defeat chapter bosses", Tiers: tiers("Boss Slayer", 1, 5, 15)},
	{ID: "flawless_chapter", Family: "mastery", Name: "Flawless Chapter", Stat: "flawless_chapters", Unit: "chapters", Description: "Finish a chapter with every lesson on the first try", Tiers: tiers("Flawless", 1, 3, 10)},
	{ID: "course_complete", Family: "mastery", Name: "Course Complete", Stat: "courses", Unit: "courses", Description: "Finish every lesson of a course", Tiers: tiers("Graduate", 1, 3, 5)},
	{ID: "scholar", Family: "mastery", Name: "Scholar", Stat: "lessons", Unit: "lessons", Description: "Pass lessons", Tiers: tiers("Scholar", 10, 50, 200)},

	{ID: "unaided", Family: "craft", Name: "Unaided", Stat: "no_hint", Unit: "lessons without hints", Description: "Pass lessons without a hint", Tiers: tiers("Unaided", 10, 50, 200)},
	{ID: "sharpshooter", Family: "craft", Name: "Sharpshooter", Stat: "first_try", Unit: "first-try passes", Description: "Pass lessons on the first try", Tiers: tiers("Sharpshooter", 10, 50, 200)},
	{ID: "bug_hunter", Family: "craft", Name: "Bug Hunter", Stat: "comebacks", Unit: "fixed attempts", Description: "Pass a lesson after a failed submit", Tiers: tiers("Bug Hunter", 5, 25, 100)},
	{ID: "big_earner", Family: "craft", Name: "Big Earner", Stat: "total_xp", Unit: "XP", Description: "Earn XP", Tiers: tiers("Big Earner", 500, 2500, 10000)},

	{ID: "promoted", Family: "league", Name: "Promoted", Stat: "tier_max", Unit: "tier", Description: "Reach Silver, Gold, then Diamond league", Tiers: tiers("Promoted", 2, 3, 4)},
	{ID: "podium", Family: "league", Name: "Podium", Stat: "podium", Unit: "top-3 weeks", Description: "Finish a league week in the top 3", Tiers: tiers("Podium", 1, 5, 20)},
	{ID: "contender", Family: "league", Name: "Contender", Stat: "top5", Unit: "top-5 weeks", Description: "Finish a league week in the top 5", Tiers: tiers("Contender", 1, 5, 20)},
	{ID: "diamond_league", Family: "league", Name: "Diamond League", Stat: "tier_max", Unit: "tier", Description: "Reach Diamond", Tiers: one("Diamond", "diamond", 4)},

	{ID: "reviewer", Family: "memory", Name: "Reviewer", Stat: "reviews", Unit: "reviews", Description: "Finish spaced reviews", ComingSoon: true, Tiers: tiers("Reviewer", 25, 100, 300)},

	{ID: "persistent", Family: "secret", Name: "Persistent", Stat: "persistent", Unit: "lessons", Secret: true, Description: "Pass a lesson after five failed submits", Tiers: one("Persistent", "bronze", 1)},
	{ID: "high_roller", Family: "secret", Name: "High Roller", Stat: "bets_won", Unit: "bets won", Secret: true, Description: "Win Double or Nothing", Tiers: tiers("High Roller", 1, 3, 10)},
	{ID: "shielded", Family: "secret", Name: "Shielded", Stat: "shield_saves", Unit: "saves", Secret: true, Description: "A League Shield saved you from demotion", Tiers: one("Shielded", "silver", 1)},
}

// Quest is a daily or weekly quest. Event names what counts toward it.
type Quest struct {
	ID      string `json:"id"`
	Cadence string `json:"cadence"` // daily | weekly
	Name    string `json:"name"`
	Event   string `json:"event"` // lessons | first_try | xp | code | quiz | goal_day
	Target  int    `json:"target"`
	Bytes   int    `json:"bytes"` // weekly quests pay Bytes; daily ones fill the chest
}

// DailyPool is drawn from (3 a day per user); WeeklyPool gives one a week.
var (
	DailyPool = []Quest{
		{ID: "d.lessons2", Cadence: "daily", Name: "Pass 2 lessons", Event: "lessons", Target: 2},
		{ID: "d.firsttry1", Cadence: "daily", Name: "Pass one lesson on the first try", Event: "first_try", Target: 1},
		{ID: "d.xp50", Cadence: "daily", Name: "Earn 50 XP", Event: "xp", Target: 50},
		{ID: "d.code1", Cadence: "daily", Name: "Pass a coding lesson", Event: "code", Target: 1},
		{ID: "d.quiz1", Cadence: "daily", Name: "Pass a quiz", Event: "quiz", Target: 1},
	}
	WeeklyPool = []Quest{
		{ID: "w.lessons10", Cadence: "weekly", Name: "Pass 10 lessons this week", Event: "lessons", Target: 10, Bytes: WeeklyQuestByte},
		{ID: "w.goal5", Cadence: "weekly", Name: "Hit your daily goal on 5 days", Event: "goal_day", Target: 5, Bytes: WeeklyQuestByte},
		{ID: "w.xp300", Cadence: "weekly", Name: "Earn 300 XP this week", Event: "xp", Target: 300, Bytes: WeeklyQuestByte},
	}
)

// QuestsPerDay is how many daily quests a user gets.
const QuestsPerDay = 3
