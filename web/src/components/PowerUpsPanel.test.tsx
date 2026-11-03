import { screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import type { GradeResult, Shop } from '../api/client';
import { nothing } from '../test/fixtures';
import { renderRoute, stubApi } from '../test/render';
import { PowerUpsPanel } from './PowerUpsPanel';

const shop = (owned: Record<string, number>, effects: Shop['effects'] = []): Shop => ({
  balance: 0, earned_today: 0, bet: null, earn: [], effects,
  items: Object.entries(owned).map(([id, n]) => ({ id, name: id, description: '', note: '', cost: 1, owned: n, state: 'buy', use: 'lessons', usable_now: false })),
});
const failed: GradeResult = {
  submission_id: 's', xp_awarded: 0, outcome: nothing,
  result: { status: 'fail', score: 0, stdout: '', runtime_ms: 1, checks: [{ name: 'Hidden test 1', passed: false, visible: false }] },
};

describe('PowerUpsPanel', () => {
  it('peeks at a failed hidden check and shows only its input', async () => {
    stubApi({
      'GET /api/v1/shop': [200, shop({ peek: 2, second_wind: 1 })],
      'POST /api/v1/powerups/peek/use': [200, { powerup: 'peek', remaining: 1, balance: 0, peek: { check: 0, name: 'Hidden test 1', input: 'fizzbuzz(15)[14]' } }],
    });
    renderRoute('/x', '/x', <PowerUpsPanel exerciseId="fd.py.fizzbuzz-boss" result={failed} />);
    await userEvent.click(await screen.findByRole('button', { name: /Peek/ }));
    expect(await screen.findByText('fizzbuzz(15)[14]')).toBeInTheDocument();
  });

  it('only enables Peek and Second Wind after a failed submit', async () => {
    stubApi({ 'GET /api/v1/shop': [200, shop({ peek: 1, second_wind: 1 })] });
    renderRoute('/x', '/x', <PowerUpsPanel exerciseId="e" result={undefined} />);
    expect(await screen.findByRole('button', { name: /Peek/ })).toBeDisabled();
    expect(screen.getByRole('button', { name: /Second Wind/ })).toBeDisabled();
  });

  it('shows an active Second Wind and links to the shop when empty', async () => {
    stubApi({ 'GET /api/v1/shop': [200, shop({}, [{ effect: 'second_wind', ref_id: 'e', expires_at: '9999-01-01T00:00:00Z' }])] });
    renderRoute('/x', '/x', <PowerUpsPanel exerciseId="e" result={failed} />);
    expect(await screen.findByRole('button', { name: /Second Wind\s*active/ })).toBeDisabled();
  });
});
