-- name: InsertByteEvent :one
-- Returns no row when (user, reason, ref) already exists: the reward was given.
INSERT INTO byte_events (user_id, amount, reason, ref_id)
VALUES (@user_id, @amount, @reason, @ref_id)
ON CONFLICT (user_id, reason, ref_id) DO NOTHING
RETURNING id;

-- name: AddBytes :one
-- Applies a ledger amount to the cached balance; no row if it would go negative.
UPDATE user_stats SET bytes = bytes + @amount::int
WHERE user_id = @user_id AND bytes + @amount::int >= 0
RETURNING bytes;

-- name: ByteLedgerSum :one
SELECT COALESCE(sum(amount), 0)::int AS total FROM byte_events WHERE user_id = @user_id;

-- name: BytesEarnedSince :one
SELECT COALESCE(sum(amount), 0)::int AS total FROM byte_events
WHERE user_id = @user_id AND amount > 0 AND created_at >= @since;

-- name: LockIdempotencyKey :exec
SELECT pg_advisory_xact_lock(hashtext(@user_id::text || ':' || @request_id::text));

-- name: GetIdempotent :one
SELECT action, response FROM idempotency WHERE user_id = @user_id AND request_id = @request_id;

-- name: PutIdempotent :exec
INSERT INTO idempotency (user_id, request_id, action, response) VALUES (@user_id, @request_id, @action, @response);

-- name: ListInventory :many
SELECT powerup, count FROM inventory WHERE user_id = @user_id AND count > 0 ORDER BY powerup;

-- name: InventoryCount :one
SELECT COALESCE((SELECT count FROM inventory WHERE user_id = @user_id AND powerup = @powerup), 0)::int AS count;

-- name: AddInventory :exec
INSERT INTO inventory (user_id, powerup, count) VALUES (@user_id, @powerup, @n::int)
ON CONFLICT (user_id, powerup) DO UPDATE SET count = inventory.count + EXCLUDED.count;

-- name: TakeInventory :one
-- Removes one; no row when none are held.
UPDATE inventory SET count = count - 1
WHERE user_id = @user_id AND powerup = @powerup AND count > 0
RETURNING count;

-- name: AddFreezes :one
-- Adds held Streak Freezes up to the cap; no row when it would exceed it.
UPDATE user_stats SET streak_freezes = streak_freezes + @n::int
WHERE user_id = @user_id AND streak_freezes + @n::int <= @cap::int
RETURNING streak_freezes;

-- name: RecordPowerupUse :exec
INSERT INTO powerup_uses (user_id, powerup, ref_id) VALUES (@user_id, @powerup, @ref_id);

-- name: CountPowerupUses :one
SELECT count(*)::int AS n FROM powerup_uses
WHERE user_id = @user_id AND powerup = @powerup AND created_at >= @since
  AND (@ref_id::text = '' OR ref_id = @ref_id::text);

-- name: SetEffect :exec
INSERT INTO active_effects (user_id, effect, ref_id, expires_at) VALUES (@user_id, @effect, @ref_id, @expires_at)
ON CONFLICT (user_id, effect, ref_id) DO UPDATE SET expires_at = EXCLUDED.expires_at;

-- name: EffectActive :one
SELECT EXISTS (SELECT 1 FROM active_effects
    WHERE user_id = @user_id AND effect = @effect AND ref_id = @ref_id AND expires_at > @now)::bool AS active;

-- name: ClearEffect :exec
DELETE FROM active_effects WHERE user_id = @user_id AND effect = @effect AND ref_id = @ref_id;

-- name: ListEffects :many
SELECT effect, ref_id, expires_at FROM active_effects WHERE user_id = @user_id AND expires_at > @now ORDER BY effect;

-- name: GetProgress :one
SELECT status, attempts, used_hint, peeked, first_try FROM progress WHERE user_id = @user_id AND exercise_id = @exercise_id;

-- name: SetPeeked :exec
UPDATE progress SET peeked = true WHERE user_id = @user_id AND exercise_id = @exercise_id;

-- name: SetFirstTry :exec
UPDATE progress SET first_try = true WHERE user_id = @user_id AND exercise_id = @exercise_id;

-- name: LatestPeekable :one
-- The user's latest submission of the exercise, with its server-only hidden inputs.
SELECT id, hidden_inputs, result FROM submissions
WHERE user_id = @user_id AND exercise_id = @exercise_id AND status <> 'pending'
ORDER BY created_at DESC LIMIT 1;

-- name: SetHiddenInputs :exec
UPDATE submissions SET hidden_inputs = @hidden_inputs WHERE id = @id;

-- name: ExerciseChapter :one
SELECT e.chapter_id, ch.course_id, e.boss, e.type FROM exercises e JOIN chapters ch ON ch.id = e.chapter_id WHERE e.id = @id;

-- name: ChapterProgress :many
-- Lessons of a chapter with the user's progress.
SELECT e.id, e.boss, COALESCE(p.status, '')::text AS status, COALESCE(p.first_try, false)::bool AS first_try
FROM exercises e
LEFT JOIN progress p ON p.exercise_id = e.id AND p.user_id = @user_id
WHERE e.chapter_id = @chapter_id AND e.published
ORDER BY e.order_index;

-- name: CourseChaptersDone :one
-- Chapters of the course, and how many the user has completed.
SELECT count(*)::int AS chapters,
       count(*) FILTER (WHERE NOT EXISTS (
           SELECT 1 FROM exercises e LEFT JOIN progress p ON p.exercise_id = e.id AND p.user_id = @user_id
           WHERE e.chapter_id = ch.id AND e.published AND COALESCE(p.status, '') <> 'passed'))::int AS done
FROM chapters ch WHERE ch.course_id = @course_id AND ch.published;

-- name: LockedChapters :many
-- Chapters behind an undefeated boss: every chapter after one whose boss the user has not
-- passed, unless a Skeleton Key waived that boss (chapter_keys.gate_chapter_id).
SELECT ch.id FROM chapters ch
WHERE ch.course_id = @course_id AND ch.published AND EXISTS (
    SELECT 1 FROM chapters prev
    JOIN exercises b ON b.chapter_id = prev.id AND b.boss AND b.published
    LEFT JOIN progress p ON p.exercise_id = b.id AND p.user_id = sqlc.narg('user_id')
    WHERE prev.course_id = ch.course_id AND prev.published AND prev.order_index < ch.order_index
      AND COALESCE(p.status, '') <> 'passed'
      AND NOT EXISTS (SELECT 1 FROM chapter_keys k WHERE k.user_id = sqlc.narg('user_id') AND k.gate_chapter_id = prev.id));

-- name: ChapterGate :one
-- The latest chapter before @chapter_id whose boss the user has not passed and no key waived:
-- the boss that locks it.
SELECT prev.id FROM chapters ch
JOIN chapters prev ON prev.course_id = ch.course_id AND prev.published AND prev.order_index < ch.order_index
JOIN exercises b ON b.chapter_id = prev.id AND b.boss AND b.published
LEFT JOIN progress p ON p.exercise_id = b.id AND p.user_id = @user_id
WHERE ch.id = @chapter_id AND COALESCE(p.status, '') <> 'passed'
  AND NOT EXISTS (SELECT 1 FROM chapter_keys k WHERE k.user_id = @user_id AND k.gate_chapter_id = prev.id)
ORDER BY prev.order_index DESC LIMIT 1;

-- name: GetChapter :one
SELECT id, course_id, order_index FROM chapters WHERE id = @id AND published;

-- name: CourseKey :one
SELECT chapter_id FROM chapter_keys WHERE user_id = @user_id AND course_id = @course_id;

-- name: InsertChapterKey :exec
INSERT INTO chapter_keys (user_id, course_id, gate_chapter_id, chapter_id)
VALUES (@user_id, @course_id, @gate_chapter_id, @chapter_id);

-- name: BumpCounter :one
INSERT INTO counters (user_id, key, value) VALUES (@user_id, @key, @n::int)
ON CONFLICT (user_id, key) DO UPDATE SET value = counters.value + EXCLUDED.value
RETURNING value;

-- name: MaxCounter :exec
INSERT INTO counters (user_id, key, value) VALUES (@user_id, @key, @n::int)
ON CONFLICT (user_id, key) DO UPDATE SET value = greatest(counters.value, EXCLUDED.value);

-- name: ListCounters :many
SELECT key, value FROM counters WHERE user_id = @user_id;

-- name: UpsertMedal :exec
INSERT INTO medals (id, family, name, secret, tiers, ord) VALUES (@id, @family, @name, @secret, @tiers, @ord)
ON CONFLICT (id) DO UPDATE SET family = EXCLUDED.family, name = EXCLUDED.name, secret = EXCLUDED.secret,
    tiers = EXCLUDED.tiers, ord = EXCLUDED.ord;

-- name: AwardMedal :one
INSERT INTO user_medals (user_id, medal_id, tier) VALUES (@user_id, @medal_id, @tier)
ON CONFLICT DO NOTHING
RETURNING earned_at;

-- name: ListUserMedals :many
SELECT medal_id, tier, earned_at FROM user_medals WHERE user_id = @user_id ORDER BY earned_at;

-- name: UpsertQuest :exec
INSERT INTO quests (id, cadence, rule, reward) VALUES (@id, @cadence, @rule, @reward)
ON CONFLICT (id) DO UPDATE SET cadence = EXCLUDED.cadence, rule = EXCLUDED.rule, reward = EXCLUDED.reward;

-- name: BumpQuest :one
-- Adds to a quest's progress; done_at is set once, when target is first reached.
INSERT INTO user_quests (user_id, quest_id, period, progress, done_at)
VALUES (@user_id, @quest_id, @period, @n::int, CASE WHEN @n::int >= @target::int THEN now() END)
ON CONFLICT (user_id, quest_id, period) DO UPDATE SET
    progress = user_quests.progress + EXCLUDED.progress,
    done_at = COALESCE(user_quests.done_at,
        CASE WHEN user_quests.progress + EXCLUDED.progress >= @target::int THEN now() END)
RETURNING progress, done_at;

-- name: ListUserQuests :many
SELECT quest_id, period, progress, done_at FROM user_quests
WHERE user_id = @user_id AND period IN (@day::date, @week::date);

-- name: OpenChest :one
-- Records a chest once; no row when it was already opened.
INSERT INTO chests (user_id, kind, ref_id, bytes, powerup) VALUES (@user_id, @kind, @ref_id, @bytes, @powerup)
ON CONFLICT DO NOTHING
RETURNING opened_at;

-- name: GetChest :one
SELECT bytes, powerup, opened_at FROM chests WHERE user_id = @user_id AND kind = @kind AND ref_id = @ref_id;

-- name: OpenBet :one
SELECT id, stake, goal, starts_on FROM bets WHERE user_id = @user_id AND status = 'open';

-- name: CreateBet :one
INSERT INTO bets (user_id, stake, goal, starts_on) VALUES (@user_id, @stake, @goal, @starts_on) RETURNING id;

-- name: SettleBet :exec
UPDATE bets SET status = @status, settled_at = now() WHERE id = @id AND status = 'open';

-- name: SetStreakState :exec
UPDATE user_stats SET streak_days = @streak_days, streak_freezes = @streak_freezes, last_active_date = @last_active_date,
    lost_streak_days = @lost_streak_days, lost_on = @lost_on
WHERE user_id = @user_id;

-- name: GetEconomyStats :one
SELECT s.bytes, s.streak_days, s.streak_freezes, s.last_active_date, s.longest_streak, s.lost_streak_days, s.lost_on,
       s.total_xp, s.level, s.league_tier, s.title, s.daily_goal, s.rolled_over_on, u.timezone
FROM user_stats s JOIN users u ON u.id = s.user_id
WHERE s.user_id = @user_id
FOR UPDATE OF s;

-- name: SetRolledOver :exec
UPDATE user_stats SET rolled_over_on = @day WHERE user_id = @user_id;

-- name: SetLongestStreak :exec
UPDATE user_stats SET longest_streak = greatest(longest_streak, @days::int) WHERE user_id = @user_id;

-- name: SetTitle :exec
UPDATE user_stats SET title = @title WHERE user_id = @user_id;

-- name: SetLostStreak :exec
UPDATE user_stats SET lost_streak_days = @lost_streak_days, lost_on = @lost_on WHERE user_id = @user_id;

-- name: DueRollovers :many
-- Users active in the last 14 days (or with an open bet) whose local date moved past their last rollover.
SELECT s.user_id FROM user_stats s JOIN users u ON u.id = s.user_id
WHERE (s.last_active_date >= timezone('UTC', sqlc.arg(now)::timestamptz)::date - 14
       OR EXISTS (SELECT 1 FROM bets b WHERE b.user_id = s.user_id AND b.status = 'open'))
  AND (s.rolled_over_on IS NULL OR s.rolled_over_on < timezone(u.timezone, sqlc.arg(now)::timestamptz)::date)
LIMIT 500;

-- name: DayLessons :many
SELECT day, count(*)::int AS lessons FROM day_passes
WHERE user_id = @user_id AND day BETWEEN @from_day AND @to_day
GROUP BY day ORDER BY day;

-- name: IsLeagueShielded :one
SELECT EXISTS (SELECT 1 FROM active_effects WHERE user_id = @user_id AND effect = 'league_shield' AND ref_id = @week::text)::bool AS shielded;

-- name: NextChapter :one
SELECT n.id FROM chapters c
JOIN chapters n ON n.course_id = c.course_id AND n.published AND n.order_index > c.order_index
WHERE c.id = @id
ORDER BY n.order_index
LIMIT 1;

-- name: CountPowerupUsesInChapter :one
SELECT count(*)::int AS n FROM powerup_uses u JOIN exercises e ON e.id = u.ref_id
WHERE u.user_id = @user_id AND u.powerup = @powerup AND e.chapter_id = @chapter_id;

-- name: PeekedChecks :many
-- Check indexes of a submission already revealed by Peek.
SELECT split_part(ref_id, ':', 2)::text AS check_index FROM powerup_uses
WHERE user_id = @user_id AND powerup = 'peek_shown' AND split_part(ref_id, ':', 1) = @submission_id::text;

-- name: RecordPeek :exec
INSERT INTO powerup_uses (user_id, powerup, ref_id) VALUES (@user_id, 'peek_shown', @submission_id::text || ':' || @check_index::text);

-- name: GetCounter :one
SELECT COALESCE((SELECT value FROM counters WHERE user_id = @user_id AND key = @key), 0)::int AS value;
