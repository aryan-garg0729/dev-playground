// Phase 4 in the browser: sign-in redirect, a first pass with its celebration
// and XP, no XP on a re-solve, and the league join.
import { expect, type Page, test } from '@playwright/test';
import { lesson } from './helpers';

const ORIGIN = { Origin: 'http://localhost:5173' };
const QUIZ = lesson('langchain-basics', 'lc.chat.message-roles'); // 5 XP, answer: SystemMessage

/** Signs this page's context in as a brand-new user. */
async function freshUser(page: Page) {
  const email = `e2e-${Date.now()}-${Math.floor(Math.random() * 1e6)}@example.com`;
  const res = await page.request.post('/api/v1/auth/dev', { headers: ORIGIN, data: { email, name: 'E2E Learner', timezone: 'UTC' } });
  expect(res.status()).toBe(200);
  // Skip onboarding (Phase 6): `/` would otherwise open /welcome. onboarding.spec.ts covers it.
  const onboard = await page.request.patch('/api/v1/me', { headers: ORIGIN, data: { onboarded: true } });
  expect(onboard.status()).toBe(200);
}

test.describe('signed out', () => {
  test.use({ storageState: { cookies: [], origins: [] } });

  test('pages send you to sign in and back', async ({ page }) => {
    await page.goto('/league');
    await expect(page).toHaveURL(/\/sign-in\?next=%2Fleague$/);
    await page.getByRole('button', { name: /Dev sign-in/ }).click();
    await expect(page).toHaveURL(/\/league$/);
    await expect(page.getByRole('heading', { level: 1, name: /league$/ })).toBeVisible();
  });
});

test('first pass: celebration, XP once, streak and league', async ({ page }) => {
  await freshUser(page);
  await page.addInitScript(() => localStorage.setItem('cq-theme', 'dark'));
  await page.goto(QUIZ);
  await expect(page.getByLabel('0-day streak')).toBeVisible();

  await page.getByLabel('SystemMessage').check();
  await page.getByRole('button', { name: 'Submit' }).click();

  const dialog = page.getByRole('dialog', { name: 'All checks passed' });
  await expect(dialog).toBeVisible();
  // 5 base + 3 first try (50%, rounded) + 1 no hint (25%, rounded).
  await expect(dialog).toContainText('+9 XP');
  await expect(dialog).toContainText('first try');
  await expect(dialog).toContainText('1-day streak');
  await expect(page.getByRole('link', { name: /^Next: / })).toBeFocused();
  await page.screenshot({ path: 'e2e/screenshots/pass-celebration-dark.png' });

  await page.keyboard.press('Escape');
  await expect(dialog).toBeHidden();
  await expect(page.getByRole('img', { name: '9 of 100 XP to level 2' })).toBeVisible();
  await expect(page.getByLabel('1-day streak')).toBeVisible();

  // Re-solving: graded again, but no XP and no celebration.
  const again = page.waitForResponse((r) => r.url().endsWith('/submit'));
  await page.getByRole('button', { name: 'Submit' }).click();
  const second = await (await again).json();
  expect(second.result.status).toBe('pass');
  expect(second.xp_awarded).toBe(0);
  await expect(page.getByText('Correct')).toBeVisible();
  await expect(dialog).toBeHidden();
  const me = await (await page.request.get('/api/v1/me')).json();
  expect(me.stats.total_xp).toBe(9);

  // The first XP of the week joined a league group.
  await page.goto('/');
  const league = page.getByRole('region', { name: /league$/ });
  await expect(league).toContainText(/#\d+of \d+ · 9 XP this week/);
});
