import { act, render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter } from 'react-router';
import { useRail } from '../stores/rail';
import { LessonRail, type RailLesson } from './LessonRail';

const lessons: RailLesson[] = [
  { id: 'a', title: 'For loops', status: 'done', boss: false },
  { id: 'b', title: 'Sum a list', status: 'current', boss: false },
  { id: 'c', title: 'Boss: FizzBuzz', status: 'todo', boss: true },
];

function renderRail() {
  return render(
    <MemoryRouter>
      <LessonRail courseId="py" chapterNumber={3} chapterTitle="Loops" lessons={lessons} />
    </MemoryRouter>,
  );
}

describe('LessonRail', () => {
  beforeEach(() => act(() => useRail.setState({ open: false })));

  it('is closed by default with one labelled segment per lesson', () => {
    renderRail();
    expect(screen.getByRole('button', { name: '2. Sum a list, current' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: '1. For loops, done' })).toBeInTheDocument();
    expect(screen.queryByText('For loops')).not.toBeInTheDocument();
  });

  it('opens from a segment, shows names and progress, and remembers the choice', async () => {
    renderRail();
    await userEvent.click(screen.getByRole('button', { name: '3. Boss: FizzBuzz, not started' }));
    expect(screen.getByText('Chapter 3 · Loops · 1/3')).toBeInTheDocument();
    const current = screen.getByRole('link', { name: '2. Sum a list, current' });
    expect(current).toHaveAttribute('aria-current', 'page');
    expect(current).toHaveAttribute('href', '/learn/py/b');
    expect(localStorage.getItem('cq-rail-open')).toBe('1');

    await userEvent.click(screen.getByRole('button', { name: 'Close outline' }));
    expect(screen.queryByText('Chapter 3 · Loops · 1/3')).not.toBeInTheDocument();
    expect(localStorage.getItem('cq-rail-open')).toBe('0');
  });
});
