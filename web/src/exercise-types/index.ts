import type { ComponentType } from 'react';
import { CodeExercise, codeSandboxJob } from './code/CodeExercise';
import { McqExercise } from './mcq/McqExercise';
import { CliExercise } from './project_cli/CliExercise';
import { SqlExercise, sqlSandboxJob } from './sql/SqlExercise';
import type { ExerciseRendererProps, ToSandboxJob } from './types';

export type { ExerciseRendererProps, RunMode, RunOutput } from './types';

// The one place `any` is allowed: each renderer narrows its own view/answer types.
/* eslint-disable @typescript-eslint/no-explicit-any */
export const renderers: Record<string, ComponentType<ExerciseRendererProps<any, any>>> = {
  code: CodeExercise,
  mcq: McqExercise,
  project_cli: CliExercise,
  sql: SqlExercise,
};

/** Types that can Run in the browser sandbox. */
export const sandboxJobs: Record<string, ToSandboxJob<any, any>> = {
  code: codeSandboxJob,
  sql: sqlSandboxJob,
};
/* eslint-enable @typescript-eslint/no-explicit-any */
