import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { useState } from 'react';
import type { RunOutput } from '../types';
import { type SqlAnswer, SqlExercise, type SqlPublicView, sqlSandboxJob } from './SqlExercise';

vi.mock('../code/CodeEditor', () => ({
  default: ({ value, onChange, ariaLabel }: { value: string; onChange: (v: string) => void; ariaLabel: string }) => (
    <textarea aria-label={ariaLabel} value={value} onChange={(e) => onChange(e.target.value)} />
  ),
}));

const view: SqlPublicView = {
  prompt: 'Return every name.', setup: 'CREATE TABLE t (name TEXT);', starter: '', run: 'browser', timeout_ms: 2000,
  visible_checks: [{ name: 'names', ordered: false, expected: { columns: ['name'], rows: [['a']] } }],
};

describe('SqlExercise', () => {
  it('runs the query and shows the result table', async () => {
    const onRun = vi.fn<(a: SqlAnswer) => Promise<RunOutput>>().mockResolvedValue({
      stdout: '', stderr: '', runtimeMs: 3, timedOut: false,
      checks: [{ name: 'names', passed: true, visible: true }],
      table: { columns: ['name'], rows: [['a'], [null]], truncated: false },
    });
    function H() {
      const [d, setD] = useState<SqlAnswer>();
      return <SqlExercise header={null} publicView={view} draft={d} onChange={setD} onRun={onRun} onSubmit={vi.fn()} result={undefined} />;
    }
    render(<H />);
    await userEvent.type(await screen.findByLabelText('query.sql editor'), 'SELECT name FROM t');
    await userEvent.click(screen.getByRole('button', { name: 'Run' }));
    expect(onRun).toHaveBeenCalledWith({ query: 'SELECT name FROM t' });
    expect(await screen.findByRole('columnheader', { name: 'name' })).toBeInTheDocument();
    expect(screen.getByText('NULL')).toBeInTheDocument();
    expect(screen.getByText('2 rows')).toBeInTheDocument();
  });

  it('builds an sql.js job with the public setup and visible checks', () => {
    expect(sqlSandboxJob(view, { query: 'SELECT 1' })).toEqual({
      runner: 'sqljs', entry: 'query.sql', files: { 'query.sql': 'SELECT 1' }, tests: [],
      sql: { setup: view.setup, checks: view.visible_checks }, timeoutMs: 2000,
    });
  });
});
