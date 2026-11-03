-- name: CreateSession :exec
INSERT INTO sessions (id, user_id, expires_at) VALUES (@id, @user_id, @expires_at);

-- name: GetSessionUser :one
SELECT user_id FROM sessions WHERE id = @id AND expires_at > now();

-- name: DeleteSession :exec
DELETE FROM sessions WHERE id = @id;

-- name: DeleteExpiredAuth :exec
WITH s AS (DELETE FROM sessions WHERE expires_at < now())
DELETE FROM magic_links WHERE expires_at < now() - interval '1 day';

-- name: CreateMagicLink :exec
INSERT INTO magic_links (token_hash, email, next, expires_at)
VALUES (@token_hash, @email, @next, @expires_at);

-- name: UseMagicLink :one
-- Single use: only the first call for an unexpired token returns a row.
UPDATE magic_links SET used_at = now()
WHERE token_hash = @token_hash AND used_at IS NULL AND expires_at > now()
RETURNING email, next;

-- name: GetOAuthUser :one
SELECT user_id FROM oauth_identities WHERE provider = @provider AND provider_user_id = @provider_user_id;

-- name: LinkOAuth :exec
INSERT INTO oauth_identities (provider, provider_user_id, user_id)
VALUES (@provider, @provider_user_id, @user_id)
ON CONFLICT (provider, provider_user_id) DO NOTHING;
