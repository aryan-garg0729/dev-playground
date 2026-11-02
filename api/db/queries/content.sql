-- name: UpsertCourse :exec
INSERT INTO courses (id, title, description, order_index, published)
VALUES (@id, @title, @description, @order_index, @published)
ON CONFLICT (id) DO UPDATE SET
    title = EXCLUDED.title,
    description = EXCLUDED.description,
    order_index = EXCLUDED.order_index,
    published = EXCLUDED.published;

-- name: UpsertChapter :exec
INSERT INTO chapters (id, course_id, title, order_index, published)
VALUES (@id, @course_id, @title, @order_index, true)
ON CONFLICT (id) DO UPDATE SET
    course_id = EXCLUDED.course_id,
    title = EXCLUDED.title,
    order_index = EXCLUDED.order_index,
    published = true;

-- name: UpsertExercise :exec
INSERT INTO exercises (id, chapter_id, type, version, title, xp, difficulty, tags, free, boss,
                       order_index, public_view, private_spec, published)
VALUES (@id, @chapter_id, @type, @version, @title, @xp, @difficulty, @tags, @free, @boss,
        @order_index, @public_view, @private_spec, true)
ON CONFLICT (id) DO UPDATE SET
    chapter_id = EXCLUDED.chapter_id,
    type = EXCLUDED.type,
    version = EXCLUDED.version,
    title = EXCLUDED.title,
    xp = EXCLUDED.xp,
    difficulty = EXCLUDED.difficulty,
    tags = EXCLUDED.tags,
    free = EXCLUDED.free,
    boss = EXCLUDED.boss,
    order_index = EXCLUDED.order_index,
    public_view = EXCLUDED.public_view,
    private_spec = EXCLUDED.private_spec,
    published = true;

-- name: UnpublishCoursesNotIn :exec
UPDATE courses SET published = false WHERE NOT (id = ANY (@ids::text[]));

-- name: UnpublishChaptersNotIn :exec
UPDATE chapters SET published = false WHERE NOT (id = ANY (@ids::text[]));

-- name: UnpublishExercisesNotIn :exec
UPDATE exercises SET published = false WHERE NOT (id = ANY (@ids::text[]));

-- name: ListCourses :many
SELECT c.id, c.title, c.description,
    (SELECT count(*) FROM chapters ch
      WHERE ch.course_id = c.id AND ch.published)::int AS chapter_count,
    (SELECT count(*) FROM exercises e JOIN chapters ch ON ch.id = e.chapter_id
      WHERE ch.course_id = c.id AND ch.published AND e.published)::int AS exercise_count,
    (SELECT count(*) FROM progress p
       JOIN exercises e ON e.id = p.exercise_id
       JOIN chapters ch ON ch.id = e.chapter_id
      WHERE ch.course_id = c.id AND ch.published AND e.published
        AND p.status = 'passed' AND p.user_id = sqlc.narg('user_id'))::int AS passed_count
FROM courses c
WHERE c.published
ORDER BY c.order_index, c.id;

-- name: GetCourse :one
SELECT id, title, description FROM courses WHERE id = @id AND published;

-- name: ListChapters :many
SELECT id, title FROM chapters
WHERE course_id = @course_id AND published
ORDER BY order_index;

-- name: ListCourseExercises :many
-- Browser-facing: never selects private_spec.
SELECT e.id, e.chapter_id, e.type, e.title, e.xp, e.difficulty, e.free, e.boss,
       COALESCE(p.status, '')::text AS progress_status
FROM exercises e
JOIN chapters ch ON ch.id = e.chapter_id
LEFT JOIN progress p ON p.exercise_id = e.id AND p.user_id = sqlc.narg('user_id')
WHERE ch.course_id = @course_id AND ch.published AND e.published
ORDER BY ch.order_index, e.order_index;

-- name: GetExercise :one
-- Browser-facing: never selects private_spec.
SELECT e.id, e.chapter_id, e.type, e.version, e.title, e.xp, e.difficulty, e.tags, e.free, e.boss,
       e.public_view, ch.course_id, ch.title AS chapter_title
FROM exercises e
JOIN chapters ch ON ch.id = e.chapter_id
WHERE e.id = @id AND e.published AND ch.published;

-- name: GetExerciseForGrading :one
-- Server-only: includes private_spec. Never return this row to the browser.
SELECT id, type, version, title, xp, difficulty, tags, free, boss, private_spec
FROM exercises
WHERE id = @id AND published;
