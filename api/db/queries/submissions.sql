-- name: CreateSubmission :one
INSERT INTO submissions (user_id, exercise_id, exercise_version, answer, status, score, result, outcome)
VALUES (@user_id, @exercise_id, @exercise_version, @answer, @status, @score, @result, @outcome)
RETURNING id, created_at;

-- name: RecordAttempt :one
-- Returns the progress after this attempt (attempts counts it).
INSERT INTO progress (user_id, exercise_id, status, first_passed_at, attempts)
VALUES (@user_id, @exercise_id,
        CASE WHEN @passed::bool THEN 'passed' ELSE 'started' END,
        CASE WHEN @passed::bool THEN now() END,
        1)
ON CONFLICT (user_id, exercise_id) DO UPDATE SET
    attempts = progress.attempts + 1,
    status = CASE WHEN progress.status = 'passed' OR @passed::bool THEN 'passed' ELSE 'started' END,
    first_passed_at = COALESCE(progress.first_passed_at, EXCLUDED.first_passed_at)
RETURNING attempts, used_hint, peeked, first_try;

-- name: GetLatestSubmission :one
SELECT id, exercise_version, answer, status, score, result, created_at
FROM submissions
WHERE user_id = @user_id AND exercise_id = @exercise_id
ORDER BY created_at DESC
LIMIT 1;

-- name: FinishSubmission :exec
UPDATE submissions SET status = @status, score = @score, result = @result, outcome = @outcome
WHERE id = @id AND status = 'pending';

-- name: GetSubmission :one
SELECT id, exercise_id, exercise_version, status, score, result, outcome, created_at
FROM submissions
WHERE id = @id AND user_id = @user_id;

-- name: FailStalePending :execrows
UPDATE submissions
SET status = 'error',
    result = '{"status":"error","score":0,"checks":[{"name":"Grader","passed":false,"visible":true,"feedback":"Grading was interrupted; submit again"}],"stdout":"","runtime_ms":0}'
WHERE status = 'pending' AND created_at < now() - make_interval(secs => @older_than_secs::float8);
