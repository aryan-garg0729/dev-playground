import type { GradeResult } from '../api/client';

export type { GradeResult };

/** Output of a client-side (ungraded) run. Filled in with the sandbox in Phase 3. */
export interface RunOutput {
  stdout: string;
  stderr: string;
  runtimeMs: number;
  timedOut: boolean;
}

export interface ExerciseRendererProps<V = unknown, A = unknown> {
  publicView: V; // from GET /exercises/:id
  draft: A | undefined; // last answer
  onChange: (answer: A) => void;
  onRun?: (answer: A) => Promise<RunOutput>; // client-side run, optional
  onSubmit: (answer: A) => Promise<GradeResult>;
  result: GradeResult | undefined;
}
