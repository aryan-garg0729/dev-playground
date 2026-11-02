import type { ReactNode } from 'react';
import { Link, NavLink } from 'react-router';
import { useMe } from '../api/queries';
import { CodeIcon } from './icons';
import { ThemeToggle } from './ThemeToggle';

function initials(name: string, email: string): string {
  const src = name.trim() || email;
  const parts = src.split(/[\s@._-]+/).filter(Boolean);
  return ((parts[0]?.[0] ?? '') + (parts[1]?.[0] ?? '')).toUpperCase() || '?';
}

/** 56px app header. `center` holds page context (e.g. lesson name and prev/next). */
export function AppHeader({ center }: { center?: ReactNode }) {
  const me = useMe();
  return (
    <header className="sticky top-0 z-10 grid h-[var(--header-height)] grid-cols-[1fr_auto_1fr] items-center border-b border-border bg-surface px-5">
      <div className="flex items-center gap-6">
        <Link to="/" className="flex items-center gap-2" aria-label="CodeQuest home">
          <span className="flex h-7 w-7 items-center justify-center rounded-md bg-accent text-accent-ink">
            <CodeIcon />
          </span>
          <span className="text-brand text-text max-sm:hidden">CodeQuest</span>
        </Link>
        {!center && (
          <nav aria-label="Main">
            <NavLink
              to="/"
              className={({ isActive }) =>
                `text-small rounded-md px-3 py-1.5 ${isActive ? 'bg-surface-raised text-text' : 'text-text-muted hover:text-text'}`
              }
            >
              Courses
            </NavLink>
          </nav>
        )}
      </div>
      <div className="min-w-0">{center}</div>
      <div className="flex items-center justify-end gap-[14px]">
        {/* GameStats (streak, level, XP) arrives with Phase 4. */}
        <ThemeToggle />
        {me.data && (
          <span
            className="text-caption flex h-8 w-8 items-center justify-center rounded-pill bg-surface-raised text-text-secondary"
            title={me.data.user.email}
            aria-label={`Signed in as ${me.data.user.name || me.data.user.email}`}
          >
            {initials(me.data.user.name, me.data.user.email)}
          </span>
        )}
      </div>
    </header>
  );
}
