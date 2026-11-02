import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { useState } from 'react';
import type { GradeResult } from '../../api/client';
import { McqExercise, type McqAnswer, type McqPublicView } from './McqExercise';

const view: McqPublicView = {
  question: 'Which message carries the **persona**?',
  multiple: false,
  options: [
    { id: 'a', text: 'SystemMessage' },
    { id: 'b', text: 'HumanMessage' },
  ],
};

const graded = (status: 'pass' | 'fail'): GradeResult => ({
  submission_id: '00000000-0000-0000-0000-000000000000',
  xp_awarded: 0,
  result: {
    status, score: status === 'pass' ? 1 : 0, stdout: '', runtime_ms: 0,
    checks: [{ name: 'Correct answer', passed: status === 'pass', visible: true, feedback: 'SystemMessage holds instructions.' }],
  },
});

function Harness({ onSubmit }: { onSubmit: (a: McqAnswer) => Promise<GradeResult> }) {
  const [draft, setDraft] = useState<McqAnswer>();
  const [result, setResult] = useState<GradeResult>();
  return (
    <McqExercise
      publicView={view}
      draft={draft}
      result={result}
      onChange={(a) => { setDraft(a); setResult(undefined); }}
      onSubmit={async (a) => { const r = await onSubmit(a); setResult(r); return r; }}
    />
  );
}

describe('McqExercise', () => {
  it('renders the question and disables Submit until an option is chosen', () => {
    render(<Harness onSubmit={vi.fn()} />);
    expect(screen.getByText('persona').tagName).toBe('STRONG');
    expect(screen.getByRole('button', { name: 'Submit' })).toBeDisabled();
  });

  it('shows fail with the explanation, then pass after changing the answer', async () => {
    const onSubmit = vi.fn<(a: McqAnswer) => Promise<GradeResult>>()
      .mockResolvedValueOnce(graded('fail'))
      .mockResolvedValueOnce(graded('pass'));
    render(<Harness onSubmit={onSubmit} />);

    await userEvent.click(screen.getByLabelText('HumanMessage'));
    await userEvent.click(screen.getByRole('button', { name: 'Submit' }));
    expect(onSubmit).toHaveBeenLastCalledWith({ choices: ['b'] });
    expect(await screen.findByText('Not quite')).toBeInTheDocument();
    expect(screen.getByText('SystemMessage holds instructions.')).toBeInTheDocument();

    await userEvent.click(screen.getByLabelText('SystemMessage'));
    expect(screen.queryByText('Not quite')).not.toBeInTheDocument();
    await userEvent.click(screen.getByRole('button', { name: 'Submit' }));
    expect(onSubmit).toHaveBeenLastCalledWith({ choices: ['a'] });
    expect(await screen.findByText('Correct')).toBeInTheDocument();
  });

  it('collects several choices when multiple is true', async () => {
    const onSubmit = vi.fn<(a: McqAnswer) => Promise<GradeResult>>().mockResolvedValue(graded('pass'));
    function Multi() {
      const [draft, setDraft] = useState<McqAnswer>();
      return <McqExercise publicView={{ ...view, multiple: true }} draft={draft} result={undefined} onChange={setDraft} onSubmit={onSubmit} />;
    }
    render(<Multi />);
    expect(screen.getByText('Select all that apply')).toBeInTheDocument();
    await userEvent.click(screen.getByLabelText('SystemMessage'));
    await userEvent.click(screen.getByLabelText('HumanMessage'));
    await userEvent.click(screen.getByRole('button', { name: 'Submit' }));
    expect(onSubmit).toHaveBeenCalledWith({ choices: ['a', 'b'] });
  });
});
