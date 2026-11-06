import { act, render, screen } from '@testing-library/react';
import { vi } from 'vitest';
import type { Result } from '../../api/client';
import { setReducedMotion } from '../../test/motion';
import { BossHP, BossResults, useAttack } from './Boss';

const play = vi.hoisted(() => vi.fn());
vi.mock('../../effects/sound', () => ({ play, resetSoundForTests: () => {} }));

const result = (passed: boolean[]): Result => ({
  status: passed.every(Boolean) ? 'pass' : 'fail', score: 0, stdout: '', runtime_ms: 1,
  checks: passed.map((p, i) => ({ name: `fizzbuzz(${i})`, passed: p, visible: i < 3, got: p ? undefined : '"Fizz"', expected: p ? undefined : '"FizzBuzz"' })),
});

function Fight({ r, fresh }: { r: Result; fresh: boolean }) {
  const n = useAttack(r, fresh);
  return (
    <>
      <BossHP maxHP={r.checks.length} result={r} revealed={fresh ? n : undefined} />
      <BossResults result={r} revealed={fresh ? n : undefined} />
    </>
  );
}

describe('Boss', () => {
  beforeEach(() => play.mockClear());

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

  it('plays a fresh attack check by check: HP drops one per hit, a sound per check, then the fanfare', () => {
    vi.useFakeTimers();
    setReducedMotion(false);
    render(<Fight r={result([true, false, true])} fresh />);
    expect(screen.getByRole('img', { name: '3 of 3 HP left' })).toBeInTheDocument();
    expect(screen.getByRole('heading')).toHaveTextContent('Attacking…');
    act(() => vi.advanceTimersByTime(150));
    expect(screen.getByRole('img', { name: '2 of 3 HP left' })).toBeInTheDocument();
    act(() => vi.advanceTimersByTime(300));
    expect(screen.getByRole('img', { name: '1 of 3 HP left' })).toBeInTheDocument();
    expect(screen.getByRole('heading')).toHaveTextContent('2 hits · boss blocked 1');
    act(() => vi.advanceTimersByTime(500));
    expect(play.mock.calls).toEqual([['bossHit'], ['bossBlock'], ['bossHit']]); // not defeated: no fanfare
    vi.useRealTimers();
  });

  it('plays the defeat fanfare after a full clear', () => {
    vi.useFakeTimers();
    setReducedMotion(false);
    render(<Fight r={result([true, true])} fresh />);
    act(() => vi.advanceTimersByTime(1000));
    expect(play.mock.calls.at(-1)).toEqual(['bossDefeat']);
    expect(screen.getByRole('heading')).toHaveTextContent('Boss defeated');
    vi.useRealTimers();
  });

  it('shows an old result (page load) at once and silently', () => {
    vi.useFakeTimers();
    setReducedMotion(false);
    render(<Fight r={result([true, true, false])} fresh={false} />);
    expect(screen.getByRole('img', { name: '1 of 3 HP left' })).toBeInTheDocument();
    act(() => vi.advanceTimersByTime(1000));
    expect(play).not.toHaveBeenCalled();
    vi.useRealTimers();
  });
});
