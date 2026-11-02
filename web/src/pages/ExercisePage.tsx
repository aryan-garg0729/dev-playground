import { type ReactNode, useState } from 'react';
import { Link, useParams } from 'react-router';
import type { CourseDetail, ExerciseDetail, GradeResult } from '../api/client';
import { useCourse, useExercise, useSubmit } from '../api/queries';
import { AppHeader } from '../components/AppHeader';
import { ChevronLeftIcon, ChevronRightIcon } from '../components/icons';
import { LessonRail, type RailLesson } from '../components/LessonRail';
import { DifficultyPill, Pill } from '../components/Pill';
import { renderers } from '../exercise-types';
import { loadDraft, saveDraft } from './drafts';
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

function Rail({ course, ex }: { course: CourseDetail; ex: ExerciseDetail }) {
  const index = course.chapters.findIndex((c) => c.id === ex.chapter.id);
  const chapter = course.chapters[index];
  if (!chapter) return null;
  const lessons: RailLesson[] = chapter.exercises.map((e) => ({
    id: e.id,
    title: e.title,
    boss: e.boss,
    status: e.id === ex.id ? 'current' : e.status === 'passed' ? 'done' : 'todo',
  }));
  return <LessonRail courseId={course.id} chapterNumber={index + 1} chapterTitle={chapter.title} lessons={lessons} />;
}

function Exercise({ ex }: { ex: ExerciseDetail }) {
  const last = ex.last_submission;
  const sameVersion = last?.exercise_version === ex.version;
  const [draft, setDraft] = useState<unknown>(() => loadDraft(ex.id, ex.version) ?? (sameVersion ? last?.answer : undefined));
  const [result, setResult] = useState<GradeResult | undefined>(
    sameVersion && last?.result ? { submission_id: last.submission_id, result: last.result, xp_awarded: 0 } : undefined,
  );
  const submit = useSubmit(ex.id, ex.course_id);
  const Renderer = renderers[ex.type];

  const header = (
    <div className="flex flex-col gap-4">
      <div className="flex flex-wrap gap-1.5">
        <DifficultyPill difficulty={ex.difficulty} />
        {ex.tags[0] && <Pill variant="topic">{ex.tags[0]}</Pill>}
        <Pill variant="xp">+{ex.xp} XP</Pill>
      </div>
      <h1 className="text-page-title text-text">{ex.title}</h1>
    </div>
  );

  if (!Renderer) {
    return <p className="text-body text-text-muted">This lesson type ({ex.type}) is not supported yet.</p>;
  }
  return (
    <Renderer
      header={header}
      publicView={ex.public_view}
      draft={draft}
      result={result}
      onChange={(a) => {
        setDraft(a);
        saveDraft(ex.id, ex.version, a);
        if (ex.type === 'mcq') setResult(undefined); // a changed choice makes the shown verdict stale
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
  const { courseId = '', exerciseId = '' } = useParams();
  const { data, error, isPending } = useExercise(exerciseId);
  const course = useCourse(courseId);
  return (
    <div className="flex h-dvh flex-col">
      <AppHeader center={data && <LessonNav ex={data} />} />
      <main className="flex min-h-0 flex-1 gap-2 p-2 max-[900px]:h-auto max-[900px]:overflow-y-auto">
        {data && course.data && <Rail course={course.data} ex={data} />}
        {isPending && <Loading label="lesson" />}
        {error && <LoadError error={error} />}
        {data && <Exercise key={data.id} ex={data} />}
      </main>
    </div>
  );
}
