import type { Result } from '../api/client';

/**
 * The result panel under the editor: a summary line and one row per check.
 * Everything here is grader output about learner code, so it is rendered as
 * text nodes only.
 */
export function CheckList({ result }: { result: Result }) {
  const total = result.checks.length;
  const failing = result.checks.filter((c) => !c.passed).length;
  const ok = failing === 0 && result.status === 'pass';
  const summary =
    result.status === 'timeout'
      ? 'Time limit exceeded'
      : result.status === 'error'
        ? 'Could not grade this submission'
        : ok
          ? `${total} of ${total} checks passing`
          : `${failing} of ${total} checks failing`;

  return (
    <section aria-label="Checks" className="rounded-lg border border-border bg-surface px-4 py-[14px]">
      <div className="flex items-baseline justify-between gap-4">
        <h2 className={`text-section ${ok ? 'text-success' : 'text-danger'}`}>{summary}</h2>
        {result.runtime_ms > 0 && <p className="text-caption shrink-0 text-text-muted">{result.runtime_ms} ms</p>}
      </div>
      <ul className="mt-2 flex flex-col gap-1.5">
        {result.checks.map((c, i) => (
          <li key={i} className="text-code-sm flex items-start gap-3 text-text">
            <span
              role="img"
              aria-label={c.passed ? 'Passed' : 'Failed'}
              className={`mt-[8px] h-2 w-2 shrink-0 rounded-pill ${c.passed ? 'bg-success' : 'bg-danger'}`}
            />
            <span className="min-w-0 [overflow-wrap:anywhere]">
              {c.name}
              {c.got !== undefined && (
                <>
                  {' → '}
                  <span className="text-danger">{c.got}</span>
                </>
              )}
              {c.expected !== undefined && <span className="text-text-muted"> (expected {c.expected})</span>}
              {c.feedback && (c.passed || c.got === undefined) && (
                <span className={c.passed ? 'text-text-muted' : 'text-danger'}> · {c.feedback}</span>
              )}
            </span>
          </li>
        ))}
      </ul>
    </section>
  );
}
