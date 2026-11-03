import type { ReactNode } from 'react';
import type { GradeResult } from '../api/client';
import type { CheckRow, Table } from '../runner-client/protocol';
import type { SandboxJob } from '../runner-client/sandbox';

export type { GradeResult };

/** Output of an ungraded Run, in the browser sandbox or on the server. */
export interface RunOutput {
  stdout: string;
  stderr: string;
  checks: CheckRow[];
  runtimeMs: number;
  timedOut: boolean;
  table?: Table;
  error?: string;
}

/** Where Run executes, from the manifest's `run:` (in every runnable public view). */
export type RunMode = 'browser' | 'server' | 'off';

export interface ExerciseRendererProps<V = unknown, A = unknown> {
  /** Pills and title, rendered by the page; the renderer places it at the top of its problem pane. */
  header: ReactNode;
  publicView: V; // from GET /exercises/:id
  draft: A | undefined; // last answer
  onChange: (answer: A) => void;
  onRun?: (answer: A) => Promise<RunOutput>; // ungraded run; absent when run is off
  onSubmit: (answer: A) => Promise<GradeResult>;
  result: GradeResult | undefined;
  /** Page content for the bottom of the problem pane (the lesson's power-ups). */
  aside?: ReactNode;
  /** Set when the lesson is a chapter boss: renderers that support it show HP and "Attack". */
  boss?: BossInfo;
  /**
   * The "Ask the tutor" control for the problem pane. Renderers that support the
   * tutor call it with the learner's current answer as text (code or SQL).
   */
  tutor?: (answerText: string) => ReactNode;
}

/** A boss lesson: what defeating it gives. Its HP is its number of checks. */
export interface BossInfo {
  reward: string; // e.g. "Boss chest · Boss Slayer medal · unlocks chapter 2"
}

/** Builds the sandbox job for a browser Run of this type. */
export type ToSandboxJob<V = unknown, A = unknown> = (publicView: V, answer: A) => SandboxJob;
