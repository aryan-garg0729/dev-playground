import type { Result } from '../../api/client';
import { CheckList } from '../../components/CheckList';
import { ResultTable } from '../../components/ResultTable';
import type { GradeResult, RunOutput } from '../types';

function runAsResult(out: RunOutput): Result {
  const failing = out.checks.some((c) => !c.passed);
  const status = out.timedOut ? 'timeout' : out.error && out.checks.length === 0 ? 'error' : failing ? 'fail' : 'pass';
  return { status, score: 0, checks: out.checks, stdout: out.stdout, runtime_ms: out.runtimeMs };
}

function Output({ label, text }: { label: string; text: string }) {
  return (
    <details open className="rounded-lg border border-border bg-surface px-4 py-[14px]">
      <summary className="text-section cursor-pointer text-text">{label}</summary>
      <pre className="text-code-sm mt-2 whitespace-pre-wrap text-text-secondary [overflow-wrap:anywhere]">{text}</pre>
    </details>
  );
}

/** What the panel under the editor shows: the last Run, or the last Submit. */
export type Shown = { kind: 'run'; out: RunOutput } | { kind: 'submit'; graded: GradeResult } | undefined;

export function Results({ shown }: { shown: Shown }) {
  if (!shown) return null;
  if (shown.kind === 'submit') {
    const r = shown.graded.result;
    return (
      <>
        <CheckList result={r} />
        {r.stdout !== '' && <Output label="Output" text={r.stdout} />}
      </>
    );
  }
  const out = shown.out;
  return (
    <>
      <p className="text-caption text-text-muted">Run · visible checks only, not graded</p>
      {out.error && (
        <p role="status" className="text-code-sm rounded-lg border border-border bg-danger-soft px-4 py-2 text-danger">
          {out.error}
        </p>
      )}
      {out.checks.length > 0 && <CheckList result={runAsResult(out)} />}
      {out.table && <ResultTable table={out.table} />}
      {out.stdout !== '' && <Output label="Output" text={out.stdout} />}
      {out.stderr !== '' && <Output label="Errors" text={out.stderr} />}
      {!out.error && out.checks.length === 0 && !out.table && out.stdout === '' && out.stderr === '' && (
        <p className="text-small text-text-muted">Ran in {out.runtimeMs} ms with no output.</p>
      )}
    </>
  );
}
