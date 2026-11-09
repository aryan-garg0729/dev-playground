import type { ReactNode } from 'react';
import { Link, useNavigate, useSearchParams } from 'react-router';
import { browserTimezone, useDevSignIn, useProviders } from '../api/queries';
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

/** Sign-in failures come back from the provider as ?error=github or ?error=google. */
const errors: Record<string, string> = {
  github: 'GitHub sign-in did not complete. Try again.',
  google: 'Google sign-in did not complete. Try again.',
};

export function SignInPage() {
  const [params] = useSearchParams();
  const next = safeNext(params.get('next'));
  const providers = useProviders();
  const dev = useDevSignIn();
  const navigate = useNavigate();
  const error = errors[params.get('error') ?? ''];
  const start = (p: string) =>
    `/api/v1/auth/${p}/start?${new URLSearchParams({ next, tz: browserTimezone() }).toString()}`;
  const none = providers.data && !providers.data.github && !providers.data.google && !providers.data.dev;

  return (
    <AuthShell>
      <div className="flex flex-col gap-2">
        <h1 className="text-page-title text-text">Sign in to CodeQuest</h1>
        <p className="text-body text-text-secondary">Use your Google or GitHub account. No password needed.</p>
      </div>
      {providers.data?.google && (
        <a href={start('google')} className="text-control flex h-10 items-center justify-center rounded-md bg-accent text-accent-ink">
          Continue with Google
        </a>
      )}
      {providers.data?.github && (
        <a href={start('github')} className="text-control flex h-10 items-center justify-center rounded-md bg-surface-raised text-text">
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
      {none && <p className="text-small text-text-muted">Sign-in is not set up on this server yet.</p>}
      {error && (
        <p role="alert" className="text-small text-danger">
          {error}
        </p>
      )}
    </AuthShell>
  );
}
