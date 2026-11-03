import type { ReactNode } from 'react';
import { Link } from 'react-router';
import type { Home, LeagueSummary } from '../../api/client';
import { useOpenChest, useQuests } from '../../api/queries';
import { Button } from '../../components/Button';
import { powerupName } from '../../components/powerups';
import { CheckIcon, ChestIcon } from '../../components/icons';
import { timeLeft, tierAbove, tierBelow, tierTitle } from './format';

function Card({ title, aside, children, label }: { title: string; aside?: ReactNode; children: ReactNode; label?: string }) {
  return (
    <section aria-label={label ?? title} className="flex flex-col gap-4 rounded-lg border border-border bg-surface p-6">
      <div className="flex items-baseline justify-between gap-2">
        <h2 className="text-section text-text">{title}</h2>
        {aside && <span className="text-small text-text-muted">{aside}</span>}
      </div>
      {children}
    </section>
  );
}

function Ring({ done, goal }: { done: number; goal: number }) {
  const r = 34;
  const c = 2 * Math.PI * r;
  const pct = Math.min(1, goal > 0 ? done / goal : 0);
  return (
    <svg width="84" height="84" viewBox="0 0 84 84" aria-hidden="true" className="-rotate-90">
      <circle cx="42" cy="42" r={r} fill="none" strokeWidth="8" className="stroke-border" />
      <circle cx="42" cy="42" r={r} fill="none" strokeWidth="8" strokeLinecap="round" className="stroke-accent"
        strokeDasharray={c} strokeDashoffset={c * (1 - pct)} />
    </svg>
  );
}

const dayLetters = ['M', 'T', 'W', 'T', 'F', 'S', 'S'];

export function TodayCard({ today }: { today: Home['today'] }) {
  return (
    <Card title="Today">
      <div className="flex items-center gap-5">
        <Ring done={today.done} goal={today.goal} />
        <div>
          <p className="text-game text-[28px] leading-[32px] text-text">
            {Math.min(today.done, today.goal)} / {today.goal}
          </p>
          <p className="text-small text-text-secondary">lessons for your daily goal</p>
        </div>
      </div>
      <p className="text-small text-text-muted">{today.streak_days}-day streak</p>
      <ol className="grid grid-cols-7 gap-2" aria-label="This week">
        {today.week.map((d, i) => (
          <li key={d.date} className="flex flex-col items-center gap-1" aria-label={`${d.date}${d.active ? ', practised' : ''}${d.today ? ', today' : ''}`}>
            <span
              className={`h-7 w-7 rounded-sm ${d.active ? 'bg-accent' : d.today ? 'border-2 border-dashed border-accent' : 'bg-surface-raised'}`}
            />
            <span className={`text-caption ${d.today ? 'font-semibold text-text' : 'text-text-muted'}`}>{dayLetters[i]}</span>
          </li>
        ))}
      </ol>
    </Card>
  );
}

const soon = <span className="text-caption rounded-pill bg-surface-raised px-2 py-0.5 text-text-muted">Coming soon</span>;

/** Today's three quests and the daily chest they open (screens/home.png). */
export function QuestsCard() {
  const { data } = useQuests();
  const open = useOpenChest();
  if (!data) return null;
  const done = data.daily.filter((q) => q.done).length;
  const chest = data.chest;
  return (
    <Card title="Daily quests" aside={`resets in ${timeLeft(data.resets_at)}`}>
      <ul className="flex flex-col gap-3">
        {[...data.daily, data.weekly].map((q) => (
          <li key={q.id} className="text-body flex items-center gap-3" aria-label={`${q.name}: ${q.progress} of ${q.target}${q.done ? ', done' : ''}`}>
            <span
              aria-hidden="true"
              className={`flex h-5 w-5 shrink-0 items-center justify-center rounded-sm ${q.done ? 'bg-success text-bg' : 'border border-border-strong'}`}
            >
              {q.done && <CheckIcon size={12} />}
            </span>
            <span className={`flex-1 ${q.done ? 'text-text-muted line-through' : 'text-text'}`}>
              {q.name}
              {q.cadence === 'weekly' && <span className="text-caption ml-2 text-text-muted">weekly · +{q.bytes}</span>}
            </span>
            <span className="text-game-sm text-text-secondary">
              {q.progress}/{q.target}
            </span>
          </li>
        ))}
      </ul>
      <div className="flex items-center gap-3 rounded-md border border-accent-line bg-accent-soft p-4">
        <ChestIcon size={22} className="text-accent-text" />
        <div className="flex flex-1 flex-col">
          <span className="text-section text-text">Daily chest</span>
          <span className="text-small text-text-secondary">
            {chest.opened
              ? `Opened: +${chest.opened.bytes} Bytes and ${powerupName(chest.opened.powerup)}`
              : chest.ready
                ? 'Ready to open'
                : `${3 - done} quest${3 - done === 1 ? '' : 's'} left · Bytes and a random power-up`}
          </span>
        </div>
        {chest.ready && !chest.opened ? (
          <Button variant="primary" disabled={open.isPending} onClick={() => open.mutate()}>
            Open
          </Button>
        ) : (
          <span className="text-game-sm text-accent-text">{done}/3</span>
        )}
      </div>
      <p className="text-caption text-text-muted">
        Chest odds: {chest.odds.map((o) => `${powerupName(o.powerup)} ${o.weight}%`).join(' · ')}
      </p>
    </Card>
  );
}

/** Placeholder until spaced review ships. */
export function ReviewCard() {
  return (
    <Card title="Review" aside={soon}>
      <p className="text-body text-text-muted">Short reviews of concepts you have learned, when they are due.</p>
      <Button disabled className="h-10 w-full">
        Start review
      </Button>
    </Card>
  );
}

export function LeagueCard({ league }: { league: LeagueSummary }) {
  const title = `${tierTitle(league.tier)} league`;
  return (
    <Card title={title} aside={`ends in ${timeLeft(league.ends_at)}`}>
      {league.joined ? (
        <>
          <p className="flex items-baseline gap-2">
            <span className="text-game text-[32px] leading-[36px] text-text">#{league.rank}</span>
            <span className="text-body text-text-secondary">
              of {league.of} · {league.weekly_xp} XP this week
            </span>
          </p>
          <p className={`text-small ${league.zone === 'demote' ? 'text-danger' : 'text-success'}`}>
            {league.zone === 'promote' && league.tier !== 'diamond' && `In the top ${league.promote}: on track for ${tierTitle(tierAbove(league.tier))}`}
            {league.zone === 'promote' && league.tier === 'diamond' && 'Top of the top tier'}
            {league.zone === 'stay' && league.promote > 0 && `Top ${league.promote} promote to ${tierTitle(tierAbove(league.tier))} · ${league.xp_to_promote} XP to go`}
            {league.zone === 'stay' && league.promote === 0 && 'Holding your place'}
            {league.zone === 'demote' && `Bottom ${league.demote} drop to ${tierTitle(tierBelow(league.tier))} · pass a lesson to climb`}
          </p>
        </>
      ) : (
        <p className="text-body text-text-secondary">Earn XP this week to join a league group.</p>
      )}
      <Link to="/league" className="text-small text-accent-text underline">
        See the standings
      </Link>
    </Card>
  );
}
