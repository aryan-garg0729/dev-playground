import { FlameIcon } from './icons';

export interface GameStatsProps {
  streakDays: number;
  level: number;
  /** Cumulative XP, and where the current and next levels start. */
  totalXP: number;
  levelStartXP: number;
  nextLevelXP: number;
}

/** XP gained inside the current level, and the level's size. */
export function levelProgress(p: Pick<GameStatsProps, 'totalXP' | 'levelStartXP' | 'nextLevelXP'>) {
  const of = Math.max(1, p.nextLevelXP - p.levelStartXP);
  const into = Math.min(of, Math.max(0, p.totalXP - p.levelStartXP));
  return { into, of, pct: (into / of) * 100 };
}

/** The header's reward cluster: streak, level badge and XP bar (docs/design/components/GameStats). */
export function GameStats(p: GameStatsProps) {
  const { into, of, pct } = levelProgress(p);
  return (
    <div className="flex items-center gap-[14px]">
      <span className="flex items-center gap-1" aria-label={`${p.streakDays}-day streak`} title={`${p.streakDays}-day streak`}>
        <FlameIcon className="text-accent" />
        <span className="text-game text-accent-text" aria-hidden="true">
          {p.streakDays}
        </span>
      </span>
      <span className="text-game-sm rounded-sm bg-level-bg px-[7px] py-[2px] text-level-ink">Lv {p.level}</span>
      <span
        role="img"
        aria-label={`${into} of ${of} XP to level ${p.level + 1}`}
        className="block h-[6px] w-[80px] overflow-hidden rounded-pill bg-border max-sm:hidden"
      >
        <span className="block h-full rounded-pill bg-accent" style={{ width: `${pct}%` }} />
      </span>
    </div>
  );
}
