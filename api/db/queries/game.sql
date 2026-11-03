-- name: EnsureUserStats :exec
INSERT INTO user_stats (user_id) VALUES (@user_id) ON CONFLICT (user_id) DO NOTHING;

-- name: LockUserStats :one
-- Serialises game updates for one user (concurrent submissions).
SELECT s.total_xp, s.level, s.streak_days, s.streak_freezes, s.last_active_date, s.league_tier, u.timezone
FROM user_stats s JOIN users u ON u.id = s.user_id
WHERE s.user_id = @user_id
FOR UPDATE OF s;

-- name: GetUserStats :one
SELECT s.total_xp, s.level, s.bytes, s.streak_days, s.streak_freezes, s.last_active_date,
       s.league_tier, s.title, s.daily_goal, u.timezone
FROM user_stats s JOIN users u ON u.id = s.user_id
WHERE s.user_id = @user_id;

-- name: UpdateUserStats :exec
UPDATE user_stats SET
    total_xp = @total_xp,
    level = @level,
    streak_days = @streak_days,
    streak_freezes = @streak_freezes,
    last_active_date = @last_active_date
WHERE user_id = @user_id;

-- name: SetDailyGoal :exec
UPDATE user_stats SET daily_goal = @daily_goal WHERE user_id = @user_id;

-- name: SetLeagueTier :exec
UPDATE user_stats SET league_tier = @league_tier WHERE user_id = @user_id;

-- name: AwardXP :one
-- Inserts the event unless (user, reason, ref) already has one; returns no row then.
INSERT INTO xp_events (user_id, amount, reason, ref_id)
VALUES (@user_id, @amount, @reason, @ref_id)
ON CONFLICT (user_id, reason, ref_id) DO NOTHING
RETURNING id;

-- name: RecordDayPass :exec
INSERT INTO day_passes (user_id, day, exercise_id)
VALUES (@user_id, @day, @exercise_id)
ON CONFLICT DO NOTHING;

-- name: ListDayPasses :many
-- Lessons passed per day in [from, to].
SELECT day, count(*)::int AS lessons
FROM day_passes
WHERE user_id = @user_id AND day BETWEEN @from_day AND @to_day
GROUP BY day
ORDER BY day;

-- name: CourseXP :one
-- XP earned in a course, and the base XP on offer.
SELECT
    COALESCE((SELECT sum(x.amount) FROM xp_events x
              JOIN exercises e ON e.id = x.ref_id
              JOIN chapters ch ON ch.id = e.chapter_id
              WHERE x.user_id = sqlc.narg('user_id') AND x.reason = 'pass' AND ch.course_id = @course_id), 0)::int AS earned,
    COALESCE((SELECT sum(e.xp) FROM exercises e
              JOIN chapters ch ON ch.id = e.chapter_id
              WHERE ch.course_id = @course_id AND ch.published AND e.published), 0)::int AS total;

-- name: CourseTags :many
-- Distinct tags in course order (first appearance).
SELECT tag::text AS tag FROM (
    SELECT t.tag, min(ch.order_index * 100000 + e.order_index * 100 + t.ord) AS first
    FROM exercises e
    JOIN chapters ch ON ch.id = e.chapter_id
    CROSS JOIN LATERAL unnest(e.tags) WITH ORDINALITY AS t(tag, ord)
    WHERE ch.course_id = @course_id AND ch.published AND e.published
    GROUP BY t.tag
) s
ORDER BY first
LIMIT 12;

-- name: LatestActivity :one
-- The exercise the user most recently submitted.
SELECT s.exercise_id, s.status, s.answer, s.result, e.type, ch.course_id
FROM submissions s
JOIN exercises e ON e.id = s.exercise_id
JOIN chapters ch ON ch.id = e.chapter_id
WHERE s.user_id = @user_id AND e.published AND ch.published
ORDER BY s.created_at DESC
LIMIT 1;
