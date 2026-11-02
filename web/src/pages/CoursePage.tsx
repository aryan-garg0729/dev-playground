import { Link, useParams } from 'react-router';
import type { Schemas } from '../api/client';
import { useCourse } from '../api/queries';
import { AppHeader } from '../components/AppHeader';
import { CheckIcon } from '../components/icons';
import { LoadError, Loading } from './states';

type Status = Schemas['ProgressStatus'];

const segment: Record<Status, string> = {
  passed: 'bg-success',
  started: 'bg-progress-current',
  not_started: 'bg-border-strong',
};

const statusLabel: Record<Status, string> = { passed: 'done', started: 'started', not_started: 'not started' };

const typeLabel: Record<string, string> = { mcq: 'Quiz', code: 'Exercise', sql: 'SQL' };

export function CoursePage() {
  const { courseId = '' } = useParams();
  const { data, error, isPending } = useCourse(courseId);
  const lessons = data?.chapters.flatMap((c) => c.exercises) ?? [];
  const done = lessons.filter((e) => e.status === 'passed').length;

  return (
    <>
      <AppHeader />
      <main className="mx-auto flex max-w-[960px] flex-col gap-6 px-5 py-8">
        {isPending && <Loading label="course" />}
        {error && <LoadError error={error} />}
        {data && (
          <>
            <div className="flex flex-col gap-2">
              <nav aria-label="Breadcrumb" className="text-small text-text-muted">
                <Link to="/" className="hover:text-text">
                  Courses
                </Link>{' '}
                › {data.title}
              </nav>
              <h1 className="text-[32px] font-bold leading-[40px] text-text">{data.title}</h1>
              <p className="text-body max-w-[640px] text-text-secondary">{data.description}</p>
              <p className="text-small text-text-muted">
                {data.chapters.length} chapters · {lessons.length} lessons · {done} done
              </p>
            </div>

            <ol className="flex flex-col gap-2">
              {data.chapters.map((ch, ci) => (
                <li key={ch.id} className="rounded-lg border border-border bg-surface p-5">
                  <h2 className="text-section text-text">
                    {ci + 1}. {ch.title}
                  </h2>
                  <p className="text-caption mt-1 text-text-muted">{ch.exercises.length} lessons</p>
                  <ol className="mt-3 flex flex-col">
                    {ch.exercises.map((ex, ei) => (
                      <li key={ex.id} className="border-t border-border first:border-t-0">
                        <Link
                          to={`/learn/${data.id}/${ex.id}`}
                          className="flex min-h-[var(--touch-min)] items-center gap-3 rounded-md px-2 hover:bg-surface-raised"
                          aria-label={`${ei + 1}. ${ex.title}, ${statusLabel[ex.status]}`}
                        >
                          <span className={`h-4 w-1 shrink-0 rounded-pill ${segment[ex.status]}`} aria-hidden="true" />
                          <span className="text-body flex-1 text-text">
                            {ei + 1}. {ex.title}
                          </span>
                          {ex.status === 'passed' && <CheckIcon className="text-success" />}
                          <span className="text-caption text-text-muted">
                            {typeLabel[ex.type] ?? ex.type} · {ex.xp} XP
                          </span>
                        </Link>
                      </li>
                    ))}
                  </ol>
                </li>
              ))}
            </ol>
          </>
        )}
      </main>
    </>
  );
}
