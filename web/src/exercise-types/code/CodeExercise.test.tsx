import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { useState } from 'react';
import type { GradeResult } from '../../api/client';
import { type CodeAnswer, CodeExercise, type CodePublicView } from './CodeExercise';

// Monaco needs a real browser; a textarea stands in for it here.
vi.mock('./CodeEditor', () => ({
  default: ({ value, onChange, ariaLabel, errorLine }: { value: string; onChange: (v: string) => void; ariaLabel: string; errorLine?: number }) => (
    <textarea aria-label={ariaLabel} data-error-line={errorLine ?? ''} value={value} onChange={(e) => onChange(e.target.value)} />
  ),
}));

const view: CodePublicView = {
  language: 'python', runner: 'docker-python', prompt: 'Return the total of every integer in `nums`.',
  entry: 'main.py', starter: 'def sum_list(nums):\n    pass\n', timeout_ms: 3000,
  examples: [{ input: 'nums = [1, 2, 3]', output: '6', explanation: '1 + 2 + 3 = 6' }],
  constraints: ['0 ≤ len(nums) ≤ 10⁴'], hints: ['What should total start at?', 'Use a for loop.'],
  visible_tests: [],
};

const failing: GradeResult = {
  submission_id: 's1', xp_awarded: 0,
  result: {
    status: 'fail', score: 0.5, stdout: 'debug print\n', runtime_ms: 84,
    checks: [
      { name: 'sum_list([5])', passed: true, visible: true },
      { name: 'sum_list([1, 2, 3])', passed: false, visible: true, got: '1', expected: '6', line: 4 },
    ],
  },
};

function Harness({ onSubmit }: { onSubmit: (a: CodeAnswer) => Promise<GradeResult> }) {
  const [draft, setDraft] = useState<CodeAnswer>();
  const [result, setResult] = useState<GradeResult>();
  return (
    <CodeExercise
      header={<h1>Sum a list</h1>}
      publicView={view}
      draft={draft}
      result={result}
      onChange={setDraft}
      onSubmit={async (a) => { const r = await onSubmit(a); setResult(r); return r; }}
    />
  );
}

describe('CodeExercise', () => {
  it('renders the problem pane: statement, examples, constraints', async () => {
    render(<Harness onSubmit={vi.fn()} />);
    expect(screen.getByRole('heading', { name: 'Sum a list' })).toBeInTheDocument();
    expect(screen.getByText('Example 1:')).toBeInTheDocument();
    expect(screen.getByText('nums = [1, 2, 3]')).toBeInTheDocument();
    expect(screen.getByText('Constraints:')).toBeInTheDocument();
    expect(screen.getByText('main.py')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Run' })).toBeDisabled();
    expect(await screen.findByLabelText('main.py editor')).toHaveValue(view.starter);
  });

  it('reveals hints one at a time', async () => {
    render(<Harness onSubmit={vi.fn()} />);
    expect(screen.queryByText(/Hint 1/)).not.toBeInTheDocument();
    await userEvent.click(screen.getByRole('button', { name: 'Show a hint' }));
    expect(screen.getByText('Hint 1 · What should total start at?')).toBeInTheDocument();
    await userEvent.click(screen.getByRole('button', { name: 'Next hint' }));
    expect(screen.getByText('Hint 2 · Use a for loop.')).toBeInTheDocument();
    expect(screen.queryByRole('button', { name: 'Next hint' })).not.toBeInTheDocument();
  });

  it('submits the edited entry file and shows the checks and failing line', async () => {
    const onSubmit = vi.fn<(a: CodeAnswer) => Promise<GradeResult>>().mockResolvedValue(failing);
    render(<Harness onSubmit={onSubmit} />);
    const editor = await screen.findByLabelText('main.py editor');
    await userEvent.clear(editor);
    await userEvent.type(editor, 'x = 1');
    await userEvent.click(screen.getByRole('button', { name: 'Submit' }));
    expect(onSubmit).toHaveBeenCalledWith({ files: { 'main.py': 'x = 1' } });
    expect(await screen.findByRole('heading', { name: '1 of 2 checks failing' })).toBeInTheDocument();
    expect(screen.getByLabelText('main.py editor')).toHaveAttribute('data-error-line', '4');
    expect(screen.getByText('debug print')).toBeInTheDocument();
  });
});
