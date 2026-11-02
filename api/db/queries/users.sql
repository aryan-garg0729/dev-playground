-- name: EnsureUserByEmail :one
INSERT INTO users (email, name)
VALUES (@email, @name)
ON CONFLICT (email) DO UPDATE SET email = EXCLUDED.email
RETURNING id, email, name, avatar, timezone;

-- name: GetUser :one
SELECT id, email, name, avatar, timezone FROM users WHERE id = @id;
