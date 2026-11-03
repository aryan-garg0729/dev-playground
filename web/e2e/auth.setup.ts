// Signs in once through dev auth (DEV_AUTH=1 in Compose) and saves the session
// for the other tests. Browser tests never need email.
import { expect, test as setup } from '@playwright/test';

export const STATE = 'e2e/.auth/dev.json';

setup('sign in (dev)', async ({ request }) => {
  const res = await request.post('/api/v1/auth/dev', { headers: { Origin: 'http://localhost:5173' }, data: { timezone: 'UTC' } });
  expect(res.status(), await res.text()).toBe(200);
  await request.storageState({ path: STATE });
});
