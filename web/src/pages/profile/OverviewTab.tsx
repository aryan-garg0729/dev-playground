import { Link } from 'react-router';
import type { Medals, Profile } from '../../api/client';
import { Heatmap, Legend } from '../../components/Heatmap';
import { MedalIcon } from '../../components/icons';
import { cap, tierBg, tierText } from './MedalsTab';
import { ProjectRow } from './ProjectList';

const plural = (n: number, one: string) => `${n} ${n === 1 ? one : `${one}s`}`;

function Card({ title, aside, children }: { title: string; aside?: React.ReactNode; children: React.ReactNode }) {
  return (
    <section aria-label={title} className="flex flex-col gap-4 rounded-lg border border-border bg-surface p-6">
      <div className="flex items-baseline justify-between gap-2">
        <h2 className="text-section text-text">{title}</h2>
        {aside && <span className="text-small text-text-muted">{aside}</span>}
      </div>
      {children}
    </section>
  );
}

export function YearCard({ profile }: { profile: Profile }) {
  const a = profile.activity;
  const lessons = a.days.reduce((n, d) => n + d.lessons, 0);
  return (
    <Card title={`${plural(lessons, 'lesson')} in the last year`} aside={`${plural(a.days.length, 'active day')} · longest streak ${plural(a.longest_streak, 'day')}`}>
      <Heatmap activity={a} weeks={53} cell={16} label="Lessons passed, last year" />
      <div className="text-caption flex flex-wrap items-center justify-between gap-2 text-text-muted">
        <span>Each square is a day. Color shows lessons passed.</span>
        <Legend />
      </div>
    </Card>
  );
}

function ProjectsCard({ profile, icons }: { profile: Profile; icons: Record<string, string> }) {
  const done = profile.projects
    .filter((p) => p.status === 'defeated')
    .sort((a, b) => (b.defeated_at ?? '').localeCompare(a.defeated_at ?? ''))
    .slice(0, 3);
  const next = profile.projects.find((p) => p.status === 'next');
  const rows = next ? [...done, next] : done;
  return (
    <Card title="Projects completed" aside={<Link to="/profile/projects" className="hover:text-text">All projects</Link>}>
      {rows.length === 0 ? (
        <p className="text-small text-text-muted">Defeat a chapter boss to ship your first project.</p>
      ) : (
        <ul>
          {rows.map((p) => (
            <ProjectRow key={p.exercise_id} p={p} icon={icons[p.course_id] ?? ''} />
          ))}
        </ul>
      )}
    </Card>
  );
}

const skillText: Record<string, string> = { gold: 'text-accent-ink', silver: 'text-accent-ink', bronze: 'text-accent-ink' };

function SkillsCard({ skills }: { skills: NonNullable<Profile['skills']> }) {
  return (
    <Card title="Skills" aside={skills.course_title}>
      <ul className="flex flex-col gap-3">
        {skills.chapters.map((c) => (
          <li key={c.title} className="flex items-center justify-between gap-3">
            <span className={`text-small ${c.tier ? 'text-text' : 'text-text-muted'}`}>{c.title}</span>
            {c.tier ? (
              <span className={`text-game-sm rounded-sm px-1.5 py-0.5 ${tierBg[c.tier]} ${skillText[c.tier]}`}>{cap(c.tier)}</span>
            ) : (
              <span className="text-small text-text-muted">Not started</span>
            )}
          </li>
        ))}
      </ul>
      <p className="text-caption text-text-muted">Bronze once you pass a lesson, Silver when the chapter is done, Gold when every lesson passed on the first try.</p>
    </Card>
  );
}

function RecentMedals({ medals }: { medals: Medals }) {
  const earned = medals.medals
    .flatMap((m) => m.tiers.filter((t) => t.earned_at).map((t) => ({ m, t })))
    .sort((a, b) => (b.t.earned_at ?? '').localeCompare(a.t.earned_at ?? ''))
    .slice(0, 3);
  return (
    <Card title="Recent medals" aside={<Link to="/profile/medals" className="hover:text-text">All {medals.earned}</Link>}>
      {earned.length === 0 ? (
        <p className="text-small text-text-muted">No medals yet. Keep a streak or beat a boss to earn one.</p>
      ) : (
        <ul className="flex flex-wrap gap-6">
          {earned.map(({ m, t }) => (
            <li key={`${m.id}:${t.tier}`} aria-label={`${m.name}, ${t.tier}`} className="flex flex-col items-center gap-1">
              <MedalIcon size={36} className={tierText[t.tier]} />
              <span className="text-small text-text-secondary">{m.name}</span>
            </li>
          ))}
        </ul>
      )}
    </Card>
  );
}

export function OverviewTab({ profile, medals, icons }: { profile: Profile; medals: Medals; icons: Record<string, string> }) {
  return (
    <>
      <YearCard profile={profile} />
      <div className="grid items-start gap-6 lg:grid-cols-[1fr_390px]">
        <ProjectsCard profile={profile} icons={icons} />
        <div className="flex flex-col gap-6">
          {profile.skills && <SkillsCard skills={profile.skills} />}
          <RecentMedals medals={medals} />
        </div>
      </div>
    </>
  );
}
