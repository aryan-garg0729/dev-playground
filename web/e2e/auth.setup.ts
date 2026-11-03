// Signs in once through dev auth (DEV_AUTH=1 in Compose) and saves the session
// for the other tests. Browser tests never need email. The dev user also
// defeats the Foundations boss, so chapters behind it are open for content CI.
import { readFileSync } from 'node:fs';
import { resolve } from 'node:path';
import { expect, test as setup } from '@playwright/test';

export const STATE = 'e2e/.auth/dev.json';
const ORIGIN = { Origin: 'http://localhost:5173' };
const BOSS = 'fd.py.fizzbuzz-boss';

setup('sign in (dev)', async ({ request }) => {
  const res = await request.post('/api/v1/auth/dev', { headers: ORIGIN, data: { timezone: 'UTC' } });
  expect(res.status(), await res.text()).toBe(200);
  const solution = readFileSync(resolve(import.meta.dirname, '../../content/foundations/01-python-warm-up/03-fizzbuzz-boss/solution.py'), 'utf8');
  const boss = await request.post(`/api/v1/exercises/${BOSS}/submit`, {
    headers: ORIGIN, data: { exercise_id: BOSS, version: 1, answer: { files: { 'main.py': solution } } },
  });
  expect(boss.status(), await boss.text()).toBe(200);
  expect((await boss.json()).result.status).toBe('pass');
  await request.storageState({ path: STATE });
});
