import type { ReactNode } from 'react';

type Variant = 'difficulty' | 'topic' | 'xp';

const variants: Record<Variant, string> = {
  difficulty: 'text-caption bg-success-soft text-success',
  topic: 'text-caption bg-surface-raised text-text-secondary',
  xp: 'text-game-sm bg-accent-soft text-accent-text',
};

export function Pill({ variant, children }: { variant: Variant; children: ReactNode }) {
  return <span className={`inline-flex items-center rounded-pill px-[9px] py-[3px] ${variants[variant]}`}>{children}</span>;
}

/**
 * Difficulty pill. Only Easy (1-2) has tokens so far; per the Pill README, Medium and
 * Hard need their own tokens before they get a color, so they render as neutral pills.
 */
export function DifficultyPill({ difficulty }: { difficulty: number }) {
  if (difficulty <= 2) return <Pill variant="difficulty">Easy</Pill>;
  return <Pill variant="topic">{difficulty === 3 ? 'Medium' : 'Hard'}</Pill>;
}
