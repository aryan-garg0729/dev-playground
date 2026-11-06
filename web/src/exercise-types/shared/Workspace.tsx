import { type CSSProperties, lazy, type ReactNode, type RefObject, Suspense, useEffect, useRef, useState } from 'react';
import { ApiError } from '../../api/client';
import { Button } from '../../components/Button';
import { Splitter } from '../../components/Splitter';
import { PANE_LIMITS, usePanes } from '../../stores/panes';
import type { BossInfo, GradeResult, RunOutput } from '../types';
import { BossHP, BossResults, DefeatReward, useAttack } from './Boss';
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
  /** Page content for the bottom of the problem pane (power-ups). */
  aside?: ReactNode;
  /** Boss lesson: maxHP is its number of checks. */
  boss?: BossInfo & { maxHP: number };
}

function message(e: unknown, what: string): string {
  if (e instanceof ApiError) return e.message;
  if (e instanceof Error && e.message) return e.message;
  return `Could not ${what}. Check your connection and try again.`;
}

/** Room the editor always keeps when a neighbouring pane grows (px). */
const EDITOR_MIN_WIDTH = 360;
const EDITOR_MIN_HEIGHT = 160;

/** Tracks an element's size (0 until measured, or where ResizeObserver is missing). */
function useSize(ref: RefObject<HTMLElement | null>) {
  const [size, setSize] = useState({ width: 0, height: 0 });
  useEffect(() => {
    const el = ref.current;
    if (!el || typeof ResizeObserver === 'undefined') return;
    const ro = new ResizeObserver(([entry]) => {
      if (entry) setSize({ width: entry.contentRect.width, height: entry.contentRect.height });
    });
    ro.observe(el);
    return () => ro.disconnect();
  }, [ref]);
  return size;
}

/**
 * The lesson workspace: problem pane | editor with Run and Submit, then results.
 * Drag handles resize the problem pane and the checks pane; sizes are kept per
 * device (stores/panes). Below 900px the panes stack and keep natural sizes.
 */
export function Workspace(p: Props) {
  const [busy, setBusy] = useState<'run' | 'submit'>();
  const [error, setError] = useState<string>();
  const [run, setRun] = useState<RunOutput>();
  const [showRun, setShowRun] = useState(false);
  // The result shown when the latest Attack was clicked (null: none). Any newer
  // result is that attack's, and only it plays the boss animation.
  const [attackFrom, setAttackFrom] = useState<string | null>();
  const shown: Shown = showRun && run ? { kind: 'run', out: run } : p.result ? { kind: 'submit', graded: p.result } : undefined;
  const failing = shown?.kind === 'submit' ? shown.graded.result.checks : shown?.out.checks;
  const errorLine = failing?.find((c) => !c.passed && (c.line ?? 0) > 0)?.line;
  const fresh = p.boss !== undefined && p.result !== undefined && attackFrom !== undefined && p.result.submission_id !== attackFrom;
  const revealed = useAttack(p.boss ? p.result?.result : undefined, fresh);
  const reveal = fresh ? revealed : undefined;

  async function act(kind: 'run' | 'submit') {
    setBusy(kind);
    setError(undefined);
    try {
      if (kind === 'run' && p.onRun) {
        setRun(await p.onRun());
        setShowRun(true);
      } else {
        setAttackFrom(p.result?.submission_id ?? null);
        await p.onSubmit();
        setShowRun(false);
      }
    } catch (e) {
      setError(message(e, kind));
    } finally {
      setBusy(undefined);
    }
  }

  const panes = usePanes();
  const row = useRef<HTMLDivElement>(null);
  const column = useRef<HTMLDivElement>(null);
  const rowSize = useSize(row);
  const columnSize = useSize(column);
  const maxProblem = Math.max(
    PANE_LIMITS.problemWidth.min,
    Math.min(PANE_LIMITS.problemWidth.max, rowSize.width ? rowSize.width - EDITOR_MIN_WIDTH : Infinity),
  );
  const maxResults = Math.max(
    PANE_LIMITS.resultsHeight.min,
    Math.min(PANE_LIMITS.resultsHeight.max, columnSize.height ? columnSize.height - EDITOR_MIN_HEIGHT : Infinity),
  );
  const problemWidth = Math.min(panes.problemWidth, maxProblem);
  const resultsHeight = Math.min(panes.resultsHeight, maxResults);
  const sizes = { '--problem-w': `${problemWidth}px`, '--results-h': `${resultsHeight}px` } as CSSProperties;

  return (
    <div ref={row} style={sizes} className="flex min-h-0 min-w-0 flex-1 max-[900px]:flex-col max-[900px]:gap-2">
      <section
        aria-label="Problem"
        className="flex w-[var(--problem-w)] shrink-0 flex-col gap-4 overflow-y-auto rounded-lg border border-border bg-surface px-6 py-5 max-[900px]:w-auto"
      >
        {p.header}
        {p.boss && <BossHP maxHP={p.boss.maxHP} result={p.result?.result} revealed={reveal} />}
        {p.problem}
        {p.aside}
        {(p.problemFooter || p.boss) && (
          <div className="mt-auto flex flex-col gap-3 pt-4">
            {p.problemFooter}
            {p.boss && <DefeatReward reward={p.boss.reward} />}
          </div>
        )}
      </section>

      <Splitter
        orientation="vertical"
        label="Resize problem and editor"
        value={problemWidth}
        min={PANE_LIMITS.problemWidth.min}
        max={maxProblem}
        direction={1}
        onChange={(px) => panes.setSize('problemWidth', Math.min(px, maxProblem))}
        onReset={() => panes.reset('problemWidth')}
      />

      <div ref={column} className="flex min-h-0 min-w-0 flex-1 flex-col max-[900px]:gap-2">
        <section
          aria-label="Editor"
          className="flex min-h-[160px] flex-1 flex-col overflow-hidden rounded-lg border border-border bg-surface max-[900px]:min-h-[320px]"
        >
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
              {busy === 'submit' ? (p.boss ? 'Attacking…' : 'Checking…') : p.boss ? 'Attack' : 'Submit'}
            </Button>
          </div>
          <div className="min-h-0 flex-1">
            <Suspense fallback={<p className="text-small p-4 text-text-muted">Loading editor…</p>}>
              <CodeEditor value={p.code} language={p.language} errorLine={errorLine} ariaLabel={`${p.file} editor`} onChange={p.onCode} />
            </Suspense>
          </div>
        </section>
        {error && (
          <p role="alert" className="text-small mt-2 rounded-lg border border-border bg-danger-soft px-4 py-2 text-danger">
            {error}
          </p>
        )}
        {shown && (
          <Splitter
            orientation="horizontal"
            label="Resize editor and checks"
            value={resultsHeight}
            min={PANE_LIMITS.resultsHeight.min}
            max={maxResults}
            direction={-1}
            onChange={(px) => panes.setSize('resultsHeight', Math.min(px, maxResults))}
            onReset={() => panes.reset('resultsHeight')}
          />
        )}
        <div
          aria-live="polite"
          className={`flex shrink-0 flex-col gap-2 overflow-y-auto ${shown ? 'h-[var(--results-h)] max-[900px]:h-auto max-[900px]:max-h-[60vh]' : ''}`}
        >
          {p.boss && shown?.kind === 'submit' ? <BossResults result={shown.graded.result} revealed={reveal} /> : <Results shown={shown} />}
        </div>
      </div>
    </div>
  );
}
