import { type FormEvent, type ReactNode, useState } from 'react';
import { Link, useNavigate, useSearchParams } from 'react-router';
import { ApiError } from '../api/client';
import { sendMagicLink, useDevSignIn, useProviders } from '../api/queries';
import { Button } from '../components/Button';
import { CodeIcon } from '../components/icons';

/** Keeps only same-site paths (the server checks again). */
export function safeNext(next: string | null): string {
  return next && next.startsWith('/') && !next.startsWith('//') && !next.includes('\\') ? next : '/';
}

export function AuthShell({ children }: { children: ReactNode }) {
  return (
    <div className="min-h-dvh bg-bg">
      <header className="flex h-[var(--header-height)] items-center px-5">
        <Link to="/" className="flex items-center gap-2" aria-label="CodeQuest home">
          <span className="flex h-7 w-7 items-center justify-center rounded-md bg-accent text-accent-ink">
            <CodeIcon />
          </span>
          <span className="text-brand text-text">CodeQuest</span>
        </Link>
      </header>
      <main className="mx-auto flex max-w-[420px] flex-col gap-6 px-5 pt-16">{children}</main>
    </div>
  );
}

export function SignInPage() {
  const [params] = useSearchParams();
  const next = safeNext(params.get('next'));
  const providers = useProviders();
  const dev = useDevSignIn();
  const navigate = useNavigate();
  const [email, setEmail] = useState('');
  const [sentTo, setSentTo] = useState<string>();
  const [error, setError] = useState<string | undefined>(
    params.get('error') === 'github' ? 'GitHub sign-in did not complete. Try again.' : undefined,
  );
  const [sending, setSending] = useState(false);

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    setError(undefined);
    setSending(true);
    try {
      await sendMagicLink(email, next);
      setSentTo(email.trim());
    } catch (err) {
      setError(err instanceof ApiError ? err.message : 'Could not send the link. Try again.');
    } finally {
      setSending(false);
    }
  }

  if (sentTo) {
    return (
      <AuthShell>
        <h1 className="text-page-title text-text">Check your email</h1>
        <p className="text-body text-text-secondary">
          We sent a sign-in link to <span className="font-semibold text-text">{sentTo}</span>. It works once and expires in 15 minutes.
        </p>
        <Button onClick={() => setSentTo(undefined)}>Use a different email</Button>
      </AuthShell>
    );
  }

  return (
    <AuthShell>
      <div className="flex flex-col gap-2">
        <h1 className="text-page-title text-text">Sign in to CodeQuest</h1>
        <p className="text-body text-text-secondary">We'll email you a link. No password needed.</p>
      </div>
      <form onSubmit={onSubmit} className="flex flex-col gap-3">
        <label className="text-small flex flex-col gap-1.5 text-text-secondary">
          Email
          <input
            type="email"
            required
            autoComplete="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            className="text-body h-10 rounded-md border border-border bg-surface px-3 text-text outline-none focus:border-accent"
          />
        </label>
        <Button variant="primary" type="submit" disabled={sending} className="h-10">
          {sending ? 'Sending…' : 'Email me a link'}
        </Button>
      </form>
      {providers.data?.github && (
        <a
          href={`/api/v1/auth/github/start?next=${encodeURIComponent(next)}`}
          className="text-control flex h-10 items-center justify-center rounded-md bg-surface-raised text-text"
        >
          Continue with GitHub
        </a>
      )}
      {providers.data?.dev && (
        <Button
          className="h-10 border border-dashed border-border"
          disabled={dev.isPending}
          onClick={() => dev.mutate(undefined, { onSuccess: () => navigate(next, { replace: true }) })}
        >
          Dev sign-in (development only)
        </Button>
      )}
      {error && (
        <p role="alert" className="text-small text-danger">
          {error}
        </p>
      )}
    </AuthShell>
  );
}
