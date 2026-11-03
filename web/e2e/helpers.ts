import type { Page } from '@playwright/test';

export const SANDBOX = process.env.SANDBOX_URL ?? 'http://localhost:5174';

/** Puts an unsent answer in place before the lesson loads (the page restores drafts). */
export async function withDraft(page: Page, exerciseId: string, version: number, answer: unknown) {
  await page.addInitScript(
    ([key, value]) => localStorage.setItem(key, value),
    [`cq-draft:${exerciseId}:v${version}`, JSON.stringify(answer)] as const,
  );
}

export function lesson(courseId: string, exerciseId: string) {
  return `/learn/${courseId}/${exerciseId}`;
}

/**
 * Runs a sandbox job from the app's origin through a fresh sandboxed iframe, the
 * same way web/src/runner-client/sandbox.ts does.
 */
export async function runInSandbox(page: Page, job: Record<string, unknown>) {
  return page.evaluate(
    async ([sandbox, job]) => {
      const frame = document.createElement('iframe');
      frame.setAttribute('sandbox', 'allow-scripts');
      frame.src = `${sandbox}/`;
      const ready = new Promise<void>((resolve) => {
        window.addEventListener('message', function on(e) {
          if (e.source === frame.contentWindow && (e.data as { type?: string }).type === 'ready') {
            window.removeEventListener('message', on);
            resolve();
          }
        });
      });
      document.body.appendChild(frame);
      await ready;
      const runId = crypto.randomUUID();
      const result = new Promise<Record<string, unknown>>((resolve) => {
        window.addEventListener('message', function on(e) {
          const d = e.data as { type?: string; runId?: string };
          if (e.source === frame.contentWindow && d.type === 'result' && d.runId === runId) {
            window.removeEventListener('message', on);
            resolve(e.data as Record<string, unknown>);
          }
        });
      });
      frame.contentWindow?.postMessage({ type: 'run', runId, ...job }, '*'); // opaque origin
      const r = await result;
      frame.remove();
      return r;
    },
    [SANDBOX, job] as const,
  );
}
