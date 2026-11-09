CREATE TABLE magic_links (
    token_hash text PRIMARY KEY,
    email      text NOT NULL,
    next       text NOT NULL DEFAULT '/',
    expires_at timestamptz NOT NULL,
    used_at    timestamptz,
    created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX magic_links_expires_at_idx ON magic_links (expires_at);
