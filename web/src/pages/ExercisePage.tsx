import { type ReactNode, useState } from 'react';
import { Link, useParams } from 'react-router';
import type { ExerciseDetail, GradeResult } from '../api/client';
import { useExercise, useSubmit } from '../api/queries';
import { AppHeader } from '../components/AppHeader';
import { ChevronLeftIcon, ChevronRightIcon } from '../components/icons';
import { DifficultyPill, Pill } from '../components/Pill';
import { renderers } from '../exercise-types';
import { LoadError, Loading } from './states';

function LessonNav({ ex }: { ex: ExerciseDetail }) {
  const nav = (id: string | null, label: string, icon: ReactNode) =>
    id ? (
      <Link
        to={`/learn/${ex.course_id}/${id}`}
        aria-label={label}
        className="inline-flex h-8 w-8 items-center justify-center rounded-md border border-border text-text-secondary hover:bg-surface-raised"
      >
        {icon}
      </Link>
    ) : (
      <span aria-hidden="true" className="inline-flex h-8 w-8 items-center justify-center rounded-md border border-border text-text-faint">
        {icon}
      </span>
    );
  return (
    <div className="flex items-center justify-center gap-4">
      {nav(ex.prev_id, 'Previous lesson', <ChevronLeftIcon />)}
      <div className="min-w-0 text-center">
        <p className="text-small truncate font-semibold text-text">{ex.title}</p>
        <p className="text-caption truncate text-text-muted">
          {ex.chapter.title} · lesson {ex.position.index} of {ex.position.total}
        </p>
      </div>
      {nav(ex.next_id, 'Next lesson', <ChevronRightIcon />)}
    </div>
  );
}

function Exercise({ ex }: { ex: ExerciseDetail }) {
  const last = ex.last_submission;
  const sameVersion = last?.exercise_version === ex.version;
  const [draft, setDraft] = useState<unknown>(sameVersion ? last?.answer : undefined);
  const [result, setResult] = useState<GradeResult | undefined>(
    sameVersion && last ? { submission_id: last.submission_id, result: last.result, xp_awarded: 0 } : undefined,
  );
  const submit = useSubmit(ex.id, ex.course_id);
  const Renderer = renderers[ex.type];

  if (!Renderer) {
    return <p className="text-body text-text-muted">This lesson type ({ex.type}) is not supported yet.</p>;
  }
  return (
    <Renderer
      publicView={ex.public_view}
      draft={draft}
      result={result}
      onChange={(a) => {
        setDraft(a);
        setResult(undefined); // the shown result belongs to the previous answer
      }}
      onSubmit={async (answer) => {
        const r = await submit.mutateAsync({ version: ex.version, answer });
        setResult(r);
        return r;
      }}
    />
  );
}

export function ExercisePage() {
  const { exerciseId = '' } = useParams();
  const { data, error, isPending } = useExercise(exerciseId);
  return (
    <>
      <AppHeader center={data && <LessonNav ex={data} />} />
      <main className="mx-auto max-w-[760px] p-2">
        {isPending && <Loading label="lesson" />}
        {error && <LoadError error={error} />}
        {data && (
          <article className="flex flex-col gap-4 rounded-lg border border-border bg-surface px-6 py-5">
            <div className="flex flex-wrap gap-1.5">
              <DifficultyPill difficulty={data.difficulty} />
              {data.tags[0] && <Pill variant="topic">{data.tags[0]}</Pill>}
              <Pill variant="xp">+{data.xp} XP</Pill>
            </div>
            <h1 className="text-page-title text-text">{data.title}</h1>
            <Exercise key={data.id} ex={data} />
          </article>
        )}
      </main>
    </>
  );
}
