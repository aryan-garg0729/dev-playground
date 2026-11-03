import { render, screen } from '@testing-library/react';
import type { Result } from '../../api/client';
import { BossHP, BossResults } from './Boss';

const result = (passed: boolean[]): Result => ({
  status: passed.every(Boolean) ? 'pass' : 'fail', score: 0, stdout: '', runtime_ms: 1,
  checks: passed.map((p, i) => ({ name: `fizzbuzz(${i})`, passed: p, visible: i < 3, got: p ? undefined : '"Fizz"', expected: p ? undefined : '"FizzBuzz"' })),
});

describe('Boss', () => {
  it('has HP = checks before any attack, and loses one per passing check', () => {
    const { rerender } = render(<BossHP maxHP={8} result={undefined} />);
    expect(screen.getByRole('img', { name: '8 of 8 HP left' })).toBeInTheDocument();
    rerender(<BossHP maxHP={8} result={result([true, true, true, false, false, false, false, false])} />);
    expect(screen.getByRole('img', { name: '5 of 8 HP left' })).toBeInTheDocument();
    expect(screen.getByLabelText('Boss HP')).toHaveTextContent('3 visible, 5 hidden');
    rerender(<BossHP maxHP={8} result={result(Array(8).fill(true))} />);
    expect(screen.getByRole('img', { name: '0 of 8 HP left' })).toBeInTheDocument();
  });

  it('shows hits and blocks', () => {
    render(<BossResults result={result([true, true, false])} />);
    expect(screen.getByRole('heading')).toHaveTextContent('2 hits · boss blocked 1');
    expect(screen.getAllByText('HIT')).toHaveLength(2);
    expect(screen.getByText('BLOCK')).toBeInTheDocument();
  });
});
