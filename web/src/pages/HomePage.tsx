import { Link } from 'react-router';
import type { Home } from '../api/client';
import { useHome, useMe } from '../api/queries';
import { AppHeader } from '../components/AppHeader';
import { ContinueCard } from './home/ContinueCard';
import { greeting, nudge } from './home/format';
import { LeagueCard, QuestsCard, ReviewCard, TodayCard } from './home/SideCards';
import { LoadError, Loading } from './states';

const statusPill = {
  in_progress: 'bg-success-soft text-success',
  not_started: 'bg-surface-raised text-text-secondary',
  done: 'bg-success-soft text-success',
} as const;
const statusText = { in_progress: 'In progress', not_started: 'Not started', done: 'Done' } as const;

function PathCard({ p }: { p: Home['paths'][number] }) {
  const pct = p.exercise_count > 0 ? (p.passed_count / p.exercise_count) * 100 : 0;
  return (
    <li>
      <Link to={`/courses/${p.id}`} className="flex h-full flex-col gap-3 rounded-lg border border-border bg-surface p-6 hover:border-border-strong">
        <span className="flex items-start justify-between gap-3">
          <span className="text-section text-text">{p.title}</span>
          <span className={`text-caption shrink-0 rounded-pill px-2 py-0.5 ${statusPill[p.status]}`}>{statusText[p.status]}</span>
        </span>
        <span className="text-small text-text-secondary">
          {p.chapter_count} chapters · {p.passed_count} of {p.exercise_count} lessons done
        </span>
        <span
          role="img"
          aria-label={`${Math.round(pct)}% complete`}
          className="mt-auto block h-[6px] overflow-hidden rounded-pill bg-border"
        >
          <span className="block h-full rounded-pill bg-success" style={{ width: `${pct}%` }} />
        </span>
      </Link>
    </li>
  );
}

export function HomePage() {
  const home = useHome();
  const me = useMe();
  const first = (me.data?.user.name || '').split(' ')[0];
  const data = home.data;
  return (
    <>
      <AppHeader />
      <main className="mx-auto grid max-w-[1200px] gap-6 px-5 py-8 lg:grid-cols-[1fr_360px]">
        {home.isPending && <Loading label="home" />}
        {home.error && <LoadError error={home.error} />}
        {data && (
          <>
            <div className="flex min-w-0 flex-col gap-6">
              <div className="flex flex-col gap-1">
                <h1 className="text-[28px] font-semibold leading-[36px] text-text">
                  {greeting(new Date().getHours())}
                  {first ? `, ${first}` : ''}
                </h1>
                <p className="text-body text-text-secondary">{nudge(data.today)}</p>
              </div>
              {data.continue && <ContinueCard card={data.continue} />}
              <section aria-labelledby="paths-title" className="flex flex-col gap-3">
                <div className="flex items-baseline justify-between">
                  <h2 id="paths-title" className="text-section text-text">
                    Your paths
                  </h2>
                  <Link to="/courses" className="text-small text-text-muted hover:text-text">
                    Browse all courses
                  </Link>
                </div>
                <ul className="grid gap-4 sm:grid-cols-2">
                  {data.paths.map((p) => (
                    <PathCard key={p.id} p={p} />
                  ))}
                </ul>
              </section>
            </div>
            <aside className="flex flex-col gap-4" aria-label="Your day">
              <TodayCard today={data.today} />
              <QuestsCard />
              <ReviewCard />
              <LeagueCard league={data.league} />
            </aside>
          </>
        )}
      </main>
    </>
  );
}
