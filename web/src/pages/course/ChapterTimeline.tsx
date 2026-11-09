import { useState } from 'react';
import { Link } from 'react-router';
import type { CourseDetail, Schemas } from '../../api/client';
import { useShop, useSkeletonKey } from '../../api/queries';
import { Button } from '../../components/Button';
import { CheckIcon, KeyIcon, LockIcon, StarIcon } from '../../components/icons';

type Chapter = CourseDetail['chapters'][number];
type Status = Schemas['ProgressStatus'];
type ChapterState = 'done' | 'current' | 'todo' | 'locked';

const typeLabel: Record<string, string> = { mcq: 'Quiz', code: 'Exercise', sql: 'SQL' };
const bar: Record<Status, string> = { passed: 'bg-success', started: 'bg-border-strong', not_started: 'bg-border-strong' };

/** The chapter to work on: the first with an unpassed lesson. */
export function currentChapter(chapters: Chapter[]): number {
  const i = chapters.findIndex((c) => c.exercises.some((e) => e.status !== 'passed'));
  return i < 0 ? chapters.length - 1 : i;
}

function Node({ state }: { state: ChapterState }) {
  if (state === 'locked') {
    return (
      <span className="relative z-[1] flex h-8 w-8 items-center justify-center rounded-pill border-[3px] border-border-strong bg-bg text-text-muted">
        <LockIcon size={12} />
      </span>
    );
  }
  if (state === 'done') {
    return (
      <span className="relative z-[1] flex h-8 w-8 items-center justify-center rounded-pill bg-success-soft text-success">
        <CheckIcon size={14} />
      </span>
    );
  }
  return (
    <span
      className={`relative z-[1] block h-8 w-8 rounded-pill border-[3px] bg-bg ${state === 'current' ? 'border-accent' : 'border-border-strong'}`}
    />
  );
}

function Lessons({ courseId, chapter, chapterNo }: { courseId: string; chapter: Chapter; chapterNo: number }) {
  const next = chapter.exercises.findIndex((e) => e.status !== 'passed');
  return (
    <ol className="flex flex-col" aria-label={`Chapter ${chapterNo} lessons`}>
      {chapter.exercises.map((e, i) => (
        <li key={e.id} className="border-t border-border">
          <Link
            to={`/learn/${courseId}/${e.id}`}
            className={`flex min-h-[var(--touch-min)] items-center gap-3 rounded-md px-2 py-2 hover:bg-surface-raised ${e.boss ? 'my-2 border border-dashed border-border' : ''}`}
            aria-label={`${i + 1}. ${e.title}${e.boss ? ', boss' : ''}, ${e.status === 'passed' ? 'done' : i === next ? 'next' : 'not started'}`}
          >
            <span className={`h-4 w-1 shrink-0 rounded-pill ${i === next ? 'bg-accent' : bar[e.status]}`} aria-hidden="true" />
            {e.boss && <StarIcon className="text-accent" />}
            <span className={`text-body flex-1 ${e.status === 'passed' ? 'text-text-secondary' : 'text-text'}`}>
              {e.boss ? 'Boss: ' : `${i + 1}. `}
              {e.title}
            </span>
            {e.status === 'passed' && <CheckIcon className="text-success" />}
            <span className="text-caption text-text-muted">
              {e.boss ? '' : `${typeLabel[e.type] ?? e.type} · `}
              {e.xp} XP
            </span>
          </Link>
        </li>
      ))}
    </ol>
  );
}

/** On the first locked chapter: open it with a held Skeleton Key (the boss before it stays undefeated). */
function KeyButton({ courseId, chapterId }: { courseId: string; chapterId: string }) {
  const { data } = useShop();
  const key = useSkeletonKey(courseId);
  const [problem, setProblem] = useState<string>();
  const held = data?.items.find((i) => i.id === 'skeleton_key')?.owned ?? 0;
  if (held === 0) return null;
  return (
    <div className="mt-2 flex flex-col items-start gap-1">
      <Button disabled={key.isPending} onClick={() => key.open(chapterId, (e) => setProblem(e.message))}>
        <KeyIcon size={14} /> Use Skeleton Key
      </Button>
      {problem && (
        <p role="alert" className="text-small text-danger">
          {problem}
        </p>
      )}
    </div>
  );
}

/** Chapters as a vertical path: done ✓, current (open), then what is ahead. */
export function ChapterTimeline({ course }: { course: CourseDetail }) {
  const current = currentChapter(course.chapters);
  const firstLocked = course.chapters.findIndex((c) => c.locked);
  const [open, setOpen] = useState<Set<number>>(() => new Set([current]));
  const toggle = (i: number) =>
    setOpen((s) => {
      const n = new Set(s);
      if (n.has(i)) n.delete(i);
      else n.add(i);
      return n;
    });

  return (
    <ol className="relative flex flex-col gap-2" aria-label="Chapters">
      <span aria-hidden="true" className="absolute bottom-4 left-[15px] top-4 w-0.5 bg-border" />
      {course.chapters.map((ch, i) => {
        const done = ch.exercises.filter((e) => e.status === 'passed').length;
        const state: ChapterState = ch.locked
          ? 'locked'
          : done === ch.exercises.length && ch.exercises.length > 0
            ? 'done'
            : i === current
              ? 'current'
              : 'todo';
        const xp = ch.exercises.reduce((n, e) => n + e.xp, 0);
        const isOpen = open.has(i) && !ch.locked;
        const summary = ch.locked
          ? `${ch.exercises.length} lessons · defeat chapter ${i}'s boss to unlock`
          : state === 'done'
            ? `${ch.exercises.length} lessons · all done`
            : `${ch.exercises.length} lessons · ${done} done`;
        return (
          <li key={ch.id} className="flex gap-5">
            <Node state={state} />
            <div className={`min-w-0 flex-1 ${isOpen ? 'rounded-lg border border-border bg-surface p-5' : 'py-0.5'}`}>
              <button
                type="button"
                aria-expanded={isOpen}
                disabled={ch.locked}
                onClick={() => toggle(i)}
                className="flex w-full items-start justify-between gap-3 text-left"
              >
                <span className="flex flex-col">
                  <span className="text-[17px] font-semibold leading-[24px] text-text">
                    {i + 1}. {ch.title}
                  </span>
                  <span className="text-small text-text-muted">{summary}</span>
                </span>
                {state === 'current' ? (
                  <span className="text-caption rounded-pill bg-accent-soft px-2 py-0.5 text-accent-text">Current</span>
                ) : (
                  <span className="text-game-sm text-text-muted">{xp} XP</span>
                )}
              </button>
              {i === firstLocked && <KeyButton courseId={course.id} chapterId={ch.id} />}
              {isOpen && (
                <div className="mt-4">
                  <Lessons courseId={course.id} chapter={ch} chapterNo={i + 1} />
                </div>
              )}
            </div>
          </li>
        );
      })}
    </ol>
  );
}
