import { Link } from 'react-router';
import type { Home } from '../../api/client';
import { Heatmap, heatmapWeeks, Legend, totals } from '../../components/Heatmap';

const WEEKS = 27; // about six months; the profile shows the full year

function Stat({ value, label, accent }: { value: number; label: string; accent?: boolean }) {
  return (
    <div className="flex flex-col">
      <dd className={`text-game order-first text-[24px] leading-[30px] ${accent ? 'text-accent-text' : 'text-text'}`}>{value}</dd>
      <dt className="text-small text-text-secondary">{label}</dt>
    </div>
  );
}

export function ActivityCard({ activity, streak }: { activity: Home['activity']; streak: number }) {
  const { lessons, days } = totals(heatmapWeeks(activity, WEEKS));
  const best = streak > 0 && streak >= activity.longest_streak;
  return (
    <section aria-labelledby="activity-title" className="flex flex-col gap-4 rounded-lg border border-border bg-surface p-6">
      <div className="flex items-baseline justify-between gap-2">
        <h2 id="activity-title" className="text-section text-text">
          Activity
        </h2>
        <Link to="/profile/activity" className="text-small text-text-muted hover:text-text">
          Full year on your profile
        </Link>
      </div>
      <div className="flex flex-wrap items-start gap-x-8 gap-y-4">
        <div className="flex min-w-0 max-w-full flex-1 flex-col items-end gap-2">
          <div className="w-full">
            <Heatmap activity={activity} weeks={WEEKS} label="Lessons passed, last six months" />
          </div>
          <Legend />
        </div>
        <dl className="flex flex-row flex-wrap gap-x-8 gap-y-4 pt-4 sm:flex-col">
          <Stat value={days} label="active days" />
          <Stat value={lessons} label="lessons passed" />
          <Stat value={streak} label={best ? 'day streak, best yet' : `day streak · best ${activity.longest_streak}`} accent />
        </dl>
      </div>
    </section>
  );
}
