// Phase 7 in the browser: a boss fight unlocks the next chapter, Bytes buy a
// Peek that shows only a hidden check's input, and medals reach the profile.
import { readFileSync } from 'node:fs';
import { resolve } from 'node:path';
import { expect, type Page, test } from '@playwright/test';
import { lesson, withDraft } from './helpers';

const ORIGIN = { Origin: 'http://localhost:5173' };
const DIR = resolve(import.meta.dirname, '../../content/foundations/01-python-warm-up');
const file = (p: string) => readFileSync(resolve(DIR, p), 'utf8');

/** Sets an unsent answer once (unlike withDraft, which re-applies on every load). */
async function setDraft(page: Page, id: string, src: string) {
  await page.evaluate(([key, value]) => localStorage.setItem(key, value), [`cq-draft:${id}:v1`, JSON.stringify({ files: { 'main.py': src } })] as const);
}

async function freshUser(page: Page) {
  const email = `e2e-eco-${Date.now()}-${Math.floor(Math.random() * 1e6)}@example.com`;
  const res = await page.request.post('/api/v1/auth/dev', { headers: ORIGIN, data: { email, name: 'Ada Lovelace', timezone: 'UTC' } });
  expect(res.status()).toBe(200);
  await page.addInitScript(() => localStorage.setItem('cq-theme', 'dark'));
}

test('boss fight, Peek and medals', async ({ page }) => {
  await freshUser(page);

  // SQL basics sits behind the Python warm-up boss.
  await page.goto(lesson('foundations', 'fd.sql.select-columns'));
  await expect(page.getByText('Chapter locked')).toBeVisible();

  // A wrong attack: 6 of 8 checks hit (FizzBuzz is checked too late).
  await setDraft(page, 'fd.py.fizzbuzz-boss', file('03-fizzbuzz-boss/wrong/fizz_first.py'));
  await page.goto(lesson('foundations', 'fd.py.fizzbuzz-boss'));
  await expect(page.getByRole('img', { name: '8 of 8 HP left' })).toBeVisible();
  await page.getByRole('button', { name: 'Attack', exact: true }).click();
  await expect(page.getByRole('heading', { name: '6 hits · boss blocked 2' })).toBeVisible({ timeout: 30_000 });
  await expect(page.getByRole('img', { name: '2 of 8 HP left' })).toBeVisible();
  await page.screenshot({ path: 'e2e/screenshots/boss-fight-dark.png' });

  // The killing blow: celebration with the boss bonus and chest; the next chapter opens.
  await setDraft(page, 'fd.py.fizzbuzz-boss', file('03-fizzbuzz-boss/solution.py'));
  await page.reload();
  await page.getByRole('button', { name: 'Attack', exact: true }).click();
  const dialog = page.getByRole('dialog', { name: 'All checks passed' });
  await expect(dialog).toBeVisible({ timeout: 30_000 });
  await expect(dialog).toContainText('Boss defeated+100');
  await expect(dialog).toContainText('Boss chest: +300 Bytes');
  await expect(dialog).toContainText('Next chapter unlocked');
  await page.keyboard.press('Escape');
  await page.goto(lesson('foundations', 'fd.sql.select-columns'));
  await expect(page.getByRole('region', { name: 'Editor' })).toBeVisible();
  await expect(page.getByText('Chapter locked')).toBeHidden();

  // Spend Bytes on Peek, fail a lesson's hidden checks, and peek at one input.
  await page.goto('/shop');
  await expect(page.getByRole('region', { name: 'Balance' })).toContainText('350'); // boss chest 300 + Boss Slayer bronze 50
  await page.getByRole('listitem', { name: 'Peek' }).getByRole('button', { name: 'Buy' }).click();
  await expect(page.getByRole('listitem', { name: 'Peek' })).toContainText('Owned 1');
  await page.screenshot({ path: 'e2e/screenshots/power-ups-shop-dark.png', fullPage: true });

  await withDraft(page, 'fd.py.sum-list', 1, { files: { 'main.py': file('01-sum-list/wrong/return_in_loop.py') } });
  await page.goto(lesson('foundations', 'fd.py.sum-list'));
  await page.getByRole('button', { name: 'Submit' }).click();
  await expect(page.getByRole('heading', { name: /checks failing/ })).toBeVisible({ timeout: 30_000 });
  const used = page.waitForResponse((r) => r.url().endsWith('/powerups/peek/use'));
  await page.getByRole('button', { name: /^Peek/ }).click();
  const peek = await (await used).json();
  expect(Object.keys(peek.peek).sort()).toEqual(['check', 'input', 'name']);
  await expect(page.getByText(`input: ${peek.peek.input}`)).toBeVisible();
  expect(peek.peek.input).toMatch(/^sum_list\(/);
  expect(JSON.stringify(peek)).not.toMatch(/expected|== |assert/);

  // The boss earned Boss Slayer: on the profile, wearable as a title, and the boss is a shipped project.
  await page.goto('/profile');
  await expect(page.getByRole('region', { name: 'Projects completed' }).getByRole('link', { name: 'View code' })).toBeVisible();
  await page.screenshot({ path: 'e2e/screenshots/profile-overview-dark.png', fullPage: true });
  await page.getByRole('button', { name: 'Edit profile' }).click();
  await page.getByRole('combobox', { name: 'Title' }).selectOption('Boss Slayer I');
  await page.getByRole('button', { name: 'Save' }).click();
  await expect(page.getByRole('heading', { level: 1 })).toContainText('Boss Slayer I');
  await page.getByRole('link', { name: 'Medals' }).click();
  await expect(page.getByRole('listitem', { name: 'Boss Slayer, bronze' })).toBeVisible();
  await page.screenshot({ path: 'e2e/screenshots/profile-medals-dark.png', fullPage: true });
});
