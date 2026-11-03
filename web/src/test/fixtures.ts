import type { Me, Outcome } from '../api/client';

export const me: Me = {
  user: { id: '00000000-0000-0000-0000-000000000001', email: 'ada@example.com', name: 'Ada Lovelace', avatar: '', timezone: 'UTC' },
  stats: {
    total_xp: 1765, level: 7, level_start_xp: 1470, next_level_xp: 1852, streak_days: 12, streak_freezes: 1,
    league_tier: 'gold', daily_goal: 3, today_lessons: 2, bytes: 1240, title: 'Boss Slayer I', longest_streak: 13,
  },
};

/** A first-try pass of a 20 XP lesson that keeps a 13-day streak. */
export const firstPass: Outcome = {
  xp_awarded: 35, breakdown: { base: 20, first_try: 10, no_hint: 5, boss: 0, boost: 0 }, first_try: true, total_xp: 1727, level: 7,
  level_up: false, level_start_xp: 1470, next_level_xp: 1852, streak_days: 13, streak_extended: true, freezes_used: 0,
  freeze_earned: false, second_wind_kept: false, focus_boost: false, rewards: null,
};

export const nothing: Outcome = {
  ...firstPass, xp_awarded: 0, breakdown: null, first_try: false, streak_extended: false,
};
