import { Markdown } from '../../components/Markdown';
import type { SqlCheck } from '../../runner-client/protocol';
import { Workspace } from '../shared/Workspace';
import type { ExerciseRendererProps, RunMode, ToSandboxJob } from '../types';

export interface SqlPublicView {
  prompt: string;
  setup: string;
  starter: string;
  run: RunMode;
  timeout_ms: number;
  visible_checks: SqlCheck[];
}

export interface SqlAnswer {
  query: string;
}

/** Browser Run: the query in sql.js on the public setup data. */
export const sqlSandboxJob: ToSandboxJob<SqlPublicView, SqlAnswer> = (pv, answer) => ({
  runner: 'sqljs',
  entry: 'query.sql',
  files: { 'query.sql': answer.query },
  tests: [],
  sql: { setup: pv.setup, checks: pv.visible_checks },
  timeoutMs: pv.timeout_ms,
});

export function SqlExercise({ header, publicView: pv, draft, onChange, onRun, onSubmit, result }: ExerciseRendererProps<SqlPublicView, SqlAnswer>) {
  const query = draft?.query ?? pv.starter;
  return (
    <Workspace
      header={header}
      problem={
        <>
          <Markdown>{pv.prompt}</Markdown>
          <details className="flex flex-col gap-2">
            <summary className="text-section cursor-pointer text-text">Tables and data</summary>
            <pre className="text-code-sm mt-2 overflow-x-auto border-l-2 border-border-strong pl-[14px] text-text-secondary">{pv.setup}</pre>
          </details>
        </>
      }
      file="query.sql"
      language="sql"
      code={query}
      onCode={(v) => onChange({ query: v })}
      onRun={onRun && (() => onRun({ query }))}
      onSubmit={() => onSubmit({ query })}
      result={result}
    />
  );
}
