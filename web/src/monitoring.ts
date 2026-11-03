/**
 * Error monitoring: uncaught errors, unhandled promise rejections and React
 * render errors go to POST /api/v1/client-errors, where they are grouped and
 * logged next to server errors. Self-hosted; no third party.
 */

const MAX_PER_PAGE = 20;
const seen = new Set<string>();
let sent = 0;

export const release = import.meta.env.VITE_RELEASE ?? 'dev';

/** Reports an error once per page load (by message), at most 20 per page. */
export function reportError(error: unknown, extra?: { componentStack?: string }): void {
  if (import.meta.env.MODE === 'test') return;
  const err = error instanceof Error ? error : new Error(typeof error === 'string' ? error : 'Non-Error thrown');
  const message = `${err.name}: ${err.message}`.slice(0, 1000);
  if (seen.has(message) || sent >= MAX_PER_PAGE) return;
  seen.add(message);
  sent++;
  const stack = [err.stack ?? '', extra?.componentStack ? `Component stack:${extra.componentStack}` : ''].filter(Boolean).join('\n');
  void globalThis
    .fetch(`${window.location.origin}/api/v1/client-errors`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        message,
        stack: stack.slice(0, 8000),
        url: window.location.href.split('#')[0], // the fragment can hold a sign-in token
        user_agent: navigator.userAgent,
        release,
      }),
      keepalive: true,
      credentials: 'same-origin',
    })
    .catch(() => undefined);
}

/** Catches errors that nothing else handled. Call once at startup. */
export function installErrorMonitoring(): void {
  window.addEventListener('error', (e) => {
    // Resource load errors (img, script) have no error object; skip them.
    if (e.error !== undefined || e.message) reportError(e.error ?? e.message);
  });
  window.addEventListener('unhandledrejection', (e) => reportError(e.reason));
}
