/**
 * Self-hosted product analytics: events are batched and sent to POST /api/v1/events
 * (no third party, no extra cookies). Visitors with Do Not Track or Global
 * Privacy Control on are not tracked. Event names must be on the server's list.
 */

export type EventName =
  | 'page_view'
  | 'landing_cta'
  | 'sign_in'
  | 'onboarding_step'
  | 'onboarding_completed'
  | 'lesson_opened'
  | 'run_clicked'
  | 'submit_clicked'
  | 'tutor_opened'
  | 'tutor_asked'
  | 'shop_bought';

type Props = Record<string, string | number | boolean>;
interface QueuedEvent {
  name: EventName;
  props: Props;
  path: string;
  at: string;
}

const FLUSH_MS = 5000;
const MAX_BATCH = 20;
let queue: QueuedEvent[] = [];
let timer: ReturnType<typeof setTimeout> | undefined;

function optedOut(): boolean {
  const nav = navigator as Navigator & { globalPrivacyControl?: boolean };
  return nav.doNotTrack === '1' || nav.globalPrivacyControl === true;
}

function enabled(): boolean {
  return import.meta.env.MODE !== 'test' && !optedOut();
}

/** A random id kept in this browser, so signed-out visits can be counted as one visitor. */
function anonId(): string | undefined {
  try {
    let id = localStorage.getItem('cq-anon');
    if (!id) {
      id = crypto.randomUUID();
      localStorage.setItem('cq-anon', id);
    }
    return id;
  } catch {
    return undefined;
  }
}

/** Sends queued events now. keepalive lets the request outlive a closing tab. */
export function flush(): void {
  if (timer) clearTimeout(timer);
  timer = undefined;
  while (queue.length > 0) {
    const events = queue.splice(0, MAX_BATCH);
    const body = JSON.stringify({ anon_id: anonId(), events });
    void globalThis
      .fetch(`${window.location.origin}/api/v1/events`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body,
        keepalive: true,
        credentials: 'same-origin',
      })
      .catch(() => undefined); // analytics never breaks the app
  }
}

/** Records an event; it is sent within a few seconds (or when the tab is hidden). */
export function track(name: EventName, props: Props = {}): void {
  if (!enabled()) return;
  queue.push({ name, props, path: window.location.pathname, at: new Date().toISOString() });
  if (queue.length >= MAX_BATCH) flush();
  else timer ??= setTimeout(flush, FLUSH_MS);
}

if (typeof document !== 'undefined') {
  document.addEventListener('visibilitychange', () => {
    if (document.visibilityState === 'hidden') flush();
  });
}
