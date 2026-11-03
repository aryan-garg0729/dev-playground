import { lazy, type ReactNode, Suspense, useState } from 'react';
import { ApiError } from '../../api/client';
import { Button } from '../../components/Button';
import type { GradeResult, RunOutput } from '../types';
import { Results, type Shown } from './Results';

const CodeEditor = lazy(() => import('../code/CodeEditor'));

interface Props {
  header: ReactNode;
  /** Problem pane content under the header. */
  problem: ReactNode;
  /** Bottom of the problem pane (hints, tutor). */
  problemFooter?: ReactNode;
  file: string;
  language: string;
  code: string;
  onCode: (code: string) => void;
  onRun?: () => Promise<RunOutput>;
  onSubmit: () => Promise<GradeResult>;
  result: GradeResult | undefined;
}

function message(e: unknown, what: string): string {
  if (e instanceof ApiError) return e.message;
  if (e instanceof Error && e.message) return e.message;
  return `Could not ${what}. Check your connection and try again.`;
}

/** The lesson workspace: problem pane | editor with Run and Submit, then results. */
export function Workspace(p: Props) {
  const [busy, setBusy] = useState<'run' | 'submit'>();
  const [error, setError] = useState<string>();
  const [run, setRun] = useState<RunOutput>();
  const [showRun, setShowRun] = useState(false);
  const shown: Shown = showRun && run ? { kind: 'run', out: run } : p.result ? { kind: 'submit', graded: p.result } : undefined;
  const failing = shown?.kind === 'submit' ? shown.graded.result.checks : shown?.out.checks;
  const errorLine = failing?.find((c) => !c.passed && (c.line ?? 0) > 0)?.line;

  async function act(kind: 'run' | 'submit') {
    setBusy(kind);
    setError(undefined);
    try {
      if (kind === 'run' && p.onRun) {
        setRun(await p.onRun());
        setShowRun(true);
      } else {
        await p.onSubmit();
        setShowRun(false);
      }
    } catch (e) {
      setError(message(e, kind));
    } finally {
      setBusy(undefined);
    }
  }

  return (
    <div className="flex min-h-0 min-w-0 flex-1 gap-2 max-[900px]:flex-col">
      <section
        aria-label="Problem"
        className="flex w-[320px] shrink-0 flex-col gap-4 overflow-y-auto rounded-lg border border-border bg-surface px-6 py-5 max-[900px]:w-auto"
      >
        {p.header}
        {p.problem}
        {p.problemFooter && <div className="mt-auto flex flex-col gap-3 pt-4">{p.problemFooter}</div>}
      </section>

      <div className="flex min-h-0 min-w-0 flex-1 flex-col gap-2">
        <section aria-label="Editor" className="flex min-h-[320px] flex-1 flex-col overflow-hidden rounded-lg border border-border bg-surface">
          <div className="flex h-[52px] shrink-0 items-center gap-2 border-b border-border pl-4 pr-2">
            <span className="text-code-sm flex-1 truncate text-text">{p.file}</span>
            <Button
              onClick={() => act('run')}
              disabled={!p.onRun || busy !== undefined}
              title={p.onRun ? 'Run with the visible checks (not graded)' : 'Run is off for this lesson'}
            >
              {busy === 'run' ? 'Running…' : 'Run'}
            </Button>
            <Button variant="primary" onClick={() => act('submit')} disabled={busy !== undefined}>
              {busy === 'submit' ? 'Checking…' : 'Submit'}
            </Button>
          </div>
          <div className="min-h-0 flex-1">
            <Suspense fallback={<p className="text-small p-4 text-text-muted">Loading editor…</p>}>
              <CodeEditor value={p.code} language={p.language} errorLine={errorLine} ariaLabel={`${p.file} editor`} onChange={p.onCode} />
            </Suspense>
          </div>
        </section>
        {error && (
          <p role="alert" className="text-small rounded-lg border border-border bg-danger-soft px-4 py-2 text-danger">
            {error}
          </p>
        )}
        <div aria-live="polite" className="flex max-h-[45%] shrink-0 flex-col gap-2 overflow-y-auto">
          <Results shown={shown} />
        </div>
      </div>
    </div>
  );
}
