-- Phase 4: accounts, XP ledger, streaks and leagues.

-- Sessions store only a SHA-256 of the cookie value (id = hex hash).
ALTER TABLE sessions ADD COLUMN created_at timestamptz NOT NULL DEFAULT now();
CREATE INDEX sessions_expires_at_idx ON sessions (expires_at);

-- Magic links: single use, short-lived, stored hashed.
CREATE TABLE magic_links (
    token_hash text PRIMARY KEY,
    email      text NOT NULL,
    next       text NOT NULL DEFAULT '/',
    expires_at timestamptz NOT NULL,
    used_at    timestamptz,
    created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX magic_links_expires_at_idx ON magic_links (expires_at);

CREATE TABLE oauth_identities (
    provider         text NOT NULL,
    provider_user_id text NOT NULL,
    user_id          uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at       timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (provider, provider_user_id)
);
CREATE INDEX oauth_identities_user_id_idx ON oauth_identities (user_id);

-- XP ledger, append-only. The unique key makes "first pass awards XP once" a
-- database guarantee, even for concurrent submissions.
CREATE TABLE xp_events (
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    amount     int  NOT NULL CHECK (amount >= 0),
    reason     text NOT NULL,
    ref_id     text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (user_id, reason, ref_id)
);
CREATE INDEX xp_events_user_created_idx ON xp_events (user_id, created_at);

-- Cache of the ledgers plus streak state; one row per user, locked per submission.
CREATE TABLE user_stats (
    user_id          uuid PRIMARY KEY REFERENCES users (id) ON DELETE CASCADE,
    total_xp         int  NOT NULL DEFAULT 0,
    level            int  NOT NULL DEFAULT 1,
    bytes            int  NOT NULL DEFAULT 0,
    streak_days      int  NOT NULL DEFAULT 0,
    streak_freezes   int  NOT NULL DEFAULT 0,
    last_active_date date,
    league_tier      text NOT NULL DEFAULT 'bronze' CHECK (league_tier IN ('bronze', 'silver', 'gold', 'diamond')),
    title            text NOT NULL DEFAULT '',
    daily_goal       int  NOT NULL DEFAULT 2 CHECK (daily_goal BETWEEN 1 AND 10)
);

-- Exercises passed per local day (the user's timezone): streak week and daily goal.
CREATE TABLE day_passes (
    user_id     uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    day         date NOT NULL,
    exercise_id text NOT NULL REFERENCES exercises (id),
    PRIMARY KEY (user_id, day, exercise_id)
);

-- Weekly leagues. week = the Monday (UTC) the week starts. Users join the
-- week's group when they first earn XP in it.
CREATE TABLE league_members (
    week     date NOT NULL,
    user_id  uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    tier     text NOT NULL CHECK (tier IN ('bronze', 'silver', 'gold', 'diamond')),
    group_no int  NOT NULL,
    final_rank int,
    result   text CHECK (result IN ('promoted', 'demoted', 'stayed')),
    PRIMARY KEY (week, user_id)
);
CREATE INDEX league_members_group_idx ON league_members (week, tier, group_no);

CREATE TABLE league_weeks (
    week      date PRIMARY KEY,
    closed_at timestamptz NOT NULL DEFAULT now()
);

-- What the game engine did for a submission (XP, level, streak), for the result screen.
ALTER TABLE submissions ADD COLUMN outcome jsonb;
