import type { Schemas } from '../api/client';

export type Activity = Schemas['Activity'];

const DAY = 86400000;
const months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
const rowLabels = ['Mon', '', 'Wed', '', 'Fri', '', ''];

/** Dates are local calendar days ("2026-11-04"); do the math in UTC so no timezone shifts them. */
const parse = (iso: string) => Date.parse(`${iso}T00:00:00Z`);
const format = (t: number) => new Date(t).toISOString().slice(0, 10);

// Fill per lesson count: none, 1, 2, 3, 4+ (tokens only; opacity steps the amber).
const levels = ['bg-surface-raised', 'bg-accent-line', 'bg-accent opacity-60', 'bg-accent opacity-80', 'bg-accent'];
export const level = (lessons: number) => levels[Math.min(lessons, 4)];

export interface Cell {
  date: string;
  lessons: number;
}

/** Columns of Monday..Sunday ending with the week of `activity.to`; days after today are null. */
export function heatmapWeeks(activity: Activity, weeks: number): (Cell | null)[][] {
  const counts = new Map(activity.days.map((d) => [d.date, d.lessons]));
  const today = parse(activity.to);
  const weekday = (new Date(today).getUTCDay() + 6) % 7; // Monday = 0
  const start = today - (weekday + (weeks - 1) * 7) * DAY;
  return Array.from({ length: weeks }, (_, w) =>
    Array.from({ length: 7 }, (_, d) => {
      const t = start + (w * 7 + d) * DAY;
      if (t > today) return null;
      const date = format(t);
      return { date, lessons: counts.get(date) ?? 0 };
    }),
  );
}

/** Lessons and active days in the shown columns. */
export function totals(cols: (Cell | null)[][]) {
  let lessons = 0;
  let days = 0;
  for (const c of cols.flat()) {
    if (c && c.lessons > 0) {
      lessons += c.lessons;
      days++;
    }
  }
  return { lessons, days };
}

export function Legend() {
  return (
    <span className="text-caption flex items-center gap-1 text-text-muted" aria-hidden="true">
      Less
      {levels.map((l) => (
        <span key={l} className={`h-3 w-3 rounded-[3px] ${l}`} />
      ))}
      More
    </span>
  );
}

/**
 * GitHub-style activity grid: one column per week, one square per day, colored
 * by lessons passed. Today is outlined. Scrolls sideways on narrow screens.
 */
export function Heatmap({ activity, weeks, label, cell = 14 }: { activity: Activity; weeks: number; label: string; cell?: number }) {
  const cols = heatmapWeeks(activity, weeks);
  const { lessons, days } = totals(cols);
  // A month is labelled at its first week; the first column only when the next one doesn't start a month (no overlap).
  const month = (i: number) => cols[i]?.[0]?.date.slice(5, 7) ?? '';
  const monthAt = cols.map((_, i) => {
    const starts = month(i) !== month(i - 1);
    if (!starts || (i === 0 && month(1) !== month(0))) return '';
    return months[Number(month(i)) - 1] ?? '';
  });
  return (
    <div className="overflow-x-auto">
      <div
        role="img"
        aria-label={`${label}: ${lessons} ${lessons === 1 ? 'lesson' : 'lessons'} on ${days} active ${days === 1 ? 'day' : 'days'}`}
        className="grid w-max gap-[3px]"
        style={{ gridTemplateColumns: `28px repeat(${weeks}, ${cell}px)`, gridAutoRows: `${cell}px` }}
      >
        <span />
        {monthAt.map((m, i) => (
          <span key={i} className="text-caption overflow-visible whitespace-nowrap leading-none text-text-muted">
            {m}
          </span>
        ))}
        {rowLabels.map((r, d) => (
          <Row key={d} label={r} cells={cols.map((c) => c[d] ?? null)} today={activity.to} />
        ))}
      </div>
    </div>
  );
}

function Row({ label, cells, today }: { label: string; cells: (Cell | null)[]; today: string }) {
  return (
    <>
      <span className="text-caption leading-none text-text-muted">{label}</span>
      {cells.map((c, i) =>
        c ? (
          <span
            key={i}
            data-date={c.date}
            title={`${c.date}: ${c.lessons} ${c.lessons === 1 ? 'lesson' : 'lessons'}`}
            className={`rounded-[3px] ${level(c.lessons)} ${c.date === today ? 'outline outline-2 -outline-offset-2 outline-text-secondary' : ''}`}
          />
        ) : (
          <span key={i} />
        ),
      )}
    </>
  );
}
