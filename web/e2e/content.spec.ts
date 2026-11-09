// Content CI, browser half: every exercise's reference solution passes its
// visible checks wherever its `run:` points (browser sandbox or server Run).
// The server half (all checks on the grading runner) is `cq test`.
import { readdirSync, readFileSync, statSync } from 'node:fs';
import { join, resolve } from 'node:path';
import { expect, test } from '@playwright/test';
import { runInSandbox } from './helpers';

const CONTENT = resolve(import.meta.dirname, '../../content');

interface Found {
  id: string;
  type: string;
  solution: string;
}

/** Walks content/ for exercise.yaml files (only the few header fields needed here). */
function exercises(dir = CONTENT): Found[] {
  const out: Found[] = [];
  for (const name of readdirSync(dir)) {
    const p = join(dir, name);
    if (!statSync(p).isDirectory() || name.startsWith('.') || name === 'wrong') continue;
    const manifest = join(p, 'exercise.yaml');
    try {
      const y = readFileSync(manifest, 'utf8');
      const field = (re: RegExp) => re.exec(y)?.[1]?.trim() ?? '';
      const type = field(/^type:\s*(\S+)/m);
      const sol = field(/^\s+solution:\s*(\S+)/m);
      if (type === 'code' || type === 'sql') {
        out.push({ id: field(/^id:\s*(\S+)/m), type, solution: readFileSync(join(p, sol), 'utf8') });
      }
      continue;
    } catch {
      out.push(...exercises(p));
    }
  }
  return out;
}

const all = exercises();

for (const ex of all) {
  test(`${ex.id}: solution passes its visible checks on Run`, async ({ page, request }) => {
    const detail = await (await request.get(`/api/v1/exercises/${ex.id}`)).json();
    const pv = detail.public_view;
    const answer = ex.type === 'code' ? { files: { [pv.entry]: ex.solution } } : { query: ex.solution };
    test.skip(pv.run === 'off', 'run: off');

    let checks: { name: string; passed: boolean; feedback?: string }[];
    let error: string | undefined;
    if (pv.run === 'server') {
      const res = await request.post(`/api/v1/exercises/${ex.id}/run`, {
        headers: { Origin: 'http://localhost:5173' },
        data: { exercise_id: ex.id, version: detail.version, answer },
      });
      expect(res.status(), await res.text()).toBe(200);
      ({ checks, error } = await res.json());
    } else {
      await page.goto('/');
      const job =
        ex.type === 'code'
          ? { runner: 'pyodide', entry: pv.entry, files: answer.files, tests: pv.visible_tests, timeoutMs: pv.timeout_ms }
          : { runner: 'sqljs', entry: 'query.sql', files: { 'query.sql': ex.solution }, tests: [], sql: { setup: pv.setup, checks: pv.visible_checks }, timeoutMs: pv.timeout_ms };
      const r = (await runInSandbox(page, job)) as { checks: typeof checks; error?: string; timedOut: boolean };
      expect(r.timedOut).toBe(false);
      ({ checks, error } = r);
    }
    expect(error ?? '').toBe('');
    expect(checks.length, 'a runnable exercise needs at least one visible check').toBeGreaterThan(0);
    expect(checks.filter((c) => !c.passed)).toEqual([]);
  });
}
