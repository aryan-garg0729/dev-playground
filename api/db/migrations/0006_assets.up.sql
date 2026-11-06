-- Lesson images (Phase 9). Content-addressed: one row per distinct file, shared
-- by every exercise that uses it. Upserted by cq build, never deleted (old
-- exercise versions and cached pages may still point at them).
CREATE TABLE content_assets (
    sha256 text PRIMARY KEY CHECK (sha256 ~ '^[0-9a-f]{64}$'),
    mime   text NOT NULL CHECK (mime IN ('image/png', 'image/jpeg', 'image/webp', 'image/gif', 'image/svg+xml')),
    size   int  NOT NULL CHECK (size > 0 AND size <= 1048576),
    bytes  bytea NOT NULL
);

-- Public name -> served file ("assets/loop.svg" -> "<sha256>.svg").
ALTER TABLE exercises ADD COLUMN assets jsonb NOT NULL DEFAULT '{}';
