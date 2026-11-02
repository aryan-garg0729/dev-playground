// Package grading is the server-side submission pipeline: it loads the
// exercise, grades it with the type plugin, and stores the result. It is the
// only code that decides pass/fail.
package grading

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aryangarg/codequest/api/internal/exercise"
	"github.com/aryangarg/codequest/api/internal/store"
)

// Errors the API maps to HTTP statuses.
var (
	ErrNotFound     = errors.New("exercise not found")
	ErrStaleVersion = errors.New("exercise version changed")
)

// Graded is the outcome of one submission.
type Graded struct {
	SubmissionID uuid.UUID
	Result       exercise.Result
	CreatedAt    time.Time
}

// Service grades and records submissions.
type Service struct {
	pool *pgxpool.Pool
}

// NewService returns a grading service.
func NewService(pool *pgxpool.Pool) *Service { return &Service{pool: pool} }

// Submit grades sub for user and stores the submission and progress in one
// transaction. Invalid answers return exercise.ErrInvalidAnswer and are not stored.
func (s *Service) Submit(ctx context.Context, user uuid.UUID, sub exercise.Submission) (Graded, error) {
	q := store.New(s.pool)
	row, err := q.GetExerciseForGrading(ctx, sub.ExerciseID)
	if errors.Is(err, pgx.ErrNoRows) {
		return Graded{}, ErrNotFound
	}
	if err != nil {
		return Graded{}, fmt.Errorf("load exercise: %w", err)
	}
	if int(row.Version) != sub.Version {
		return Graded{}, fmt.Errorf("%w: submitted v%d, current v%d", ErrStaleVersion, sub.Version, row.Version)
	}
	typ, ok := exercise.Get(row.Type)
	if !ok {
		return Graded{}, fmt.Errorf("exercise %s has unregistered type %q", row.ID, row.Type)
	}
	spec, err := typ.LoadSpec(row.PrivateSpec)
	if err != nil {
		return Graded{}, fmt.Errorf("load spec %s: %w", row.ID, err)
	}
	ex := exercise.Exercise{
		ID: row.ID, Type: row.Type, Version: int(row.Version), Title: row.Title, XP: int(row.Xp),
		Difficulty: int(row.Difficulty), Tags: row.Tags, Free: row.Free, Boss: row.Boss, Spec: spec,
	}
	result, err := typ.Grade(ctx, ex, sub)
	if err != nil {
		return Graded{}, fmt.Errorf("grade %s: %w", row.ID, err)
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		return Graded{}, fmt.Errorf("encode result: %w", err)
	}
	var g Graded
	err = pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		tq := store.New(tx)
		created, err := tq.CreateSubmission(ctx, store.CreateSubmissionParams{
			UserID: user, ExerciseID: row.ID, ExerciseVersion: row.Version,
			Answer: sub.Answer, Status: result.Status, Score: result.Score, Result: resultJSON,
		})
		if err != nil {
			return fmt.Errorf("store submission: %w", err)
		}
		if err := tq.RecordAttempt(ctx, store.RecordAttemptParams{
			UserID: user, ExerciseID: row.ID, Passed: result.Status == exercise.StatusPass,
		}); err != nil {
			return fmt.Errorf("record progress: %w", err)
		}
		// Phase 4: game.Engine.OnGraded runs here, in the same transaction.
		g = Graded{SubmissionID: created.ID, Result: result, CreatedAt: created.CreatedAt}
		return nil
	})
	if err != nil {
		return Graded{}, err
	}
	return g, nil
}
