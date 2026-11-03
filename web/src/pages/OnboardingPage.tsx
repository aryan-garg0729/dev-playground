import { type ReactNode, useState } from 'react';
import { Link, useNavigate } from 'react-router';
import { track } from '../analytics';
import type { CourseSummary } from '../api/client';
import { type MeUpdate, useCourses, useFinishOnboarding } from '../api/queries';
import { CodeIcon } from '../components/icons';
import { LoadError, Loading } from './states';

// screens/onboarding.png. Its "5-question check to skip what you know" has no
// design or question bank yet, so onboarding is two steps: goal and time, then
// the first course (SPEC Phase 6).

type Goal = NonNullable<MeUpdate['goal']>;

const goals: { id: Goal; title: string; body: string; match?: RegExp }[] = [
  { id: 'zero', title: 'Learn to code from zero', body: 'Never written a program? Start with the basics and build small tools.' },
  { id: 'backend', title: 'Get job-ready for backend', body: 'Python, SQL and the habits of writing tested code, with real exercises.' },
  { id: 'interviews', title: 'Prepare for interviews', body: 'Patterns, data structures and timed practice problems.' },
  { id: 'ai', title: 'Build with AI', body: 'LLM APIs, prompts, chains and evaluation, hands-on.', match: /lang|ai|llm/i },
];

const times = [
  { minutes: 10, goal: 1, label: 'About 1 lesson a day' },
  { minutes: 20, goal: 2, label: 'About 2 lessons a day' },
  { minutes: 40, goal: 4, label: 'A chapter a week' },
];

/** The course a goal starts with: a matching one, else the first. */
function recommend(goal: Goal, courses: CourseSummary[]): CourseSummary | undefined {
  const m = goals.find((g) => g.id === goal)?.match;
  return (m && courses.find((c) => m.test(c.id) || m.test(c.title))) || courses[0];
}

function Choice({ selected, onSelect, children, name }: { selected: boolean; onSelect: () => void; children: ReactNode; name: string }) {
  return (
    <label
      className={`flex cursor-pointer gap-3 rounded-lg border bg-surface p-6 has-[:focus-visible]:outline has-[:focus-visible]:outline-2 has-[:focus-visible]:outline-accent ${
        selected ? 'border-accent' : 'border-border hover:border-border-strong'
      }`}
    >
      <input type="radio" name={name} checked={selected} onChange={onSelect} className="sr-only" />
      <div className="flex min-w-0 flex-1 flex-col gap-2">{children}</div>
      <span aria-hidden className={`mt-1 h-5 w-5 shrink-0 rounded-pill border-2 ${selected ? 'border-accent bg-accent shadow-[inset_0_0_0_3px_var(--surface)]' : 'border-border-strong'}`} />
    </label>
  );
}

export function OnboardingPage() {
  const courses = useCourses();
  const finish = useFinishOnboarding();
  const navigate = useNavigate();
  const [step, setStep] = useState(1);
  const [goal, setGoal] = useState<Goal>('zero');
  const [minutes, setMinutes] = useState(20);
  const [course, setCourse] = useState<string>();

  if (courses.isPending) return <Loading label="courses" />;
  if (courses.error) return <LoadError error={courses.error} />;
  const list = courses.data.courses;
  const first = list[0];
  const picked = course ?? recommend(goal, list)?.id;

  const done = (body: MeUpdate, to: string, skipped: boolean) =>
    finish.mutate(body, {
      onSuccess: () => {
        track('onboarding_completed', { goal: body.goal ?? '', skipped, course: to });
        navigate(to, { replace: true });
      },
    });

  return (
    <div className="min-h-dvh bg-bg">
      <header className="flex h-16 items-center justify-between px-4 sm:px-8">
        <Link to="/" className="flex items-center gap-2" aria-label="CodeQuest home">
          <span className="flex h-8 w-8 items-center justify-center rounded-md bg-accent text-accent-ink"><CodeIcon /></span>
          <span className="text-brand text-text">CodeQuest</span>
        </Link>
        <div className="flex items-center gap-3" aria-label={`Step ${step} of 2`}>
          <span className="text-small text-text-muted">Step {step} of 2</span>
          {[1, 2].map((s) => <span key={s} aria-hidden className={`h-1 w-8 rounded-pill ${s <= step ? 'bg-accent' : 'bg-border-strong'}`} />)}
        </div>
      </header>

      <main className="mx-auto flex max-w-[880px] flex-col gap-6 px-4 py-10 sm:px-6">
        {step === 1 ? (
          <>
            <div>
              <h1 className="text-[32px] font-bold leading-10 text-text">What do you want CodeQuest to help you do?</h1>
              <p className="text-body mt-1 text-text-muted">We'll pick your first course. You can switch any time.</p>
            </div>
            <fieldset className="grid gap-4 md:grid-cols-2">
              <legend className="sr-only">Your goal</legend>
              {goals.map((g) => (
                <Choice key={g.id} name="goal" selected={goal === g.id} onSelect={() => { setGoal(g.id); setCourse(undefined); }}>
                  <span className="text-[17px] font-semibold text-text">{g.title}</span>
                  <span className="text-small text-text-secondary">{g.body}</span>
                  <span className="text-caption text-text-muted">Starts with {recommend(g.id, list)?.title ?? 'the first course'}</span>
                </Choice>
              ))}
            </fieldset>
            <fieldset className="flex flex-col gap-3">
              <legend className="text-section mb-3 text-text">How much time per day?</legend>
              <div className="grid gap-3 sm:grid-cols-3">
                {times.map((t) => (
                  <label key={t.minutes}
                    className={`flex cursor-pointer flex-col gap-1 rounded-lg border bg-surface px-5 py-4 has-[:focus-visible]:outline has-[:focus-visible]:outline-2 has-[:focus-visible]:outline-accent ${minutes === t.minutes ? 'border-accent' : 'border-border hover:border-border-strong'}`}>
                    <input type="radio" name="time" checked={minutes === t.minutes} onChange={() => setMinutes(t.minutes)} className="sr-only" />
                    <span className="text-game text-text">{t.minutes} min</span>
                    <span className="text-small text-text-secondary">{t.label}</span>
                  </label>
                ))}
              </div>
            </fieldset>
          </>
        ) : (
          <>
            <div>
              <h1 className="text-[32px] font-bold leading-10 text-text">Where do you want to start?</h1>
              <p className="text-body mt-1 text-text-muted">We picked one for your goal. The first chapters of every course are free.</p>
            </div>
            <fieldset className="grid gap-4 md:grid-cols-2">
              <legend className="sr-only">First course</legend>
              {list.map((c) => (
                <Choice key={c.id} name="course" selected={picked === c.id} onSelect={() => setCourse(c.id)}>
                  <span className="text-[17px] font-semibold text-text">{c.title}</span>
                  <span className="text-small text-text-secondary">{c.description}</span>
                  <span className="text-caption text-text-muted">{c.chapter_count} chapters · {c.exercise_count} lessons</span>
                </Choice>
              ))}
            </fieldset>
          </>
        )}

        {finish.error && <p className="text-small text-danger" role="alert">{finish.error.message}</p>}
        <div className="flex flex-wrap items-center justify-between gap-4 border-t border-border pt-6">
          {step === 1 ? (
            <button type="button" className="text-small text-text-secondary underline disabled:opacity-60" disabled={finish.isPending}
              onClick={() => done({}, first ? `/courses/${first.id}` : '/', true)}>
              Skip{first ? ` and start ${first.title}` : ''}
            </button>
          ) : (
            <button type="button" className="text-small text-text-secondary underline" onClick={() => setStep(1)}>Back</button>
          )}
          <div className="flex items-center gap-4">
            {step === 1 && <span className="text-small text-text-muted max-sm:hidden">Next: pick your first course</span>}
            <button
              type="button"
              disabled={finish.isPending || (step === 2 && !picked)}
              className="text-control h-12 rounded-md bg-accent px-8 text-[16px] text-accent-ink hover:brightness-105 disabled:opacity-60"
              onClick={() => {
                if (step === 1) {
                  track('onboarding_step', { step: 1, goal, minutes });
                  setStep(2);
                  return;
                }
                const dailyGoal = times.find((t) => t.minutes === minutes)?.goal ?? 2;
                done({ goal, daily_goal: dailyGoal }, picked ? `/courses/${picked}` : '/', false);
              }}
            >
              {step === 1 ? 'Continue' : 'Start learning'}
            </button>
          </div>
        </div>
      </main>
    </div>
  );
}
