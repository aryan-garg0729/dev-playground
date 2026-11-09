// Signed-out pages send you to sign in and bring you back.
import { expect, test } from '@playwright/test';

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
