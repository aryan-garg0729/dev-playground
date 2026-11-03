// Message protocol between the app (web/) and the sandbox iframe.
// Kept identical in sandbox/src/protocol.ts and web/src/runner-client/protocol.ts (a web test checks).
import { z } from 'zod';

export const MAX_FILE_CHARS = 64 * 1024;
const name = z.string().min(1).max(64);
const source = z.string().max(MAX_FILE_CHARS);
const cell = z.union([z.string(), z.number(), z.null()]);

export const Table = z.strictObject({
  columns: z.array(z.string()).max(64),
  rows: z.array(z.array(cell).max(64)).max(1000),
  truncated: z.boolean(),
});
export type Table = z.infer<typeof Table>;

export const SqlCheck = z.strictObject({
  name: z.string().max(200),
  ordered: z.boolean(),
  expected: z.strictObject({ columns: z.array(z.string()).max(64), rows: z.array(z.array(cell).max(64)).max(1000) }),
});
export type SqlCheck = z.infer<typeof SqlCheck>;

/** App -> sandbox. */
export const RunRequest = z.strictObject({
  type: z.literal('run'),
  runId: name,
  runner: z.enum(['pyodide', 'sqljs']),
  entry: name,
  files: z.record(name, source),
  tests: z.array(z.strictObject({ file: name, source })).max(20),
  sql: z.strictObject({ setup: z.string().max(1 << 20), checks: z.array(SqlCheck).max(20) }).optional(),
  timeoutMs: z.number().int().min(100).max(10000),
});
export type RunRequest = z.infer<typeof RunRequest>;

export const CheckRow = z.strictObject({
  name: z.string().max(500),
  passed: z.boolean(),
  visible: z.boolean(),
  feedback: z.string().max(2000).optional(),
  got: z.string().max(2000).optional(),
  expected: z.string().max(2000).optional(),
  line: z.number().int().min(0).optional(),
});
export type CheckRow = z.infer<typeof CheckRow>;

/** Sandbox -> app. */
export const RunResult = z.strictObject({
  type: z.literal('result'),
  runId: name,
  stdout: z.string().max(64 * 1024),
  stderr: z.string().max(64 * 1024),
  checks: z.array(CheckRow).max(200),
  runtimeMs: z.number().int().min(0),
  timedOut: z.boolean(),
  table: Table.optional(),
  error: z.string().max(2000).optional(),
});
export type RunResult = z.infer<typeof RunResult>;

export const Ready = z.strictObject({ type: z.literal('ready') });

export const SandboxMessage = z.discriminatedUnion('type', [RunResult, Ready]);
export type SandboxMessage = z.infer<typeof SandboxMessage>;
