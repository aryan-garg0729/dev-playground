import { useQuery } from '@tanstack/react-query';
import { api, ApiError, type Schemas, unwrap } from './client';

export type TutorThread = Schemas['TutorThread'];
export type TutorMessage = Schemas['TutorMessage'];

export const tutorKey = (exerciseId: string) => ['tutor', exerciseId] as const;

export function useTutorThread(exerciseId: string, enabled: boolean) {
  return useQuery({
    queryKey: tutorKey(exerciseId),
    queryFn: () => unwrap(api.GET('/exercises/{exerciseId}/tutor', { params: { path: { exerciseId } } })),
    enabled,
  });
}

interface SSEEvent {
  event: string;
  data: string;
}

/** Splits a text/event-stream into events; `rest` is an unfinished tail. */
export function parseSSE(buf: string): { events: SSEEvent[]; rest: string } {
  const parts = buf.split('\n\n');
  const rest = parts.pop() ?? '';
  const events: SSEEvent[] = [];
  for (const block of parts) {
    let event = 'message';
    const data: string[] = [];
    for (const line of block.split('\n')) {
      if (line.startsWith('event: ')) event = line.slice(7);
      else if (line.startsWith('data: ')) data.push(line.slice(6));
    }
    events.push({ event, data: data.join('\n') });
  }
  return { events, rest };
}

/**
 * Asks the tutor and streams the reply: onDelta gets each piece of text.
 * Resolves with today's hints left; rejects with ApiError (limit, validation,
 * or a failed reply, which the server refunds).
 */
export async function askTutor(
  exerciseId: string,
  body: { message: string; code: string },
  onDelta: (text: string) => void,
  signal?: AbortSignal,
): Promise<{ remaining: number }> {
  const res = await globalThis.fetch(`${window.location.origin}/api/v1/exercises/${encodeURIComponent(exerciseId)}/hint`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', Accept: 'text/event-stream' },
    body: JSON.stringify(body),
    credentials: 'same-origin',
    signal,
  });
  if (!res.ok || !res.headers.get('Content-Type')?.startsWith('text/event-stream') || !res.body) {
    const env = (await res.json().catch(() => undefined)) as Schemas['Error'] | undefined;
    throw new ApiError(res.status, env?.error?.code ?? 'http_error', env?.error?.message ?? `Request failed (${res.status})`);
  }
  const reader = res.body.pipeThrough(new TextDecoderStream()).getReader();
  let buf = '';
  for (;;) {
    const { value, done } = await reader.read();
    if (done) break;
    const { events, rest } = parseSSE(buf + value);
    buf = rest;
    for (const ev of events) {
      const data = JSON.parse(ev.data) as Record<string, unknown>;
      if (ev.event === 'delta') onDelta(String(data.text ?? ''));
      else if (ev.event === 'done') return { remaining: Number(data.remaining ?? 0) };
      else if (ev.event === 'error') throw new ApiError(502, String(data.code ?? 'tutor_failed'), String(data.message ?? 'The tutor could not answer.'));
    }
  }
  throw new ApiError(502, 'tutor_failed', 'The tutor stopped before finishing. Try again.');
}
