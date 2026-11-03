// Screenshots of the lesson screen in both themes, for comparison with
// docs/design/screens/lesson-*.png. Written to e2e/screenshots/.
import { expect, test } from '@playwright/test';
import { lesson, withDraft } from './helpers';

const failing = 'def sum_list(nums):\n    total = 0\n    for n in nums:\n        total += n\n        return total\n';

for (const theme of ['dark', 'light'] as const) {
  for (const rail of ['closed', 'open'] as const) {
    test(`lesson screen, ${theme}, rail ${rail}`, async ({ page }) => {
      await page.addInitScript(
        ([t, r]) => {
          localStorage.setItem('cq-theme', t);
          localStorage.setItem('cq-rail-open', r === 'open' ? '1' : '0');
        },
        [theme, rail] as const,
      );
      await withDraft(page, 'fd.py.sum-list', 1, { files: { 'main.py': failing } });
      await page.goto(lesson('foundations', 'fd.py.sum-list'));
      await page.getByRole('button', { name: 'Show a hint' }).click();
      await page.getByRole('button', { name: 'Submit' }).click();
      await expect(page.getByRole('heading', { name: /checks (failing|passing)/ })).toBeVisible({ timeout: 30_000 });
      await expect(page.locator('.monaco-editor')).toBeVisible();
      await page.screenshot({ path: `e2e/screenshots/lesson-${theme}${rail === 'open' ? '-outline-open' : ''}.png` });
    });
  }
}
