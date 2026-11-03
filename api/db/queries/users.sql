-- name: EnsureUserByEmail :one
INSERT INTO users (email, name)
VALUES (@email, @name)
ON CONFLICT (email) DO UPDATE SET email = EXCLUDED.email
RETURNING id, email, name, avatar, timezone, goal, onboarded_at;

-- name: GetUser :one
SELECT id, email, name, avatar, timezone, goal, onboarded_at FROM users WHERE id = @id;

-- name: GetUserByEmail :one
SELECT id, email, name, avatar, timezone, goal, onboarded_at FROM users WHERE email = @email;

-- name: CreateUser :one
INSERT INTO users (email, name, avatar)
VALUES (@email, @name, @avatar)
RETURNING id, email, name, avatar, timezone, goal, onboarded_at;

-- name: UpdateUserProfile :exec
-- Fills name and avatar only where they are still empty (sign-in never overwrites edits).
UPDATE users SET
    name = CASE WHEN name = '' THEN @name ELSE name END,
    avatar = CASE WHEN avatar = '' THEN @avatar ELSE avatar END
WHERE id = @id;

-- name: SetUserTimezone :exec
UPDATE users SET timezone = @timezone WHERE id = @id;

-- name: SetUserName :exec
UPDATE users SET name = @name WHERE id = @id;

-- name: SetUserGoal :exec
UPDATE users SET goal = @goal WHERE id = @id;

-- name: SetUserOnboarded :exec
UPDATE users SET onboarded_at = COALESCE(onboarded_at, now()) WHERE id = @id;
