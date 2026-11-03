import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter } from 'react-router';
import { vi } from 'vitest';
import type { Result } from '../api/client';
import { firstPass } from '../test/fixtures';
import { PassCelebration } from './PassCelebration';

const result: Result = {
  status: 'pass', score: 1, stdout: '', runtime_ms: 41,
  checks: [1, 2, 3].map((i) => ({ name: `c${i}`, passed: true, visible: true })),
};

const findTheMax = { to: '/learn/c/max', title: 'Find the max' };

function show(onClose = vi.fn(), next: { to: string; title: string } | null = findTheMax) {
  render(
    <MemoryRouter>
      <PassCelebration outcome={firstPass} result={result} next={next ?? undefined} courseTo="/courses/c" onClose={onClose} />
    </MemoryRouter>,
  );
  return onClose;
}

describe('PassCelebration', () => {
  it('shows the XP breakdown, level progress and streak', () => {
    show();
    const dialog = screen.getByRole('dialog', { name: 'All checks passed' });
    expect(dialog).toHaveTextContent('3 of 3 checks · 41 ms · first try');
    expect(dialog).toHaveTextContent('Lesson complete+20');
    expect(dialog).toHaveTextContent('First-try bonus+10');
    expect(dialog).toHaveTextContent('No hints used+5');
    expect(dialog).toHaveTextContent('+35 XP');
    // 1727 - 1470 = 257 of 382 in level 7; 125 to go.
    expect(screen.getByRole('img', { name: '257 of 382 XP in level 7' })).toBeInTheDocument();
    expect(dialog).toHaveTextContent('125 XP to Lv 8');
    expect(dialog).toHaveTextContent('13-day streak');
  });

  it('focuses Next, which links to the next lesson', () => {
    show();
    const next = screen.getByRole('link', { name: 'Next: Find the max' });
    expect(next).toHaveFocus();
    expect(next).toHaveAttribute('href', '/learn/c/max');
  });

  it('is skippable with Escape, the backdrop and Keep working', async () => {
    const onClose = show();
    await userEvent.keyboard('{Escape}');
    await userEvent.click(screen.getByTestId('backdrop'));
    await userEvent.click(screen.getByRole('button', { name: 'Keep working' }));
    expect(onClose).toHaveBeenCalledTimes(3);
  });

  it('goes back to the course after the last lesson', () => {
    show(vi.fn(), null);
    expect(screen.getByRole('link', { name: 'Back to the course' })).toHaveAttribute('href', '/courses/c');
  });

  it('animates the XP bar only with motion allowed (motion-safe classes)', () => {
    show();
    const gained = screen.getByRole('img', { name: /XP in level/ }).lastElementChild!;
    expect(gained.className).toMatch(/motion-safe:transition-\[width\]/);
    expect(gained.className).not.toMatch(/(^|\s)transition/);
  });
});
