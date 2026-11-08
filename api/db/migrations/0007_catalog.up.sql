-- Course catalog metadata and career paths (courses.png, home.png, profile-overview.png).
-- Written by cq build from course.yaml and content/paths.yaml.
ALTER TABLE courses
    ADD COLUMN icon   text   NOT NULL DEFAULT '',
    ADD COLUMN level  text   NOT NULL DEFAULT 'beginner' CHECK (level IN ('beginner', 'intermediate', 'advanced')),
    ADD COLUMN tags   text[] NOT NULL DEFAULT '{}',
    ADD COLUMN topics text   NOT NULL DEFAULT '';

-- A career path: courses in order. steps is [{course, note}]. goal links the
-- path to an onboarding goal (users.goal), which makes it "your path".
CREATE TABLE career_paths (
    id          text PRIMARY KEY,
    title       text NOT NULL,
    description text NOT NULL DEFAULT '',
    goal        text NOT NULL DEFAULT '',
    order_index int  NOT NULL DEFAULT 0,
    steps       jsonb NOT NULL DEFAULT '[]',
    published   boolean NOT NULL DEFAULT true
);
