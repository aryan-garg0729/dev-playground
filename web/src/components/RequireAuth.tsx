import type { ReactNode } from 'react';
import { Navigate, useLocation } from 'react-router';
import { ApiError } from '../api/client';
import { useMe } from '../api/queries';
import { LoadError, Loading } from '../pages/states';

/** Renders children for a signed-in user; otherwise sends them to sign in and back. */
export function RequireAuth({ children }: { children: ReactNode }) {
  const me = useMe();
  const loc = useLocation();
  if (me.isPending) return <Loading label="your account" />;
  if (me.error instanceof ApiError && me.error.status === 401) {
    const next = loc.pathname + loc.search;
    return <Navigate to={next === '/' ? '/sign-in' : `/sign-in?next=${encodeURIComponent(next)}`} replace />;
  }
  if (me.error) return <LoadError error={me.error} />;
  return <>{children}</>;
}
