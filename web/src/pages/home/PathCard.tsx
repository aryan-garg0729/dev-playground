import { Link } from 'react-router';
import type { Home } from '../../api/client';

type Path = Home['paths'][number];

const pill = {
  in_progress: ['In progress', 'bg-success-soft text-success'],
  done: ['Done', 'bg-success-soft text-success'],
  not_started: ['Not started', 'bg-surface-raised text-text-secondary'],
} as const;

/** "8 chapters · 2 done · next boss: FizzBuzz" or "After Foundations · Chat models, prompts". */
export function pathLine(p: Path): string {
  if (p.up_next) return [p.after && `After ${p.after}`, p.topics || p.description].filter(Boolean).join(' · ');
  const parts = [`${p.chapter_count} ${p.chapter_count === 1 ? 'chapter' : 'chapters'}`, `${p.chapters_done} done`];
  if (p.next_boss) parts.push(`next boss: ${p.next_boss}`);
  return parts.join(' · ');
}

export function PathCard({ p }: { p: Path }) {
  const pct = p.exercise_count > 0 ? (p.passed_count / p.exercise_count) * 100 : 0;
  const [text, cls] = p.up_next ? ['Up next', 'bg-surface-raised text-text-secondary'] : pill[p.status];
  return (
    <li>
      <Link to={`/courses/${p.id}`} className="flex h-full flex-col gap-3 rounded-lg border border-border bg-surface p-5 hover:border-border-strong">
        <span className="flex items-start justify-between gap-3">
          <span className="text-section text-text">{p.title}</span>
          <span className={`text-caption shrink-0 rounded-pill px-2 py-0.5 ${cls}`}>{text}</span>
        </span>
        <span className="text-small text-text-secondary">{pathLine(p)}</span>
        <span
          role="img"
          aria-label={`${p.passed_count} of ${p.exercise_count} lessons done`}
          className="mt-auto block h-[6px] overflow-hidden rounded-pill bg-border"
        >
          <span className="block h-full rounded-pill bg-success" style={{ width: `${pct}%` }} />
        </span>
      </Link>
    </li>
  );
}
