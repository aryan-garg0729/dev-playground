import { useEffect, useRef, useState } from 'react';
import { Link } from 'react-router';
import type { Outcome, Result } from '../api/client';
import { Button } from './Button';
import { levelProgress } from './GameStats';
import { CheckIcon, FlameIcon } from './icons';
import { powerupName } from './powerups';

interface Props {
  outcome: Outcome;
  result: Result;
  /** Where "Next" goes; omitted on the last lesson. */
  next?: { to: string; title: string };
  courseTo: string;
  onClose: () => void;
}

/** Bytes, medals, quests and boss loot from the economy, if any. */
function RewardsList({ outcome }: { outcome: Outcome }) {
  const r = outcome.rewards;
  const lines: string[] = [];
  if (outcome.second_wind_kept) lines.push('Second Wind kept your first-try bonus');
  if (r?.boss?.defeated && r.boss.chest) lines.push(`Boss chest: +${r.boss.chest.bytes} Bytes and ${powerupName(r.boss.chest.powerup)}`);
  if (r?.boss?.unlocked_chapter) lines.push('Next chapter unlocked');
  for (const m of r?.medals ?? []) lines.push(`Medal: ${m.name} (${m.tier}) · title "${m.title}"`);
  for (const q of r?.quests_done ?? []) lines.push(`Quest done: ${q}`);
  if (r?.chest_ready) lines.push('Daily chest ready to open on Home');
  for (const g of r?.bytes ?? []) if (g.reason !== 'boss' && g.reason !== 'medal') lines.push(`${g.label}: +${g.amount} Bytes`);
  if (lines.length === 0) return null;
  return (
    <ul aria-label="Rewards" className="flex flex-col gap-1.5">
      {lines.map((l) => (
        <li key={l} className="text-small text-text-secondary">
          {l}
        </li>
      ))}
    </ul>
  );
}

function Row({ label, xp }: { label: string; xp: number }) {
  return (
    <div className="flex items-center justify-between">
      <span className="text-body text-text-secondary">{label}</span>
      <span className="text-game-sm text-text">+{xp}</span>
    </div>
  );
}

/**
 * "Lesson clear" dialog after a first pass (screens/pass-celebration.png).
 * Skippable: Escape, the backdrop, or "Keep working" close it. The XP bar
 * animates only when the user allows motion.
 */
export function PassCelebration({ outcome, result, next, courseTo, onClose }: Props) {
  const primary = useRef<HTMLAnchorElement>(null);
  const [filled, setFilled] = useState(false);
  useEffect(() => {
    primary.current?.focus();
    const t = requestAnimationFrame(() => setFilled(true));
    const onKey = (e: KeyboardEvent) => e.key === 'Escape' && onClose();
    document.addEventListener('keydown', onKey);
    return () => {
      cancelAnimationFrame(t);
      document.removeEventListener('keydown', onKey);
    };
  }, [onClose]);

  const b = outcome.breakdown;
  const passed = result.checks.filter((c) => c.passed).length;
  const now = levelProgress({ totalXP: outcome.total_xp, levelStartXP: outcome.level_start_xp, nextLevelXP: outcome.next_level_xp });
  const before = outcome.level_up ? 0 : Math.max(0, now.into - outcome.xp_awarded);
  const toNext = now.of - now.into;
  const detail = [`${passed} of ${result.checks.length} checks`, result.runtime_ms > 0 ? `${result.runtime_ms} ms` : '', outcome.first_try ? 'first try' : '']
    .filter(Boolean)
    .join(' · ');

  return (
    <div className="fixed inset-0 z-30 flex items-center justify-center p-4">
      <div aria-hidden="true" data-testid="backdrop" className="absolute inset-0 bg-bg opacity-80" onMouseDown={onClose} />
      <div
        role="dialog"
        aria-modal="true"
        aria-labelledby="clear-title"
        className="relative flex w-full max-w-[420px] flex-col gap-5 rounded-lg border border-border bg-surface p-6"
      >
        <div className="flex flex-col items-center gap-3 text-center">
          <p className="text-game-sm tracking-[0.2em] text-accent-text">LESSON CLEAR</p>
          <span className="flex h-11 w-11 items-center justify-center rounded-pill bg-success-soft text-success">
            <CheckIcon size={22} />
          </span>
          <h2 id="clear-title" className="text-page-title text-text">
            All checks passed
          </h2>
          <p className="text-small text-text-muted">{detail}</p>
        </div>

        {b && (
          <div className="flex flex-col gap-2 rounded-md border border-border p-4">
            <Row label="Lesson complete" xp={b.base} />
            {b.first_try > 0 && <Row label="First-try bonus" xp={b.first_try} />}
            {b.no_hint > 0 && <Row label="No hints used" xp={b.no_hint} />}
            {b.boss > 0 && <Row label="Boss defeated" xp={b.boss} />}
            {b.boost > 0 && <Row label="Focus Boost ×2" xp={b.boost} />}
            <div className="mt-1 flex items-center justify-between border-t border-border pt-3">
              <span className="text-section text-text">Total</span>
              <span className="text-game text-[24px] leading-[28px] text-accent-text">+{outcome.xp_awarded} XP</span>
            </div>
          </div>
        )}

        <div className="flex flex-col gap-2">
          <div className="flex items-center justify-between gap-2">
            <span className="flex items-center gap-2">
              <span className="text-game-sm rounded-sm bg-level-bg px-[7px] py-[2px] text-level-ink">Lv {outcome.level}</span>
              <span className="text-small text-text-muted">
                {outcome.level_up ? `Level up! ${toNext} XP to Lv ${outcome.level + 1}` : `${toNext} XP to Lv ${outcome.level + 1}`}
              </span>
            </span>
            {outcome.streak_days > 0 && (
              <span className="text-game-sm flex items-center gap-1 text-accent-text">
                <FlameIcon size={14} className="text-accent" />
                {outcome.streak_days}-day streak
              </span>
            )}
          </div>
          <div
            role="img"
            aria-label={`${now.into} of ${now.of} XP in level ${outcome.level}`}
            className="flex h-[6px] overflow-hidden rounded-pill bg-border"
          >
            <span className="h-full bg-accent" style={{ width: `${(before / now.of) * 100}%` }} />
            <span
              className="h-full bg-accent-line motion-safe:transition-[width] motion-safe:duration-700 motion-safe:ease-out"
              style={{ width: filled ? `${((now.into - before) / now.of) * 100}%` : '0%' }}
            />
          </div>
        </div>

        <RewardsList outcome={outcome} />

        <div className="flex flex-col gap-2">
          <Link
            ref={primary}
            to={next?.to ?? courseTo}
            onClick={onClose}
            className="text-control flex h-11 items-center justify-center rounded-md bg-accent text-accent-ink"
          >
            {next ? `Next: ${next.title}` : 'Back to the course'}
          </Link>
          <Button className="h-10" onClick={onClose}>
            Keep working
          </Button>
        </div>
      </div>
    </div>
  );
}
