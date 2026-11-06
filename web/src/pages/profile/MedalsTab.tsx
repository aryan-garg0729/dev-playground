import type { Medal, Medals } from '../../api/client';
import { MedalIcon } from '../../components/icons';

export const tierText: Record<string, string> = {
  bronze: 'text-medal-bronze',
  silver: 'text-medal-silver',
  gold: 'text-medal-gold',
  diamond: 'text-medal-diamond',
};
export const tierBg: Record<string, string> = {
  bronze: 'bg-medal-bronze',
  silver: 'bg-medal-silver',
  gold: 'bg-medal-gold',
  diamond: 'bg-medal-diamond',
};
export const cap = (s: string) => s.charAt(0).toUpperCase() + s.slice(1);

const leagues = ['', 'Bronze', 'Silver', 'Gold', 'Diamond'];
const singular: Record<string, string> = {
  days: 'day', lessons: 'lesson', chapters: 'chapter', courses: 'course', bosses: 'boss', saves: 'save',
  'late lessons': 'late lesson', 'fixed attempts': 'fixed attempt', 'bets won': 'bet won',
  'top-3 weeks': 'top-3 week', 'top-5 weeks': 'top-5 week', 'first-try passes': 'first-try pass',
};

/** "30 days", "1 boss", "Reach Gold league". */
export function goal(n: number, unit: string): string {
  if (unit === 'tier') return `Reach ${leagues[n] ?? n} league`;
  return `${n} ${n === 1 ? (singular[unit] ?? unit) : unit}`;
}

/** "Bronze · Silver at 30 days", "10 lessons before 8 am", "Gold · all tiers earned". */
export function medalLine(m: Medal): string {
  if (m.coming_soon) return 'Coming soon';
  if (m.name === '???') return 'Hidden until earned';
  const target = goal(m.next_gte, m.unit);
  const next = m.next_tier ? (m.unit === 'tier' ? `${cap(m.next_tier)}: ${target.replace("Reach", "reach")}` : `${cap(m.next_tier)} at ${target}`) : 'all tiers earned';
  return m.tier ? `${cap(m.tier)} · ${next}` : target;
}

function MedalCard({ m }: { m: Medal }) {
  const earned = m.tier !== '';
  return (
    <li
      aria-label={`${m.name}${earned ? `, ${m.tier}` : ', not earned'}`}
      className={`flex items-center gap-4 rounded-lg border border-border bg-surface p-4 ${earned ? '' : 'opacity-70'}`}
    >
      <MedalIcon size={40} className={earned ? tierText[m.tier] : 'text-border-strong'} />
      <div className="flex min-w-0 flex-col gap-1">
        <span className={`text-section ${earned ? 'text-text' : 'text-text-secondary'}`}>{m.name}</span>
        <span className="text-small truncate text-text-muted">{medalLine(m)}</span>
        <span className="flex gap-1" aria-hidden="true">
          {(m.tiers.length ? m.tiers : [{ tier: '', earned_at: null }]).map((t, i) => (
            <span key={i} className={`h-1 w-5 rounded-pill ${t.earned_at ? tierBg[t.tier] : 'bg-border'}`} />
          ))}
        </span>
      </div>
    </li>
  );
}

function Closest({ medals }: { medals: Medal[] }) {
  const near = medals
    .filter((m) => !m.coming_soon && m.next_tier && m.next_gte > 0 && m.name !== '???')
    .map((m) => ({ m, pct: Math.min(1, m.value / m.next_gte) }))
    .sort((a, b) => b.pct - a.pct)
    .slice(0, 3);
  if (near.length === 0) return null;
  return (
    <section aria-labelledby="closest" className="flex flex-col gap-3">
      <h2 id="closest" className="text-section text-text">
        Closest to unlocking
      </h2>
      <ul className="grid gap-4 md:grid-cols-3">
        {near.map(({ m, pct }) => (
          <li key={m.id} className="flex flex-col gap-3 rounded-lg border border-border bg-surface p-5">
            <span className="flex justify-between">
              <span className="text-section text-text">{m.name}</span>
              <span className={`text-game-sm ${tierText[m.next_tier]}`}>{cap(m.next_tier)}</span>
            </span>
            <span role="img" aria-label={`${m.value} of ${m.next_gte}`} className="block h-[6px] overflow-hidden rounded-pill bg-border">
              <span className={`block h-full rounded-pill ${tierBg[m.next_tier]}`} style={{ width: `${pct * 100}%` }} />
            </span>
            <span className="text-small text-text-muted">
              {m.unit === 'tier' ? goal(m.next_gte, m.unit) : `${Math.min(m.value, m.next_gte)} of ${goal(m.next_gte, m.unit)}`}
            </span>
          </li>
        ))}
      </ul>
    </section>
  );
}

/** The Medals tab: closest to unlocking, then every family. */
export function MedalsTab({ data }: { data: Medals }) {
  return (
    <>
      <Closest medals={data.medals} />
      {data.families.map((f) => {
        const list = data.medals.filter((m) => m.family === f.id);
        if (list.length === 0) return null;
        const earned = list.filter((m) => m.tier).length;
        return (
          <section key={f.id} aria-labelledby={`fam-${f.id}`} className="flex flex-col gap-3">
            <h2 id={`fam-${f.id}`} className="text-section flex items-baseline gap-2 text-text">
              {f.name} <span className="text-small font-normal text-text-muted">{earned} of {list.length}</span>
            </h2>
            <ul className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
              {list.map((m) => (
                <MedalCard key={m.id} m={m} />
              ))}
            </ul>
          </section>
        );
      })}
    </>
  );
}
