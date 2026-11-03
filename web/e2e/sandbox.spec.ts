import { expect, type Page, test } from '@playwright/test';
import { lesson, SANDBOX, withDraft } from './helpers';

/** The Run button by its title: its name changes to "Running…" while a run is in progress. */
function runButton(page: Page) {
  return page.getByTitle('Run with the visible checks (not graded)');
}

const SUM = 'fd.py.sum-list';
const SUM_URL = lesson('foundations', SUM);
// The results pane (Monaco adds its own aria-live regions, so exclude those).
const resultsPane = (page: Page) => page.locator('[aria-live="polite"]:not(.monaco-status)');

test('Run shows Python output in under 1 s once the runtime is loaded', async ({ page }) => {
  await withDraft(page, SUM, 1, { files: { 'main.py': 'print("hello from pyodide")\ndef sum_list(nums):\n    return 0\n' } });
  await page.goto(SUM_URL);
  const run = runButton(page);
  await run.click(); // first Run may wait for Pyodide to finish loading
  await expect(resultsPane(page).getByText('hello from pyodide')).toBeVisible({ timeout: 60_000 });

  for (let i = 0; i < 3; i++) {
    const start = Date.now();
    await run.click();
    await expect(run).toHaveText('Run');
    await expect(resultsPane(page).getByText('hello from pyodide')).toBeVisible();
    const ms = Date.now() - start;
    console.log(`warm Run ${i + 1}: ${ms} ms`);
    expect(ms).toBeLessThan(1000);
  }
});

test('code in the sandbox cannot read the app cookies, DOM or storage, or reach the network', async ({ page, context }) => {
  await context.addCookies([{ name: 'cq_secret', value: 's3cr3t-cookie', url: 'http://localhost:5173' }]);
  const learner = [
    'import js',
    'out = []',
    'for name, f in [',
    '    ("doc", lambda: js.document.cookie),',
    '    ("parent", lambda: js.parent.document.cookie),',
    '    ("storage", lambda: js.localStorage.getItem("cq-theme")),',
    ']:',
    '    try:',
    '        out.append(name + ":" + str(f()))',
    '    except Exception as e:',
    '        out.append(name + "-blocked")',
    'try:',
    '    x = js.XMLHttpRequest.new()',
    '    x.open("GET", "https://example.com/", False)',
    '    x.send()',
    '    out.append("net:" + str(x.status))',
    'except Exception as e:',
    '    out.append("net-blocked")',
    'print(" ".join(out))',
    'def sum_list(nums):',
    '    return 0',
  ].join('\n');
  await withDraft(page, SUM, 1, { files: { 'main.py': learner } });
  await page.goto(SUM_URL);
  expect(await page.evaluate(() => document.cookie)).toContain('s3cr3t-cookie');

  await runButton(page).click();
  const line = resultsPane(page).getByText(/doc-blocked/);
  await expect(line).toBeVisible({ timeout: 60_000 });
  await expect(line).toHaveText('doc-blocked parent-blocked storage-blocked net-blocked');
  await expect(resultsPane(page)).not.toContainText('s3cr3t');

  // The frame itself: opaque origin, no cookie or storage access, no reach into the app.
  const iframe = page.locator(`iframe[src^="${SANDBOX}"]`);
  await expect(iframe).toHaveAttribute('sandbox', 'allow-scripts');
  const frame = page.frames().find((f) => f.url().startsWith(SANDBOX));
  expect(frame).toBeDefined();
  const probe = await frame!.evaluate(() => {
    const tryIt = (f: () => unknown) => {
      try {
        return `read:${String(f())}`;
      } catch (e) {
        return `blocked:${(e as Error).name}`;
      }
    };
    return {
      origin: window.origin,
      cookie: tryIt(() => document.cookie),
      storage: tryIt(() => window.localStorage.length),
      parentDom: tryIt(() => window.parent.document.title),
      parentCookie: tryIt(() => window.parent.document.cookie),
    };
  });
  expect(probe.origin).toBe('null');
  expect(probe.cookie).toMatch(/^blocked/);
  expect(probe.storage).toMatch(/^blocked/);
  expect(probe.parentDom).toMatch(/^blocked/);
  expect(probe.parentCookie).toMatch(/^blocked/);
});

test('an infinite loop is stopped at timeout_ms and the page stays responsive', async ({ page }) => {
  await withDraft(page, SUM, 1, { files: { 'main.py': 'def sum_list(nums):\n    while True:\n        pass\n' } });
  await page.goto(SUM_URL);
  const run = runButton(page);
  const stopped = resultsPane(page).getByText('Stopped after 3000 ms');

  // First Run also loads Pyodide; after it, the spare worker is warm.
  await run.click();
  await expect(stopped).toBeVisible({ timeout: 70_000 });
  await expect(run).toHaveText('Run');

  const start = Date.now();
  await run.click();
  await expect(run).toHaveText('Running…');
  // While the learner's loop spins in the worker, the app reacts at once.
  const themeBefore = await page.evaluate(() => document.documentElement.dataset.theme);
  const t0 = Date.now();
  await page.getByRole('button', { name: /Switch to (light|dark) theme/ }).click();
  await expect.poll(() => page.evaluate(() => document.documentElement.dataset.theme)).not.toBe(themeBefore);
  expect(Date.now() - t0).toBeLessThan(1000);

  await expect(run).toHaveText('Run', { timeout: 10_000 });
  await expect(stopped).toBeVisible();
  const took = Date.now() - start;
  console.log(`loop stopped after ${took} ms (limit 3000 ms)`);
  expect(took).toBeGreaterThanOrEqual(3000);
  expect(took).toBeLessThan(4000);
});
