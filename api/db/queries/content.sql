-- name: UpsertCourse :exec
INSERT INTO courses (id, title, description, order_index, published, icon, level, tags, topics)
VALUES (@id, @title, @description, @order_index, @published, @icon, @level, @tags, @topics)
ON CONFLICT (id) DO UPDATE SET
    title = EXCLUDED.title,
    description = EXCLUDED.description,
    order_index = EXCLUDED.order_index,
    published = EXCLUDED.published,
    icon = EXCLUDED.icon,
    level = EXCLUDED.level,
    tags = EXCLUDED.tags,
    topics = EXCLUDED.topics;

-- name: UpsertCareerPath :exec
INSERT INTO career_paths (id, title, description, goal, order_index, steps, published)
VALUES (@id, @title, @description, @goal, @order_index, @steps, true)
ON CONFLICT (id) DO UPDATE SET
    title = EXCLUDED.title,
    description = EXCLUDED.description,
    goal = EXCLUDED.goal,
    order_index = EXCLUDED.order_index,
    steps = EXCLUDED.steps,
    published = true;

-- name: UnpublishCareerPathsNotIn :exec
UPDATE career_paths SET published = false WHERE NOT (id = ANY (@ids::text[]));

-- name: ListCareerPaths :many
SELECT id, title, description, goal, steps FROM career_paths
WHERE published
ORDER BY order_index, id;

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
                       order_index, public_view, private_spec, assets, published)
VALUES (@id, @chapter_id, @type, @version, @title, @xp, @difficulty, @tags, @free, @boss,
        @order_index, @public_view, @private_spec, @assets, true)
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
    assets = EXCLUDED.assets,
    published = true;

-- name: UpsertContentAsset :exec
-- Content-addressed: the same bytes always have the same row.
INSERT INTO content_assets (sha256, mime, size, bytes)
VALUES (@sha256, @mime, @size, @bytes)
ON CONFLICT (sha256) DO NOTHING;

-- name: GetContentAsset :one
-- Public: lesson images are served to anyone who has the hash.
SELECT mime, bytes FROM content_assets WHERE sha256 = @sha256;

-- name: UnpublishCoursesNotIn :exec
UPDATE courses SET published = false WHERE NOT (id = ANY (@ids::text[]));

-- name: UnpublishChaptersNotIn :exec
UPDATE chapters SET published = false WHERE NOT (id = ANY (@ids::text[]));

-- name: UnpublishExercisesNotIn :exec
UPDATE exercises SET published = false WHERE NOT (id = ANY (@ids::text[]));

-- name: ListCourses :many
SELECT c.id, c.title, c.description, c.icon, c.level, c.tags, c.topics,
    (SELECT count(*) FROM chapters ch
      WHERE ch.course_id = c.id AND ch.published)::int AS chapter_count,
    (SELECT count(*) FROM exercises e JOIN chapters ch ON ch.id = e.chapter_id
      WHERE ch.course_id = c.id AND ch.published AND e.published)::int AS exercise_count,
    (SELECT count(*) FROM progress p
       JOIN exercises e ON e.id = p.exercise_id
       JOIN chapters ch ON ch.id = e.chapter_id
      WHERE ch.course_id = c.id AND ch.published AND e.published
        AND p.status = 'passed' AND p.user_id = sqlc.narg('user_id'))::int AS passed_count,
    -- Chapters with lessons, every one of them passed.
    (SELECT count(*) FROM chapters ch
      WHERE ch.course_id = c.id AND ch.published
        AND EXISTS (SELECT 1 FROM exercises e WHERE e.chapter_id = ch.id AND e.published)
        AND NOT EXISTS (
            SELECT 1 FROM exercises e
            LEFT JOIN progress p ON p.exercise_id = e.id AND p.user_id = sqlc.narg('user_id')
            WHERE e.chapter_id = ch.id AND e.published AND COALESCE(p.status, '') <> 'passed'))::int AS chapters_done,
    (SELECT count(*) FROM exercises e JOIN chapters ch ON ch.id = e.chapter_id
      WHERE ch.course_id = c.id AND ch.published AND e.published AND e.free)::int AS free_count,
    EXISTS (SELECT 1 FROM exercises e JOIN chapters ch ON ch.id = e.chapter_id
      WHERE ch.course_id = c.id AND ch.published AND e.published AND e.type = 'project_cli')::bool AS has_cli,
    EXISTS (SELECT 1 FROM exercises e JOIN chapters ch ON ch.id = e.chapter_id
      WHERE ch.course_id = c.id AND ch.published AND e.published AND e.type <> 'project_cli')::bool AS has_browser
FROM courses c
WHERE c.published
ORDER BY c.order_index, c.id;

-- name: CourseCurrentChapters :many
-- Per course, the chapter of the user's first unpassed lesson (1-based index).
SELECT DISTINCT ON (ch.course_id) ch.course_id, ch.title,
    (SELECT count(*) FROM chapters c2
      WHERE c2.course_id = ch.course_id AND c2.published AND c2.order_index <= ch.order_index)::int AS chapter_index
FROM exercises e
JOIN chapters ch ON ch.id = e.chapter_id
LEFT JOIN progress p ON p.exercise_id = e.id AND p.user_id = @user_id
WHERE ch.published AND e.published AND COALESCE(p.status, '') <> 'passed'
ORDER BY ch.course_id, ch.order_index, e.order_index;

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
       e.public_view, e.assets, ch.course_id, ch.title AS chapter_title
FROM exercises e
JOIN chapters ch ON ch.id = e.chapter_id
WHERE e.id = @id AND e.published AND ch.published;

-- name: GetExerciseForGrading :one
-- Server-only: includes private_spec. Never return this row to the browser.
SELECT id, type, version, title, xp, difficulty, tags, free, boss, private_spec
FROM exercises
WHERE id = @id AND published;
