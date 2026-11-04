import { useEffect, useRef, useState } from 'react';
import { api, ApiError, unwrap } from '../api/client';
import { Button } from '../components/Button';
import { AuthShell } from './SignInPage';

interface Code {
  code: string;
  expiresAt: number; // ms since epoch
}

function remaining(expiresAt: number, now: number): string {
  const s = Math.max(0, Math.round((expiresAt - now) / 1000));
  return `${Math.floor(s / 60)}:${String(s % 60).padStart(2, '0')}`;
}

/**
 * /cli/login: shows a signed-in learner a one-time code to paste into
 * `codequest login`. Single use, 10 minutes (the server enforces both).
 */
export function CliLoginPage() {
  const [code, setCode] = useState<Code>();
  const [error, setError] = useState<string>();
  const [busy, setBusy] = useState(false);
  const [copied, setCopied] = useState(false);
  const [now, setNow] = useState(() => Date.now());
  const started = useRef(false);

  async function create() {
    setBusy(true);
    setError(undefined);
    setCopied(false);
    try {
      const c = await unwrap(api.POST('/cli/codes', { body: {} }));
      setCode({ code: c.code, expiresAt: Date.parse(c.expires_at) });
      setNow(Date.now());
    } catch (e) {
      setError(e instanceof ApiError ? e.message : 'Could not get a code. Check your connection and try again.');
    } finally {
      setBusy(false);
    }
  }

  useEffect(() => {
    // Once per visit (StrictMode runs effects twice in development).
    if (started.current) return;
    started.current = true;
    void create();
  }, []);

  useEffect(() => {
    if (!code) return;
    const t = setInterval(() => setNow(Date.now()), 1000);
    return () => clearInterval(t);
  }, [code]);

  async function copy() {
    if (!code) return;
    try {
      await navigator.clipboard.writeText(code.code);
      setCopied(true);
    } catch {
      setCopied(false);
    }
  }

  const expired = code !== undefined && now >= code.expiresAt;

  return (
    <AuthShell>
      <div className="flex flex-col gap-2">
        <h1 className="text-page-title text-text">Sign in to the CLI</h1>
        <p className="text-body text-text-secondary">
          Paste this code into your terminal where <code className="font-mono text-text">codequest login</code> is waiting.
        </p>
      </div>

      {code && !expired && (
        <section aria-label="Login code" className="flex flex-col gap-3 rounded-lg border border-border bg-surface px-5 py-4">
          <p data-testid="cli-code" className="select-all text-center font-mono text-[28px] leading-[40px] tracking-[0.12em] text-text">
            {code.code}
          </p>
          <div className="flex items-center justify-between gap-3">
            <p className="text-small text-text-muted" aria-live="polite">
              Expires in {remaining(code.expiresAt, now)} · works once
            </p>
            <Button onClick={copy}>{copied ? 'Copied' : 'Copy'}</Button>
          </div>
        </section>
      )}

      {expired && (
        <section className="flex flex-col gap-3 rounded-lg border border-border bg-surface px-5 py-4">
          <p className="text-body text-text-secondary">This code expired.</p>
          <Button variant="primary" onClick={create} disabled={busy}>
            Get a new code
          </Button>
        </section>
      )}

      {!code && busy && <p className="text-body text-text-muted">Getting a code…</p>}

      {error && (
        <div className="flex flex-col gap-3">
          <p role="alert" className="text-small text-danger">
            {error}
          </p>
          <Button onClick={create} disabled={busy}>
            Try again
          </Button>
        </div>
      )}

      <p className="text-small text-text-muted">
        Didn't start from the terminal? Install the CLI and run <code className="font-mono text-text">codequest login</code>. Never
        share this code: it signs the CLI in as you.
      </p>
    </AuthShell>
  );
}
