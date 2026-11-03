// Bridge to the sandbox runner: a hidden <iframe sandbox="allow-scripts"> served from a
// separate origin. Without allow-same-origin the frame's origin is opaque ("null"), so
// replies are accepted only when they come from our frame's window with origin "null",
// and every message is validated with zod.
import { type RunRequest, type RunResult, SandboxMessage } from './protocol';

export const SANDBOX_ORIGIN = import.meta.env.VITE_SANDBOX_ORIGIN ?? 'http://localhost:5174';
const LOAD_GRACE_MS = 60_000; // first run may wait for Pyodide to load

export type SandboxJob = Omit<RunRequest, 'type' | 'runId'>;

interface Pending {
  resolve: (r: RunResult) => void;
  reject: (e: Error) => void;
  timer: ReturnType<typeof setTimeout>;
}

let frame: HTMLIFrameElement | undefined;
let ready: Promise<void> | undefined;
let markReady: (() => void) | undefined;
const pending = new Map<string, Pending>();

function onMessage(event: MessageEvent) {
  if (!frame || event.source !== frame.contentWindow || event.origin !== 'null') return;
  const parsed = SandboxMessage.safeParse(event.data);
  if (!parsed.success) return;
  const msg = parsed.data;
  if (msg.type === 'ready') {
    markReady?.();
    return;
  }
  const p = pending.get(msg.runId);
  if (!p) return;
  clearTimeout(p.timer);
  pending.delete(msg.runId);
  p.resolve(msg);
}

/** Creates the sandbox frame (once) so the runtime starts loading before the first Run. */
export function warmSandbox(): Promise<void> {
  if (ready) return ready;
  ready = new Promise<void>((resolve) => (markReady = resolve));
  window.addEventListener('message', onMessage);
  frame = document.createElement('iframe');
  frame.setAttribute('sandbox', 'allow-scripts');
  frame.setAttribute('referrerpolicy', 'no-referrer');
  frame.setAttribute('aria-hidden', 'true');
  frame.title = 'Code runner';
  frame.tabIndex = -1;
  frame.style.cssText = 'position:absolute;width:0;height:0;border:0;visibility:hidden';
  frame.src = `${SANDBOX_ORIGIN}/`;
  document.body.appendChild(frame);
  return ready;
}

/** Runs a job in the sandbox. Rejects only if the sandbox never answers. */
export async function runInSandbox(job: SandboxJob): Promise<RunResult> {
  await warmSandbox();
  const runId = crypto.randomUUID();
  const req: RunRequest = { type: 'run', runId, ...job };
  return new Promise<RunResult>((resolve, reject) => {
    const timer = setTimeout(() => {
      pending.delete(runId);
      reject(new Error('The code runner did not respond. Reload the page and try again.'));
    }, job.timeoutMs + LOAD_GRACE_MS);
    pending.set(runId, { resolve, reject, timer });
    // '*' because the frame's origin is opaque; see SPEC "Message protocol".
    frame?.contentWindow?.postMessage(req, '*');
  });
}
