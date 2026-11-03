import { useState } from 'react';
import type { Shop, ShopItem } from '../api/client';
import { useBuy, useShop, useUsePowerup } from '../api/queries';
import { AppHeader } from '../components/AppHeader';
import { Button } from '../components/Button';
import { Bytes, PowerUpIcon } from '../components/powerups';
import { LoadError, Loading } from './states';

const stateLabel: Record<string, string> = {
  max_held: 'Max held',
  streak_safe: 'Streak is safe',
  coming_soon: 'Coming soon',
  pro: 'Pro',
  bet_open: 'Bet placed',
};

function Card({ item, balance, onError }: { item: ShopItem; balance: number; onError: (m: string) => void }) {
  const buy = useBuy();
  const use = useUsePowerup();
  const [stake, setStake] = useState(200);
  const bet = item.id === 'double_or_nothing';
  const busy = buy.isPending || use.isPending;
  const fail = (e: Error) => onError(e.message);
  let action;
  if (item.state === 'buy') {
    action = (
      <Button disabled={busy || balance < item.cost} onClick={() => buy.mutate(item.id, { onError: fail })}>
        Buy
      </Button>
    );
  } else if (item.state === 'repair') {
    action = (
      <Button variant="primary" disabled={busy || balance < item.cost} onClick={() => use.mutate({ powerup: item.id }, { onError: fail })}>
        Repair
      </Button>
    );
  } else if (item.state === 'place_bet') {
    action = (
      <span className="flex items-center gap-2">
        <label className="sr-only" htmlFor="stake">
          Stake
        </label>
        <select
          id="stake"
          value={stake}
          onChange={(e) => setStake(Number(e.target.value))}
          className="text-control h-[34px] rounded-md border border-border bg-surface px-2 text-text"
        >
          {[100, 200, 300, 500].map((s) => (
            <option key={s} value={s}>
              {s}
            </option>
          ))}
        </select>
        <Button variant="primary" disabled={busy || balance < stake} onClick={() => use.mutate({ powerup: item.id, stake }, { onError: fail })}>
          Place bet
        </Button>
      </span>
    );
  } else {
    action = <Button disabled>{stateLabel[item.state] ?? item.state}</Button>;
  }
  return (
    <li
      className={`flex flex-col gap-3 rounded-lg border bg-surface p-5 ${bet ? 'border-dashed border-accent-line' : 'border-border'}`}
      aria-label={item.name}
    >
      <div className="flex items-start justify-between">
        <span className={`flex h-10 w-10 items-center justify-center rounded-md ${bet ? 'bg-accent-soft text-accent-text' : 'bg-surface-raised text-accent'}`}>
          <PowerUpIcon id={item.id} />
        </span>
        {item.state === 'pro' ? (
          <span className="text-caption rounded-pill bg-accent-soft px-2 py-0.5 text-accent-text">Pro</span>
        ) : bet ? (
          <span className="text-game-sm text-accent-text">×2 payout</span>
        ) : (
          <span className="text-game-sm text-text-muted">Owned {item.owned}</span>
        )}
      </div>
      <h2 className="text-section text-text">{item.name}</h2>
      <p className="text-body text-text-secondary">{item.description}</p>
      <p className="text-small text-text-muted">{item.note}</p>
      <div className="mt-auto flex items-center justify-between gap-2 pt-2">
        {item.cost > 0 ? (
          <Bytes n={item.cost} className="text-text-secondary" />
        ) : (
          <span className="text-small text-text-muted">{bet ? 'Bet 100–500' : 'Included'}</span>
        )}
        {action}
      </div>
    </li>
  );
}

function Inventory({ shop }: { shop: Shop }) {
  const use = useUsePowerup();
  const held = shop.items.filter((i) => i.owned > 0 && i.use !== '');
  return (
    <section aria-label="Inventory" className="flex flex-col gap-3 rounded-lg border border-border bg-surface p-6">
      <h2 className="text-section text-text">Inventory</h2>
      {held.length === 0 && <p className="text-small text-text-muted">Nothing held yet.</p>}
      <ul className="flex flex-col">
        {held.map((i) => (
          <li key={i.id} className="flex items-center gap-3 border-t border-border py-2.5 first:border-t-0">
            <span className="text-body flex-1 text-text">{i.name}</span>
            {i.use === 'now' ? (
              <Button
                className="border border-accent-line text-accent-text"
                disabled={!i.usable_now || use.isPending}
                onClick={() => use.mutate({ powerup: i.id })}
              >
                {i.usable_now ? 'Use now' : 'Used'}
              </Button>
            ) : (
              <span className="text-small text-text-muted">{i.use === 'auto' ? 'used automatically' : 'in lessons'}</span>
            )}
            <span className="text-game-sm w-8 text-right text-text">×{i.owned}</span>
          </li>
        ))}
      </ul>
      {shop.bet && (
        <p className="text-small text-text-secondary">
          Double or Nothing: {shop.bet.days_met} of 7 days on goal ({shop.bet.stake} staked)
        </p>
      )}
    </section>
  );
}

export function ShopPage() {
  const { data, error, isPending } = useShop();
  const [problem, setProblem] = useState<string>();
  return (
    <>
      <AppHeader />
      <main className="mx-auto grid max-w-[1200px] gap-6 px-5 py-8 lg:grid-cols-[1fr_360px]">
        {isPending && <Loading label="power-ups" />}
        {error && <LoadError error={error} />}
        {data && (
          <>
            <div className="flex flex-col gap-5">
              <div className="flex flex-col gap-1">
                <h1 className="text-[28px] font-semibold leading-[36px] text-text">Power-ups</h1>
                <p className="text-body text-text-secondary">
                  Bytes are earned, never bought. Power-ups change your rewards and safety nets, never whether your code passes.
                </p>
              </div>
              {problem && (
                <p role="alert" className="text-small rounded-md bg-danger-soft px-4 py-2 text-danger">
                  {problem}
                </p>
              )}
              <ul className="grid gap-4 sm:grid-cols-2 xl:grid-cols-3">
                {data.items.map((i) => (
                  <Card key={i.id} item={i} balance={data.balance} onError={setProblem} />
                ))}
              </ul>
            </div>
            <aside className="flex flex-col gap-4">
              <section aria-label="Balance" className="flex flex-col gap-2 rounded-lg border border-border bg-surface p-6">
                <h2 className="text-section text-text-secondary">Balance</h2>
                <p className="flex items-baseline gap-2">
                  <Bytes n={data.balance} className="text-[36px] leading-[40px] text-text" />
                  <span className="text-body text-text-muted">Bytes</span>
                </p>
                {data.earned_today > 0 && <p className="text-small text-success">+{data.earned_today} today</p>}
              </section>
              <Inventory shop={data} />
              <section aria-label="Earn Bytes" className="flex flex-col gap-2 rounded-lg border border-border bg-surface p-6">
                <h2 className="text-section text-text">Earn Bytes</h2>
                {data.earn.map((e) => (
                  <p key={e.label} className="text-body flex justify-between text-text-secondary">
                    {e.label}
                    <span className="text-game-sm text-text">{e.amount}</span>
                  </p>
                ))}
              </section>
            </aside>
          </>
        )}
      </main>
    </>
  );
}
