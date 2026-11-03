import { screen, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import type { Shop } from '../api/client';
import { me } from '../test/fixtures';
import { renderRoute, stubApi } from '../test/render';
import { ShopPage } from './ShopPage';

const item = (id: string, name: string, cost: number, owned: number, state: string, use = '', usable = false) => ({
  id, name, description: `${name} does a thing.`, note: 'note', cost, owned, state, use, usable_now: usable,
}) as Shop['items'][number];

const shop: Shop = {
  balance: 300, earned_today: 60, effects: [], bet: null,
  earn: [{ label: 'Boss defeated', amount: '300' }],
  items: [
    item('streak_freeze', 'Streak Freeze', 200, 2, 'max_held', 'auto'),
    item('focus_boost', 'Focus Boost', 150, 1, 'buy', 'now', true),
    item('league_shield', 'League Shield', 400, 0, 'buy', 'now'),
    item('holiday_mode', 'Holiday Mode', 0, 0, 'pro'),
    item('extra_hint', 'Extra Hint', 100, 0, 'coming_soon'),
    item('double_or_nothing', 'Double or Nothing', 0, 0, 'place_bet'),
  ],
};

describe('ShopPage', () => {
  it('shows states, disables what cannot be afforded, and buys with a request id', async () => {
    const spy = stubApi({ 'GET /api/v1/shop': [200, shop], 'GET /api/v1/me': [200, me], 'POST /api/v1/shop/focus_boost/buy': [200, { powerup: 'focus_boost', owned: 2, balance: 150 }] });
    renderRoute('/shop', '/shop', <ShopPage />);
    const card = (name: string) => within(screen.getByRole('listitem', { name }));
    expect(await screen.findByRole('listitem', { name: 'Streak Freeze' })).toBeInTheDocument();
    expect(card('Streak Freeze').getByRole('button', { name: 'Max held' })).toBeDisabled();
    expect(card('League Shield').getByRole('button', { name: 'Buy' })).toBeDisabled(); // 400 > 300
    expect(card('Holiday Mode').getByRole('button', { name: 'Pro' })).toBeDisabled();
    expect(card('Extra Hint').getByRole('button', { name: 'Coming soon' })).toBeDisabled();
    expect(screen.getByRole('region', { name: 'Balance' })).toHaveTextContent('+60 today');
    expect(within(screen.getByRole('region', { name: 'Inventory' })).getByRole('button', { name: 'Use now' })).toBeEnabled();

    await userEvent.click(card('Focus Boost').getByRole('button', { name: 'Buy' }));
    const req = spy.mock.calls.map(([r]) => r as Request).find((r) => r.url.endsWith('/shop/focus_boost/buy'))!;
    expect((await req.json()).request_id).toMatch(/^[0-9a-f-]{36}$/);
  });

  it('shows why a buy was refused', async () => {
    stubApi({
      'GET /api/v1/shop': [200, { ...shop, balance: 1000 }], 'GET /api/v1/me': [200, me],
      'POST /api/v1/shop/league_shield/buy': [409, { error: { code: 'insufficient_bytes', message: 'not enough Bytes' } }],
    });
    renderRoute('/shop', '/shop', <ShopPage />);
    await userEvent.click(within(await screen.findByRole('listitem', { name: 'League Shield' })).getByRole('button', { name: 'Buy' }));
    expect(await screen.findByRole('alert')).toHaveTextContent('not enough Bytes');
  });
});
