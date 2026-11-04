import { useState } from 'react';
import { useExercise } from '../../api/queries';
import { Button } from '../../components/Button';
import { CheckList } from '../../components/CheckList';
import { Markdown } from '../../components/Markdown';
import type { ExerciseRendererProps } from '../types';
import { type CliPublicView, stepAction, stepTests } from './steps';

/** The CLI posts the answer; the page never builds one. */
type CliAnswer = unknown;

function CopyCommand({ command }: { command: string }) {
  const [copied, setCopied] = useState(false);
  return (
    <div className="flex items-center justify-between gap-3 rounded-md border border-border bg-surface-raised px-3 py-2">
      <code className="text-code-sm min-w-0 font-mono text-text [overflow-wrap:anywhere]">{command}</code>
      <Button
        onClick={() => {
          navigator.clipboard.writeText(command).then(
            () => setCopied(true),
            () => setCopied(false),
          );
        }}
      >
        {copied ? 'Copied' : 'Copy'}
      </Button>
    </div>
  );
}

/**
 * A command-line lesson: the steps run on the learner's machine through the
 * `codequest` CLI, which submits their output. This page shows the prompt,
 * the steps and tests, the commands to run, and the latest result.
 */
export function CliExercise({ header, publicView, result, aside }: ExerciseRendererProps<CliPublicView, CliAnswer>) {
  const { lesson_id: id, prompt, steps, base_url: baseURL } = publicView;
  // Same query as the page: refetching picks up a submit made from the terminal.
  const ex = useExercise(id);
  const last = ex.data?.last_submission;
  const latest = last?.result && last.exercise_version === ex.data?.version ? last.result : result?.result;

  return (
    <div className="flex min-w-0 flex-1 flex-col gap-4 lg:flex-row lg:items-start">
      <article className="flex min-w-0 flex-1 flex-col gap-4 rounded-lg border border-border bg-surface px-6 py-5 lg:max-w-[560px]">
        {header}
        <Markdown>{prompt}</Markdown>
        {aside}
      </article>

      <div className="flex min-w-0 flex-1 flex-col gap-4">
        <section aria-label="Run in your terminal" className="flex flex-col gap-3 rounded-lg border border-border bg-surface px-5 py-4">
          <h2 className="text-section text-text">Run in your terminal</h2>
          <p className="text-small text-text-secondary">
            This lesson runs on your computer. From the folder with your code, run it to see the output, then submit it to be graded.
          </p>
          <CopyCommand command={`codequest run ${id}`} />
          <CopyCommand command={`codequest submit ${id}`} />
          <p className="text-small text-text-muted">
            New to the CLI? Install it, then run <code className="font-mono text-text">codequest login</code>.
            {baseURL && (
              <>
                {' '}Requests go to <code className="font-mono text-text">{baseURL}</code> unless you set{' '}
                <code className="font-mono text-text">codequest config base_url</code>.
              </>
            )}
          </p>
        </section>

        <section aria-label="Steps" className="flex flex-col gap-3 rounded-lg border border-border bg-surface px-5 py-4">
          <h2 className="text-section text-text">What gets checked</h2>
          <ol className="flex flex-col gap-3">
            {steps.map((s, i) => (
              <li key={i} className="flex flex-col gap-1">
                <p className="text-body text-text">
                  {i + 1}. {s.description}
                </p>
                <code className="text-code-sm font-mono text-text-secondary [overflow-wrap:anywhere]">{stepAction(s)}</code>
                <ul className="flex flex-col gap-0.5 pl-4">
                  {stepTests(s).map((t, j) => (
                    <li key={j} className="text-small list-disc text-text-muted">
                      {t}
                    </li>
                  ))}
                </ul>
              </li>
            ))}
          </ol>
        </section>

        <div aria-live="polite" className="flex flex-col gap-2">
          {latest ? (
            <CheckList result={latest} />
          ) : (
            <p className="text-small text-text-muted">No submission yet. Results from the CLI show up here.</p>
          )}
          <div>
            <Button onClick={() => void ex.refetch()} disabled={ex.isFetching}>
              {ex.isFetching ? 'Checking…' : 'Check for new results'}
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}
