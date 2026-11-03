// Lighthouse audit of the landing page (SPEC Phase 6: score above 90).
// Runs in the Playwright image (`make lighthouse`) against a production build
// served by `vite preview`, using the image's Chromium. Fails if any category
// scores below 90. The HTML report goes to e2e/.results/lighthouse.html.
import { mkdirSync, writeFileSync } from 'node:fs';
import { chromium } from '@playwright/test';
import * as chromeLauncher from 'chrome-launcher';
import lighthouse from 'lighthouse';

const url = process.argv[2] ?? 'http://localhost:4173/';
const out = process.argv[3] ?? 'e2e/.results';
const MIN = 0.9;

const chrome = await chromeLauncher.launch({
  chromePath: chromium.executablePath(),
  chromeFlags: ['--headless=new', '--no-sandbox', '--disable-gpu'],
});
try {
  const result = await lighthouse(url, {
    port: chrome.port,
    output: 'html',
    onlyCategories: ['performance', 'accessibility', 'best-practices', 'seo'],
    formFactor: 'mobile', // Lighthouse's default (and stricter) profile
  });
  if (!result) throw new Error('lighthouse returned no result');
  mkdirSync(out, { recursive: true });
  writeFileSync(`${out}/lighthouse.html`, result.report);
  let failed = false;
  for (const c of Object.values(result.lhr.categories)) {
    const score = Math.round((c.score ?? 0) * 100);
    const ok = (c.score ?? 0) >= MIN;
    failed ||= !ok;
    console.log(`${ok ? 'ok  ' : 'FAIL'} ${c.title}: ${score}`);
  }
  {
    // List what costs points, failing or not, so a drop is easy to trace.
    for (const a of Object.values(result.lhr.audits)) {
      if (a.score !== null && a.score < 0.9 && a.scoreDisplayMode !== 'informative' && a.scoreDisplayMode !== 'notApplicable') {
        console.log(`  - ${a.title}${a.displayValue ? ` (${a.displayValue})` : ''}`);
        const items = a.details?.items ?? [];
        for (const it of items.slice(0, 5)) {
          const where = it.node?.selector ?? it.url ?? it.source?.url ?? '';
          const why = it.node?.explanation?.split('\n')[0] ?? it.description ?? '';
          if (where || why) console.log(`      ${where} ${why}`.trimEnd());
        }
      }
    }
  }
  if (failed) process.exitCode = 1;
} finally {
  await chrome.kill();
}
