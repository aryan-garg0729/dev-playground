import { useState } from 'react';
import { Link, useNavigate } from 'react-router';
import { useVerifyMagicLink } from '../api/queries';
import { Button } from '../components/Button';
import { AuthShell, safeNext } from './SignInPage';

/**
 * Landing page of a magic link (/sign-in/verify#token=...). The token is in the
 * fragment so it never reaches a server log; it is read once and removed from
 * the address bar. Signing in takes a click, so mail scanners that open links
 * cannot use the token up.
 */
export function VerifyPage() {
  const [token] = useState(() => {
    const t = new URLSearchParams(window.location.hash.slice(1)).get('token') ?? '';
    if (window.location.hash) window.history.replaceState(null, '', window.location.pathname);
    return t;
  });
  const verify = useVerifyMagicLink();
  const navigate = useNavigate();

  if (!token) {
    return (
      <AuthShell>
        <h1 className="text-page-title text-text">This link is incomplete</h1>
        <p className="text-body text-text-secondary">Open the link from your email again, or ask for a new one.</p>
        <Link to="/sign-in" className="text-small text-accent-text underline">
          Back to sign in
        </Link>
      </AuthShell>
    );
  }
  return (
    <AuthShell>
      <h1 className="text-page-title text-text">Finish signing in</h1>
      <p className="text-body text-text-secondary">Continue to CodeQuest on this device.</p>
      <Button
        variant="primary"
        className="h-10"
        disabled={verify.isPending || verify.isSuccess}
        onClick={() => verify.mutate(token, { onSuccess: (r) => navigate(safeNext(r.next), { replace: true }) })}
      >
        {verify.isPending ? 'Signing in…' : 'Continue'}
      </Button>
      {verify.error && (
        <p role="alert" className="text-small text-danger">
          {verify.error.message}{' '}
          <Link to="/sign-in" className="text-accent-text underline">
            Get a new link
          </Link>
        </p>
      )}
    </AuthShell>
  );
}
