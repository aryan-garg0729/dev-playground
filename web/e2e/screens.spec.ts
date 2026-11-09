// Screenshots for comparison with docs/design/screens/*.png, written to
// e2e/screenshots/: Home and the course list in both themes.
import { expect, test } from '@playwright/test';

for (const theme of ['dark', 'light'] as const) {
  test(`home, ${theme}`, async ({ page }) => {
    await page.addInitScript((t) => localStorage.setItem('cq-theme', t), theme);
    await page.goto('/');
    await expect(page.getByRole('region', { name: 'Today' })).toBeVisible();
    await page.screenshot({ path: `e2e/screenshots/home-${theme}.png`, fullPage: true });
  });

  test(`courses, ${theme}`, async ({ page }) => {
    await page.addInitScript((t) => localStorage.setItem('cq-theme', t), theme);
    await page.goto('/courses');
    await expect(page.getByRole('region', { name: /All courses/ })).toBeVisible();
    await page.screenshot({ path: `e2e/screenshots/courses-${theme}.png`, fullPage: true });
  });
}
