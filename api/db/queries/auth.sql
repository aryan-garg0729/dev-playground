-- name: CreateSession :exec
INSERT INTO sessions (id, user_id, expires_at) VALUES (@id, @user_id, @expires_at);

-- name: GetSessionUser :one
SELECT user_id FROM sessions WHERE id = @id AND expires_at > now();

-- name: DeleteSession :exec
DELETE FROM sessions WHERE id = @id;

-- name: DeleteExpiredAuth :exec
DELETE FROM sessions WHERE expires_at < now();

-- name: GetOAuthUser :one
SELECT user_id FROM oauth_identities WHERE provider = @provider AND provider_user_id = @provider_user_id;

-- name: LinkOAuth :exec
INSERT INTO oauth_identities (provider, provider_user_id, user_id)
VALUES (@provider, @provider_user_id, @user_id)
ON CONFLICT (provider, provider_user_id) DO NOTHING;
