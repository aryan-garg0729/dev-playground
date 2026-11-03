import { Link, useParams } from 'react-router';
import { useCourse } from '../api/queries';
import { AppHeader } from '../components/AppHeader';
import { ChapterTimeline } from './course/ChapterTimeline';
import { LoadError, Loading } from './states';

/** The course map (screens/course-map.png). */
export function CoursePage() {
  const { courseId = '' } = useParams();
  const { data, error, isPending } = useCourse(courseId);
  const lessons = data?.chapters.flatMap((c) => c.exercises) ?? [];
  const done = lessons.filter((e) => e.status === 'passed').length;
  const pct = lessons.length ? Math.round((done / lessons.length) * 100) : 0;
  const next = lessons.find((e) => e.status !== 'passed');

  return (
    <>
      <AppHeader />
      {isPending && <Loading label="course" />}
      {error && <LoadError error={error} />}
      {data && (
        <>
          <section className="border-b border-border">
            <div className="mx-auto grid max-w-[1200px] gap-6 px-5 py-8 lg:grid-cols-[1fr_380px] lg:items-end">
              <div className="flex flex-col gap-3">
                <nav aria-label="Breadcrumb" className="text-small text-text-muted">
                  <Link to="/courses" className="hover:text-text">
                    Courses
                  </Link>{' '}
                  › {data.title}
                </nav>
                <h1 className="text-[32px] font-bold leading-[40px] text-text">{data.title}</h1>
                <p className="text-body max-w-[640px] text-text-secondary">{data.description}</p>
                <p className="text-small flex gap-5 text-text-muted">
                  <span>{data.chapters.length} chapters</span>
                  <span>{lessons.length} lessons</span>
                </p>
              </div>
              <div className="flex flex-col gap-3">
                <p className="text-small flex justify-between text-text-secondary">
                  <span>
                    {done} of {lessons.length} lessons
                  </span>
                  <span className="font-semibold text-success">{pct}%</span>
                </p>
                <span role="img" aria-label={`${pct}% complete`} className="block h-[6px] overflow-hidden rounded-pill bg-border">
                  <span className="block h-full rounded-pill bg-success" style={{ width: `${pct}%` }} />
                </span>
                {next && (
                  <Link
                    to={`/learn/${data.id}/${next.id}`}
                    className="text-control flex h-11 items-center justify-center rounded-md bg-accent text-accent-ink"
                  >
                    {done === 0 ? 'Start' : 'Continue'}: {next.title}
                  </Link>
                )}
              </div>
            </div>
          </section>
          <main className="mx-auto grid max-w-[1200px] gap-8 px-5 py-8 lg:grid-cols-[1fr_380px]">
            <ChapterTimeline course={data} />
            <aside className="flex flex-col gap-4">
              {data.skills.length > 0 && (
                <section className="flex flex-col gap-4 rounded-lg border border-border bg-surface p-6">
                  <h2 className="text-section text-text">Skills you'll leave with</h2>
                  <ul className="flex flex-wrap gap-2">
                    {data.skills.map((s) => (
                      <li key={s} className="text-small rounded-pill bg-surface-raised px-3 py-1 text-text-secondary">
                        {s}
                      </li>
                    ))}
                  </ul>
                </section>
              )}
              <section className="flex flex-col gap-3 rounded-lg border border-border bg-surface p-6">
                <h2 className="text-section text-text">XP in this course</h2>
                <p className="flex items-baseline gap-2">
                  <span className="text-game text-[32px] leading-[36px] text-accent-text">{data.xp_earned}</span>
                  <span className="text-body text-text-secondary">of about {data.xp_total}</span>
                </p>
                <p className="text-small text-text-muted">
                  First-try passes give +50%, and passes without a hint +25%. Re-solving a lesson gives no XP.
                </p>
              </section>
            </aside>
          </main>
        </>
      )}
    </>
  );
}
