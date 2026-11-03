import { Link } from 'react-router';
import type { Home } from '../../api/client';

type Card = NonNullable<Home['continue']>;

const segment = { passed: 'bg-success', started: 'bg-border-strong', not_started: 'bg-border-strong' } as const;

/** Lesson progress in a chapter as thin segments; the current one is accent. */
export function Segments({ lessons, current }: { lessons: Card['lessons']; current: number }) {
  return (
    <span className="flex gap-1" aria-hidden="true">
      {lessons.map((s, i) => (
        <span key={i} className={`h-4 w-1 rounded-pill ${i + 1 === current ? 'bg-accent' : segment[s]}`} />
      ))}
    </span>
  );
}

const typeLabel: Record<string, string> = { mcq: 'quiz', code: 'exercise', sql: 'SQL exercise' };

export function ContinueCard({ card }: { card: Card }) {
  const done = card.lessons.filter((s) => s === 'passed').length;
  const last = card.last;
  return (
    <section aria-labelledby="continue-title" className="grid overflow-hidden rounded-lg border border-border bg-surface md:grid-cols-2">
      <div className="flex flex-col gap-4 p-6">
        <p className="text-caption uppercase tracking-[0.08em] text-text-muted">Continue · {card.course_title}</p>
        <div className="flex flex-col gap-1">
          <h2 id="continue-title" className="text-[22px] font-semibold leading-[30px] text-text">
            {card.exercise_title}
          </h2>
          <p className="text-small text-text-muted">
            Chapter {card.chapter_index} · {card.chapter_title} · lesson {card.position.index} of {card.position.total}
          </p>
        </div>
        <span role="img" aria-label={`${done} of ${card.lessons.length} lessons done in this chapter`}>
          <Segments lessons={card.lessons} current={card.position.index} />
        </span>
        <div className="mt-auto flex flex-wrap gap-2">
          <Link
            to={`/learn/${card.course_id}/${card.exercise_id}`}
            className="text-control flex h-10 items-center rounded-md bg-accent px-5 text-accent-ink"
          >
            {last ? 'Resume lesson' : 'Start lesson'}
          </Link>
          <Link to={`/courses/${card.course_id}`} className="text-control flex h-10 items-center rounded-md bg-surface-raised px-5 text-text">
            View chapter
          </Link>
        </div>
      </div>
      <div className="flex flex-col gap-3 border-border bg-bg p-6 max-md:border-t md:border-l">
        {last && last.lines.length > 0 ? (
          <>
            <p className="text-code-sm text-text-muted">{last.file} · where you left off</p>
            <pre className="text-code-sm -mx-6 overflow-x-auto" aria-label="Your last answer">
              {last.lines.map((l) => (
                <div key={l.n} className={`flex gap-4 px-6 ${l.n === last.error_line ? 'bg-danger-soft' : ''}`}>
                  <span className={`w-6 shrink-0 text-right ${l.n === last.error_line ? 'text-danger' : 'text-text-faint'}`}>{l.n}</span>
                  <code className="whitespace-pre text-text">{l.text}</code>
                </div>
              ))}
            </pre>
            {last.total > 0 && (
              <p className="text-small flex items-center gap-2 text-text-secondary">
                <span className="h-2 w-2 rounded-pill bg-danger" aria-hidden="true" />
                {last.total - last.passed} of {last.total} checks were failing
              </p>
            )}
          </>
        ) : (
          <p className="text-small text-text-muted">
            Up next in {card.chapter_title}: {card.exercise_title} ({typeLabel[card.exercise_type] ?? 'lesson'}).
          </p>
        )}
      </div>
    </section>
  );
}
