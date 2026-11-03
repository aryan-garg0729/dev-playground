// Resizable lesson panes: drag the handles in a real browser (over Monaco), and
// the sizes stick after a reload.
import { expect, type Page, test } from '@playwright/test';
import { lesson } from './helpers';

const width = async (page: Page, name: string) => (await page.getByRole('region', { name }).boundingBox())!.width;
const checksHeight = async (page: Page) => (await page.locator('[aria-live="polite"]:not(.monaco-status)').boundingBox())!.height;

async function drag(page: Page, name: string, dx: number, dy: number) {
  const box = (await page.getByRole('separator', { name }).boundingBox())!;
  const x = box.x + box.width / 2;
  const y = box.y + box.height / 2;
  await page.mouse.move(x, y);
  await page.mouse.down();
  await page.mouse.move(x + dx / 2, y + dy / 2, { steps: 4 });
  await page.mouse.move(x + dx, y + dy, { steps: 4 });
  await page.mouse.up();
}

test('problem, editor and checks panes can be resized', async ({ page }) => {
  await page.goto(lesson('foundations', 'fd.py.sum-list'));
  await page.evaluate(() => localStorage.removeItem('cq-panes')); // start from the defaults, once
  await page.reload();
  await page.getByRole('button', { name: 'Submit' }).click();
  await expect(page.getByRole('heading', { name: /checks (failing|passing)/ })).toBeVisible({ timeout: 30_000 });

  expect(Math.round(await width(page, 'Problem'))).toBe(320);
  const editorBefore = await width(page, 'Editor');
  await drag(page, 'Resize problem and editor', 120, 0);
  expect(Math.round(await width(page, 'Problem'))).toBe(440);
  expect(Math.round(await width(page, 'Editor'))).toBe(Math.round(editorBefore - 120));

  const checksBefore = await checksHeight(page);
  await drag(page, 'Resize editor and checks', 0, -100); // dragging up grows the checks
  expect(Math.round(await checksHeight(page))).toBe(Math.round(checksBefore + 100));

  // The editor always keeps room: dragging far right stops at its minimum width.
  await drag(page, 'Resize problem and editor', 2000, 0);
  expect(await width(page, 'Editor')).toBeGreaterThanOrEqual(359);
  await page.screenshot({ path: 'e2e/screenshots/lesson-resized.png' });

  // Keyboard: Enter on the handle resets it.
  await page.getByRole('separator', { name: 'Resize problem and editor' }).focus();
  await page.keyboard.press('Enter');
  expect(Math.round(await width(page, 'Problem'))).toBe(320);

  // Sizes are remembered on this device.
  await drag(page, 'Resize problem and editor', 80, 0);
  await page.reload();
  await expect(page.getByRole('region', { name: 'Editor' })).toBeVisible();
  expect(Math.round(await width(page, 'Problem'))).toBe(400);
});
