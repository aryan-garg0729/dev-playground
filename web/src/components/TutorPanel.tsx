import { useQueryClient } from '@tanstack/react-query';
import { type FormEvent, type KeyboardEvent, useEffect, useRef, useState } from 'react';
import { ApiError } from '../api/client';
import { useMe } from '../api/queries';
import { askTutor, type TutorMessage, tutorKey, useTutorThread } from '../api/tutor';
import { track } from '../analytics';
import { Button } from './Button';
import { Markdown } from './Markdown';
import { Spinner } from './Spinner';

function Bubble({ m }: { m: Pick<TutorMessage, 'role' | 'content'> }) {
  if (m.role === 'user') {
    return (
      <p className="text-small ml-6 self-end whitespace-pre-wrap rounded-md bg-surface-raised px-3 py-2 text-text">{m.content}</p>
    );
  }
  return (
    <div className="mr-6 border-l-2 border-accent-line pl-3">
      <Markdown>{m.content}</Markdown>
    </div>
  );
}

function Conversation({ exerciseId, code, onClose }: { exerciseId: string; code: string; onClose: () => void }) {
  const qc = useQueryClient();
  const thread = useTutorThread(exerciseId, true);
  const [question, setQuestion] = useState('');
  const [asked, setAsked] = useState('');
  const [streaming, setStreaming] = useState<string>();
  const [error, setError] = useState<string>();
  const abort = useRef<AbortController>(undefined);
  const end = useRef<HTMLDivElement>(null);
  useEffect(() => () => abort.current?.abort(), []);
  // Braces matter: scrollIntoView returns a Promise in current Chromium, and an
  // effect must return nothing or a cleanup function.
  useEffect(() => {
    end.current?.scrollIntoView?.({ block: 'nearest' });
  }, [streaming, thread.data]);

  const busy = streaming !== undefined;
  const data = thread.data;
  const remaining = data?.remaining ?? 0;

  async function send(e?: FormEvent) {
    e?.preventDefault();
    const message = question.trim();
    if (!message || busy) return;
    setError(undefined);
    setAsked(message);
    setQuestion('');
    setStreaming('');
    track('tutor_asked', { exercise: exerciseId });
    abort.current = new AbortController();
    try {
      await askTutor(exerciseId, { message, code }, (t) => setStreaming((s) => (s ?? '') + t), abort.current.signal);
    } catch (err) {
      if (abort.current.signal.aborted) return;
      setError(err instanceof ApiError ? err.message : 'The tutor could not answer. Try again.');
      if (!(err instanceof ApiError && err.code === 'limit_reached')) setQuestion(message);
    } finally {
      await qc.invalidateQueries({ queryKey: tutorKey(exerciseId) });
      setStreaming(undefined);
      setAsked('');
    }
  }

  function onKey(e: KeyboardEvent<HTMLTextAreaElement>) {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      void send();
    }
  }

  return (
    <section aria-label="AI tutor" className="flex flex-col gap-3 rounded-md border border-dashed border-accent-line p-3">
      <div className="flex items-center justify-between gap-2">
        <h2 className="text-small font-semibold text-accent-text">Tutor</h2>
        <span className="text-caption text-text-muted">
          {data ? `${remaining} of ${data.daily_limit} hints left today` : ''}
        </span>
        <button type="button" onClick={onClose} aria-label="Close the tutor" className="text-caption text-text-muted hover:text-text">
          Close
        </button>
      </div>
      {thread.isPending && <Spinner size="sm" label="Loading the tutor…" />}
      {thread.error && <p className="text-caption text-danger" role="alert">{thread.error.message}</p>}
      {data && data.messages.length === 0 && !busy && (
        <p className="text-caption text-text-muted">
          The tutor asks questions and points at the problem; it won't write the answer. Asking removes this lesson's
          no-hint bonus (+25% XP).
        </p>
      )}
      <div className="flex max-h-80 flex-col gap-3 overflow-y-auto" aria-live="polite">
        {data?.messages.map((m, i) => <Bubble key={i} m={m} />)}
        {busy && asked && <Bubble m={{ role: 'user', content: asked }} />}
        {busy && (streaming ? <Bubble m={{ role: 'assistant', content: streaming }} /> : <p className="text-caption text-text-muted">Thinking…</p>)}
        <div ref={end} />
      </div>
      {error && <p className="text-caption text-danger" role="alert">{error}</p>}
      <form onSubmit={send} className="flex flex-col gap-2">
        <textarea
          aria-label="Ask the tutor"
          value={question}
          onChange={(e) => setQuestion(e.target.value)}
          onKeyDown={onKey}
          maxLength={2000}
          rows={2}
          disabled={busy || (data !== undefined && remaining === 0)}
          placeholder={data && remaining === 0 ? 'No hints left today' : 'What are you stuck on?'}
          className="text-small resize-y rounded-md border border-border bg-surface-raised px-3 py-2 text-text placeholder:text-text-faint"
        />
        <Button type="submit" disabled={busy || !question.trim() || remaining === 0}>
          {busy ? 'Answering…' : 'Ask'}
        </Button>
      </form>
    </section>
  );
}

/** "Ask the tutor" at the bottom of the problem pane. Off unless the server enables it (development only for now). */
export function TutorPanel({ exerciseId, code }: { exerciseId: string; code: string }) {
  const me = useMe();
  const [open, setOpen] = useState(false);
  if (!me.data?.features.tutor) {
    return (
      <Button variant="tutor" disabled title="Coming soon">
        Ask the tutor
      </Button>
    );
  }
  if (!open) {
    return (
      <Button
        variant="tutor"
        onClick={() => {
          setOpen(true);
          track('tutor_opened', { exercise: exerciseId });
        }}
      >
        Ask the tutor
      </Button>
    );
  }
  return <Conversation exerciseId={exerciseId} code={code} onClose={() => setOpen(false)} />;
}
