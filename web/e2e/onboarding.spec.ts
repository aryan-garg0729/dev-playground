// Phase 6: the landing page for visitors, onboarding for a new learner, and
// analytics events reaching the API.
import { expect, test } from '@playwright/test';

const ORIGIN = { Origin: 'http://localhost:5173' };

test.describe('signed out', () => {
  test.use({ storageState: { cookies: [], origins: [] } });

  test('the landing page shows the hero, real courses and the try-it card', async ({ page }) => {
    await page.goto('/');
    await expect(page.getByRole('heading', { level: 1 })).toContainText('Learn to code');
    await expect(page.getByRole('heading', { name: 'Pick a path' })).toBeVisible();
    await expect(page.getByRole('link', { name: 'Start the path →' }).first()).toBeVisible();
    await page.getByRole('button', { name: 'Run checks' }).click();
    await expect(page.getByRole('status').filter({ hasText: 'checks passing' })).toBeVisible();
    await page.screenshot({ path: 'e2e/screenshots/landing-dark.png', fullPage: true });
    await page.getByRole('link', { name: 'Start free' }).first().click();
    await expect(page).toHaveURL(/\/sign-in/);
  });

  test('a new learner is onboarded and lands on the course for their goal', async ({ page }) => {
    const email = `new-${Date.now()}@codequest.local`;
    const res = await page.request.post('/api/v1/auth/dev', { headers: ORIGIN, data: { email, timezone: 'UTC' } });
    expect(res.status(), await res.text()).toBe(200);

    const events = page.waitForRequest((r) => r.url().endsWith('/api/v1/events') && r.method() === 'POST');
    await page.goto('/');
    await expect(page).toHaveURL(/\/welcome$/);
    await expect(page.getByText('Step 1 of 2')).toBeVisible();
    await page.screenshot({ path: 'e2e/screenshots/onboarding-dark.png', fullPage: true });
    await page.getByText('Learn to code from zero').click();
    await page.getByText('10 min').click();
    await page.getByRole('button', { name: 'Continue' }).click();
    await expect(page.getByText('Step 2 of 2')).toBeVisible();
    await page.getByRole('button', { name: 'Start learning' }).click();
    await expect(page).toHaveURL(/\/courses\/[^/]+$/);

    const me = await (await page.request.get('/api/v1/me')).json();
    expect(me.user).toMatchObject({ onboarded: true, goal: 'zero' });
    expect(me.stats.daily_goal).toBe(1);
    // Page views are batched and sent within a few seconds.
    expect((await events).postDataJSON().events.length).toBeGreaterThan(0);
  });
});
