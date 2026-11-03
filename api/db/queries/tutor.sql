-- name: InsertTutorMessage :exec
-- created_at is set by the caller: a question and its reply are saved in one
-- transaction, where now() would give both the same time and an unstable order.
INSERT INTO tutor_messages (user_id, exercise_id, role, content, input_tokens, output_tokens, model, created_at)
VALUES (@user_id, @exercise_id, @role, @content, @input_tokens, @output_tokens, @model, @created_at);

-- name: ListTutorMessages :many
-- The most recent messages for one exercise, newest first.
SELECT role, content, created_at FROM tutor_messages
WHERE user_id = @user_id AND exercise_id = @exercise_id
ORDER BY created_at DESC, id DESC
LIMIT @lim;

-- name: MarkHintUsed :exec
-- A tutor hint removes the no-hint bonus for this exercise.
INSERT INTO progress (user_id, exercise_id, status, used_hint)
VALUES (@user_id, @exercise_id, 'started', true)
ON CONFLICT (user_id, exercise_id) DO UPDATE SET used_hint = true;

-- name: GetLatestResult :one
SELECT result FROM submissions
WHERE user_id = @user_id AND exercise_id = @exercise_id AND result IS NOT NULL
ORDER BY created_at DESC
LIMIT 1;
