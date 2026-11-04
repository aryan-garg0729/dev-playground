-- name: CreateCLILoginCode :exec
INSERT INTO cli_login_codes (code_hash, user_id, expires_at) VALUES (@code_hash, @user_id, @expires_at);

-- name: UseCLILoginCode :one
-- Single use: only the first call for an unexpired code returns a row.
UPDATE cli_login_codes SET used_at = now()
WHERE code_hash = @code_hash AND used_at IS NULL AND expires_at > now()
RETURNING user_id;

-- name: CreateCLIToken :exec
INSERT INTO cli_tokens (id, user_id, label, expires_at) VALUES (@id, @user_id, @label, @expires_at);

-- name: TouchCLIToken :one
-- Looks up an unexpired token and slides its expiry forward.
UPDATE cli_tokens SET last_used_at = now(), expires_at = @expires_at
WHERE id = @id AND expires_at > now()
RETURNING user_id;

-- name: DeleteCLIToken :exec
DELETE FROM cli_tokens WHERE id = @id;

-- name: DeleteExpiredCLIAuth :exec
WITH t AS (DELETE FROM cli_tokens WHERE expires_at < now())
DELETE FROM cli_login_codes WHERE expires_at < now() - interval '1 day';

-- name: NextCLILesson :one
-- The learner's first unpassed project_cli lesson, in course order.
SELECT e.id, e.title, c.id AS course_id, c.title AS course_title,
    ch.order_index AS chapter_index, e.order_index AS lesson_index
FROM exercises e
JOIN chapters ch ON ch.id = e.chapter_id
JOIN courses c ON c.id = ch.course_id
WHERE e.type = 'project_cli' AND e.published AND ch.published AND c.published
  AND NOT EXISTS (SELECT 1 FROM progress p
                  WHERE p.user_id = @user_id AND p.exercise_id = e.id AND p.status = 'passed')
ORDER BY c.order_index, c.id, ch.order_index, e.order_index
LIMIT 1;
