-- Skeleton Key uses: one per user and course. gate_chapter_id is the chapter
-- whose undefeated boss no longer locks the chapters after it; that boss stays
-- unpassed and pays nothing. chapter_id is the chapter the key opened.
CREATE TABLE chapter_keys (
    user_id         uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    course_id       text NOT NULL,
    gate_chapter_id text NOT NULL,
    chapter_id      text NOT NULL,
    used_at         timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, course_id)
);
