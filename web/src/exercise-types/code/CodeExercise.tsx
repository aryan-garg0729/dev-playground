import { lazy, Suspense, useState } from 'react';
import { ApiError } from '../../api/client';
import { Button } from '../../components/Button';
import { CheckList } from '../../components/CheckList';
import { Markdown } from '../../components/Markdown';
import { type Example, ProblemExamples } from '../../components/ProblemExample';
import type { ExerciseRendererProps } from '../types';

const CodeEditor = lazy(() => import('./CodeEditor'));

export interface CodePublicView {
  language: string;
  runner: string;
  prompt: string;
  entry: string;
  starter: string;
  timeout_ms: number;
  examples: Example[];
  constraints: string[];
  hints: string[];
  visible_tests: { file: string; source: string }[];
}

export interface CodeAnswer {
  files: Record<string, string>;
}

function Hints({ hints }: { hints: string[] }) {
  const [shown, setShown] = useState(0);
  if (hints.length === 0) return null;
  return (
    <div className="flex flex-col gap-2">
      {hints.slice(0, shown).map((h, i) => (
        <p key={i} className="text-body rounded-md border border-border bg-surface-raised px-4 py-2 text-text-secondary">
          Hint {i + 1} · {h}
        </p>
      ))}
      {shown < hints.length && (
        <Button onClick={() => setShown(shown + 1)} className="w-full">
          {shown === 0 ? 'Show a hint' : 'Next hint'}
        </Button>
      )}
    </div>
  );
}

export function CodeExercise({ header, publicView: pv, draft, onChange, onSubmit, result }: ExerciseRendererProps<CodePublicView, CodeAnswer>) {
  const code = draft?.files[pv.entry] ?? pv.starter;
  const [pending, setPending] = useState(false);
  const [error, setError] = useState<string>();
  const errorLine = result?.result.checks.find((c) => !c.passed && (c.line ?? 0) > 0)?.line;

  async function submit() {
    setPending(true);
    setError(undefined);
    try {
      await onSubmit({ files: { [pv.entry]: code } });
    } catch (e) {
      setError(e instanceof ApiError ? e.message : 'Could not submit. Check your connection and try again.');
    } finally {
      setPending(false);
    }
  }

  return (
    <div className="flex min-h-0 min-w-0 flex-1 gap-2 max-[900px]:flex-col">
      <section
        aria-label="Problem"
        className="flex w-[320px] shrink-0 flex-col gap-4 overflow-y-auto rounded-lg border border-border bg-surface px-6 py-5 max-[900px]:w-auto"
      >
        {header}
        <Markdown>{pv.prompt}</Markdown>
        <ProblemExamples examples={pv.examples} constraints={pv.constraints} />
        <div className="mt-auto flex flex-col gap-3 pt-4">
          <Hints hints={pv.hints} />
          <Button variant="tutor" disabled title="The AI tutor arrives in a later release">
            Ask the tutor
          </Button>
        </div>
      </section>

      <div className="flex min-h-0 min-w-0 flex-1 flex-col gap-2">
        <section aria-label="Editor" className="flex min-h-[320px] flex-1 flex-col overflow-hidden rounded-lg border border-border bg-surface">
          <div className="flex h-[52px] shrink-0 items-center gap-2 border-b border-border pl-4 pr-2">
            <span className="text-code-sm flex-1 truncate text-text">{pv.entry}</span>
            <Button disabled title="Run in the browser arrives with the sandbox (Phase 3)">
              Run
            </Button>
            <Button variant="primary" onClick={submit} disabled={pending}>
              {pending ? 'Checking…' : 'Submit'}
            </Button>
          </div>
          <div className="min-h-0 flex-1">
            <Suspense fallback={<p className="text-small p-4 text-text-muted">Loading editor…</p>}>
              <CodeEditor
                value={code}
                language={pv.language}
                errorLine={errorLine}
                ariaLabel={`${pv.entry} editor`}
                onChange={(v) => onChange({ files: { [pv.entry]: v } })}
              />
            </Suspense>
          </div>
        </section>
        {error && (
          <p role="alert" className="text-small rounded-lg border border-border bg-danger-soft px-4 py-2 text-danger">
            {error}
          </p>
        )}
        <div aria-live="polite" className="flex max-h-[45%] shrink-0 flex-col gap-2 overflow-y-auto">
          {result && <CheckList result={result.result} />}
          {result && result.result.stdout !== '' && (
            <details className="rounded-lg border border-border bg-surface px-4 py-[14px]">
              <summary className="text-section cursor-pointer text-text">Output</summary>
              <pre className="text-code-sm mt-2 whitespace-pre-wrap text-text-secondary [overflow-wrap:anywhere]">{result.result.stdout}</pre>
            </details>
          )}
        </div>
      </div>
    </div>
  );
}
