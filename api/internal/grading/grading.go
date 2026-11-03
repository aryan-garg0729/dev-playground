// Package grading is the server-side submission pipeline: it loads the
// exercise, grades it with the type plugin, and stores the result. It is the
// only code that decides pass/fail.
package grading

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aryangarg/codequest/api/internal/exercise"
	"github.com/aryangarg/codequest/api/internal/game"
	"github.com/aryangarg/codequest/api/internal/runner"
	"github.com/aryangarg/codequest/api/internal/store"
)

// Errors the API maps to HTTP statuses.
var (
	ErrNotFound     = errors.New("exercise not found")
	ErrStaleVersion = errors.New("exercise version changed")
)

// Graded is the outcome of one submission. When Pending is true, grading is
// still running; poll the submission by id.
type Graded struct {
	SubmissionID uuid.UUID
	Result       exercise.Result
	Outcome      game.Outcome
	CreatedAt    time.Time
	Pending      bool
}

// Options tune the service. Zero values take the defaults.
type Options struct {
	// Wait is how long Submit waits for a result before answering "pending" (15s).
	Wait time.Duration
	// MaxGrade bounds one grading run, queue time included (2m).
	MaxGrade time.Duration
	Log      *slog.Logger
	// Game applies XP, levels and streaks in the grading transaction (a
	// default Engine when nil). Board, when set, gets the league update after commit.
	Game  *game.Engine
	Board *game.Board
}

// Service grades and records submissions.
type Service struct {
	pool   *pgxpool.Pool
	runner runner.Runner
	opt    Options
	wg     sync.WaitGroup
}

// NewService returns a grading service. run executes code for runner-based
// types; it may be nil when only in-process types (mcq) are graded.
func NewService(pool *pgxpool.Pool, run runner.Runner, opt Options) *Service {
	if opt.Wait <= 0 {
		opt.Wait = 15 * time.Second
	}
	if opt.MaxGrade <= 0 {
		opt.MaxGrade = 2 * time.Minute
	}
	if opt.Log == nil {
		opt.Log = slog.Default()
	}
	if opt.Game == nil {
		opt.Game = &game.Engine{}
	}
	return &Service{pool: pool, runner: run, opt: opt}
}

// Close waits for background grading to finish.
func (s *Service) Close() { s.wg.Wait() }

type outcome struct {
	res exercise.Result
	err error
}

func (s *Service) load(ctx context.Context, sub exercise.Submission) (store.GetExerciseForGradingRow, exercise.Type, exercise.Exercise, error) {
	q := store.New(s.pool)
	row, err := q.GetExerciseForGrading(ctx, sub.ExerciseID)
	if errors.Is(err, pgx.ErrNoRows) {
		return row, nil, exercise.Exercise{}, ErrNotFound
	}
	if err != nil {
		return row, nil, exercise.Exercise{}, fmt.Errorf("load exercise: %w", err)
	}
	if int(row.Version) != sub.Version {
		return row, nil, exercise.Exercise{}, fmt.Errorf("%w: submitted v%d, current v%d", ErrStaleVersion, sub.Version, row.Version)
	}
	typ, ok := exercise.Get(row.Type)
	if !ok {
		return row, nil, exercise.Exercise{}, fmt.Errorf("exercise %s has unregistered type %q", row.ID, row.Type)
	}
	spec, err := typ.LoadSpec(row.PrivateSpec)
	if err != nil {
		return row, nil, exercise.Exercise{}, fmt.Errorf("load spec %s: %w", row.ID, err)
	}
	ex := exercise.Exercise{
		ID: row.ID, Type: row.Type, Version: int(row.Version), Title: row.Title, XP: int(row.Xp),
		Difficulty: int(row.Difficulty), Tags: row.Tags, Free: row.Free, Boss: row.Boss, Spec: spec,
	}
	return row, typ, ex, nil
}

// Run executes an ungraded server run (`run: server`): visible checks only,
// nothing stored, no progress. exercise.ErrRunUnavailable when the exercise
// does not run on the server.
func (s *Service) Run(ctx context.Context, sub exercise.Submission) (exercise.RunOutput, error) {
	_, typ, ex, err := s.load(ctx, sub)
	if err != nil {
		return exercise.RunOutput{}, err
	}
	r, ok := typ.(exercise.Runnable)
	if !ok {
		return exercise.RunOutput{}, exercise.ErrRunUnavailable
	}
	if s.runner != nil {
		ctx = runner.NewContext(ctx, s.runner)
	}
	out, err := r.Run(ctx, ex, sub.Answer)
	if err != nil {
		return exercise.RunOutput{}, fmt.Errorf("run %s: %w", ex.ID, err)
	}
	return out, nil
}

// Submit grades sub for user. A result within Options.Wait is stored with the
// progress update in one transaction and returned. Otherwise a pending
// submission is stored and returned, and grading finishes in the background.
// Invalid answers return exercise.ErrInvalidAnswer and are not stored; so are
// runner outages (runner.ErrUnavailable) that happen before the wait ends.
func (s *Service) Submit(ctx context.Context, user uuid.UUID, sub exercise.Submission) (Graded, error) {
	q := store.New(s.pool)
	row, typ, ex, err := s.load(ctx, sub)
	if err != nil {
		return Graded{}, err
	}

	// Grade on a context that outlives the request, so a slow run is not lost
	// when the browser stops waiting.
	gctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), s.opt.MaxGrade)
	if s.runner != nil {
		gctx = runner.NewContext(gctx, s.runner)
	}
	done := make(chan outcome, 1)
	s.wg.Go(func() {
		defer cancel()
		res, err := typ.Grade(gctx, ex, sub)
		done <- outcome{res, err}
	})

	timer := time.NewTimer(s.opt.Wait)
	defer timer.Stop()
	select {
	case o := <-done:
		if o.err != nil {
			return Graded{}, fmt.Errorf("grade %s: %w", row.ID, o.err)
		}
		return s.store(ctx, user, ex, sub.Answer, o.res)
	case <-timer.C:
	case <-ctx.Done():
	}

	// Still grading: store a pending submission and finish it later.
	created, err := q.CreateSubmission(context.WithoutCancel(ctx), store.CreateSubmissionParams{
		UserID: user, ExerciseID: row.ID, ExerciseVersion: row.Version,
		Answer: sub.Answer, Status: "pending", Score: 0, Result: nil,
	})
	if err != nil {
		return Graded{}, fmt.Errorf("store pending submission: %w", err)
	}
	s.wg.Go(func() { s.finish(created.ID, user, ex, <-done) })
	return Graded{SubmissionID: created.ID, CreatedAt: created.CreatedAt, Pending: true}, nil
}

func (s *Service) store(ctx context.Context, user uuid.UUID, ex exercise.Exercise, answer json.RawMessage, res exercise.Result) (Graded, error) {
	resultJSON, err := json.Marshal(res)
	if err != nil {
		return Graded{}, fmt.Errorf("encode result: %w", err)
	}
	var g Graded
	err = pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		out, outJSON, err := s.applyGame(ctx, tx, user, ex, res)
		if err != nil {
			return err
		}
		created, err := store.New(tx).CreateSubmission(ctx, store.CreateSubmissionParams{
			UserID: user, ExerciseID: ex.ID, ExerciseVersion: int32(ex.Version),
			Answer: answer, Status: res.Status, Score: res.Score, Result: resultJSON, Outcome: outJSON,
		})
		if err != nil {
			return fmt.Errorf("store submission: %w", err)
		}
		if err := storeHiddenInputs(ctx, tx, created.ID, res); err != nil {
			return err
		}
		g = Graded{SubmissionID: created.ID, Result: res, Outcome: out, CreatedAt: created.CreatedAt}
		return nil
	})
	if err != nil {
		return Graded{}, err
	}
	s.afterCommit(ctx, g.Outcome)
	return g, nil
}

// applyGame records the attempt and runs the game engine, in tx. The engine
// reads the Result; nothing here changes it.
func (s *Service) applyGame(ctx context.Context, tx pgx.Tx, user uuid.UUID, ex exercise.Exercise, res exercise.Result) (game.Outcome, json.RawMessage, error) {
	p, err := store.New(tx).RecordAttempt(ctx, store.RecordAttemptParams{
		UserID: user, ExerciseID: ex.ID, Passed: res.Status == exercise.StatusPass,
	})
	if err != nil {
		return game.Outcome{}, nil, fmt.Errorf("record progress: %w", err)
	}
	out, err := s.opt.Game.OnGraded(ctx, tx, user, ex, res, game.Progress{Attempts: int(p.Attempts), UsedHint: p.UsedHint, Peeked: p.Peeked})
	if err != nil {
		return game.Outcome{}, nil, fmt.Errorf("game: %w", err)
	}
	outJSON, err := json.Marshal(out)
	if err != nil {
		return game.Outcome{}, nil, fmt.Errorf("encode outcome: %w", err)
	}
	return out, outJSON, nil
}

// afterCommit updates the league board. It is idempotent (the board is
// rebuilt from Postgres), so a retry or a duplicate call is harmless.
func (s *Service) afterCommit(ctx context.Context, out game.Outcome) {
	if s.opt.Board == nil || out.League == nil {
		return
	}
	if err := s.opt.Board.Sync(context.WithoutCancel(ctx), s.pool, *out.League); err != nil {
		s.opt.Log.Error("league board", "err", err)
	}
}

// finish completes a pending submission. Grading failures are stored as an
// error result and do not count as an attempt.
func (s *Service) finish(id, user uuid.UUID, ex exercise.Exercise, o outcome) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	res, attempt := o.res, true
	if o.err != nil {
		s.opt.Log.Error("background grading failed", "submission", id, "exercise", ex.ID, "err", o.err)
		res = ErrorResult("Grading failed; submit again")
		attempt = false
	}
	resultJSON, err := json.Marshal(res)
	if err != nil {
		s.opt.Log.Error("encode result", "submission", id, "err", err)
		return
	}
	var out game.Outcome
	err = pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		var outJSON json.RawMessage
		if attempt {
			var err error
			if out, outJSON, err = s.applyGame(ctx, tx, user, ex, res); err != nil {
				return err
			}
		}
		if err := store.New(tx).FinishSubmission(ctx, store.FinishSubmissionParams{
			ID: id, Status: res.Status, Score: res.Score, Result: resultJSON, Outcome: outJSON,
		}); err != nil {
			return fmt.Errorf("finish submission: %w", err)
		}
		return storeHiddenInputs(ctx, tx, id, res)
	})
	if err != nil {
		s.opt.Log.Error("store background result", "submission", id, "err", err)
		return
	}
	s.afterCommit(ctx, out)
}

// storeHiddenInputs keeps the declared inputs of failed hidden checks
// ({"<row index>": input}) in a server-only column, for Peek.
func storeHiddenInputs(ctx context.Context, tx pgx.Tx, id uuid.UUID, res exercise.Result) error {
	inputs := map[string]string{}
	for i, c := range res.Checks {
		if !c.Passed && !c.Visible && c.Input != "" {
			inputs[strconv.Itoa(i)] = c.Input
		}
	}
	if len(inputs) == 0 {
		return nil
	}
	b, err := json.Marshal(inputs)
	if err != nil {
		return fmt.Errorf("encode hidden inputs: %w", err)
	}
	if err := store.New(tx).SetHiddenInputs(ctx, store.SetHiddenInputsParams{ID: id, HiddenInputs: b}); err != nil {
		return fmt.Errorf("store hidden inputs: %w", err)
	}
	return nil
}

// ErrorResult is a Result for a grading failure that is not the learner's.
func ErrorResult(msg string) exercise.Result {
	return exercise.Result{
		Status: exercise.StatusError,
		Checks: []exercise.CheckResult{{Name: "Grader", Passed: false, Visible: true, Feedback: msg}},
	}
}

// Submission is a stored submission as the owner may see it.
type Submission struct {
	ID        uuid.UUID
	Exercise  string
	Version   int
	Status    string
	Result    json.RawMessage // nil while pending
	Outcome   json.RawMessage // game.Outcome; nil while pending or after an error
	CreatedAt time.Time
}

// ErrNoSubmission means the id is unknown or belongs to someone else.
var ErrNoSubmission = errors.New("submission not found")

// Get returns user's submission id.
func (s *Service) Get(ctx context.Context, user, id uuid.UUID) (Submission, error) {
	r, err := store.New(s.pool).GetSubmission(ctx, store.GetSubmissionParams{ID: id, UserID: user})
	if errors.Is(err, pgx.ErrNoRows) {
		return Submission{}, ErrNoSubmission
	}
	if err != nil {
		return Submission{}, fmt.Errorf("get submission: %w", err)
	}
	return Submission{
		ID: r.ID, Exercise: r.ExerciseID, Version: int(r.ExerciseVersion), Status: r.Status,
		Result: r.Result, Outcome: r.Outcome, CreatedAt: r.CreatedAt,
	}, nil
}

// FailStale marks submissions left pending for longer than age (e.g. after a
// restart) as errors. Run it at startup.
func (s *Service) FailStale(ctx context.Context, age time.Duration) (int64, error) {
	n, err := store.New(s.pool).FailStalePending(ctx, age.Seconds())
	if err != nil {
		return 0, fmt.Errorf("fail stale submissions: %w", err)
	}
	return n, nil
}
