import { useEffect, useMemo, useRef, useState } from 'react';
import type { CourseSummary } from '../api/client';
import { useCourses } from '../api/queries';
import { AppHeader } from '../components/AppHeader';
import { CourseCard } from './courses/CourseCard';
import { CareerPathCard } from './courses/PathCard';
import { LoadError, Loading } from './states';

/** Courses whose title, description, topics or tags contain every word of the query. */
export function filterCourses(courses: CourseSummary[], query: string, tag: string): CourseSummary[] {
  const words = query.toLowerCase().split(/\s+/).filter(Boolean);
  return courses.filter((c) => {
    if (tag && !c.tags.includes(tag)) return false;
    const text = [c.title, c.description, c.topics, ...c.tags].join(' ').toLowerCase();
    return words.every((w) => text.includes(w));
  });
}

function SearchBox({ value, onChange }: { value: string; onChange: (v: string) => void }) {
  const ref = useRef<HTMLInputElement>(null);
  useEffect(() => {
    // "/" focuses search unless the learner is already typing somewhere.
    const onKey = (e: KeyboardEvent) => {
      const t = e.target instanceof Element ? e.target : null;
      if (e.key !== '/' || e.metaKey || e.ctrlKey || t?.closest('input, textarea, select, [contenteditable="true"]')) return;
      e.preventDefault();
      ref.current?.focus();
    };
    window.addEventListener('keydown', onKey);
    return () => window.removeEventListener('keydown', onKey);
  }, []);
  return (
    <label className="flex h-10 w-full items-center gap-2 rounded-md border border-border bg-surface px-3 focus-within:border-border-strong md:w-[340px]">
      <span className="sr-only">Search courses and topics</span>
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" className="shrink-0 text-text-muted" aria-hidden="true">
        <circle cx="11" cy="11" r="7" />
        <path d="m20 20-3.5-3.5" />
      </svg>
      <input
        ref={ref}
        type="search"
        value={value}
        onChange={(e) => onChange(e.target.value)}
        placeholder="Search courses and topics"
        className="text-body min-w-0 flex-1 bg-transparent text-text outline-none placeholder:text-text-muted"
      />
      <kbd className="text-caption rounded-sm border border-border px-1.5 text-text-muted" aria-hidden="true">
        /
      </kbd>
    </label>
  );
}

export function CoursesPage() {
  const { data, error, isPending } = useCourses();
  const [query, setQuery] = useState('');
  const [tag, setTag] = useState('');
  const courses = useMemo(() => data?.courses ?? [], [data]);
  const tags = useMemo(() => [...new Set(courses.flatMap((c) => c.tags))], [courses]);
  const shown = filterCourses(courses, query, tag);
  return (
    <>
      <AppHeader />
      <main className="mx-auto flex max-w-[1200px] flex-col gap-8 px-5 py-8">
        <div className="flex flex-wrap items-end justify-between gap-4">
          <div className="flex flex-col gap-1">
            <h1 className="text-page-title text-text">Courses</h1>
            <p className="text-body text-text-secondary">Follow a career path, or pick any course. The first chapters of every course are free.</p>
          </div>
          <SearchBox value={query} onChange={setQuery} />
        </div>
        {isPending && <Loading label="courses" />}
        {error && <LoadError error={error} />}
        {data && data.paths.length > 0 && (
          <section aria-labelledby="paths-title" className="flex flex-col gap-3">
            <h2 id="paths-title" className="text-section flex flex-wrap items-baseline gap-x-3 text-text">
              Career paths <span className="text-small font-normal text-text-muted">Courses in order, from first steps to real projects</span>
            </h2>
            <ul className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
              {data.paths.map((p) => (
                <CareerPathCard key={p.id} p={p} />
              ))}
            </ul>
          </section>
        )}
        {data && (
          <section aria-labelledby="all-title" className="flex flex-col gap-3">
            <div className="flex flex-wrap items-center justify-between gap-3">
              <h2 id="all-title" className="text-section flex items-baseline gap-2 text-text">
                All courses <span className="text-small font-normal text-text-muted">{shown.length} {shown.length === 1 ? 'course' : 'courses'}</span>
              </h2>
              {tags.length > 1 && (
                <div role="group" aria-label="Filter by topic" className="flex flex-wrap gap-2">
                  {['', ...tags].map((t) => (
                    <button
                      key={t || 'all'}
                      type="button"
                      aria-pressed={tag === t}
                      onClick={() => setTag(t)}
                      className={`text-small h-8 rounded-pill border px-3 ${tag === t ? 'border-text bg-text text-bg' : 'border-border text-text-secondary hover:border-border-strong hover:text-text'}`}
                    >
                      {t || 'All'}
                    </button>
                  ))}
                </div>
              )}
            </div>
            {courses.length === 0 && <p className="text-body text-text-muted">No courses yet. Run `make content`.</p>}
            {courses.length > 0 && shown.length === 0 && <p className="text-body text-text-muted">No courses match. Try another word or topic.</p>}
            <ul className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
              {shown.map((c) => (
                <CourseCard key={c.id} c={c} />
              ))}
            </ul>
          </section>
        )}
      </main>
    </>
  );
}
