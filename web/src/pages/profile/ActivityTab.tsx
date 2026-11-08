import { Link } from 'react-router';
import type { Profile } from '../../api/client';
import { YearCard } from './OverviewTab';

const dayLabel = (iso: string) =>
  new Date(`${iso}T00:00:00Z`).toLocaleDateString('en-US', { weekday: 'short', month: 'short', day: 'numeric', timeZone: 'UTC' });

/** The year heatmap and the latest passed lessons, by day. */
export function ActivityTab({ profile }: { profile: Profile }) {
  const days = [...new Set(profile.recent.map((r) => r.date))];
  return (
    <>
      <YearCard profile={profile} />
      <section aria-labelledby="recent-title" className="flex flex-col gap-4 rounded-lg border border-border bg-surface p-6">
        <h2 id="recent-title" className="text-section text-text">
          Recently passed
        </h2>
        {days.length === 0 && <p className="text-small text-text-muted">Nothing yet. Pass a lesson and it shows up here.</p>}
        {days.map((d) => (
          <div key={d} className="flex flex-col gap-2">
            <h3 className="text-caption uppercase tracking-[0.08em] text-text-muted">{dayLabel(d)}</h3>
            <ul className="flex flex-col gap-1">
              {profile.recent
                .filter((r) => r.date === d)
                .map((r) => (
                  <li key={r.exercise_id} className="text-small flex justify-between gap-3">
                    <Link to={`/learn/${r.course_id}/${r.exercise_id}`} className="text-text hover:underline">
                      {r.title}
                      {r.boss && <span className="text-caption ml-2 rounded-pill bg-danger-soft px-2 py-0.5 text-danger">Boss</span>}
                    </Link>
                    <span className="shrink-0 text-text-muted">{r.course_title}</span>
                  </li>
                ))}
            </ul>
          </div>
        ))}
      </section>
    </>
  );
}
