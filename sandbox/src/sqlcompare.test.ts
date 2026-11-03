import { describe, expect, it } from 'vitest';
import { RunRequest, SandboxMessage } from './protocol';
import { sameRows, sqlCheckRow } from './sqlcompare';

describe('sameRows', () => {
  it('compares in order when ordered', () => {
    expect(sameRows([[1], [2]], [[1], [2]], true)).toBe(true);
    expect(sameRows([[2], [1]], [[1], [2]], true)).toBe(false);
  });

  it('compares as a multiset when unordered', () => {
    expect(sameRows([[2], [1]], [[1], [2]], false)).toBe(true);
    expect(sameRows([[1], [1]], [[1], [2]], false)).toBe(false); // duplicates count
    expect(sameRows([[1]], [[1], [1]], false)).toBe(false);
  });

  it('tolerates float rounding but not type changes', () => {
    expect(sameRows([[0.1 + 0.2]], [[0.3]], true)).toBe(true);
    expect(sameRows([[1.5]], [[1.6]], true)).toBe(false);
    expect(sameRows([['1']], [[1]], true)).toBe(false);
    expect(sameRows([[null]], [[null]], true)).toBe(true);
    expect(sameRows([[null]], [[0]], true)).toBe(false);
  });
});

describe('sqlCheckRow', () => {
  const check = { name: 'all rows', ordered: false, expected: { columns: ['a', 'b'], rows: [[1, 'x']] } };

  it('passes on a matching table', () => {
    expect(sqlCheckRow(check, { columns: ['a', 'b'], rows: [[1, 'x']], truncated: false }, undefined)).toEqual({
      name: 'all rows', passed: true, visible: true,
    });
  });

  it('reports a column-count mismatch', () => {
    const row = sqlCheckRow(check, { columns: ['a'], rows: [[1]], truncated: false }, undefined);
    expect(row.passed).toBe(false);
    expect(row.feedback).toBe('Expected 2 column(s) (a, b), got 1');
  });

  it('fails a truncated result even when the kept rows match', () => {
    expect(sqlCheckRow(check, { columns: ['a', 'b'], rows: [[1, 'x']], truncated: true }, undefined).passed).toBe(false);
  });

  it('shows got and expected on wrong rows, and the error when the query failed', () => {
    const row = sqlCheckRow(check, { columns: ['a', 'b'], rows: [[2, 'y']], truncated: false }, undefined);
    expect(row).toMatchObject({ passed: false, got: '[[2,"y"]]', expected: '[[1,"x"]]' });
    expect(sqlCheckRow(check, undefined, 'no such table: t')).toMatchObject({ passed: false, feedback: 'no such table: t' });
  });
});

describe('protocol', () => {
  it('drops unknown fields and message types', () => {
    const run = { type: 'run', runId: 'r1', runner: 'pyodide', entry: 'main.py', files: {}, tests: [], timeoutMs: 3000 };
    expect(RunRequest.safeParse(run).success).toBe(true);
    expect(RunRequest.safeParse({ ...run, extra: 1 }).success).toBe(false);
    expect(SandboxMessage.safeParse({ type: 'eval', code: '1' }).success).toBe(false);
  });
});
