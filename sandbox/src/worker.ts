// Runs learner code. One worker per run slot; the page terminates it on timeout.
// Built as a single classic script and started from a blob URL (the page has an opaque origin).
import { loadPyodide, type PyodideAPI } from 'pyodide';
import initSqlJs, { type Database, type SqlJsStatic } from 'sql.js';
import type { RunRequest, RunResult, Table } from './protocol';
import { sqlCheckRow } from './sqlcompare';

type Cell = string | number | null;
const ROW_CAP = 1000;

let base = '';
let pyodide: Promise<PyodideAPI> | undefined;
let sqljs: Promise<SqlJsStatic> | undefined;

async function startPyodide(): Promise<PyodideAPI> {
  const py = await loadPyodide({ indexURL: `${base}vendor/pyodide/`, stdout: () => {}, stderr: () => {} });
  const [harness, cq] = await Promise.all([
    fetch(`${base}vendor/py/browser_harness.py`).then((r) => r.text()),
    fetch(`${base}vendor/py/cq/__init__.py`).then((r) => r.text()),
  ]);
  py.FS.mkdirTree('/cqlib/cq');
  py.FS.writeFile('/cqlib/browser_harness.py', harness);
  py.FS.writeFile('/cqlib/cq/__init__.py', cq);
  py.runPython("import sys; sys.path.append('/cqlib'); import browser_harness");
  return py;
}

function startSql(): Promise<SqlJsStatic> {
  return initSqlJs({ locateFile: (f: string) => `${base}vendor/sqljs/${f}` });
}

function toCell(v: unknown): Cell {
  if (v === null || typeof v === 'number' || typeof v === 'string') return v;
  if (v instanceof Uint8Array) return Array.from(v, (b) => b.toString(16).padStart(2, '0')).join('');
  return String(v);
}

function query(db: Database, sql: string): Table {
  let count = 0;
  for (const stmt of db.iterateStatements(sql)) {
    stmt.free();
    count++;
  }
  if (count !== 1) throw new Error(count === 0 ? 'Write one SQL statement' : 'Write a single SQL statement');
  const stmt = db.prepare(sql);
  try {
    const table: Table = { columns: stmt.getColumnNames(), rows: [], truncated: false };
    while (stmt.step()) {
      if (table.rows.length >= ROW_CAP) {
        table.truncated = true;
        break;
      }
      table.rows.push(stmt.get().map(toCell));
    }
    return table;
  } finally {
    stmt.free();
  }
}

async function runSql(req: RunRequest, t0: number): Promise<RunResult> {
  const SQL = await (sqljs ??= startSql());
  const db = new SQL.Database();
  const base = { type: 'result' as const, runId: req.runId, stdout: '', stderr: '', timedOut: false };
  try {
    db.run(req.sql?.setup ?? '');
    let table: Table | undefined;
    let error: string | undefined;
    try {
      table = query(db, req.files[req.entry] ?? '');
    } catch (e) {
      error = e instanceof Error ? e.message : String(e);
    }
    const checks = (req.sql?.checks ?? []).map((c) => sqlCheckRow(c, table, error));
    return { ...base, checks, table, error, runtimeMs: Math.round(performance.now() - t0) };
  } finally {
    db.close();
  }
}

async function runPython(req: RunRequest, t0: number): Promise<RunResult> {
  const py = await (pyodide ??= startPyodide());
  const job = JSON.stringify({ files: req.files, entry: req.entry, tests: req.tests });
  py.globals.set('cq_job', job);
  const out = JSON.parse(py.runPython('browser_harness.run(cq_job)') as string) as Omit<RunResult, 'type' | 'runId' | 'timedOut'>;
  return { type: 'result', runId: req.runId, timedOut: false, ...out, runtimeMs: Math.round(performance.now() - t0) };
}

self.onmessage = async (e: MessageEvent) => {
  const msg = e.data as { type: 'init'; base: string; warm: ('pyodide' | 'sqljs')[] } | { type: 'run'; req: RunRequest };
  if (msg.type === 'init') {
    base = msg.base;
    try {
      if (msg.warm.includes('pyodide')) await (pyodide ??= startPyodide());
      if (msg.warm.includes('sqljs')) await (sqljs ??= startSql());
      self.postMessage({ type: 'ready' });
    } catch (err) {
      self.postMessage({ type: 'failed', error: String(err) });
    }
    return;
  }
  const t0 = performance.now();
  try {
    const res = msg.req.runner === 'sqljs' ? await runSql(msg.req, t0) : await runPython(msg.req, t0);
    self.postMessage(res);
  } catch (err) {
    self.postMessage({
      type: 'result', runId: msg.req.runId, stdout: '', stderr: '', checks: [], timedOut: false,
      runtimeMs: Math.round(performance.now() - t0), error: err instanceof Error ? err.message : String(err),
    } satisfies RunResult);
  }
};
