import { type ReactNode, useEffect, useRef, useState } from 'react';
import { Link, NavLink, useNavigate } from 'react-router';
import type { Me } from '../api/client';
import { useMe, useSignOut } from '../api/queries';
import { GameStats } from './GameStats';
import { Bytes } from './powerups';
import { CodeIcon } from './icons';
import { SoundToggle } from './SoundToggle';
import { ThemeToggle } from './ThemeToggle';

function initials(name: string, email: string): string {
  const src = name.trim() || email;
  const parts = src.split(/[\s@._-]+/).filter(Boolean);
  return ((parts[0]?.[0] ?? '') + (parts[1]?.[0] ?? '')).toUpperCase() || '?';
}

const navItems = [
  { to: '/', label: 'Home', end: true },
  { to: '/courses', label: 'Courses', end: false },
  { to: '/league', label: 'League', end: false },
];

function AccountMenu({ me }: { me: Me }) {
  const [open, setOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);
  const signOut = useSignOut();
  const navigate = useNavigate();
  useEffect(() => {
    if (!open) return;
    const close = (e: MouseEvent | KeyboardEvent) => {
      if (e instanceof KeyboardEvent ? e.key === 'Escape' : !ref.current?.contains(e.target as Node)) setOpen(false);
    };
    document.addEventListener('mousedown', close);
    document.addEventListener('keydown', close);
    return () => {
      document.removeEventListener('mousedown', close);
      document.removeEventListener('keydown', close);
    };
  }, [open]);
  const who = me.user.name || me.user.email;
  return (
    <div ref={ref} className="relative">
      <button
        type="button"
        aria-haspopup="menu"
        aria-expanded={open}
        aria-label={`Account: ${who}`}
        onClick={() => setOpen((o) => !o)}
        className="text-caption flex h-8 w-8 items-center justify-center rounded-pill bg-surface-raised text-text-secondary"
      >
        {initials(me.user.name, me.user.email)}
      </button>
      {open && (
        <div role="menu" className="absolute right-0 top-10 z-20 w-56 rounded-md border border-border bg-surface p-2">
          <p className="text-caption truncate px-2 py-1 text-text-muted">{me.user.email}</p>
          <Link role="menuitem" to="/profile" className="text-small block rounded-md px-2 py-1.5 text-text hover:bg-surface-raised">
            Profile and medals
          </Link>
          <Link role="menuitem" to="/shop" className="text-small block rounded-md px-2 py-1.5 text-text hover:bg-surface-raised">
            Power-ups
          </Link>
          <button
            type="button"
            role="menuitem"
            className="text-small w-full rounded-md px-2 py-1.5 text-left text-text hover:bg-surface-raised"
            onClick={() => signOut.mutate(undefined, { onSuccess: () => navigate('/sign-in') })}
          >
            Sign out
          </button>
        </div>
      )}
    </div>
  );
}

/** 56px app header. `center` holds page context (e.g. lesson name and prev/next) and replaces the nav. */
export function AppHeader({ center }: { center?: ReactNode }) {
  const me = useMe();
  const s = me.data?.stats;
  return (
    <header className="sticky top-0 z-10 grid h-[var(--header-height)] grid-cols-[1fr_auto_1fr] items-center border-b border-border bg-surface px-5">
      <div className="flex min-w-0 items-center gap-6">
        <Link to="/" className="flex items-center gap-2" aria-label="CodeQuest home">
          <span className="flex h-7 w-7 items-center justify-center rounded-md bg-accent text-accent-ink">
            <CodeIcon />
          </span>
          <span className="text-brand text-text max-sm:hidden">CodeQuest</span>
        </Link>
        {!center && (
          <nav aria-label="Main" className="flex gap-1 max-md:hidden">
            {navItems.map((n) => (
              <NavLink
                key={n.to}
                to={n.to}
                end={n.end}
                className={({ isActive }) =>
                  `text-small rounded-md px-3 py-1.5 ${isActive ? 'bg-surface-raised text-text' : 'text-text-muted hover:text-text'}`
                }
              >
                {n.label}
              </NavLink>
            ))}
          </nav>
        )}
      </div>
      <div className="min-w-0">{center}</div>
      <div className="flex items-center justify-end gap-[14px]">
        {s && (
          <Link
            to="/shop"
            aria-label={`${s.bytes} Bytes: open the power-up shop`}
            className="flex h-8 items-center rounded-md border border-border px-2.5 text-text hover:bg-surface-raised max-sm:hidden"
          >
            <Bytes n={s.bytes} className="text-game-sm" />
          </Link>
        )}
        {s && (
          <GameStats
            streakDays={s.streak_days}
            level={s.level}
            totalXP={s.total_xp}
            levelStartXP={s.level_start_xp}
            nextLevelXP={s.next_level_xp}
          />
        )}
        <SoundToggle />
        <ThemeToggle />
        {me.data && <AccountMenu me={me.data} />}
      </div>
    </header>
  );
}
