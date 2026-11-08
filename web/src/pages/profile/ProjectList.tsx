import { Link } from 'react-router';
import type { Profile } from '../../api/client';
import { LockIcon } from '../../components/icons';

type Project = Profile['projects'][number];

const day = (iso: string) => new Date(iso).toLocaleDateString('en-US', { month: 'short', day: 'numeric' });

/** "Foundations · Chapter 1 boss · defeated first try · Oct 29". */
export function projectLine(p: Project): string {
  const parts = [p.course_title, `Chapter ${p.chapter_index} boss`];
  if (p.status === 'defeated') {
    if (p.first_try) parts.push('defeated first try');
    if (p.defeated_at) parts.push(day(p.defeated_at));
  } else {
    parts[1] = `${p.status === 'next' ? 'Next boss' : 'Boss'} · Chapter ${p.chapter_index}: ${p.chapter_title}`;
  }
  return parts.join(' · ');
}

export function ProjectRow({ p, icon }: { p: Project; icon: string }) {
  const defeated = p.status === 'defeated';
  return (
    <li className="flex items-center gap-4 border-t border-border py-4 first:border-t-0">
      <span
        aria-hidden="true"
        className="text-game flex h-10 w-10 shrink-0 items-center justify-center rounded-md border border-border bg-surface-raised text-[13px] text-text-secondary"
      >
        {defeated ? icon : <LockIcon size={14} />}
      </span>
      <span className="flex min-w-0 flex-1 flex-col">
        <span className={`text-body font-semibold ${defeated ? 'text-text' : 'text-text-secondary'}`}>{p.title}</span>
        <span className="text-small truncate text-text-muted">{projectLine(p)}</span>
      </span>
      {defeated && (
        <Link
          to={`/learn/${p.course_id}/${p.exercise_id}`}
          className="text-control flex h-[34px] shrink-0 items-center rounded-md border border-border px-3 text-text hover:border-border-strong"
        >
          View code
        </Link>
      )}
      {p.status === 'next' && (
        <Link to={`/learn/${p.course_id}/${p.exercise_id}`} className="text-small shrink-0 font-semibold text-accent-text hover:underline">
          Open
        </Link>
      )}
    </li>
  );
}
