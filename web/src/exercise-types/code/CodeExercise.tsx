import { useState } from 'react';
import { Button } from '../../components/Button';
import { Markdown } from '../../components/Markdown';
import { type Example, ProblemExamples } from '../../components/ProblemExample';
import { Workspace } from '../shared/Workspace';
import type { ExerciseRendererProps, RunMode, ToSandboxJob } from '../types';

export interface CodePublicView {
  language: string;
  runner: string;
  run: RunMode;
  prompt: string;
  entry: string;
  starter: string;
  timeout_ms: number;
  examples: Example[];
  constraints: string[];
  hints: string[];
  visible_tests: { file: string; source: string }[];
  /** How many result rows grading produces: a boss's HP. */
  check_counts?: { visible: number; hidden: number };
}

export interface CodeAnswer {
  files: Record<string, string>;
}

/** Browser Run: the entry plus the visible tests in Pyodide. */
export const codeSandboxJob: ToSandboxJob<CodePublicView, CodeAnswer> = (pv, answer) => ({
  runner: 'pyodide',
  entry: pv.entry,
  files: { [pv.entry]: answer.files[pv.entry] ?? pv.starter },
  tests: pv.visible_tests,
  timeoutMs: pv.timeout_ms,
});

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

export function CodeExercise(props: ExerciseRendererProps<CodePublicView, CodeAnswer>) {
  const { header, publicView: pv, draft, onChange, onRun, onSubmit, result, aside, boss, tutor } = props;
  const code = draft?.files[pv.entry] ?? pv.starter;
  const counts = pv.check_counts;
  const answer: CodeAnswer = { files: { [pv.entry]: code } };
  return (
    <Workspace
      header={header}
      problem={
        <>
          <Markdown>{pv.prompt}</Markdown>
          <ProblemExamples examples={pv.examples} constraints={pv.constraints} />
        </>
      }
      problemFooter={
        <>
          <Hints hints={pv.hints} />
          {tutor?.(code)}
        </>
      }
      file={pv.entry}
      language={pv.language}
      code={code}
      onCode={(v) => onChange({ files: { [pv.entry]: v } })}
      onRun={onRun && (() => onRun(answer))}
      onSubmit={() => onSubmit(answer)}
      result={result}
      aside={aside}
      boss={boss && counts ? { ...boss, maxHP: counts.visible + counts.hidden } : undefined}
    />
  );
}
