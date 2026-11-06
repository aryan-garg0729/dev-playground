import { useEffect, useState } from 'react';
import type { Result } from '../../api/client';
import { ChestIcon } from '../../components/icons';
import { play } from '../../effects/sound';
import { prefersReducedMotion } from '../../effects/useReducedMotion';

/** Time between revealed checks in an attack (ms). */
const STEP_MS = 150;

/** How long an attack's reveal takes, so the page can hold the celebration until it ends. */
export function attackDuration(result: Result): number {
  return prefersReducedMotion() ? 0 : result.checks.length * STEP_MS + 450;
}

/**
 * Plays a fresh attack: reveals its checks one by one (HIT / BLOCK sounds),
 * then the defeat fanfare when every check passed. Returns how many checks
 * are revealed; everything at once for a result that is not fresh (e.g. the
 * last submission on page load) or under reduced motion.
 */
export function useAttack(result: Result | undefined, fresh: boolean): number {
  const total = result?.checks.length ?? 0;
  const [revealed, setRevealed] = useState(total);
  useEffect(() => {
    if (!result || !fresh) {
      setRevealed(result?.checks.length ?? 0);
      return;
    }
    const defeated = result.checks.every((c) => c.passed);
    if (prefersReducedMotion()) {
      setRevealed(result.checks.length);
      if (defeated) play('bossDefeat');
      return;
    }
    setRevealed(0);
    let i = 0;
    let fanfare: ReturnType<typeof setTimeout> | undefined;
    const id = setInterval(() => {
      const c = result.checks[i];
      if (c) play(c.passed ? 'bossHit' : 'bossBlock');
      i += 1;
      setRevealed(i);
      if (i >= result.checks.length) {
        clearInterval(id);
        if (defeated) fanfare = setTimeout(() => play('bossDefeat'), STEP_MS);
      }
    }, STEP_MS);
    return () => {
      clearInterval(id);
      clearTimeout(fanfare);
    };
  }, [result, fresh]);
  return Math.min(revealed, total);
}

/** HP bar: one segment per check; each passing check is a hit (screens/boss-fight.png). */
export function BossHP({ maxHP, result, revealed }: { maxHP: number; result: Result | undefined; revealed?: number }) {
  const shown = result ? result.checks.slice(0, revealed ?? result.checks.length) : [];
  const hits = shown.filter((c) => c.passed).length;
  const hp = Math.max(0, maxHP - hits);
  const visible = result ? result.checks.filter((c) => c.visible).length : undefined;
  const animating = revealed !== undefined && result !== undefined && revealed < result.checks.length;
  return (
    <section aria-label="Boss HP" className="flex flex-col gap-3 rounded-md border border-border bg-bg p-4">
      <p className="flex items-baseline justify-between">
        <span className="text-caption uppercase tracking-[0.08em] text-text-secondary">Boss HP</span>
        <span className="text-game text-boss-hp">
          {hp} / {maxHP}
        </span>
      </p>
      <span role="img" aria-label={`${hp} of ${maxHP} HP left`} className="flex gap-1">
        {Array.from({ length: maxHP }, (_, i) => {
          const lost = i < maxHP - hp;
          // The segment lost by the latest revealed hit flashes while the attack plays.
          const fresh = animating && lost && i === maxHP - hp - 1;
          return <span key={i} className={`h-3 flex-1 rounded-sm ${lost ? (fresh ? 'cq-hit' : 'bg-border') : 'bg-boss-hp'}`} />;
        })}
      </span>
      <p className="text-small text-text-muted">
        Each passing check is a hit.
        {visible !== undefined && ` ${visible} visible, ${maxHP - visible} hidden.`}
      </p>
    </section>
  );
}

/** Checks as HIT / BLOCK rows, revealed one by one during a fresh attack. */
export function BossResults({ result, revealed }: { result: Result; revealed?: number }) {
  const n = revealed ?? result.checks.length;
  const done = n >= result.checks.length;
  const hits = result.checks.filter((c) => c.passed).length;
  const blocks = result.checks.length - hits;
  return (
    <section aria-label="Attack result" className="flex flex-col gap-2 rounded-lg border border-border bg-surface px-5 py-4">
      <h2 className={`text-section text-text ${done && blocks === 0 && revealed !== undefined ? 'cq-pop' : ''}`}>
        {!done ? 'Attacking…' : blocks === 0 ? 'Boss defeated' : `${hits} hit${hits === 1 ? '' : 's'} · boss blocked ${blocks}`}
      </h2>
      <ul className="flex flex-col gap-1.5">
        {result.checks.slice(0, n).map((c, i) => (
          <li key={i} className={`text-code-sm flex gap-4 ${revealed !== undefined ? 'cq-row-in' : ''}`}>
            <span className={`text-game-sm w-12 shrink-0 ${c.passed ? 'text-success' : 'text-danger'}`}>{c.passed ? 'HIT' : 'BLOCK'}</span>
            <span className="text-text">
              {c.name}
              {!c.passed && c.got !== undefined && c.got !== '' && (
                <>
                  {' → '}
                  <span className="text-danger">{c.got}</span>
                  <span className="text-text-muted"> (expected {c.expected})</span>
                </>
              )}
              {!c.passed && c.feedback && <span className="text-danger"> · {c.feedback}</span>}
            </span>
          </li>
        ))}
      </ul>
    </section>
  );
}

export function DefeatReward({ reward }: { reward: string }) {
  return (
    <div className="flex items-center gap-3 rounded-md border border-dashed border-accent-line p-4">
      <ChestIcon size={22} className="text-accent-text" />
      <div className="flex flex-1 flex-col">
        <span className="text-section text-text">Defeat reward</span>
        <span className="text-small text-text-muted">{reward}</span>
      </div>
      <span className="text-game text-accent-text">+100 XP</span>
    </div>
  );
}
