-- name: CreateSubmission :one
INSERT INTO submissions (user_id, exercise_id, exercise_version, answer, status, score, result)
VALUES (@user_id, @exercise_id, @exercise_version, @answer, @status, @score, @result)
RETURNING id, created_at;

-- name: RecordAttempt :exec
INSERT INTO progress (user_id, exercise_id, status, first_passed_at, attempts)
VALUES (@user_id, @exercise_id,
        CASE WHEN @passed::bool THEN 'passed' ELSE 'started' END,
        CASE WHEN @passed::bool THEN now() END,
        1)
ON CONFLICT (user_id, exercise_id) DO UPDATE SET
    attempts = progress.attempts + 1,
    status = CASE WHEN progress.status = 'passed' OR @passed::bool THEN 'passed' ELSE 'started' END,
    first_passed_at = COALESCE(progress.first_passed_at, EXCLUDED.first_passed_at);

-- name: GetLatestSubmission :one
SELECT id, exercise_version, answer, status, score, result, created_at
FROM submissions
WHERE user_id = @user_id AND exercise_id = @exercise_id
ORDER BY created_at DESC
LIMIT 1;
