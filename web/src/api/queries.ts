import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import type { RunOutput } from '../exercise-types/types';
import { api, ApiError, type GradeResult, unwrap } from './client';

export const keys = {
  me: ['me'] as const,
  courses: ['courses'] as const,
  course: (id: string) => ['course', id] as const,
  exercise: (id: string) => ['exercise', id] as const,
};

export function useMe() {
  return useQuery({
    queryKey: keys.me,
    queryFn: () => unwrap(api.GET('/me')),
    retry: false,
  });
}

export function useCourses() {
  return useQuery({
    queryKey: keys.courses,
    queryFn: () => unwrap(api.GET('/courses')),
  });
}

export function useCourse(courseId: string) {
  return useQuery({
    queryKey: keys.course(courseId),
    queryFn: () => unwrap(api.GET('/courses/{courseId}', { params: { path: { courseId } } })),
  });
}

export function useExercise(exerciseId: string) {
  return useQuery({
    queryKey: keys.exercise(exerciseId),
    queryFn: () => unwrap(api.GET('/exercises/{exerciseId}', { params: { path: { exerciseId } } })),
  });
}

const POLL_MS = 1000;
const POLL_LIMIT_MS = 3 * 60 * 1000;

const sleep = (ms: number) => new Promise((r) => setTimeout(r, ms));

/**
 * Submits an answer. If the server is still grading after 15 s it answers 202;
 * then this polls GET /submissions/{id} until the result is in.
 */
export async function submitAndWait(exerciseId: string, version: number, answer: unknown, pollMs = POLL_MS): Promise<GradeResult> {
  const first = await unwrap(
    api.POST('/exercises/{exerciseId}/submit', {
      params: { path: { exerciseId } },
      body: { exercise_id: exerciseId, version, answer },
    }),
  );
  if ('result' in first) return first;
  const id = first.submission_id;
  const deadline = Date.now() + POLL_LIMIT_MS;
  while (Date.now() < deadline) {
    await sleep(pollMs);
    const s = await unwrap(api.GET('/submissions/{submissionId}', { params: { path: { submissionId: id } } }));
    if (s.status !== 'pending' && s.result) {
      return { submission_id: s.submission_id, result: s.result, xp_awarded: s.xp_awarded };
    }
  }
  throw new ApiError(504, 'timeout', 'Grading is taking too long. Reload the page later to see the result.');
}

export function useSubmit(exerciseId: string, courseId: string | undefined) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (body: { version: number; answer: unknown }) => submitAndWait(exerciseId, body.version, body.answer),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: keys.courses });
      if (courseId) void qc.invalidateQueries({ queryKey: keys.course(courseId) });
    },
  });
}

/** Ungraded server Run (`run: server`). */
export async function runOnServer(exerciseId: string, version: number, answer: unknown): Promise<RunOutput> {
  const out = await unwrap(
    api.POST('/exercises/{exerciseId}/run', {
      params: { path: { exerciseId } },
      body: { exercise_id: exerciseId, version, answer },
    }),
  );
  return {
    stdout: out.stdout,
    stderr: out.stderr,
    checks: out.checks,
    runtimeMs: out.runtime_ms,
    timedOut: out.timed_out,
    table: out.table,
    error: out.error,
  };
}
