import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import type { RunOutput } from '../exercise-types/types';
import { api, ApiError, type GradeResult, unwrap } from './client';

export const keys = {
  me: ['me'] as const,
  home: ['home'] as const,
  league: ['league'] as const,
  providers: ['auth-providers'] as const,
  shop: ['shop'] as const,
  quests: ['quests'] as const,
  medals: ['medals'] as const,
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

/** Outcome of a submission whose grading failed (nothing changed). */
export const noOutcome: GradeResult['outcome'] = {
  xp_awarded: 0, breakdown: null, first_try: false, total_xp: 0, level: 0, level_up: false,
  level_start_xp: 0, next_level_xp: 0, streak_days: 0, streak_extended: false, freezes_used: 0, freeze_earned: false,
  second_wind_kept: false, focus_boost: false, rewards: null,
};
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
      return { submission_id: s.submission_id, result: s.result, xp_awarded: s.xp_awarded, outcome: s.outcome ?? noOutcome };
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
      void qc.invalidateQueries({ queryKey: keys.me });
      void qc.invalidateQueries({ queryKey: keys.home });
      void qc.invalidateQueries({ queryKey: keys.league });
      void qc.invalidateQueries({ queryKey: keys.shop });
      void qc.invalidateQueries({ queryKey: keys.quests });
      void qc.invalidateQueries({ queryKey: keys.medals });
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

export function useHome() {
  return useQuery({ queryKey: keys.home, queryFn: () => unwrap(api.GET('/home')) });
}

export function useLeague() {
  return useQuery({ queryKey: keys.league, queryFn: () => unwrap(api.GET('/leagues/current')) });
}

export function useProviders() {
  return useQuery({ queryKey: keys.providers, queryFn: () => unwrap(api.GET('/auth/providers')), staleTime: Infinity });
}

/** The browser's IANA timezone, saved on first sign-in (streak days follow it). */
export function browserTimezone(): string {
  try {
    return Intl.DateTimeFormat().resolvedOptions().timeZone || 'UTC';
  } catch {
    return 'UTC';
  }
}

export function sendMagicLink(email: string, next: string) {
  return unwrap(api.POST('/auth/magic-link', { body: { email, next } }));
}

/** Redeems a magic-link token; on success the session cookie is set. */
export function useVerifyMagicLink() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (token: string) => unwrap(api.POST('/auth/magic-link/verify', { body: { token, timezone: browserTimezone() } })),
    onSuccess: () => qc.invalidateQueries(),
  });
}

/** Development sign-in (DEV_AUTH=1 only). */
export function useDevSignIn() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: () => unwrap(api.POST('/auth/dev', { body: { timezone: browserTimezone() } })),
    onSuccess: () => qc.invalidateQueries(),
  });
}

export function useSignOut() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: async () => {
      const { response } = await api.POST('/auth/logout', { body: {} });
      if (!response.ok) throw new ApiError(response.status, 'http_error', 'Sign out failed');
    },
    onSuccess: () => qc.clear(),
  });
}

export function useShop() {
  return useQuery({ queryKey: keys.shop, queryFn: () => unwrap(api.GET('/shop')) });
}

export function useQuests() {
  return useQuery({ queryKey: keys.quests, queryFn: () => unwrap(api.GET('/quests')) });
}

export function useMedals() {
  return useQuery({ queryKey: keys.medals, queryFn: () => unwrap(api.GET('/medals')) });
}

/** A fresh idempotency key per action: a retried request cannot spend twice. */
const requestId = () => crypto.randomUUID();

function useEconomyMutation<T, V>(fn: (v: V, id: string) => Promise<T>) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (v: V) => fn(v, requestId()),
    onSettled: () => {
      for (const k of [keys.shop, keys.me, keys.quests, keys.home, keys.medals]) void qc.invalidateQueries({ queryKey: k });
    },
  });
}

export function useBuy() {
  return useEconomyMutation((powerup: string, id) =>
    unwrap(api.POST('/shop/{powerup}/buy', { params: { path: { powerup } }, body: { request_id: id } })),
  );
}

export function useUsePowerup() {
  return useEconomyMutation((v: { powerup: string; exerciseId?: string; stake?: number }, id) =>
    unwrap(
      api.POST('/powerups/{powerup}/use', {
        params: { path: { powerup: v.powerup } },
        body: { request_id: id, exercise_id: v.exerciseId, stake: v.stake },
      }),
    ),
  );
}

export function useOpenChest() {
  return useEconomyMutation((_: void, id) => unwrap(api.POST('/quests/chest/open', { body: { request_id: id } })));
}

export function useSetTitle() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (title: string) => unwrap(api.PATCH('/me', { body: { title } })),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: keys.me });
      void qc.invalidateQueries({ queryKey: keys.medals });
    },
  });
}
