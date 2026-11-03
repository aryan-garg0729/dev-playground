-- Phase 7: game economy. Nothing here can change a grading Result.

-- Byte ledger, append-only. Balance cached in user_stats.bytes, always the sum.
-- (user_id, reason, ref_id) is unique, so each reward (or each request) lands once.
CREATE TABLE byte_events (
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    amount     int  NOT NULL,
    reason     text NOT NULL,
    ref_id     text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (user_id, reason, ref_id)
);
CREATE INDEX byte_events_user_created_idx ON byte_events (user_id, created_at);

-- Idempotency: the first response to a (user, request id) is stored and replayed.
CREATE TABLE idempotency (
    user_id    uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    request_id text NOT NULL,
    action     text NOT NULL,
    response   json NOT NULL, -- json, not jsonb: a replay returns the exact first response
    created_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, request_id)
);

CREATE TABLE inventory (
    user_id uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    powerup text NOT NULL,
    count   int  NOT NULL CHECK (count >= 0),
    PRIMARY KEY (user_id, powerup)
);

CREATE TABLE powerup_uses (
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    powerup    text NOT NULL,
    ref_id     text NOT NULL DEFAULT '',
    created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX powerup_uses_user_idx ON powerup_uses (user_id, powerup, created_at);

-- Effects in force: focus_boost (ref ''), second_wind (ref exercise), league_shield (ref week).
CREATE TABLE active_effects (
    user_id    uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    effect     text NOT NULL,
    ref_id     text NOT NULL DEFAULT '',
    expires_at timestamptz NOT NULL,
    PRIMARY KEY (user_id, effect, ref_id)
);

-- Double or Nothing: one open bet per user.
CREATE TABLE bets (
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    stake      int  NOT NULL CHECK (stake > 0),
    goal       int  NOT NULL,
    starts_on  date NOT NULL,
    status     text NOT NULL DEFAULT 'open' CHECK (status IN ('open', 'won', 'lost')),
    settled_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX bets_one_open_idx ON bets (user_id) WHERE status = 'open';

-- Medal and quest catalogues are synced from code (economy/catalog.go) at startup.
CREATE TABLE medals (
    id     text PRIMARY KEY,
    family text NOT NULL,
    name   text NOT NULL,
    secret boolean NOT NULL DEFAULT false,
    tiers  jsonb NOT NULL,
    ord    int  NOT NULL
);

CREATE TABLE user_medals (
    user_id   uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    medal_id  text NOT NULL REFERENCES medals (id),
    tier      text NOT NULL,
    earned_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, medal_id, tier)
);

CREATE TABLE quests (
    id      text PRIMARY KEY,
    cadence text NOT NULL CHECK (cadence IN ('daily', 'weekly')),
    rule    jsonb NOT NULL,
    reward  jsonb NOT NULL
);

-- period: the local date (daily) or the Monday of the week (weekly).
CREATE TABLE user_quests (
    user_id  uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    quest_id text NOT NULL REFERENCES quests (id),
    period   date NOT NULL,
    progress int  NOT NULL DEFAULT 0,
    done_at  timestamptz,
    PRIMARY KEY (user_id, quest_id, period)
);

-- Opened chests: kind daily (ref = local date) or boss (ref = exercise id). Once each.
CREATE TABLE chests (
    user_id   uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    kind      text NOT NULL CHECK (kind IN ('daily', 'boss')),
    ref_id    text NOT NULL,
    bytes     int  NOT NULL,
    powerup   text NOT NULL DEFAULT '',
    opened_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, kind, ref_id)
);

ALTER TABLE user_stats
    ADD COLUMN longest_streak   int  NOT NULL DEFAULT 0,
    ADD COLUMN lost_streak_days int  NOT NULL DEFAULT 0,
    ADD COLUMN lost_on          date,
    ADD COLUMN rolled_over_on   date,
    ADD CONSTRAINT user_stats_bytes_nonneg CHECK (bytes >= 0),
    ADD CONSTRAINT user_stats_freezes_max CHECK (streak_freezes BETWEEN 0 AND 2);

-- Peek was used on this exercise: no first-try bonus (even with Second Wind).
-- first_try: the first pass earned the first-try bonus (Flawless Chapter medal).
ALTER TABLE progress
    ADD COLUMN peeked boolean NOT NULL DEFAULT false,
    ADD COLUMN first_try boolean NOT NULL DEFAULT false;

-- Event counters for medals (first_try, no_hint, night, early, bosses, ...).
CREATE TABLE counters (
    user_id uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    key     text NOT NULL,
    value   int  NOT NULL DEFAULT 0,
    PRIMARY KEY (user_id, key)
);

-- Server-only: inputs of failed hidden checks ({"<row index>": "<input>"}), for Peek.
-- Never selected by a browser-facing query except one Peek at a time.
ALTER TABLE submissions ADD COLUMN hidden_inputs jsonb;
