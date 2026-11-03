import { defineConfig } from 'vitest/config';

const appOrigins = (process.env.VITE_APP_ORIGINS ?? 'http://localhost:5173').split(',').map((s) => s.trim());
// The page has an opaque origin (sandbox="allow-scripts", no allow-same-origin), so its
// own origin is named explicitly where 'self' would not match it.
const self = process.env.SANDBOX_ORIGIN ?? 'http://localhost:5174';

const csp = [
  "default-src 'none'",
  `script-src 'self' ${self} 'wasm-unsafe-eval'`,
  `connect-src 'self' ${self}`,
  'worker-src blob:',
  "base-uri 'none'",
  "form-action 'none'",
  `frame-ancestors ${appOrigins.join(' ')}`,
].join('; ');

const headers = {
  'Content-Security-Policy': csp,
  // Assets are public; the opaque-origin page and its blob worker fetch them cross-origin.
  'Access-Control-Allow-Origin': '*',
  'Cross-Origin-Resource-Policy': 'cross-origin',
  'X-Content-Type-Options': 'nosniff',
  'Referrer-Policy': 'no-referrer',
};

export default defineConfig({
  // No SPA fallback: a missing asset is a 404, not index.html.
  appType: 'mpa',
  worker: { format: 'iife' },
  build: { target: 'es2022', assetsInlineLimit: 0 },
  preview: { headers },
  server: { headers },
  test: { environment: 'node' },
});
