import { useState } from 'react';
import { NavLink, useParams } from 'react-router';
import { useCourses, useMe, useMedals, useProfile } from '../api/queries';
import { AppHeader } from '../components/AppHeader';
import { Button } from '../components/Button';
import { ActivityTab } from './profile/ActivityTab';
import { EditProfile } from './profile/EditProfile';
import { MedalsTab, tierBg } from './profile/MedalsTab';
import { OverviewTab } from './profile/OverviewTab';
import { ProjectsTab } from './profile/ProjectsTab';
import { LoadError, Loading } from './states';

export { medalLine } from './profile/MedalsTab';

const tabs = [
  ['', 'Overview'],
  ['medals', 'Medals'],
  ['projects', 'Projects'],
  ['activity', 'Activity'],
] as const;

const joined = (iso: string) => new Date(iso).toLocaleDateString('en-US', { month: 'short', year: 'numeric' });

export function ProfilePage() {
  const { tab = '' } = useParams();
  const me = useMe();
  const medals = useMedals();
  const profile = useProfile();
  const courses = useCourses();
  const [editing, setEditing] = useState(false);
  const user = me.data?.user;
  const stats = me.data?.stats;
  const m = medals.data;
  const p = profile.data;
  const icons = Object.fromEntries((courses.data?.courses ?? []).map((c) => [c.id, c.icon]));
  const initials = (user?.name || user?.email || '?').split(/\s+/).map((w) => w[0]).join('').slice(0, 2).toUpperCase();
  // The worn title takes the color of the medal tier that grants it.
  const titleTier = m?.medals.flatMap((x) => x.tiers).find((t) => t.title === stats?.title)?.tier;
  const error = me.error ?? medals.error ?? profile.error;
  const known = tabs.some(([id]) => id === tab);
  return (
    <>
      <AppHeader />
      <main className="mx-auto flex max-w-[1200px] flex-col gap-6 px-5 py-8">
        {error && <LoadError error={error} />}
        {!error && (!m || !user || !stats || !p) && <Loading label="profile" />}
        {m && user && stats && p && (
          <>
            <section className="flex flex-wrap items-center gap-5">
              <span className="text-[30px] flex h-20 w-20 items-center justify-center rounded-lg border-2 border-accent bg-surface font-bold text-text">
                {initials}
              </span>
              <div className="flex min-w-0 flex-1 flex-col gap-1">
                <h1 className="flex flex-wrap items-center gap-3 text-[26px] font-semibold leading-[34px] text-text">
                  {user.name || user.email}
                  {stats.title && (
                    <span className={`text-game-sm rounded-sm px-2 py-1 text-accent-ink ${titleTier ? tierBg[titleTier] : 'bg-accent'}`}>{stats.title}</span>
                  )}
                </h1>
                <p className="text-small text-text-muted">
                  {p.path ? `${p.path.title} path · ` : ''}joined {joined(p.joined_at)}
                </p>
              </div>
              {!editing && <Button onClick={() => setEditing(true)}>Edit profile</Button>}
            </section>
            {editing && <EditProfile name={user.name} title={m.title} titles={m.titles} onDone={() => setEditing(false)} />}
            <dl className="grid grid-cols-2 gap-6 md:grid-cols-4">
              {[
                ['Level', String(stats.level), false],
                ['Total XP', stats.total_xp.toLocaleString('en-US'), false],
                ['Current streak (days)', String(stats.streak_days), true],
                ['Projects shipped', String(p.projects_shipped), false],
              ].map(([k, v, accent]) => (
                <div key={String(k)} className="flex flex-col gap-1">
                  <dt className="text-small order-last text-text-muted">{k}</dt>
                  <dd className={`text-game text-[28px] leading-[34px] ${accent ? 'text-accent-text' : 'text-text'}`}>{v}</dd>
                </div>
              ))}
            </dl>
            <nav aria-label="Profile" className="flex gap-2 overflow-x-auto border-b border-border">
              {tabs.map(([id, name]) => (
                <NavLink
                  key={id}
                  to={id ? `/profile/${id}` : '/profile'}
                  end
                  className={({ isActive }) =>
                    `text-body whitespace-nowrap border-b-2 px-3 pb-3 ${isActive ? 'border-accent text-text' : 'border-transparent text-text-muted hover:text-text'}`
                  }
                >
                  {name}
                </NavLink>
              ))}
            </nav>
            {tab === '' && <OverviewTab profile={p} medals={m} icons={icons} />}
            {tab === 'medals' && <MedalsTab data={m} />}
            {tab === 'projects' && <ProjectsTab profile={p} icons={icons} />}
            {tab === 'activity' && <ActivityTab profile={p} />}
            {!known && <LoadError error={new Error('Page not found')} />}
          </>
        )}
      </main>
    </>
  );
}
