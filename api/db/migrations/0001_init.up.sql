CREATE TABLE users (
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    email      text NOT NULL UNIQUE,
    name       text NOT NULL DEFAULT '',
    avatar     text NOT NULL DEFAULT '',
    country    text NOT NULL DEFAULT '',
    timezone   text NOT NULL DEFAULT 'UTC',
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE sessions (
    id         text PRIMARY KEY,
    user_id    uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    expires_at timestamptz NOT NULL
);
CREATE INDEX sessions_user_id_idx ON sessions (user_id);

CREATE TABLE courses (
    id          text PRIMARY KEY,
    title       text NOT NULL,
    description text NOT NULL DEFAULT '',
    order_index int  NOT NULL DEFAULT 0,
    published   boolean NOT NULL DEFAULT true
);

CREATE TABLE chapters (
    id          text PRIMARY KEY,
    course_id   text NOT NULL REFERENCES courses (id),
    title       text NOT NULL,
    order_index int  NOT NULL,
    published   boolean NOT NULL DEFAULT true
);
CREATE INDEX chapters_course_id_idx ON chapters (course_id, order_index);

CREATE TABLE exercises (
    id           text PRIMARY KEY, -- manifest id, stable forever
    chapter_id   text NOT NULL REFERENCES chapters (id),
    type         text NOT NULL,
    version      int  NOT NULL,
    title        text NOT NULL,
    xp           int  NOT NULL,
    difficulty   int  NOT NULL,
    tags         text[] NOT NULL DEFAULT '{}',
    free         boolean NOT NULL DEFAULT false,
    boss         boolean NOT NULL DEFAULT false,
    order_index  int  NOT NULL,
    public_view  jsonb NOT NULL,
    private_spec jsonb NOT NULL, -- server-only; never selected by browser-facing queries
    published    boolean NOT NULL DEFAULT true
);
CREATE INDEX exercises_chapter_id_idx ON exercises (chapter_id, order_index);

CREATE TABLE submissions (
    id               uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id          uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    exercise_id      text NOT NULL REFERENCES exercises (id),
    exercise_version int  NOT NULL,
    answer           jsonb NOT NULL,
    status           text NOT NULL CHECK (status IN ('pass', 'fail', 'error', 'timeout', 'pending')),
    score            double precision NOT NULL DEFAULT 0,
    result           jsonb,
    created_at       timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX submissions_user_exercise_idx ON submissions (user_id, exercise_id, created_at DESC);

CREATE TABLE progress (
    user_id         uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    exercise_id     text NOT NULL REFERENCES exercises (id),
    status          text NOT NULL CHECK (status IN ('started', 'passed')),
    first_passed_at timestamptz,
    attempts        int NOT NULL DEFAULT 0,
    used_hint       boolean NOT NULL DEFAULT false,
    PRIMARY KEY (user_id, exercise_id)
);
