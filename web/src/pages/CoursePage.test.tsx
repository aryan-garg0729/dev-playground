import { screen, within } from '@testing-library/react';
import { renderRoute, stubApi } from '../test/render';
import { CoursePage } from './CoursePage';

const ex = (id: string, title: string, status = 'not_started') => ({
  id, title, type: 'mcq', xp: 5, difficulty: 1, free: true, boss: false, status,
});

describe('CoursePage', () => {
  it('lists chapters and lessons in the order the API returns (manifest order)', async () => {
    stubApi({
      'GET /api/v1/courses/lc': [200, {
        id: 'lc', title: 'LangChain Basics', description: 'Blocks.',
        chapters: [
          { id: 'b', title: 'Zeta chapter', exercises: [ex('b2', 'Second'), ex('b1', 'First', 'passed')] },
          { id: 'a', title: 'Alpha chapter', exercises: [ex('a1', 'Only')] },
        ],
      }],
    });
    renderRoute('/courses/lc', '/courses/:courseId', <CoursePage />);

    const chapters = await screen.findAllByRole('heading', { level: 2 });
    expect(chapters.map((h) => h.textContent)).toEqual(['1. Zeta chapter', '2. Alpha chapter']);

    const firstChapter = chapters[0]!.closest('li')!;
    const links = within(firstChapter).getAllByRole('link');
    expect(links.map((l) => l.getAttribute('aria-label'))).toEqual(['1. Second, not started', '2. First, done']);
    expect(links[0]).toHaveAttribute('href', '/learn/lc/b2');
    expect(screen.getByText(/3 lessons · 1 done/)).toBeInTheDocument();
  });

  it('shows not found for an unknown course', async () => {
    stubApi({});
    renderRoute('/courses/nope', '/courses/:courseId', <CoursePage />);
    expect(await screen.findByText('Not found')).toBeInTheDocument();
  });
});
