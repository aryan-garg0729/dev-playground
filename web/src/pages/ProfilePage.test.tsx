import { screen, within } from '@testing-library/react';
import type { Medal, Medals } from '../api/client';
import { me } from '../test/fixtures';
import { renderRoute, stubApi } from '../test/render';
import { medalLine, ProfilePage } from './ProfilePage';

const tier = (t: string, gte: number, earned = false) => ({ tier: t, gte, title: `X ${t}`, bytes: 50, earned_at: earned ? '2026-11-01T00:00:00Z' : null });
const medal = (o: Partial<Medal>): Medal => ({
  id: 'm', family: 'consistency', name: 'Streak', description: '', unit: 'days', secret: false, coming_soon: false,
  value: 13, tier: 'bronze', next_tier: 'silver', next_gte: 30, tiers: [tier('bronze', 7, true), tier('silver', 30), tier('gold', 100)], ...o,
});

const medals: Medals = {
  earned: 1, total: 3, title: 'Streak I', titles: ['Streak I'],
  families: [{ id: 'consistency', name: 'Consistency' }, { id: 'secret', name: 'Secret' }, { id: 'memory', name: 'Memory' }],
  medals: [
    medal({ id: 'streak' }),
    medal({ id: 'early', name: 'Early Bird', value: 0, tier: '', next_tier: 'bronze', next_gte: 3, unit: 'lessons before 8 am', tiers: [tier('bronze', 3)] }),
    medal({ id: 'secret1', family: 'secret', name: '???', secret: true, tier: '', next_tier: '', next_gte: 0, tiers: [] }),
    medal({ id: 'reviewer', family: 'memory', name: 'Reviewer', coming_soon: true, tier: '', value: 0 }),
  ],
};

describe('ProfilePage', () => {
  it('shows stats, closest medals, families, masked secrets and titles', async () => {
    stubApi({ 'GET /api/v1/me': [200, me], 'GET /api/v1/medals': [200, medals] });
    renderRoute('/profile', '/profile', <ProfilePage />);
    expect(await screen.findByText('1 / 3')).toBeInTheDocument();
    expect(screen.getByRole('heading', { level: 1 })).toHaveTextContent('Boss Slayer I');
    const closest = screen.getByRole('region', { name: 'Closest to unlocking' });
    expect(within(closest).getByRole('img', { name: '13 of 30' })).toBeInTheDocument();
    expect(screen.getByRole('listitem', { name: 'Streak, bronze' })).toHaveTextContent('Bronze · Silver at 30 days');
    expect(screen.getByRole('listitem', { name: '???, not earned' })).toHaveTextContent('Hidden until earned');
    expect(screen.getByRole('listitem', { name: 'Reviewer, not earned' })).toHaveTextContent('Coming soon');
    expect(screen.getByRole('combobox', { name: 'Title' })).toHaveValue('Streak I');
  });

  it('describes progress', () => {
    expect(medalLine(medal({}))).toBe('Bronze · Silver at 30 days');
    expect(medalLine(medal({ tier: 'gold', next_tier: '' }))).toBe('Gold · all tiers earned');
    expect(medalLine(medal({ tier: '', next_tier: 'bronze', next_gte: 3, unit: 'late lessons' }))).toBe('3 late lessons');
    expect(medalLine(medal({ tier: '', next_tier: 'bronze', next_gte: 1, unit: 'bosses' }))).toBe('1 boss');
    expect(medalLine(medal({ tier: '', next_tier: 'bronze', next_gte: 2, unit: 'tier' }))).toBe('Reach Silver league');
    expect(medalLine(medal({ tier: 'bronze', next_tier: 'silver', next_gte: 3, unit: 'tier' }))).toBe('Bronze · Silver: reach Gold league');
  });
});
