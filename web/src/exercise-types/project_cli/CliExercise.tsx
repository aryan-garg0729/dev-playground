import { useRef } from 'react';
import { Link } from 'react-router';
import { useExercise } from '../../api/queries';
import { Button } from '../../components/Button';
import { CheckList } from '../../components/CheckList';
import { CopyCommand } from '../../components/CopyCommand';
import { Markdown } from '../../components/Markdown';
import { BossHP, BossResults, DefeatReward, useAttack } from '../shared/Boss';
import type { ExerciseRendererProps } from '../types';
import { type CliPublicView, stepAction, stepTests } from './steps';

/** The CLI posts the answer; the page never builds one. */
type CliAnswer = unknown;

/**
 * A command-line lesson: the steps run on the learner's machine through the
 * `codequest` CLI, which submits their output. This page shows the prompt,
 * the steps and tests, the commands to run, and the latest result.
 * As a chapter boss, every test is one HP and `codequest submit` is the attack.
 */
export function CliExercise({ header, publicView, result, aside, boss }: ExerciseRendererProps<CliPublicView, CliAnswer>) {
  const { lesson_id: id, prompt, steps, base_url: baseURL } = publicView;
  // Same query as the page: refetching picks up a submit made from the terminal.
  const ex = useExercise(id);
  const last = ex.data?.last_submission;
  const latest = last?.result && last.exercise_version === ex.data?.version ? last.result : result?.result;
  // Each test is one result row, so one HP. A result that arrives after the page
  // loaded (a submit from the terminal) plays as a fresh attack.
  const maxHP = steps.reduce((n, s) => n + stepTests(s).length, 0);
  const onLoad = useRef<{ result: typeof latest } | null>(null);
  if (onLoad.current === null && ex.data) onLoad.current = { result: latest };
  const fresh = boss !== undefined && onLoad.current !== null && latest !== undefined && latest !== onLoad.current.result;
  const revealed = useAttack(boss ? latest : undefined, fresh);

  return (
    <div className="flex min-w-0 flex-1 flex-col gap-4 lg:flex-row lg:items-start">
      <article className="flex min-w-0 flex-1 flex-col gap-4 rounded-lg border border-border bg-surface px-6 py-5 lg:max-w-[560px]">
        {header}
        <Markdown dropTitle>{prompt}</Markdown>
        {boss && <BossHP maxHP={maxHP} result={latest} revealed={revealed} />}
        {boss && <DefeatReward reward={boss.reward} />}
        {aside}
      </article>

      <div className="flex min-w-0 flex-1 flex-col gap-4">
        <section aria-label="Run in your terminal" className="flex flex-col gap-3 rounded-lg border border-border bg-surface px-5 py-4">
          <h2 className="text-section text-text">Run in your terminal</h2>
          <p className="text-small text-text-secondary">
            {boss
              ? 'This boss fight runs on your computer. Run the checks to see where you stand, then submit to attack.'
              : 'This lesson runs on your computer. From the folder with your code, run it to see the output, then submit it to be graded.'}
          </p>
          <CopyCommand command={`codequest run ${id}`} />
          <CopyCommand command={`codequest submit ${id}`} />
          <p className="text-small text-text-muted">
            New to the CLI?{' '}
            <Link to="/cli/install" className="text-accent-text underline">
              Install it
            </Link>
            , then run <code className="font-mono text-text">codequest login</code>.
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
            boss ? (
              <BossResults result={latest} revealed={fresh ? revealed : undefined} />
            ) : (
              <CheckList result={latest} />
            )
          ) : (
            <p className="text-small text-text-muted">
              {boss ? 'No attack yet. Submit from the CLI and the result shows up here.' : 'No submission yet. Results from the CLI show up here.'}
            </p>
          )}
          <div>
            <Button onClick={() => void ex.refetch()} disabled={ex.isFetching}>
              {ex.isFetching ? 'Checking…' : boss ? 'Check attack result' : 'Check for new results'}
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}
