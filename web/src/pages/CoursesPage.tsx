import { Link } from 'react-router';
import { useCourses } from '../api/queries';
import { AppHeader } from '../components/AppHeader';
import { LoadError, Loading } from './states';

export function CoursesPage() {
  const { data, error, isPending } = useCourses();
  return (
    <>
      <AppHeader />
      <main className="mx-auto flex max-w-[960px] flex-col gap-6 px-5 py-8">
        <h1 className="text-page-title text-text">Courses</h1>
        {isPending && <Loading label="courses" />}
        {error && <LoadError error={error} />}
        {data && data.courses.length === 0 && <p className="text-body text-text-muted">No courses yet. Run `make content`.</p>}
        <ul className="grid gap-2 sm:grid-cols-2">
          {data?.courses.map((c) => (
            <li key={c.id}>
              <Link
                to={`/courses/${c.id}`}
                className="flex h-full flex-col gap-2 rounded-lg border border-border bg-surface p-5 hover:border-border-strong"
              >
                <span className="text-section text-text">{c.title}</span>
                <span className="text-small text-text-secondary">{c.description}</span>
                <span className="text-caption text-text-muted">
                  {c.chapter_count} chapters · {c.exercise_count} lessons · {c.passed_count} done
                </span>
              </Link>
            </li>
          ))}
        </ul>
      </main>
    </>
  );
}
