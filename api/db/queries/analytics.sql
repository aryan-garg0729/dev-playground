-- name: InsertAnalyticsEvent :exec
INSERT INTO analytics_events (user_id, anon_id, name, props, path, created_at)
VALUES (@user_id, @anon_id, @name, @props, @path, @created_at);

-- name: UpsertClientError :exec
INSERT INTO client_errors (fingerprint, day, message, stack, url, user_agent, release, users)
VALUES (@fingerprint, @day, @message, @stack, @url, @user_agent, @release,
        CASE WHEN sqlc.narg(user_id)::uuid IS NULL THEN '{}'::uuid[] ELSE ARRAY[sqlc.narg(user_id)::uuid] END)
ON CONFLICT (fingerprint, day) DO UPDATE SET
    count = client_errors.count + 1,
    last_seen = now(),
    users = CASE
        WHEN sqlc.narg(user_id)::uuid IS NULL OR sqlc.narg(user_id)::uuid = ANY (client_errors.users)
             OR cardinality(client_errors.users) >= 100 THEN client_errors.users
        ELSE client_errors.users || sqlc.narg(user_id)::uuid
    END;

-- name: DeleteOldTelemetry :exec
-- Analytics events and client errors are kept 90 days.
WITH e AS (DELETE FROM analytics_events WHERE created_at < now() - interval '90 days')
DELETE FROM client_errors WHERE day < (now() - interval '90 days')::date;
