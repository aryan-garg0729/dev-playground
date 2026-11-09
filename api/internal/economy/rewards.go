package economy

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/aryangarg/codequest/api/internal/exercise"
	"github.com/aryangarg/codequest/api/internal/game"
	"github.com/aryangarg/codequest/api/internal/store"
)

var (
	_ game.Rewards       = (*Service)(nil)
	_ game.LeagueRewards = (*Service)(nil)
)

// OnGraded is the game engine's rewards hook, run in the grading transaction.
// It reads the Result and the engine's Outcome; it cannot change either.
func (s *Service) OnGraded(ctx context.Context, tx pgx.Tx, ev game.Graded) (*game.RewardsOutcome, error) {
	q := store.New(tx)
	out := &game.RewardsOutcome{Bytes: []game.Gain{}, Medals: []game.MedalAward{}, QuestsDone: []string{}}
	ex := ev.Exercise
	meta, err := q.ExerciseChapter(ctx, ex.ID)
	if err != nil {
		return nil, fmt.Errorf("exercise chapter: %w", err)
	}
	if ex.Boss {
		out.Boss = bossHP(ev.Result)
	}
	if ev.Result.Status != exercise.StatusPass {
		return out, nil
	}
	user := ev.User
	st, err := q.GetEconomyStats(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("stats: %w", err)
	}

	// Counters behind medals, on the first pass of an exercise only (no farming).
	if ev.FirstPass {
		prog, err := q.GetProgress(ctx, store.GetProgressParams{UserID: user, ExerciseID: ex.ID})
		if err != nil {
			return nil, fmt.Errorf("progress: %w", err)
		}
		hour := ev.LocalNow.Hour()
		for key, n := range map[string]int{
			"lessons":    1,
			"first_try":  b2i(ev.Outcome.FirstTry),
			"no_hint":    b2i(!prog.UsedHint),
			"comebacks":  b2i(prog.Attempts > 1),
			"persistent": b2i(prog.Attempts > 5),
			"night":      b2i(hour >= 22 || hour < 4),
			"early":      b2i(hour >= 5 && hour < 8),
		} {
			if err := bump(ctx, q, user, key, n); err != nil {
				return nil, err
			}
		}
		if err := s.chapterDone(ctx, q, ev, meta, out); err != nil {
			return nil, err
		}
		if ex.Boss {
			if err := s.bossDefeated(ctx, q, ev, meta, out); err != nil {
				return nil, err
			}
		}
	}

	// Quests count distinct lessons per day (re-solving the same one twice does not).
	if ev.NewToday {
		if err := s.questProgress(ctx, q, ev, meta.Type, int(st.DailyGoal), st.Timezone, out); err != nil {
			return nil, err
		}
		if err := s.settleBet(ctx, q, user, dayOf(s.now(), st.Timezone).Today, out); err != nil {
			return nil, err
		}
	}
	if err := s.evaluateMedals(ctx, q, user, out); err != nil {
		return nil, err
	}
	return out, nil
}

func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}

// bossHP: the boss's HP is its number of checks; each passing check is a hit.
func bossHP(r exercise.Result) *game.BossResult {
	b := &game.BossResult{MaxHP: len(r.Checks)}
	for _, c := range r.Checks {
		if !c.Passed {
			b.HP++
		}
	}
	b.Defeated = r.Status == exercise.StatusPass && b.HP == 0
	return b
}

// chapterDone pays a chapter's completion once and counts flawless chapters
// and finished courses.
func (s *Service) chapterDone(ctx context.Context, q *store.Queries, ev game.Graded, meta store.ExerciseChapterRow, out *game.RewardsOutcome) error {
	rows, err := q.ChapterProgress(ctx, store.ChapterProgressParams{UserID: ev.User, ChapterID: meta.ChapterID})
	if err != nil {
		return fmt.Errorf("chapter progress: %w", err)
	}
	flawless := true
	for _, r := range rows {
		if r.Status != "passed" {
			return nil
		}
		flawless = flawless && r.FirstTry
	}
	paid, err := credit(ctx, q, ev.User, ChapterBytes, "chapter", meta.ChapterID)
	if err != nil || !paid {
		return err
	}
	out.ChapterComplete = true
	out.Bytes = append(out.Bytes, game.Gain{Reason: "chapter", Label: "Chapter complete", Amount: ChapterBytes})
	if err := bump(ctx, q, ev.User, "flawless_chapters", b2i(flawless)); err != nil {
		return err
	}
	c, err := q.CourseChaptersDone(ctx, store.CourseChaptersDoneParams{UserID: ev.User, CourseID: meta.CourseID})
	if err != nil {
		return fmt.Errorf("course progress: %w", err)
	}
	if c.Chapters > 0 && c.Done == c.Chapters {
		first, err := credit(ctx, q, ev.User, 0, "course", meta.CourseID) // a zero entry marks it once
		if err != nil {
			return err
		}
		return bump(ctx, q, ev.User, "courses", b2i(first))
	}
	return nil
}

// bossDefeated opens the boss chest and reports the chapter it unlocks.
func (s *Service) bossDefeated(ctx context.Context, q *store.Queries, ev game.Graded, meta store.ExerciseChapterRow, out *game.RewardsOutcome) error {
	chest, err := s.openChest(ctx, q, ev.User, "boss", ev.Exercise.ID, BossBytes)
	if err != nil {
		return err
	}
	out.Boss.Chest = chest
	out.Bytes = append(out.Bytes, game.Gain{Reason: "boss", Label: "Boss defeated", Amount: chest.Bytes})
	if err := bump(ctx, q, ev.User, "bosses", 1); err != nil {
		return err
	}
	next, err := q.NextChapter(ctx, meta.ChapterID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("next chapter: %w", err)
	}
	out.Boss.UnlockedChapter = next
	return nil
}

// Today's quests and this week's quest for a user, the same all day: drawn by
// a hash of the user and the date, so nothing needs storing until progress.
func todaysQuests(user string, day, week time.Time) (daily []Quest, weekly Quest) {
	pool := append([]Quest{}, DailyPool...)
	seed := fnv(user + day.Format(time.DateOnly))
	for i := len(pool) - 1; i > 0; i-- { // Fisher-Yates with a deterministic seed
		seed = seed*6364136223846793005 + 1442695040888963407
		j := int((seed >> 33) % uint64(i+1))
		pool[i], pool[j] = pool[j], pool[i]
	}
	w := WeeklyPool[fnv(user+week.Format(time.DateOnly))%uint64(len(WeeklyPool))]
	return pool[:QuestsPerDay], w
}

func fnv(s string) uint64 {
	h := uint64(14695981039346656037)
	for i := 0; i < len(s); i++ {
		h ^= uint64(s[i])
		h *= 1099511628211
	}
	return h
}

func (s *Service) questProgress(ctx context.Context, q *store.Queries, ev game.Graded, typ string, goal int, tz string, out *game.RewardsOutcome) error {
	d := dayOf(s.now(), tz)
	todayLessons, err := q.DayLessons(ctx, store.DayLessonsParams{UserID: ev.User, FromDay: d.Today, ToDay: d.Today})
	if err != nil {
		return fmt.Errorf("today: %w", err)
	}
	reachedGoal := len(todayLessons) > 0 && int(todayLessons[0].Lessons) == goal
	if reachedGoal {
		if err := bump(ctx, q, ev.User, "goal_days", 1); err != nil {
			return err
		}
	}
	events := map[string]int{
		"lessons":   1,
		"first_try": b2i(ev.FirstPass && ev.Outcome.FirstTry),
		"xp":        ev.Outcome.XPAwarded,
		"code":      b2i(typ == "code" || typ == "sql" || typ == "project_cli"), // hands-on lessons
		"quiz":      b2i(typ == "mcq"),
		"goal_day":  b2i(reachedGoal),
	}
	daily, weekly := todaysQuests(ev.User.String(), d.Today, d.Week)
	allDone := true
	for _, qu := range append(daily, weekly) {
		period := d.Today
		if qu.Cadence == "weekly" {
			period = d.Week
		}
		n := events[qu.Event]
		row, err := q.BumpQuest(ctx, store.BumpQuestParams{UserID: ev.User, QuestID: qu.ID, Period: period, N: int32(n), Target: int32(qu.Target)})
		if err != nil {
			return fmt.Errorf("quest %s: %w", qu.ID, err)
		}
		done := row.DoneAt != nil
		justDone := done && int(row.Progress)-n < qu.Target
		if qu.Cadence == "daily" {
			allDone = allDone && done
		}
		if !justDone {
			continue
		}
		out.QuestsDone = append(out.QuestsDone, qu.Name)
		if qu.Bytes > 0 {
			paid, err := credit(ctx, q, ev.User, qu.Bytes, "weekly_quest", qu.ID+":"+period.Format(time.DateOnly))
			if err != nil {
				return err
			}
			if paid {
				out.Bytes = append(out.Bytes, game.Gain{Reason: "weekly_quest", Label: qu.Name, Amount: qu.Bytes})
			}
		}
	}
	out.ChestReady = allDone
	return nil
}
