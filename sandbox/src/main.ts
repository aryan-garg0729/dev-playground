// The sandbox page. Embedded by the app as <iframe sandbox="allow-scripts"> from its own
// origin, so it has an opaque origin: no cookies, no storage, no access to the app's DOM.
// Learner code runs only in Web Workers started from here.
import { RunRequest, type RunResult } from './protocol';
import RunWorker from './worker.ts?worker&inline';

const APP_ORIGINS = (import.meta.env.VITE_APP_ORIGINS ?? 'http://localhost:5173')
  .split(',')
  .map((s: string) => s.trim())
  .filter(Boolean);

// Absolute base for assets: the worker starts from a blob URL and cannot resolve relative paths.
const BASE = new URL('./', location.href).href;
const START_TIMEOUT_MS = 60_000;

type Runner = RunRequest['runner'];

interface Slot {
  worker: Worker;
  ready: Promise<void>;
}

function spawn(warm: Runner[]): Slot {
  const worker = new RunWorker();
  const ready = new Promise<void>((resolve, reject) => {
    const timer = setTimeout(() => reject(new Error('runtime did not start')), START_TIMEOUT_MS);
    worker.addEventListener('message', function onMsg(e: MessageEvent) {
      const d = e.data as { type?: string; error?: string };
      if (d.type === 'ready' || d.type === 'failed') {
        clearTimeout(timer);
        worker.removeEventListener('message', onMsg);
        if (d.type === 'ready') resolve();
        else reject(new Error(d.error ?? 'runtime failed to start'));
      }
    });
  });
  ready.catch(() => {}); // surfaced when a run awaits it
  worker.postMessage({ type: 'init', base: BASE, warm });
  return { worker, ready };
}

const WARM: Runner[] = ['pyodide', 'sqljs'];
let active = spawn(WARM);
let spare: Slot | undefined;
// Start the spare once the first runtime is up, so the two do not compete for the first load.
void active.ready.then(() => (spare ??= spawn(WARM)), () => {});

function replaceActive() {
  active.worker.terminate();
  active = spare ?? spawn(WARM);
  spare = undefined;
  void active.ready.then(() => (spare ??= spawn(WARM)), () => {});
}

function fail(req: RunRequest, error: string, timedOut = false, runtimeMs = 0): RunResult {
  return { type: 'result', runId: req.runId, stdout: '', stderr: '', checks: [], runtimeMs, timedOut, error };
}

async function execute(req: RunRequest): Promise<RunResult> {
  const slot = active;
  try {
    await slot.ready;
  } catch (e) {
    replaceActive();
    return fail(req, e instanceof Error ? e.message : 'runtime failed to start');
  }
  return new Promise<RunResult>((resolve) => {
    const started = performance.now();
    const timer = setTimeout(() => {
      slot.worker.removeEventListener('message', onMsg);
      replaceActive(); // kills the learner code, whatever it is doing
      resolve(fail(req, `Stopped after ${req.timeoutMs} ms`, true, Math.round(performance.now() - started)));
    }, req.timeoutMs);
    function onMsg(e: MessageEvent) {
      const d = e.data as Partial<RunResult>;
      if (d.type !== 'result' || d.runId !== req.runId) return;
      clearTimeout(timer);
      slot.worker.removeEventListener('message', onMsg);
      resolve(d as RunResult);
    }
    slot.worker.addEventListener('message', onMsg);
    slot.worker.postMessage({ type: 'run', req });
  });
}

// One run at a time, in arrival order.
let queue: Promise<unknown> = Promise.resolve();

window.addEventListener('message', (event: MessageEvent) => {
  if (!APP_ORIGINS.includes(event.origin) || event.source !== window.parent) return;
  const parsed = RunRequest.safeParse(event.data);
  if (!parsed.success) return; // unknown or malformed: dropped
  const req = parsed.data;
  const origin = event.origin;
  queue = queue.then(async () => {
    const res = await execute(req);
    window.parent.postMessage(res, origin);
  });
});

for (const origin of APP_ORIGINS) window.parent.postMessage({ type: 'ready' }, origin);
