ALTER TABLE submissions DROP COLUMN hidden_inputs;
DROP TABLE counters;
ALTER TABLE progress DROP COLUMN first_try, DROP COLUMN peeked;
ALTER TABLE user_stats
    DROP CONSTRAINT user_stats_freezes_max,
    DROP CONSTRAINT user_stats_bytes_nonneg,
    DROP COLUMN rolled_over_on,
    DROP COLUMN lost_on,
    DROP COLUMN lost_streak_days,
    DROP COLUMN longest_streak;
DROP TABLE chests;
DROP TABLE user_quests;
DROP TABLE quests;
DROP TABLE user_medals;
DROP TABLE medals;
DROP TABLE bets;
DROP TABLE active_effects;
DROP TABLE powerup_uses;
DROP TABLE inventory;
DROP TABLE idempotency;
DROP TABLE byte_events;
