import { Link } from 'react-router';
import { ApiError } from '../api/client';

export function Loading({ label }: { label: string }) {
  return (
    <p className="text-small p-6 text-text-muted" role="status">
      Loading {label}…
    </p>
  );
}

export function LoadError({ error }: { error: unknown }) {
  const notFound = error instanceof ApiError && error.status === 404;
  return (
    <div className="flex flex-col items-start gap-2 p-6" role="alert">
      <p className="text-section text-text">{notFound ? 'Not found' : 'Something went wrong'}</p>
      <p className="text-small text-text-muted">
        {notFound ? 'This page does not exist or is no longer published.' : error instanceof Error ? error.message : 'Try again.'}
      </p>
      <Link to="/" className="text-small text-accent-text underline">
        Back to courses
      </Link>
    </div>
  );
}
