import type { ReactNode } from 'react';
import { Link } from 'react-router';
import { track } from '../analytics';
import { useCourses } from '../api/queries';
import { CodeIcon } from '../components/icons';
import { Pill } from '../components/Pill';
import { Spinner } from '../components/Spinner';
import { ThemeToggle } from '../components/ThemeToggle';
import { TryItCard } from './landing/TryItCard';

// Pricing, "For teams" and the cancellation FAQ from screens/landing.png wait
// for Phase 5 (payments); they are left out rather than linking nowhere.

const features = [
  { mark: <span aria-hidden className="text-accent-text font-mono text-[18px]">&lt;&gt;</span>, title: 'Real code, graded instantly',
    body: 'Every lesson ends in an exercise. Run it in the browser and see which checks pass, with the exact value that went wrong.' },
  { mark: <Bars />, title: 'A path, not a pile of videos',
    body: 'Chapters build on each other and end with a boss fight that proves you can put the pieces together.' },
  { mark: <span className="text-game text-accent-text">+35 XP</span>, title: 'Progress you can feel',
    body: 'XP, levels, daily streaks and a weekly league keep you coming back, without getting in the way of the work.' },
];

const steps = [
  { n: '01', title: 'Read one idea', body: 'A short explanation with an example. Two minutes, not twenty.' },
  { n: '02', title: 'Write the code', body: 'Solve it in the editor. Stuck? The tutor asks a question to point you the right way.' },
  { n: '03', title: 'Pass the checks', body: 'Hidden tests confirm it works. Earn XP, keep your streak, move to the next one.' },
];

const faqs = [
  { q: 'Do I need to install anything?', a: 'No. Code runs in your browser. Later projects show you how to set up a local environment.' },
  { q: "I've never coded. Is this for me?", a: 'Yes. The first course starts from zero and moves one small step at a time.' },
  { q: 'Does the AI tutor just give me the answer?', a: 'No. It asks questions and points at the problem so you still solve it yourself.' },
];

function Bars() {
  return (
    <span aria-hidden className="flex h-[18px] items-end gap-[3px]">
      {['bg-success', 'bg-success', 'bg-accent', 'bg-border-strong', 'bg-border-strong'].map((c, i) => (
        <span key={i} className={`h-full w-[4px] rounded-sm ${c}`} />
      ))}
    </span>
  );
}

function Cta({ to, children, primary, where }: { to: string; children: ReactNode; primary?: boolean; where: string }) {
  return (
    <Link
      to={to}
      onClick={() => track('landing_cta', { where })}
      className={`text-control inline-flex h-11 items-center rounded-md px-5 ${
        primary ? 'bg-accent text-accent-ink hover:brightness-105' : 'border border-border bg-surface-raised text-text'
      }`}
    >
      {children}
    </Link>
  );
}

function Paths() {
  const courses = useCourses();
  const list = courses.data?.courses ?? [];
  return (
    <section id="paths" aria-labelledby="paths-title" className="mx-auto max-w-[1120px] px-4 py-16 sm:px-6">
      <h2 id="paths-title" className="text-[32px] font-bold leading-10 text-text">Pick a path</h2>
      <p className="text-body mt-1 text-text-muted">Each path is a sequence of courses with projects along the way.</p>
      <div className="mt-6 grid gap-4 md:grid-cols-3">
        {courses.isPending && (
          <div className="flex justify-center p-6 md:col-span-3">
            <Spinner label="Loading courses…" />
          </div>
        )}
        {list.map((c, i) => (
          <article key={c.id} className="flex flex-col gap-2 rounded-lg border border-border bg-surface p-6">
            <span className="w-fit">
              <Pill variant={i === 0 ? 'difficulty' : 'topic'}>{i === 0 ? 'Beginner' : 'Intermediate'}</Pill>
            </span>
            <h3 className="text-[18px] font-semibold text-text">{c.title}</h3>
            <p className="text-small text-text-secondary">{c.description}</p>
            <p className="text-caption text-text-muted">
              {c.chapter_count} chapters · {c.exercise_count} lessons
            </p>
            <Link to={`/sign-in?next=${encodeURIComponent(`/courses/${c.id}`)}`} onClick={() => track('landing_cta', { where: 'path', course: c.id })}
              className="text-control mt-auto text-accent-text hover:underline">
              Start the path →
            </Link>
          </article>
        ))}
      </div>
    </section>
  );
}

/** Public landing page for signed-out visitors (screens/landing.png). */
export function LandingPage() {
  return (
    <div className="min-h-dvh bg-bg">
      <header className="border-b border-border">
        <div className="mx-auto flex h-16 max-w-[1120px] items-center justify-between px-4 sm:px-6">
          <div className="flex items-center gap-6">
            <Link to="/" className="flex items-center gap-2" aria-label="CodeQuest home">
              <span className="flex h-7 w-7 items-center justify-center rounded-md bg-accent text-accent-ink"><CodeIcon /></span>
              <span className="text-brand text-text">CodeQuest</span>
            </Link>
            <nav aria-label="Main" className="max-sm:hidden">
              <a href="#paths" className="text-small text-text-secondary hover:text-text">Courses</a>
            </nav>
          </div>
          <div className="flex items-center gap-3">
            <ThemeToggle />
            <Link to="/sign-in" className="text-control px-2 text-text max-sm:hidden">Sign in</Link>
            <Link to="/sign-in" onClick={() => track('landing_cta', { where: 'header' })}
              className="text-control inline-flex h-9 items-center rounded-md bg-accent px-4 text-accent-ink hover:brightness-105">
              Start free
            </Link>
          </div>
        </div>
      </header>

      <main>
        <section className="mx-auto grid max-w-[1120px] items-center gap-12 px-4 py-16 sm:px-6 lg:grid-cols-[1.1fr_1fr] lg:py-24">
          <div className="flex flex-col items-start gap-6">
            <span className="text-caption rounded-pill border border-border bg-surface px-3 py-1 text-text-secondary">Python · SQL · AI engineering</span>
            <h1 className="text-[44px] font-extrabold leading-[1.05] tracking-[-0.02em] text-text sm:text-[60px]">
              Learn to code
              <br />
              by writing code.
            </h1>
            <p className="max-w-[520px] text-[18px] leading-7 text-text-secondary">
              Short lessons, real exercises graded in your browser, and XP for every check you pass. No videos to sit through.
            </p>
            <div className="flex flex-wrap gap-3">
              <Cta to="/sign-in" primary where="hero">Start free</Cta>
              <a href="#paths" className="text-control inline-flex h-11 items-center rounded-md border border-border bg-surface-raised px-5 text-text">
                See the curriculum
              </a>
            </div>
            <p className="text-caption text-text-muted">The first chapters of every course are free. No card needed.</p>
          </div>
          <TryItCard />
        </section>

        <section aria-label="Why CodeQuest" className="border-y border-border bg-surface">
          <div className="mx-auto grid max-w-[1120px] gap-10 px-4 py-14 sm:px-6 md:grid-cols-3">
            {features.map((f) => (
              <div key={f.title} className="flex flex-col gap-2">
                <div className="flex h-6 items-center">{f.mark}</div>
                <h2 className="text-[17px] font-semibold text-text">{f.title}</h2>
                <p className="text-small text-text-secondary">{f.body}</p>
              </div>
            ))}
          </div>
        </section>

        <Paths />

        <section aria-labelledby="how-title" className="mx-auto max-w-[1120px] px-4 sm:px-6">
          <div className="rounded-lg border border-border bg-surface p-8">
            <h2 id="how-title" className="text-[24px] font-bold text-text">How a lesson works</h2>
            <ol className="mt-6 grid gap-8 md:grid-cols-3">
              {steps.map((s) => (
                <li key={s.n} className="flex flex-col gap-2">
                  <span className="text-game text-accent-text">{s.n}</span>
                  <h3 className="text-[16px] font-semibold text-text">{s.title}</h3>
                  <p className="text-small text-text-secondary">{s.body}</p>
                </li>
              ))}
            </ol>
          </div>
        </section>

        <section className="mx-auto grid max-w-[1120px] gap-12 px-4 py-16 sm:px-6 md:grid-cols-2">
          <div className="flex flex-col items-start gap-4">
            <h2 className="text-[28px] font-bold leading-9 text-text">Free to start.</h2>
            <p className="text-body text-text-secondary">
              Try the first chapters of every course for free. Your XP, streak and league standing are saved from the first lesson.
            </p>
            <Cta to="/sign-in" where="footer">Create your account</Cta>
          </div>
          <div>
            <h2 className="text-section border-b border-border pb-3 text-text">Questions</h2>
            <dl>
              {faqs.map((f) => (
                <div key={f.q} className="border-b border-border py-4">
                  <dt className="text-small font-semibold text-text">{f.q}</dt>
                  <dd className="text-small mt-1 text-text-secondary">{f.a}</dd>
                </div>
              ))}
            </dl>
          </div>
        </section>
      </main>

      <footer className="border-t border-border">
        <div className="mx-auto flex max-w-[1120px] items-center justify-between px-4 py-6 sm:px-6">
          <span className="text-caption text-text-muted">© 2026 CodeQuest</span>
          <a href="#paths" className="text-caption text-text-muted hover:text-text">Courses</a>
        </div>
      </footer>
    </div>
  );
}
