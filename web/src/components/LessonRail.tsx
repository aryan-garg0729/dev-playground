import { Link } from 'react-router';
import { useRail } from '../stores/rail';
import { IconButton } from './Button';
import { CheckIcon, ChevronLeftIcon, LockIcon, PanelLeftIcon } from './icons';

export type LessonStatus = 'done' | 'current' | 'todo';

export interface RailLesson {
  id: string;
  title: string;
  status: LessonStatus;
  boss: boolean;
}

interface Props {
  courseId: string;
  chapterNumber: number;
  chapterTitle: string;
  lessons: RailLesson[];
}

const segmentColor: Record<LessonStatus, string> = {
  done: 'bg-success',
  current: 'bg-progress-current',
  todo: 'bg-border-strong',
};

const statusWord: Record<LessonStatus, string> = { done: 'done', current: 'current', todo: 'not started' };

function Segment({ status }: { status: LessonStatus }) {
  return <span aria-hidden="true" className={`block h-4 w-1 shrink-0 rounded-pill ${segmentColor[status]}`} />;
}

/**
 * The chapter outline and the only progress indicator on the lesson screen.
 * Closed (48px): one segment per lesson. Open (236px): lesson names with ticks and locks.
 */
export function LessonRail({ courseId, chapterNumber, chapterTitle, lessons }: Props) {
  const { open, setOpen } = useRail();
  const done = lessons.filter((l) => l.status === 'done').length;

  if (!open) {
    return (
      <nav
        aria-label="Lesson outline"
        className="flex w-[var(--rail-collapsed)] shrink-0 flex-col items-center gap-3 rounded-lg border border-border bg-surface py-3 max-[900px]:hidden"
      >
        <IconButton aria-label="Open outline" aria-expanded={false} className="border-0" onClick={() => setOpen(true)}>
          <PanelLeftIcon size={18} />
        </IconButton>
        <ol className="flex flex-col items-center gap-1">
          {lessons.map((l, i) => (
            <li key={l.id}>
              <button
                type="button"
                aria-label={`${i + 1}. ${l.title}, ${statusWord[l.status]}`}
                onClick={() => setOpen(true)}
                className="flex h-[22px] w-8 items-center justify-center rounded-sm hover:bg-surface-raised"
              >
                <Segment status={l.status} />
              </button>
            </li>
          ))}
        </ol>
      </nav>
    );
  }

  return (
    <nav
      aria-label="Lesson outline"
      className="flex w-[var(--rail-open)] shrink-0 flex-col gap-2 overflow-y-auto rounded-lg border border-border bg-surface p-3 max-[900px]:hidden"
    >
      <div className="flex items-center justify-between gap-2 pl-2">
        <p className="text-caption truncate text-text-muted">
          Chapter {chapterNumber} · {chapterTitle} · {done}/{lessons.length}
        </p>
        <IconButton aria-label="Close outline" aria-expanded className="shrink-0 border-0" onClick={() => setOpen(false)}>
          <ChevronLeftIcon />
        </IconButton>
      </div>
      <ol className="flex flex-col gap-0.5">
        {lessons.map((l, i) => {
          const current = l.status === 'current';
          return (
            <li key={l.id}>
              <Link
                to={`/learn/${courseId}/${l.id}`}
                aria-current={current ? 'page' : undefined}
                aria-label={`${i + 1}. ${l.title}, ${statusWord[l.status]}`}
                className={`text-small flex min-h-9 items-center gap-3 rounded-md px-2 ${
                  current
                    ? 'bg-surface-raised font-semibold text-text shadow-[inset_2px_0_0_var(--accent)] [[data-theme=light]_&]:bg-accent-soft'
                    : l.status === 'done'
                      ? 'text-text hover:bg-surface-raised'
                      : 'text-text-muted hover:bg-surface-raised'
                }`}
              >
                <Segment status={l.status} />
                <span className="min-w-0 flex-1 truncate">{l.title}</span>
                {l.status === 'done' && <CheckIcon className="shrink-0 text-success" />}
                {l.boss && l.status !== 'done' && <LockIcon className="shrink-0 text-text-faint" />}
              </Link>
            </li>
          );
        })}
      </ol>
    </nav>
  );
}
