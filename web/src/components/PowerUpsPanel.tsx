import { useState } from 'react';
import { Link } from 'react-router';
import type { GradeResult, UseResult } from '../api/client';
import { useShop, useUsePowerup } from '../api/queries';
import { PowerUpIcon } from './powerups';

type Peeked = NonNullable<UseResult['peek']>;

const minutesLeft = (iso: string) => Math.max(1, Math.ceil((new Date(iso).getTime() - Date.now()) / 60000));

/**
 * Power-ups usable in a lesson (screens/boss-fight.png). They change rewards
 * only: Second Wind keeps the first-try bonus, Peek shows a failed hidden
 * check's input, Focus Boost doubles XP. The server checks every rule.
 */
export function PowerUpsPanel({ exerciseId, result }: { exerciseId: string; result: GradeResult | undefined }) {
  const shop = useShop();
  const use = useUsePowerup();
  const [peeks, setPeeks] = useState<Peeked[]>([]);
  const [error, setError] = useState<string>();
  if (!shop.data) return null;
  const owned = (id: string) => shop.data.items.find((i) => i.id === id)?.owned ?? 0;
  const effect = (name: string, ref = '') => shop.data.effects.find((e) => e.effect === name && e.ref_id === ref);
  const failed = result?.result.status === 'fail' || result?.result.status === 'error';
  const hiddenFailed = !!result?.result.checks.some((c) => !c.visible && !c.passed);
  const secondWind = effect('second_wind', exerciseId);
  const focus = effect('focus_boost');

  const run = (powerup: string, after?: (r: UseResult) => void) => {
    setError(undefined);
    use.mutate({ powerup, exerciseId }, { onSuccess: after, onError: (e) => setError(e.message) });
  };

  const chips = [
    {
      id: 'second_wind',
      label: 'Second Wind',
      show: owned('second_wind') > 0 || !!secondWind,
      active: secondWind ? 'active' : '',
      disabled: !!secondWind || !failed,
      title: 'After one failed submit: your next pass keeps the first-try bonus',
      onClick: () => run('second_wind'),
    },
    {
      id: 'peek',
      label: 'Peek',
      show: owned('peek') > 0,
      active: '',
      disabled: !hiddenFailed,
      title: 'Show the input of one hidden check you failed (removes the first-try bonus)',
      onClick: () => run('peek', (r) => r.peek && setPeeks((p) => [...p, r.peek!])),
    },
    {
      id: 'focus_boost',
      label: 'Focus Boost',
      show: owned('focus_boost') > 0 || !!focus,
      active: focus ? `active · ${minutesLeft(focus.expires_at)} min` : '',
      disabled: !!focus,
      title: 'Double XP for 15 minutes',
      onClick: () => run('focus_boost'),
    },
  ].filter((c) => c.show);

  return (
    <section aria-labelledby="pu-title" className="flex flex-col gap-3">
      <h2 id="pu-title" className="text-section text-text">
        Power-ups
      </h2>
      {chips.length === 0 ? (
        <p className="text-small text-text-muted">
          None in hand. <Link to="/shop" className="text-accent-text underline">Get some in the shop</Link>
        </p>
      ) : (
        <div className="flex flex-wrap gap-2">
          {chips.map((c) => (
            <button
              key={c.id}
              type="button"
              title={c.title}
              disabled={c.disabled || use.isPending}
              onClick={c.onClick}
              className={`text-control flex h-10 items-center gap-2 rounded-md border px-3 disabled:cursor-not-allowed ${
                c.active ? 'border-accent bg-accent-soft text-text' : 'border-border bg-surface-raised text-text disabled:text-text-faint'
              }`}
            >
              <PowerUpIcon id={c.id} size={16} className="text-accent" />
              {c.label}
              {c.active ? (
                <span className="text-accent-text">{c.active}</span>
              ) : (
                <span className="text-game-sm text-text-muted">×{owned(c.id)}</span>
              )}
            </button>
          ))}
        </div>
      )}
      {peeks.map((p) => (
        <p key={p.check} className="text-code-sm rounded-md border border-border bg-bg px-3 py-2 text-text">
          <span className="text-text-muted">{p.name} input: </span>
          {p.input}
        </p>
      ))}
      {error && (
        <p role="alert" className="text-small text-danger">
          {error}
        </p>
      )}
    </section>
  );
}
