import { Link } from 'react-router';
import type { CourseSummary } from '../../api/client';

const levelPill = {
  beginner: 'bg-success-soft text-success',
  intermediate: 'bg-accent-soft text-accent-text',
  advanced: 'bg-danger-soft text-danger',
} as const;

const runsText = { browser: 'Runs in browser', cli: '>_ CodeQuest CLI' } as const;

const chip = 'text-caption inline-flex items-center rounded-pill px-[9px] py-[3px]';

/** Status pill and call to action, from progress and whether any lesson is free. */
export function courseState(c: CourseSummary): { pill: string; pillClass: string; action: string } {
  if (c.status === 'done') return { pill: 'Completed', pillClass: 'bg-success-soft text-success', action: 'Review' };
  if (c.status === 'in_progress') return { pill: 'In progress', pillClass: 'bg-success-soft text-success', action: 'Continue' };
  if (c.free_count > 0) return { pill: 'Free to start', pillClass: 'bg-surface-raised text-text-secondary', action: 'Start free' };
  return { pill: 'Pro', pillClass: 'bg-accent-soft text-accent-text', action: 'Preview' };
}

const plural = (n: number, one: string) => `${n} ${n === 1 ? one : `${one}s`}`;

export function CourseCard({ c }: { c: CourseSummary }) {
  const s = courseState(c);
  const started = c.status !== 'not_started';
  const pct = c.exercise_count > 0 ? (c.passed_count / c.exercise_count) * 100 : 0;
  return (
    <li>
      <Link
        to={`/courses/${c.id}`}
        aria-label={`${c.title}, ${s.pill}`}
        className="flex h-full flex-col gap-3 rounded-lg border border-border bg-surface p-5 hover:border-border-strong"
      >
        <span className="flex items-start justify-between gap-3">
          <span
            aria-hidden="true"
            className={`text-game flex h-11 w-11 items-center justify-center rounded-md border border-border bg-surface-raised text-[14px] ${c.status === 'in_progress' ? 'text-accent-text' : 'text-text'}`}
          >
            {c.icon}
          </span>
          <span className={`${chip} ${s.pillClass}`}>{s.pill}</span>
        </span>
        <span className="text-section text-text">{c.title}</span>
        <span className="text-small text-text-secondary">{c.description}</span>
        <span className="flex flex-wrap gap-2">
          <span className={`${chip} capitalize ${levelPill[c.level]}`}>{c.level}</span>
          {c.runs.map((r) => (
            <span key={r} className={`${chip} bg-surface-raised text-text-secondary`}>
              {runsText[r]}
            </span>
          ))}
        </span>
        {started && (
          <span className="flex flex-col gap-2">
            <span role="img" aria-label={`${Math.round(pct)}% complete`} className="block h-1 overflow-hidden rounded-pill bg-border">
              <span className="block h-full rounded-pill bg-success" style={{ width: `${pct}%` }} />
            </span>
            <span className="text-caption text-text-muted">
              {c.passed_count} of {plural(c.exercise_count, 'lesson')}
              {c.current_chapter && ` · Chapter ${c.current_chapter.index} · ${c.current_chapter.title}`}
            </span>
          </span>
        )}
        <span className="text-small mt-auto flex items-center justify-between gap-3 border-t border-border pt-3">
          <span className="truncate text-text-muted">
            {started || !c.topics ? `${plural(c.chapter_count, 'chapter')} · ${plural(c.exercise_count, 'lesson')}` : c.topics}
          </span>
          <span className="shrink-0 font-semibold text-accent-text">{s.action}</span>
        </span>
      </Link>
    </li>
  );
}
