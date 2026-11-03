import { screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { me } from '../test/fixtures';
import { renderRoute, stubApi } from '../test/render';
import { LandingPage } from './LandingPage';
import { OnboardingPage } from './OnboardingPage';

const courses = {
  courses: [
    { id: 'foundations', title: 'Foundations', description: 'Python and SQL from zero.', chapter_count: 2, exercise_count: 6, passed_count: 0 },
    { id: 'langchain-basics', title: 'LangChain Basics', description: 'Chat models and chains.', chapter_count: 3, exercise_count: 9, passed_count: 0 },
  ],
};

async function patchBody(spy: ReturnType<typeof stubApi>) {
  await waitFor(() => expect(spy.mock.calls.some(([r]) => (r as Request).method === 'PATCH')).toBe(true));
  const req = spy.mock.calls.map(([r]) => r as Request).find((r) => r.method === 'PATCH')!;
  return req.json();
}

describe('OnboardingPage', () => {
  it('saves the goal and daily goal and opens the course picked for that goal', async () => {
    const spy = stubApi({ 'GET /api/v1/courses': [200, courses], 'PATCH /api/v1/me': [200, me] });
    renderRoute('/welcome', '/welcome', <OnboardingPage />);
    expect(await screen.findByText('Step 1 of 2')).toBeInTheDocument();
    await userEvent.click(screen.getByRole('radio', { name: /Build with AI/ }));
    expect(screen.getByRole('radio', { name: /Build with AI/ })).toBeChecked();
    expect(screen.getByText('Starts with LangChain Basics')).toBeInTheDocument();
    await userEvent.click(screen.getByRole('radio', { name: /40 min/ }));
    await userEvent.click(screen.getByRole('button', { name: 'Continue' }));

    expect(await screen.findByText('Step 2 of 2')).toBeInTheDocument();
    expect(screen.getByRole('radio', { name: /LangChain Basics/ })).toBeChecked();
    await userEvent.click(screen.getByRole('button', { name: 'Start learning' }));
    expect(await patchBody(spy)).toEqual({ goal: 'ai', daily_goal: 4, onboarded: true });
  });

  it('can be skipped', async () => {
    const spy = stubApi({ 'GET /api/v1/courses': [200, courses], 'PATCH /api/v1/me': [200, me] });
    renderRoute('/welcome', '/welcome', <OnboardingPage />);
    await userEvent.click(await screen.findByRole('button', { name: 'Skip and start Foundations' }));
    expect(await patchBody(spy)).toEqual({ onboarded: true });
  });
});

describe('LandingPage', () => {
  it('shows the hero, real courses and the try-it card', async () => {
    stubApi({ 'GET /api/v1/courses': [200, courses] });
    renderRoute('/', '/', <LandingPage />);
    expect(screen.getByRole('heading', { level: 1 })).toHaveTextContent('Learn to codeby writing code.');
    expect(await screen.findByRole('heading', { name: 'LangChain Basics' })).toBeInTheDocument();
    expect(screen.getAllByRole('link', { name: 'Start the path →' })[0]).toHaveAttribute('href', '/sign-in?next=%2Fcourses%2Ffoundations');
    await userEvent.click(screen.getByRole('button', { name: 'Run checks' }));
    expect(screen.getByRole('status')).toHaveTextContent('2 of 2 checks passing');
  });
});
