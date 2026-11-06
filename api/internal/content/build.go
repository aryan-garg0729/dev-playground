package content

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aryangarg/codequest/api/internal/store"
)

// BuildStats summarises what Build wrote.
type BuildStats struct {
	Courses, Chapters, Exercises int
}

// Build writes a validated bundle to Postgres in one transaction. Rows are
// upserted by id; anything no longer in the bundle is unpublished, never
// deleted, so progress and submissions keyed on old ids stay intact.
func Build(ctx context.Context, pool *pgxpool.Pool, b *Bundle) (BuildStats, error) {
	var st BuildStats
	err := pgx.BeginFunc(ctx, pool, func(tx pgx.Tx) error {
		q := store.New(tx)
		var courseIDs, chapterIDs, exerciseIDs []string
		for _, c := range b.Courses {
			if err := q.UpsertCourse(ctx, store.UpsertCourseParams{
				ID: c.ID, Title: c.Title, Description: c.Description,
				OrderIndex: int32(c.Order), Published: c.Published,
				Icon: c.Icon, Level: c.Level, Tags: nonNil(c.Tags), Topics: c.Topics,
			}); err != nil {
				return fmt.Errorf("course %s: %w", c.ID, err)
			}
			courseIDs = append(courseIDs, c.ID)
			for ci, ch := range c.Chapters {
				if err := q.UpsertChapter(ctx, store.UpsertChapterParams{
					ID: ch.ID, CourseID: c.ID, Title: ch.Title, OrderIndex: int32(ci),
				}); err != nil {
					return fmt.Errorf("chapter %s: %w", ch.ID, err)
				}
				chapterIDs = append(chapterIDs, ch.ID)
				for ei, ex := range ch.Exercises {
					for _, a := range ex.Assets {
						if err := q.UpsertContentAsset(ctx, store.UpsertContentAssetParams{
							Sha256: a.SHA256, Mime: a.MIME, Size: int32(len(a.Data)), Bytes: a.Data,
						}); err != nil {
							return fmt.Errorf("exercise %s: asset %s: %w", ex.ID, a.Name, err)
						}
					}
					assets, err := json.Marshal(AssetMap(ex.Assets))
					if err != nil {
						return fmt.Errorf("exercise %s: encode assets: %w", ex.ID, err)
					}
					if err := q.UpsertExercise(ctx, store.UpsertExerciseParams{
						ID: ex.ID, ChapterID: ch.ID, Type: ex.Type, Version: int32(ex.Version),
						Title: ex.Title, Xp: int32(ex.XP), Difficulty: int32(ex.Difficulty),
						Tags: ex.Tags, Free: ex.Free, Boss: ex.Boss, OrderIndex: int32(ei),
						PublicView: ex.PublicView, PrivateSpec: ex.PrivateSpec, Assets: assets,
					}); err != nil {
						return fmt.Errorf("exercise %s: %w", ex.ID, err)
					}
					exerciseIDs = append(exerciseIDs, ex.ID)
				}
			}
		}
		var pathIDs []string
		for i, p := range b.Paths {
			steps, err := json.Marshal(p.Steps)
			if err != nil {
				return fmt.Errorf("path %s: encode steps: %w", p.ID, err)
			}
			if err := q.UpsertCareerPath(ctx, store.UpsertCareerPathParams{
				ID: p.ID, Title: p.Title, Description: p.Description, Goal: p.Goal, OrderIndex: int32(i), Steps: steps,
			}); err != nil {
				return fmt.Errorf("path %s: %w", p.ID, err)
			}
			pathIDs = append(pathIDs, p.ID)
		}
		if err := q.UnpublishCareerPathsNotIn(ctx, nonNil(pathIDs)); err != nil {
			return fmt.Errorf("unpublish paths: %w", err)
		}
		if err := q.UnpublishExercisesNotIn(ctx, nonNil(exerciseIDs)); err != nil {
			return fmt.Errorf("unpublish exercises: %w", err)
		}
		if err := q.UnpublishChaptersNotIn(ctx, nonNil(chapterIDs)); err != nil {
			return fmt.Errorf("unpublish chapters: %w", err)
		}
		if err := q.UnpublishCoursesNotIn(ctx, nonNil(courseIDs)); err != nil {
			return fmt.Errorf("unpublish courses: %w", err)
		}
		st = BuildStats{Courses: len(courseIDs), Chapters: len(chapterIDs), Exercises: len(exerciseIDs)}
		return nil
	})
	if err != nil {
		return BuildStats{}, fmt.Errorf("build content: %w", err)
	}
	return st, nil
}

func nonNil(s []string) []string {
	if s == nil {
		return []string{}
	}
	return s
}
