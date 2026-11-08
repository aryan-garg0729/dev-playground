import { fireEvent, screen, within } from '@testing-library/react';
import type { Medal, Medals, Profile } from '../api/client';
import { course, me } from '../test/fixtures';
import { renderRoute, stubApi } from '../test/render';
import { medalLine, ProfilePage } from './ProfilePage';
import { projectLine } from './profile/ProjectList';

const tier = (t: string, gte: number, earned = false) => ({ tier: t, gte, title: `X ${t}`, bytes: 50, earned_at: earned ? '2026-11-01T00:00:00Z' : null });
const medal = (o: Partial<Medal>): Medal => ({
  id: 'm', family: 'consistency', name: 'Streak', description: '', unit: 'days', secret: false, coming_soon: false,
  value: 13, tier: 'bronze', next_tier: 'silver', next_gte: 30, tiers: [tier('bronze', 7, true), tier('silver', 30), tier('gold', 100)], ...o,
});

const medals: Medals = {
  earned: 1, total: 3, title: 'Streak I', titles: ['Streak I', 'Boss Slayer I'],
  families: [{ id: 'consistency', name: 'Consistency' }, { id: 'secret', name: 'Secret' }, { id: 'memory', name: 'Memory' }],
  medals: [
    medal({ id: 'streak' }),
    medal({ id: 'early', name: 'Early Bird', value: 0, tier: '', next_tier: 'bronze', next_gte: 3, unit: 'lessons before 8 am', tiers: [tier('bronze', 3)] }),
    medal({ id: 'secret1', family: 'secret', name: '???', secret: true, tier: '', next_tier: '', next_gte: 0, tiers: [] }),
    medal({ id: 'reviewer', family: 'memory', name: 'Reviewer', coming_soon: true, tier: '', value: 0 }),
  ],
};

const project = (o: Partial<Profile['projects'][number]>): Profile['projects'][number] => ({
  exercise_id: 'fd.py.fizzbuzz-boss', course_id: 'foundations', course_title: 'Foundations', title: 'FizzBuzz', chapter_index: 1,
  chapter_title: 'Python warm-up', status: 'defeated', first_try: true, defeated_at: '2026-10-29T10:00:00Z', ...o,
});

const profile: Profile = {
  joined_at: '2026-10-02T10:00:00Z',
  path: { id: 'backend-developer', title: 'Backend Developer' },
  projects_shipped: 1,
  activity: { from: '2025-10-30', to: '2026-11-04', longest_streak: 12, days: [{ date: '2026-11-03', lessons: 3 }, { date: '2026-11-04', lessons: 1 }] },
  projects: [
    project({}),
    project({ exercise_id: 'demo.basics.sum', course_id: 'demo', course_title: 'Demo', title: 'Sum a list', status: 'next', first_try: false, defeated_at: null }),
    project({ exercise_id: 'x.boss', title: 'Later boss', status: 'upcoming', first_try: false, defeated_at: null, chapter_index: 2 }),
  ],
  skills: {
    course_id: 'foundations', course_title: 'Foundations',
    chapters: [{ title: 'Python warm-up', tier: 'gold', passed: 3, lessons: 3 }, { title: 'SQL basics', tier: '', passed: 0, lessons: 3 }],
  },
  recent: [{ date: '2026-11-04', exercise_id: 'fd.py.sum-list', title: 'Sum a list', boss: false, course_id: 'foundations', course_title: 'Foundations' }],
};

const routes = () => ({
  'GET /api/v1/me': [200, me] as [number, unknown],
  'GET /api/v1/medals': [200, medals] as [number, unknown],
  'GET /api/v1/profile': [200, profile] as [number, unknown],
  'GET /api/v1/courses': [200, { courses: [course()], paths: [] }] as [number, unknown],
});

describe('ProfilePage', () => {
  it('shows the overview: stats, year of activity, projects, skills and recent medals', async () => {
    stubApi(routes());
    renderRoute('/profile', '/profile/:tab?', <ProfilePage />);
    expect(await screen.findByText('Backend Developer path · joined Oct 2026')).toBeInTheDocument();
    expect(screen.getByRole('heading', { level: 1 })).toHaveTextContent('Boss Slayer I');
    expect(screen.getByText('Projects shipped').parentElement).toHaveTextContent('Projects shipped1');
    expect(screen.getByRole('link', { name: 'Overview' })).toHaveAttribute('aria-current', 'page');

    expect(screen.getByRole('region', { name: '4 lessons in the last year' })).toHaveTextContent('2 active days · longest streak 12 days');

    const projects = screen.getByRole('region', { name: 'Projects completed' });
    expect(within(projects).getByRole('link', { name: 'View code' })).toHaveAttribute('href', '/learn/foundations/fd.py.fizzbuzz-boss');
    expect(projects).toHaveTextContent('Foundations · Chapter 1 boss · defeated first try · Oct 29');
    expect(projects).toHaveTextContent('Sum a list');
    expect(projects).not.toHaveTextContent('Later boss');

    const skills = screen.getByRole('region', { name: 'Skills' });
    expect(skills).toHaveTextContent('Python warm-upGold');
    expect(skills).toHaveTextContent('SQL basicsNot started');

    const recent = screen.getByRole('region', { name: 'Recent medals' });
    expect(within(recent).getByRole('listitem', { name: 'Streak, bronze' })).toBeInTheDocument();
    expect(within(recent).getByRole('link', { name: 'All 1' })).toHaveAttribute('href', '/profile/medals');
  });

  it('edits the name and title', async () => {
    const spy = stubApi({ ...routes(), 'PATCH /api/v1/me': [200, me] });
    renderRoute('/profile', '/profile/:tab?', <ProfilePage />);
    fireEvent.click(await screen.findByRole('button', { name: 'Edit profile' }));
    const form = screen.getByRole('form', { name: 'Edit profile' });
    expect(within(form).getByRole('combobox', { name: 'Title' })).toHaveValue('Streak I');
    fireEvent.change(within(form).getByRole('textbox', { name: 'Name' }), { target: { value: '  Ada L ' } });
    fireEvent.change(within(form).getByRole('combobox', { name: 'Title' }), { target: { value: 'Boss Slayer I' } });
    fireEvent.click(within(form).getByRole('button', { name: 'Save' }));
    await screen.findByRole('button', { name: 'Edit profile' });
    const patch = spy.mock.calls.find(([r]) => r instanceof Request && r.method === 'PATCH')?.[0] as Request;
    expect(await patch.json()).toEqual({ name: 'Ada L', title: 'Boss Slayer I' });
  });

  it('shows closest medals, families, masked secrets on the Medals tab', async () => {
    stubApi(routes());
    renderRoute('/profile/medals', '/profile/:tab?', <ProfilePage />);
    const closest = await screen.findByRole('region', { name: 'Closest to unlocking' });
    expect(within(closest).getByRole('img', { name: '13 of 30' })).toBeInTheDocument();
    expect(screen.getByRole('listitem', { name: 'Streak, bronze' })).toHaveTextContent('Bronze · Silver at 30 days');
    expect(screen.getByRole('listitem', { name: '???, not earned' })).toHaveTextContent('Hidden until earned');
    expect(screen.getByRole('listitem', { name: 'Reviewer, not earned' })).toHaveTextContent('Coming soon');
  });

  it('lists every boss by course on the Projects tab', async () => {
    stubApi(routes());
    renderRoute('/profile/projects', '/profile/:tab?', <ProfilePage />);
    expect(await screen.findByRole('region', { name: 'Foundations 1 of 2 shipped' })).toHaveTextContent('Later boss');
    expect(screen.getByRole('region', { name: 'Demo 0 of 1 shipped' })).toHaveTextContent('Next boss · Chapter 1: Python warm-up');
  });

  it('shows recent passes on the Activity tab', async () => {
    stubApi(routes());
    renderRoute('/profile/activity', '/profile/:tab?', <ProfilePage />);
    const recent = await screen.findByRole('region', { name: 'Recently passed' });
    expect(within(recent).getByRole('link', { name: 'Sum a list' })).toHaveAttribute('href', '/learn/foundations/fd.py.sum-list');
    expect(recent).toHaveTextContent('Wed, Nov 4');
  });

  it('describes projects', () => {
    expect(projectLine(project({ first_try: false }))).toBe('Foundations · Chapter 1 boss · Oct 29');
    expect(projectLine(project({ status: 'next', defeated_at: null }))).toBe('Foundations · Next boss · Chapter 1: Python warm-up');
  });

  it('describes progress', () => {
    expect(medalLine(medal({}))).toBe('Bronze · Silver at 30 days');
    expect(medalLine(medal({ tier: 'gold', next_tier: '' }))).toBe('Gold · all tiers earned');
    expect(medalLine(medal({ tier: '', next_tier: 'bronze', next_gte: 3, unit: 'late lessons' }))).toBe('3 late lessons');
    expect(medalLine(medal({ tier: '', next_tier: 'bronze', next_gte: 1, unit: 'bosses' }))).toBe('1 boss');
    expect(medalLine(medal({ tier: '', next_tier: 'bronze', next_gte: 2, unit: 'tier' }))).toBe('Reach Silver league');
    expect(medalLine(medal({ tier: 'bronze', next_tier: 'silver', next_gte: 3, unit: 'tier' }))).toBe('Bronze · Silver: reach Gold league');
  });
});
