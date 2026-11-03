-- name: LockLeagueKey :exec
-- Transaction-scoped lock for group assignment and the weekly close.
SELECT pg_advisory_xact_lock(hashtext(@key::text));

-- name: GetLeagueMember :one
SELECT tier, group_no FROM league_members WHERE week = @week AND user_id = @user_id;

-- name: OpenLeagueGroup :one
-- The lowest-numbered group in the tier that still has room.
SELECT group_no FROM league_members
WHERE week = @week AND tier = @tier
GROUP BY group_no
HAVING count(*) < @size::int
ORDER BY group_no
LIMIT 1;

-- name: MaxLeagueGroup :one
SELECT COALESCE(max(group_no), 0)::int AS group_no FROM league_members WHERE week = @week AND tier = @tier;

-- name: InsertLeagueMember :exec
INSERT INTO league_members (week, user_id, tier, group_no)
VALUES (@week, @user_id, @tier, @group_no)
ON CONFLICT (week, user_id) DO NOTHING;

-- name: LeagueStandings :many
-- A group's members by weekly XP (ties: whoever reached it first, then id).
SELECT m.user_id, u.name, u.email, u.avatar,
       COALESCE(sum(x.amount), 0)::int AS weekly_xp,
       COALESCE(max(x.created_at), '9999-01-01'::timestamptz)::timestamptz AS last_xp_at
FROM league_members m
JOIN users u ON u.id = m.user_id
LEFT JOIN xp_events x ON x.user_id = m.user_id
    AND x.created_at >= m.week::timestamp AT TIME ZONE 'UTC'
    AND x.created_at < (m.week + 7)::timestamp AT TIME ZONE 'UTC'
WHERE m.week = @week AND m.tier = @tier AND m.group_no = @group_no
GROUP BY m.user_id, u.name, u.email, u.avatar
ORDER BY weekly_xp DESC, last_xp_at ASC, m.user_id;

-- name: UnclosedLeagueWeeks :many
SELECT DISTINCT m.week FROM league_members m
WHERE m.week < @before AND NOT EXISTS (SELECT 1 FROM league_weeks w WHERE w.week = m.week)
ORDER BY m.week;

-- name: IsLeagueWeekClosed :one
SELECT EXISTS (SELECT 1 FROM league_weeks WHERE week = @week)::bool AS closed;

-- name: ListLeagueGroups :many
SELECT DISTINCT tier, group_no FROM league_members WHERE week = @week ORDER BY tier, group_no;

-- name: SetLeagueResult :exec
UPDATE league_members SET final_rank = @final_rank::int, result = @result::text WHERE week = @week AND user_id = @user_id;

-- name: CloseLeagueWeek :exec
INSERT INTO league_weeks (week) VALUES (@week) ON CONFLICT DO NOTHING;
