import { fireEvent, screen, within } from '@testing-library/react';
import type { CareerPath } from '../api/client';
import { course } from '../test/fixtures';
import { renderRoute, stubApi } from '../test/render';
import { CoursesPage, filterCourses } from './CoursesPage';
import { courseState } from './courses/CourseCard';

const courses = [
  course(),
  course({
    id: 'langchain-basics', title: 'LangChain Basics', description: 'Chat models and prompts.', icon: 'lc', level: 'intermediate',
    tags: ['Python', 'AI'], topics: 'Chat models, prompt templates', status: 'not_started', passed_count: 0, current_chapter: null,
  }),
  course({
    id: 'command-line', title: 'Command Line', description: 'On your own machine.', icon: 'cli', tags: ['CLI'], runs: ['cli'],
    topics: 'Shell and HTTP', status: 'not_started', passed_count: 0, current_chapter: null, free_count: 0,
  }),
];

const paths: CareerPath[] = [
  {
    id: 'ai-engineer', title: 'AI Engineer', description: 'From warm-ups to chat models.', yours: true, exercise_count: 20, passed_count: 2,
    steps: [
      { course_id: 'foundations', title: 'Foundations', note: 'Python and SQL basics', status: 'in_progress' },
      { course_id: 'langchain-basics', title: 'LangChain Basics', note: 'Chat models and prompts', status: 'not_started' },
    ],
  },
  {
    id: 'backend-developer', title: 'Backend Developer', description: '', yours: false, exercise_count: 8, passed_count: 0,
    steps: [{ course_id: 'command-line', title: 'Command Line', note: 'Shell and HTTP', status: 'not_started' }],
  },
];

describe('CoursesPage', () => {
  it('shows career paths and course cards with status and progress', async () => {
    stubApi({ 'GET /api/v1/courses': [200, { courses, paths }] });
    renderRoute('/courses', '/courses', <CoursesPage />);

    const mine = await screen.findByRole('listitem', { name: 'AI Engineer, your path' });
    expect(mine).toHaveTextContent('Your path');
    expect(mine).toHaveTextContent('FoundationsIn progress');
    expect(mine).toHaveTextContent('LangChain BasicsChat models and prompts');
    expect(mine).toHaveTextContent('10% · 2 courses');
    expect(within(mine).getByRole('link', { name: 'Continue' })).toHaveAttribute('href', '/courses/foundations');
    const other = screen.getByRole('listitem', { name: 'Backend Developer' });
    expect(within(other).getByRole('link', { name: 'Start path' })).toHaveAttribute('href', '/courses/command-line');

    const fd = screen.getByRole('link', { name: 'Foundations, In progress' });
    expect(fd).toHaveTextContent('2 of 6 lessons · Chapter 1 · Python warm-up');
    expect(fd).toHaveTextContent('Continue');
    expect(screen.getByRole('link', { name: 'LangChain Basics, Free to start' })).toHaveTextContent('Chat models, prompt templatesStart free');
    expect(screen.getByRole('link', { name: 'Command Line, Pro' })).toHaveTextContent('>_ CodeQuest CLI');
    expect(screen.getByRole('heading', { name: 'All courses 3 courses' })).toBeInTheDocument();
  });

  it('filters by topic and search, and "/" focuses search', async () => {
    stubApi({ 'GET /api/v1/courses': [200, { courses, paths }] });
    renderRoute('/courses', '/courses', <CoursesPage />);
    const all = await screen.findByRole('region', { name: /All courses/ });
    const filters = within(all).getByRole('group', { name: 'Filter by topic' });
    expect(within(filters).getAllByRole('button').map((b) => b.textContent)).toEqual(['All', 'Python', 'SQL', 'AI', 'CLI']);

    fireEvent.click(within(filters).getByRole('button', { name: 'AI' }));
    expect(within(filters).getByRole('button', { name: 'AI' })).toHaveAttribute('aria-pressed', 'true');
    expect(within(all).getAllByRole('listitem')).toHaveLength(1);

    fireEvent.click(within(filters).getByRole('button', { name: 'All' }));
    fireEvent.keyDown(window, { key: '/' });
    const search = screen.getByRole('searchbox', { name: 'Search courses and topics' });
    expect(search).toHaveFocus();
    fireEvent.change(search, { target: { value: 'nothing-matches' } });
    expect(all).toHaveTextContent('No courses match.');
  });

  it('filters and labels courses', () => {
    expect(filterCourses(courses, 'chat PROMPT', '').map((c) => c.id)).toEqual(['langchain-basics']);
    expect(filterCourses(courses, '', 'Python').map((c) => c.id)).toEqual(['foundations', 'langchain-basics']);
    expect(courseState(course({ status: 'done' })).action).toBe('Review');
    expect(courseState(course({ status: 'not_started', free_count: 0 })).pill).toBe('Pro');
  });
});
