-- name: GetUserJoined :one
SELECT created_at FROM users WHERE id = @id;

-- name: ListBosses :many
-- Every published boss (a profile "project") with the user's progress, in course order.
SELECT e.id, e.title, ch.course_id, c.title AS course_title, ch.title AS chapter_title,
    (SELECT count(*) FROM chapters c2
      WHERE c2.course_id = ch.course_id AND c2.published AND c2.order_index <= ch.order_index)::int AS chapter_index,
    COALESCE(p.status, '')::text AS status,
    COALESCE(p.first_try, false)::bool AS first_try,
    p.first_passed_at
FROM exercises e
JOIN chapters ch ON ch.id = e.chapter_id
JOIN courses c ON c.id = ch.course_id
LEFT JOIN progress p ON p.exercise_id = e.id AND p.user_id = @user_id
WHERE e.boss AND e.published AND ch.published AND c.published
ORDER BY c.order_index, c.id, ch.order_index, e.order_index;

-- name: ChapterSkills :many
-- Lessons per chapter of a course: total, passed, passed on the first try.
SELECT ch.id, ch.title,
    count(e.id)::int AS lessons,
    count(*) FILTER (WHERE p.status = 'passed')::int AS passed,
    count(*) FILTER (WHERE p.status = 'passed' AND p.first_try)::int AS first_try
FROM chapters ch
JOIN exercises e ON e.chapter_id = ch.id AND e.published
LEFT JOIN progress p ON p.exercise_id = e.id AND p.user_id = @user_id
WHERE ch.course_id = @course_id AND ch.published
GROUP BY ch.id, ch.title, ch.order_index
ORDER BY ch.order_index;

-- name: RecentPasses :many
-- The user's latest passed lessons, one row per lesson per local day.
SELECT d.day, e.id, e.title, e.boss, ch.course_id, c.title AS course_title
FROM day_passes d
JOIN exercises e ON e.id = d.exercise_id
JOIN chapters ch ON ch.id = e.chapter_id
JOIN courses c ON c.id = ch.course_id
WHERE d.user_id = @user_id AND e.published AND ch.published AND c.published
ORDER BY d.day DESC, ch.order_index DESC, e.order_index DESC
LIMIT 30;
