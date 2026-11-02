import type { ComponentType } from 'react';
import { CodeExercise } from './code/CodeExercise';
import { McqExercise } from './mcq/McqExercise';
import type { ExerciseRendererProps } from './types';

export type { ExerciseRendererProps, RunOutput } from './types';

// The one place `any` is allowed: each renderer narrows its own view/answer types.
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export const renderers: Record<string, ComponentType<ExerciseRendererProps<any, any>>> = {
  code: CodeExercise,
  mcq: McqExercise,
};
