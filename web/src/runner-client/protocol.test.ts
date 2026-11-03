import { readFileSync } from 'node:fs';
import { resolve } from 'node:path';
import { RunRequest, SandboxMessage } from './protocol';

describe('sandbox protocol', () => {
  it('is identical in web and sandbox', () => {
    const web = readFileSync(resolve(__dirname, 'protocol.ts'), 'utf8');
    const sandbox = readFileSync(resolve(__dirname, '../../../sandbox/src/protocol.ts'), 'utf8');
    expect(web).toBe(sandbox);
  });

  it('accepts a valid run request and drops unknown fields or types', () => {
    const ok = { type: 'run', runId: 'r1', runner: 'pyodide', entry: 'main.py', files: { 'main.py': 'x' }, tests: [], timeoutMs: 3000 };
    expect(RunRequest.safeParse(ok).success).toBe(true);
    expect(RunRequest.safeParse({ ...ok, extra: 1 }).success).toBe(false);
    expect(RunRequest.safeParse({ ...ok, runner: 'docker-python' }).success).toBe(false);
    expect(RunRequest.safeParse({ ...ok, timeoutMs: 999_999 }).success).toBe(false);
    expect(SandboxMessage.safeParse({ type: 'hello' }).success).toBe(false);
    expect(SandboxMessage.safeParse({ type: 'ready' }).success).toBe(true);
  });
});
