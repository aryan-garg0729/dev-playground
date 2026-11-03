import { render, screen } from '@testing-library/react';
import { GameStats, levelProgress } from './GameStats';

describe('GameStats', () => {
  it('shows streak, level and XP inside the level', () => {
    render(<GameStats streakDays={12} level={7} totalXP={1765} levelStartXP={1470} nextLevelXP={1852} />);
    expect(screen.getByLabelText('12-day streak')).toHaveTextContent('12');
    expect(screen.getByText('Lv 7')).toBeInTheDocument();
    // 1765 - 1470 = 295 of 382 XP in level 7.
    const bar = screen.getByRole('img', { name: '295 of 382 XP to level 8' });
    expect(bar.firstElementChild).toHaveStyle({ width: `${(295 / 382) * 100}%` });
  });

  it('clamps progress to the level', () => {
    expect(levelProgress({ totalXP: 0, levelStartXP: 0, nextLevelXP: 100 })).toEqual({ into: 0, of: 100, pct: 0 });
    expect(levelProgress({ totalXP: 150, levelStartXP: 0, nextLevelXP: 100 }).pct).toBe(100);
  });
});
