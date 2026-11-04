-- Phase 8: the learner CLI (`codequest`). The web app shows a one-time code to
-- a signed-in learner; the CLI trades it for a bearer token. Only SHA-256
-- hashes of codes and tokens are stored.

-- Login codes: 10 minutes, single use.
CREATE TABLE cli_login_codes (
    code_hash  text PRIMARY KEY,
    user_id    uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    expires_at timestamptz NOT NULL,
    used_at    timestamptz,
    created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX cli_login_codes_user_idx ON cli_login_codes (user_id);

-- CLI tokens: 90 days, extended on use; revoked by `codequest logout`.
CREATE TABLE cli_tokens (
    id           text PRIMARY KEY,
    user_id      uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    label        text NOT NULL DEFAULT '',
    created_at   timestamptz NOT NULL DEFAULT now(),
    last_used_at timestamptz NOT NULL DEFAULT now(),
    expires_at   timestamptz NOT NULL
);
CREATE INDEX cli_tokens_user_idx ON cli_tokens (user_id);
