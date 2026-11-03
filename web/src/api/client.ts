import createClient from 'openapi-fetch';
import type { components, paths } from './schema';

export type Schemas = components['schemas'];
export type CourseSummary = Schemas['CourseSummary'];
export type CourseDetail = Schemas['CourseDetail'];
export type ExerciseDetail = Schemas['ExerciseDetail'];
export type GradeResult = Schemas['SubmitResponse'];
export type Result = Schemas['Result'];
export type Me = Schemas['Me'];
export type Stats = Schemas['Stats'];
export type Outcome = Schemas['Outcome'];
export type Home = Schemas['Home'];
export type League = Schemas['League'];
export type LeagueSummary = Schemas['LeagueSummary'];
export type Shop = Schemas['Shop'];
export type ShopItem = Schemas['ShopItem'];
export type Quests = Schemas['Quests'];
export type Medals = Schemas['Medals'];
export type Medal = Schemas['Medal'];
export type Rewards = Schemas['Rewards'];
export type UseResult = Schemas['UseResult'];

// Resolve fetch at call time so tests can stub globalThis.fetch.
export const api = createClient<paths>({
  baseUrl: `${window.location.origin}/api/v1`,
  fetch: (req) => globalThis.fetch(req),
});

/** An error response from the API, carrying the envelope's code. */
export class ApiError extends Error {
  readonly status: number;
  readonly code: string;
  constructor(status: number, code: string, message: string) {
    super(message);
    this.status = status;
    this.code = code;
  }
}

/** Returns data, or throws ApiError built from the { error: { code, message } } envelope. */
export async function unwrap<T>(
  call: Promise<{ data?: T; error?: unknown; response: Response }>,
): Promise<T> {
  const { data, error, response } = await call;
  if (data !== undefined && response.ok) return data;
  const env = error as Schemas['Error'] | undefined;
  throw new ApiError(
    response.status,
    env?.error?.code ?? 'http_error',
    env?.error?.message ?? `Request failed (${response.status})`,
  );
}
