import { useLeague } from '../api/queries';
import { AppHeader } from '../components/AppHeader';
import { timeLeft, tierAbove, tierBelow, tierTitle } from './home/format';
import { LoadError, Loading } from './states';

const zoneRow = { promote: 'border-l-success', stay: 'border-l-transparent', demote: 'border-l-danger' } as const;

export function LeaguePage() {
  const { data, error, isPending } = useLeague();
  return (
    <>
      <AppHeader />
      <main className="mx-auto flex max-w-[720px] flex-col gap-6 px-5 py-8">
        {isPending && <Loading label="league" />}
        {error && <LoadError error={error} />}
        {data && (
          <>
            <div className="flex flex-col gap-1">
              <h1 className="text-page-title text-text">{tierTitle(data.tier)} league</h1>
              <p className="text-body text-text-secondary">
                Weekly, ends in {timeLeft(data.ends_at)} (Monday 00:00 UTC).
                {data.tier !== 'diamond' && ` Top ${data.promote || 5} move up to ${tierTitle(tierAbove(data.tier))}.`}
                {data.tier !== 'bronze' && ` Bottom ${data.demote || 5} move down to ${tierTitle(tierBelow(data.tier))}.`}
              </p>
            </div>
            {!data.joined ? (
              <p className="text-body rounded-lg border border-border bg-surface p-6 text-text-secondary">
                You are not in a group yet this week. Pass a lesson to join one with up to 30 learners.
              </p>
            ) : (
              <ol className="overflow-hidden rounded-lg border border-border bg-surface" aria-label="Standings">
                {data.standings.map((s) => (
                  <li
                    key={s.rank}
                    className={`flex items-center gap-4 border-b border-l-4 border-b-border px-5 py-3 last:border-b-0 ${zoneRow[s.zone]} ${s.you ? 'bg-surface-raised' : ''}`}
                    aria-label={`Rank ${s.rank}: ${s.you ? 'you' : s.name}, ${s.weekly_xp} XP${s.zone === 'promote' ? ', promotion zone' : s.zone === 'demote' ? ', demotion zone' : ''}`}
                  >
                    <span className="text-game w-8 text-text-muted">{s.rank}</span>
                    <span className={`text-body flex-1 truncate ${s.you ? 'font-semibold text-text' : 'text-text-secondary'}`}>
                      {s.you ? `${s.name} (you)` : s.name}
                    </span>
                    <span className="text-game-sm text-text">{s.weekly_xp} XP</span>
                  </li>
                ))}
              </ol>
            )}
          </>
        )}
      </main>
    </>
  );
}
