import { Link } from 'react-router';
import type { CareerPath } from '../../api/client';

type Step = CareerPath['steps'][number];

function Dot({ status }: { status: Step['status'] }) {
  const cls = {
    in_progress: 'bg-accent border-accent',
    done: 'bg-success border-success',
    not_started: 'border-border-strong bg-surface',
  }[status];
  return <span aria-hidden="true" className={`relative z-[1] mt-1 h-3 w-3 shrink-0 rounded-pill border-2 ${cls}`} />;
}

const stepLine = (s: Step) => (s.status === 'done' ? 'Done' : s.status === 'in_progress' ? 'In progress' : s.note);

export function CareerPathCard({ p }: { p: CareerPath }) {
  const pct = p.exercise_count > 0 ? Math.round((p.passed_count / p.exercise_count) * 100) : 0;
  const started = p.passed_count > 0;
  const next = p.steps.find((s) => s.status !== 'done') ?? p.steps[0];
  if (!next) return null;
  const courses = `${p.steps.length} ${p.steps.length === 1 ? 'course' : 'courses'}`;
  return (
    <li
      aria-label={`${p.title}${p.yours ? ', your path' : ''}`}
      className={`flex flex-col gap-4 rounded-lg border bg-surface p-5 ${p.yours ? 'border-accent-line' : 'border-border'}`}
    >
      <div className="flex items-start justify-between gap-3">
        <h3 className="text-section text-text">{p.title}</h3>
        {p.yours && <span className="text-caption shrink-0 rounded-pill bg-accent-soft px-2 py-0.5 text-accent-text">Your path</span>}
      </div>
      {p.description && <p className="text-small -mt-2 text-text-secondary">{p.description}</p>}
      <ol className="flex flex-col gap-3">
        {p.steps.map((s, i) => (
          <li key={s.course_id} className="relative flex gap-3">
            {i < p.steps.length - 1 && <span aria-hidden="true" className="absolute left-[5px] top-4 h-[calc(100%+4px)] w-px bg-border" />}
            <Dot status={s.status} />
            <Link to={`/courses/${s.course_id}`} className="flex min-w-0 flex-col hover:underline">
              <span className={`text-small ${s.status === 'not_started' ? 'text-text-secondary' : 'font-semibold text-text'}`}>{s.title}</span>
              <span className="text-caption text-text-muted">{stepLine(s)}</span>
            </Link>
          </li>
        ))}
      </ol>
      <div className="text-small mt-auto flex items-center justify-between gap-3 border-t border-border pt-3">
        <span className="text-text-muted">{started ? `${pct}% · ${courses}` : courses}</span>
        <Link to={`/courses/${next.course_id}`} className="font-semibold text-accent-text hover:underline">
          {started ? 'Continue' : 'Start path'}
        </Link>
      </div>
    </li>
  );
}
