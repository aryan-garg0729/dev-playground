import { screen, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { renderRoute, stubApi } from '../test/render';
import { CoursePage } from './CoursePage';

const ex = (id: string, title: string, status = 'not_started') => ({
  id, title, type: 'mcq', xp: 5, difficulty: 1, free: true, boss: false, status,
});

const course = {
  id: 'lc', title: 'LangChain Basics', description: 'Blocks.', skills: ['messages', 'prompts'], xp_earned: 9, xp_total: 15,
  chapters: [
    { id: 'b', title: 'Zeta chapter', exercises: [ex('b2', 'Second'), ex('b1', 'First', 'passed')] },
    { id: 'a', title: 'Alpha chapter', exercises: [ex('a1', 'Only')] },
  ],
};

describe('CoursePage (course map)', () => {
  it('shows chapters in API (manifest) order with the current chapter open', async () => {
    stubApi({ 'GET /api/v1/courses/lc': [200, course] });
    renderRoute('/courses/lc', '/courses/:courseId', <CoursePage />);

    const chapters = await screen.findByRole('list', { name: 'Chapters' });
    const toggles = within(chapters).getAllByRole('button');
    expect(toggles.map((b) => b.textContent)).toEqual([
      expect.stringMatching(/^1\. Zeta chapter2 lessons · 1 done/),
      expect.stringMatching(/^2\. Alpha chapter1 lessons · 0 done/),
    ]);
    // The first chapter with an unpassed lesson is current and open.
    expect(toggles[0]).toHaveAttribute('aria-expanded', 'true');
    expect(toggles[1]).toHaveAttribute('aria-expanded', 'false');
    const lessons = within(screen.getByRole('list', { name: 'Chapter 1 lessons' })).getAllByRole('link');
    expect(lessons.map((l) => l.getAttribute('aria-label'))).toEqual(['1. Second, next', '2. First, done']);
    expect(lessons[0]).toHaveAttribute('href', '/learn/lc/b2');

    // Progress, continue and XP.
    expect(screen.getByText('1 of 3 lessons')).toBeInTheDocument();
    expect(screen.getByRole('link', { name: 'Continue: Second' })).toHaveAttribute('href', '/learn/lc/b2');
    expect(screen.getByText('9')).toBeInTheDocument();
    expect(screen.getByText('of about 15')).toBeInTheDocument();
    expect(screen.getByText('prompts')).toBeInTheDocument();
  });

  it('opens a collapsed chapter on click', async () => {
    stubApi({ 'GET /api/v1/courses/lc': [200, course] });
    renderRoute('/courses/lc', '/courses/:courseId', <CoursePage />);
    await userEvent.click(await screen.findByRole('button', { name: /2\. Alpha chapter/ }));
    expect(within(screen.getByRole('list', { name: 'Chapter 2 lessons' })).getByRole('link')).toHaveAttribute('href', '/learn/lc/a1');
  });

  it('offers a held Skeleton Key on the first locked chapter only', async () => {
    const locked = {
      ...course,
      chapters: [
        { id: 'a', title: 'Alpha chapter', exercises: [ex('a1', 'Only', 'passed')] },
        { id: 'b', title: 'Beta chapter', locked: true, exercises: [ex('b1', 'First')] },
        { id: 'c', title: 'Gamma chapter', locked: true, exercises: [ex('c1', 'First')] },
      ],
    };
    const shop = { balance: 0, earned_today: 0, effects: [], bet: null, earn: [], items: [{ id: 'skeleton_key', owned: 1 }] };
    const fetch = stubApi({
      'GET /api/v1/courses/lc': [200, locked],
      'GET /api/v1/shop': [200, shop],
      'POST /api/v1/powerups/skeleton_key/use': [200, { powerup: 'skeleton_key', remaining: 0, balance: 0, unlocked_chapter: 'b' }],
    });
    renderRoute('/courses/lc', '/courses/:courseId', <CoursePage />);
    const keys = await screen.findAllByRole('button', { name: 'Use Skeleton Key' });
    expect(keys).toHaveLength(1);
    await userEvent.click(screen.getByRole('button', { name: 'Use Skeleton Key' }));
    const use = fetch.mock.calls.map(([r]) => r as Request).find((r) => r.method === 'POST');
    expect(await use?.clone().json()).toMatchObject({ chapter_id: 'b' });
  });

  it('shows not found for an unknown course', async () => {
    stubApi({});
    renderRoute('/courses/nope', '/courses/:courseId', <CoursePage />);
    expect(await screen.findByText('Not found')).toBeInTheDocument();
  });
});
