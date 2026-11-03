/** "2d 6h", "5h 12m", "8m": time left until `iso`, from `now`. */
export function timeLeft(iso: string, now = Date.now()): string {
  const mins = Math.max(0, Math.floor((new Date(iso).getTime() - now) / 60000));
  const d = Math.floor(mins / 1440);
  const h = Math.floor((mins % 1440) / 60);
  const m = mins % 60;
  if (d > 0) return `${d}d ${h}h`;
  if (h > 0) return `${h}h ${m}m`;
  return `${m}m`;
}

export function greeting(hour: number): string {
  if (hour < 5) return 'Good evening';
  if (hour < 12) return 'Good morning';
  if (hour < 18) return 'Good afternoon';
  return 'Good evening';
}

export const tierTitle = (t: string) => t.charAt(0).toUpperCase() + t.slice(1);

const tiers = ['bronze', 'silver', 'gold', 'diamond'];
export const tierAbove = (t: string) => tiers[Math.min(tiers.indexOf(t) + 1, tiers.length - 1)] ?? t;
export const tierBelow = (t: string) => tiers[Math.max(tiers.indexOf(t) - 1, 0)] ?? t;

export const plural = (n: number, one: string, many = `${one}s`) => `${n} ${n === 1 ? one : many}`;

/** The line under the greeting: what keeps the streak and the daily goal. */
export function nudge(t: { goal: number; done: number; streak_days: number; week: { today: boolean; active: boolean }[] }): string {
  const activeToday = t.week.some((d) => d.today && d.active);
  const left = Math.max(0, t.goal - t.done);
  if (!activeToday && t.streak_days > 0) return `One lesson keeps your ${t.streak_days}-day streak.`;
  if (!activeToday) return 'Pass a lesson today to start a streak.';
  if (left > 0) return `${plural(left, 'more lesson')} to reach your daily goal.`;
  return 'Daily goal done. Anything more is a bonus.';
}
