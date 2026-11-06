import { Link } from 'react-router';
import { useHome, useMe } from '../api/queries';
import { AppHeader } from '../components/AppHeader';
import { ActivityCard } from './home/ActivityCard';
import { ContinueCard } from './home/ContinueCard';
import { PathCard } from './home/PathCard';
import { greeting, nudge } from './home/format';
import { LeagueCard, QuestsCard, ReviewCard, TodayCard } from './home/SideCards';
import { LoadError, Loading } from './states';

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
              <ActivityCard activity={data.activity} streak={data.today.streak_days} />
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
