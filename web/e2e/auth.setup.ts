// Signs in once through dev auth (DEV_AUTH=1 in Compose) and saves the session
// for the other tests. Browser tests never need email. The dev user has every
// chapter open (DEV_AUTH unlocks them), so content CI can reach every lesson.
import { expect, test as setup } from '@playwright/test';

export const STATE = 'e2e/.auth/dev.json';
const ORIGIN = { Origin: 'http://localhost:5173' };

setup('sign in (dev)', async ({ request }) => {
  const res = await request.post('/api/v1/auth/dev', { headers: ORIGIN, data: { timezone: 'UTC' } });
  expect(res.status(), await res.text()).toBe(200);
  // The shared dev user skips onboarding (onboarding.spec.ts covers it with a fresh user).
  const onboard = await request.patch('/api/v1/me', { headers: ORIGIN, data: { onboarded: true } });
  expect(onboard.status(), await onboard.text()).toBe(200);
  await request.storageState({ path: STATE });
});
