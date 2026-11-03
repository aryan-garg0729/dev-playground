// Result-set comparison for sql_result_equals. Must match runners/python/harness/sql_child.py.
import type { CheckRow, SqlCheck, Table } from './protocol';

type Cell = string | number | null;

function sameCell(a: Cell, b: Cell): boolean {
  if (typeof a === 'number' && typeof b === 'number') {
    return Math.abs(a - b) <= 1e-9 * Math.max(1, Math.abs(a), Math.abs(b));
  }
  return a === b;
}

function sameRow(a: Cell[], b: Cell[]): boolean {
  return a.length === b.length && a.every((c, i) => sameCell(c, b[i] ?? null));
}

/** True when got has the same rows as want (in order when ordered, else as a multiset). */
export function sameRows(got: Cell[][], want: Cell[][], ordered: boolean): boolean {
  if (got.length !== want.length) return false;
  if (ordered) return got.every((r, i) => sameRow(r, want[i] ?? []));
  const left = [...want];
  for (const r of got) {
    const i = left.findIndex((w) => sameRow(r, w));
    if (i < 0) return false;
    left.splice(i, 1);
  }
  return true;
}

function preview(rows: Cell[][]): string {
  const s = JSON.stringify(rows.slice(0, 5));
  return rows.length > 5 ? `${s.slice(0, -1)}, … ${rows.length} rows]` : s;
}

export function sqlCheckRow(check: SqlCheck, got: Table | undefined, error: string | undefined): CheckRow {
  if (error || !got) return { name: check.name, passed: false, visible: true, feedback: error ?? 'Did not run' };
  if (got.columns.length !== check.expected.columns.length) {
    return {
      name: check.name, passed: false, visible: true,
      feedback: `Expected ${check.expected.columns.length} column(s) (${check.expected.columns.join(', ')}), got ${got.columns.length}`,
    };
  }
  const passed = !got.truncated && sameRows(got.rows, check.expected.rows, check.ordered);
  return passed
    ? { name: check.name, passed, visible: true }
    : { name: check.name, passed, visible: true, got: preview(got.rows), expected: preview(check.expected.rows) };
}
