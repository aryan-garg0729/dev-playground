import { screen, within } from '@testing-library/react';
import type { Home } from '../api/client';
import { me } from '../test/fixtures';
import { renderRoute, stubApi } from '../test/render';
import { HomePage } from './HomePage';
import { nudge, timeLeft } from './home/format';

const week = ['2026-11-02', '2026-11-03', '2026-11-04', '2026-11-05', '2026-11-06', '2026-11-07', '2026-11-08'].map((date, i) => ({
  date, active: i < 2, today: i === 2,
}));

const home: Home = {
  continue: {
    course_id: 'foundations', course_title: 'Foundations', chapter_title: 'Python warm-up', chapter_index: 1,
    exercise_id: 'fd.py.sum-list', exercise_title: 'Sum a list', exercise_type: 'code', position: { index: 2, total: 3 },
    lessons: ['passed', 'started', 'not_started'],
    last: { status: 'fail', passed: 1, total: 3, file: 'main.py', error_line: 5, lines: [
      { n: 3, text: '    for n in nums:' }, { n: 4, text: '        total += n' }, { n: 5, text: '        return total' },
    ] },
  },
  today: { goal: 3, done: 0, streak_days: 12, week },
  paths: [{ id: 'foundations', title: 'Foundations', description: '', chapter_count: 2, exercise_count: 5, passed_count: 1, status: 'in_progress' }],
  league: {
    joined: true, tier: 'gold', rank: 7, of: 30, weekly_xp: 340, ends_at: new Date(Date.now() + 2 * 86400000 + 6 * 3600000 + 60000).toISOString(),
    promote: 5, demote: 5, xp_to_promote: 85, zone: 'stay',
  },
};

describe('HomePage', () => {
  it('shows where to continue, today, paths, placeholders and the league', async () => {
    stubApi({ 'GET /api/v1/home': [200, home], 'GET /api/v1/me': [200, me] });
    renderRoute('/', '/', <HomePage />);

    expect(await screen.findByRole('heading', { level: 1, name: /, Ada$/ })).toBeInTheDocument();
    expect(screen.getByText('One lesson keeps your 12-day streak.')).toBeInTheDocument();

    const cont = screen.getByRole('region', { name: 'Sum a list' });
    expect(within(cont).getByRole('link', { name: 'Resume lesson' })).toHaveAttribute('href', '/learn/foundations/fd.py.sum-list');
    expect(within(cont).getByText('return total', { exact: false }).closest('div')?.className).toContain('bg-danger-soft');
    expect(cont).toHaveTextContent('2 of 3 checks were failing');

    const today = screen.getByRole('region', { name: 'Today' });
    expect(today).toHaveTextContent('0 / 3');
    expect(within(today).getByLabelText('2026-11-04, today')).toBeInTheDocument();
    expect(within(today).getByLabelText('2026-11-02, practised')).toBeInTheDocument();

    expect(screen.getByRole('region', { name: 'Daily quests' })).toHaveTextContent('Coming soon');
    expect(screen.getByRole('button', { name: 'Start review' })).toBeDisabled();

    const league = screen.getByRole('region', { name: 'Gold league' });
    expect(league).toHaveTextContent('#7of 30 · 340 XP this week');
    expect(league).toHaveTextContent('Top 5 promote to Diamond · 85 XP to go');
    expect(league).toHaveTextContent('ends in 2d 6h');
  });
});

describe('home format', () => {
  it('formats time left', () => {
    const now = Date.parse('2026-11-03T00:00:00Z');
    expect(timeLeft('2026-11-05T06:30:00Z', now)).toBe('2d 6h');
    expect(timeLeft('2026-11-03T05:12:00Z', now)).toBe('5h 12m');
    expect(timeLeft('2026-11-02T00:00:00Z', now)).toBe('0m');
  });

  it('nudges toward the streak, then the goal', () => {
    const t = (activeToday: boolean, done: number, streak: number) => ({
      goal: 2, done, streak_days: streak, week: [{ today: true, active: activeToday }],
    });
    expect(nudge(t(false, 0, 4))).toBe('One lesson keeps your 4-day streak.');
    expect(nudge(t(false, 0, 0))).toBe('Pass a lesson today to start a streak.');
    expect(nudge(t(true, 1, 5))).toBe('1 more lesson to reach your daily goal.');
    expect(nudge(t(true, 2, 5))).toBe('Daily goal done. Anything more is a bonus.');
  });
});
