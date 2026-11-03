import type { Result } from '../../api/client';
import { ChestIcon } from '../../components/icons';

/** HP bar: one segment per check; each passing check is a hit (screens/boss-fight.png). */
export function BossHP({ maxHP, result }: { maxHP: number; result: Result | undefined }) {
  const hits = result ? result.checks.filter((c) => c.passed).length : 0;
  const hp = Math.max(0, maxHP - hits);
  const visible = result ? result.checks.filter((c) => c.visible).length : undefined;
  return (
    <section aria-label="Boss HP" className="flex flex-col gap-3 rounded-md border border-border bg-bg p-4">
      <p className="flex items-baseline justify-between">
        <span className="text-caption uppercase tracking-[0.08em] text-text-secondary">Boss HP</span>
        <span className="text-game text-boss-hp">
          {hp} / {maxHP}
        </span>
      </p>
      <span role="img" aria-label={`${hp} of ${maxHP} HP left`} className="flex gap-1">
        {Array.from({ length: maxHP }, (_, i) => (
          <span key={i} className={`h-3 flex-1 rounded-sm ${i < maxHP - hp ? 'bg-border' : 'bg-boss-hp'}`} />
        ))}
      </span>
      <p className="text-small text-text-muted">
        Each passing check is a hit.
        {visible !== undefined && ` ${visible} visible, ${maxHP - visible} hidden.`}
      </p>
    </section>
  );
}

/** Checks as HIT / BLOCK rows. */
export function BossResults({ result }: { result: Result }) {
  const hits = result.checks.filter((c) => c.passed).length;
  const blocks = result.checks.length - hits;
  return (
    <section aria-label="Attack result" className="flex flex-col gap-2 rounded-lg border border-border bg-surface px-5 py-4">
      <h2 className="text-section text-text">
        {blocks === 0 ? 'Boss defeated' : `${hits} hit${hits === 1 ? '' : 's'} · boss blocked ${blocks}`}
      </h2>
      <ul className="flex flex-col gap-1.5">
        {result.checks.map((c, i) => (
          <li key={i} className="text-code-sm flex gap-4">
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
