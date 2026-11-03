import type { ComponentType } from 'react';
import { BoltIcon, BulbIcon, CalendarIcon, DiceIcon, EyeIcon, FlameOffIcon, HexIcon, RotateIcon, ShieldIcon, SnowflakeIcon } from './icons';

type IconComp = ComponentType<{ size?: number; className?: string }>;

const icons: Record<string, IconComp> = {
  streak_freeze: SnowflakeIcon,
  focus_boost: BoltIcon,
  second_wind: RotateIcon,
  extra_hint: BulbIcon,
  peek: EyeIcon,
  league_shield: ShieldIcon,
  streak_repair: FlameOffIcon,
  holiday_mode: CalendarIcon,
  double_or_nothing: DiceIcon,
};

const names: Record<string, string> = {
  streak_freeze: 'Streak Freeze',
  focus_boost: 'Focus Boost',
  second_wind: 'Second Wind',
  extra_hint: 'Extra Hint',
  peek: 'Peek',
  league_shield: 'League Shield',
  streak_repair: 'Streak Repair',
  holiday_mode: 'Holiday Mode',
  double_or_nothing: 'Double or Nothing',
};

export const powerupName = (id: string) => names[id] ?? id;

export function PowerUpIcon({ id, size = 20, className = '' }: { id: string; size?: number; className?: string }) {
  const I = icons[id] ?? HexIcon;
  return <I size={size} className={className} />;
}

/** "⬡ 1,240": a Byte amount in game type. */
export function Bytes({ n, className = '' }: { n: number; className?: string }) {
  return (
    <span className={`text-game inline-flex items-center gap-1.5 ${className}`}>
      <HexIcon size={14} />
      {n.toLocaleString('en-US')}
    </span>
  );
}
