-- Phase 6: AI tutor (development only for now), onboarding, analytics and
-- error monitoring. Nothing here can change a grading Result.

-- Tutor conversation per user and exercise. Token counts are logged per reply.
CREATE TABLE tutor_messages (
    id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id       uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    exercise_id   text NOT NULL REFERENCES exercises (id),
    role          text NOT NULL CHECK (role IN ('user', 'assistant')),
    content       text NOT NULL,
    input_tokens  int  NOT NULL DEFAULT 0,
    output_tokens int  NOT NULL DEFAULT 0,
    model         text NOT NULL DEFAULT '',
    created_at    timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX tutor_messages_user_exercise_idx ON tutor_messages (user_id, exercise_id, created_at);

-- Onboarding: the learner's goal and when they finished onboarding.
ALTER TABLE users
    ADD COLUMN goal text NOT NULL DEFAULT '' CHECK (goal IN ('', 'zero', 'backend', 'interviews', 'ai')),
    ADD COLUMN onboarded_at timestamptz;
-- Accounts from before onboarding existed are not sent through it.
UPDATE users SET onboarded_at = created_at;

-- Product analytics, self-hosted. user_id is null for signed-out visitors;
-- anon_id is a random id the browser keeps. Kept 90 days.
CREATE TABLE analytics_events (
    id         bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id    uuid REFERENCES users (id) ON DELETE SET NULL,
    anon_id    text NOT NULL DEFAULT '',
    name       text NOT NULL,
    props      jsonb NOT NULL DEFAULT '{}',
    path       text NOT NULL DEFAULT '',
    created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX analytics_events_name_created_idx ON analytics_events (name, created_at);
CREATE INDEX analytics_events_created_idx ON analytics_events (created_at);

-- Browser errors, grouped by fingerprint per day. Kept 90 days.
CREATE TABLE client_errors (
    fingerprint text NOT NULL,
    day         date NOT NULL,
    message     text NOT NULL,
    stack       text NOT NULL DEFAULT '',
    url         text NOT NULL DEFAULT '',
    user_agent  text NOT NULL DEFAULT '',
    release     text NOT NULL DEFAULT '',
    count       int  NOT NULL DEFAULT 1,
    users       uuid[] NOT NULL DEFAULT '{}',
    first_seen  timestamptz NOT NULL DEFAULT now(),
    last_seen   timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (fingerprint, day)
);
