// Made-up grader results and game outcomes for the dev effects lab. Never sent
// anywhere: they only drive the real components so every moment can be seen.
import type { Outcome, Result } from '../../api/client';

/** A result with `total` checks of which exactly `hits` pass, the blocks spread through it. */
export function attack(total: number, hits: number): Result {
  const blocks = total - hits;
  const blockAt = new Set(Array.from({ length: blocks }, (_, i) => Math.floor(((i + 0.5) * total) / blocks)));
  return {
    status: blocks === 0 ? 'pass' : 'fail',
    score: total ? hits / total : 0,
    stdout: '',
    runtime_ms: 42,
    checks: Array.from({ length: total }, (_, i) =>
      blockAt.has(i)
        ? { name: `check_${i + 1}()`, passed: false, visible: i < 2, got: '"wrong"', expected: '"right"' }
        : { name: `check_${i + 1}()`, passed: true, visible: i < 2 },
    ),
  };
}

const base: Outcome = {
  xp_awarded: 35,
  breakdown: { base: 20, first_try: 10, no_hint: 5, boss: 0, boost: 0 },
  first_try: true,
  total_xp: 1727,
  level: 7,
  level_up: false,
  level_start_xp: 1470,
  next_level_xp: 1852,
  streak_days: 13,
  streak_extended: true,
  freezes_used: 0,
  freeze_earned: false,
  second_wind_kept: false,
  focus_boost: false,
  rewards: null,
};

const bossRewards = (chestPowerup: string): Outcome['rewards'] => ({
  bytes: [{ reason: 'chapter', label: 'Chapter complete', amount: 150 }],
  medals: [{ id: 'boss_slayer', name: 'Boss Slayer', tier: 'bronze', title: 'Boss Slayer I' }],
  quests_done: ['Pass a coding lesson'],
  chest_ready: false,
  chapter_complete: true,
  boss: { max_hp: 5, hp: 0, defeated: true, chest: { bytes: 300, powerup: chestPowerup }, unlocked_chapter: 'next' },
});

export type CelebrationKind = 'clear' | 'levelUp' | 'boss' | 'bossLevelUp';

export const celebrations: Record<CelebrationKind, { label: string; outcome: Outcome }> = {
  clear: { label: 'Lesson clear', outcome: base },
  levelUp: {
    label: 'Lesson clear + level up',
    outcome: { ...base, xp_awarded: 140, total_xp: 1870, level: 8, level_up: true, level_start_xp: 1852, next_level_xp: 2263 },
  },
  boss: {
    label: 'Boss defeated + chest',
    outcome: {
      ...base,
      xp_awarded: 188,
      breakdown: { base: 50, first_try: 25, no_hint: 13, boss: 100, boost: 0 },
      rewards: bossRewards('focus_boost'),
    },
  },
  bossLevelUp: {
    label: 'Boss defeated + level up + chest',
    outcome: {
      ...base,
      xp_awarded: 188,
      breakdown: { base: 50, first_try: 25, no_hint: 13, boss: 100, boost: 0 },
      total_xp: 1900,
      level: 8,
      level_up: true,
      level_start_xp: 1852,
      next_level_xp: 2263,
      rewards: bossRewards('peek'),
    },
  },
};

export const POWERUPS = ['focus_boost', 'peek', 'second_wind', 'league_shield'];
