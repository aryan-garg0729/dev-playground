// End-to-end tests against the running Compose stack (make e2e runs them in the
// official Playwright image, so no browsers are installed on the host).
import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
  testDir: './e2e',
  timeout: 120_000,
  workers: 1,
  reporter: [['list']],
  outputDir: './e2e/.results',
  use: {
    baseURL: process.env.APP_URL ?? 'http://localhost:5173',
    ...devices['Desktop Chrome'],
    viewport: { width: 1280, height: 800 },
  },
});
