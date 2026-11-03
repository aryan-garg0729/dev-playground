// The tutor panel opens on a lesson without errors (it needs no model to open).
// Regression: an effect returned scrollIntoView's Promise in Chromium and crashed the page.
import { expect, test } from '@playwright/test';
import { lesson } from './helpers';

test('the tutor panel opens on a lesson without page errors', async ({ page }) => {
  const errors: string[] = [];
  page.on('pageerror', (e) => errors.push(e.message));
  page.on('console', (m) => {
    if (m.type() === 'error' && m.text().includes('useEffect')) errors.push(m.text());
  });
  await page.goto(lesson('foundations', 'fd.py.sum-list'));
  const ask = page.getByRole('button', { name: 'Ask the tutor' });
  await expect(ask).toBeVisible();
  test.skip(await ask.isDisabled(), 'tutor is off (TUTOR_PROVIDER empty)');
  await ask.click();
  await expect(page.getByRole('region', { name: 'AI tutor' })).toBeVisible();
  await expect(page.getByText(/hints left today/)).toBeVisible();
  await expect(page.getByRole('textbox', { name: 'Ask the tutor' })).toBeEnabled();
  expect(errors).toEqual([]);
});
