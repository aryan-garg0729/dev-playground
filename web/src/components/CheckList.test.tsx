import { render, screen, within } from '@testing-library/react';
import type { Result } from '../api/client';
import { CheckList } from './CheckList';

const result = (over: Partial<Result>): Result => ({ status: 'fail', score: 0, stdout: '', runtime_ms: 84, checks: [], ...over });

describe('CheckList', () => {
  it('summarises failures and shows got / expected on failed rows', () => {
    render(
      <CheckList
        result={result({
          checks: [
            { name: 'sum_list([5])', passed: true, visible: true },
            { name: 'sum_list([1, 2, 3])', passed: false, visible: true, got: '1', expected: '6' },
            { name: 'Hidden test 1', passed: false, visible: false, feedback: 'Wrong result' },
          ],
        })}
      />,
    );
    expect(screen.getByRole('heading', { name: '2 of 3 checks failing' })).toHaveClass('text-danger');
    expect(screen.getByText('84 ms')).toBeInTheDocument();
    const rows = screen.getAllByRole('listitem');
    expect(within(rows[0]!).getByRole('img', { name: 'Passed' })).toBeInTheDocument();
    expect(within(rows[1]!).getByRole('img', { name: 'Failed' })).toBeInTheDocument();
    expect(rows[1]).toHaveTextContent('sum_list([1, 2, 3]) → 1 (expected 6)');
    expect(rows[2]).toHaveTextContent('Hidden test 1 · Wrong result');
  });

  it('says all checks pass in success color', () => {
    render(<CheckList result={result({ status: 'pass', score: 1, checks: [{ name: 'a', passed: true, visible: true }] })} />);
    expect(screen.getByRole('heading', { name: '1 of 1 checks passing' })).toHaveClass('text-success');
  });

  it('names a timeout', () => {
    render(<CheckList result={result({ status: 'timeout', checks: [{ name: 'Time limit', passed: false, visible: true, feedback: 'Your code ran longer than 3000 ms' }] })} />);
    expect(screen.getByRole('heading', { name: 'Time limit exceeded' })).toBeInTheDocument();
  });

  it('renders grader text as text, never as HTML', () => {
    const { container } = render(
      <CheckList result={result({ checks: [{ name: '<img src=x onerror=alert(1)>', passed: false, visible: true, got: '<b>x</b>' }] })} />,
    );
    expect(container.querySelector('img')).toBeNull();
    expect(container.querySelector('b')).toBeNull();
    expect(screen.getByText('<b>x</b>')).toBeInTheDocument();
  });
});
