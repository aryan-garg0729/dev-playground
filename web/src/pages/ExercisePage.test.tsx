import { screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { renderRoute, stubApi } from '../test/render';
import { ExercisePage } from './ExercisePage';
import { nothing } from '../test/fixtures';

const detail = {
  id: 'lc.q', type: 'mcq', version: 3, title: 'Message roles', xp: 5, difficulty: 1, tags: ['messages'],
  free: true, boss: false, course_id: 'lc', chapter: { id: 'ch', title: 'Chat models' },
  position: { index: 2, total: 4 }, prev_id: 'lc.p', next_id: null, last_submission: null,
  public_view: { question: 'Pick one', multiple: false, options: [{ id: 'a', text: 'Yes' }, { id: 'b', text: 'No' }] },
};

describe('ExercisePage', () => {
  it('renders lesson context and submits through the API', async () => {
    const fetchSpy = stubApi({
      'GET /api/v1/exercises/lc.q': [200, detail],
      'POST /api/v1/exercises/lc.q/submit': [200, {
        submission_id: 'x', xp_awarded: 0, outcome: nothing,
        result: { status: 'fail', score: 0, stdout: '', runtime_ms: 0, checks: [{ name: 'Correct answer', passed: false, visible: true, feedback: 'Because.' }] },
      }],
    });
    renderRoute('/learn/lc/lc.q', '/learn/:courseId/:exerciseId', <ExercisePage />);

    expect(await screen.findByRole('heading', { level: 1, name: 'Message roles' })).toBeInTheDocument();
    expect(screen.getByText('Chat models · lesson 2 of 4')).toBeInTheDocument();
    expect(screen.getByRole('link', { name: 'Previous lesson' })).toHaveAttribute('href', '/learn/lc/lc.p');
    expect(screen.queryByRole('link', { name: 'Next lesson' })).not.toBeInTheDocument();

    await userEvent.click(screen.getByLabelText('No'));
    await userEvent.click(screen.getByRole('button', { name: 'Submit' }));
    expect(await screen.findByText('Not quite')).toBeInTheDocument();
    expect(screen.getByText('Because.')).toBeInTheDocument();

    const post = fetchSpy.mock.calls.map(([req]) => req as Request).find((r) => r.method === 'POST')!;
    expect(post.headers.get('Content-Type')).toContain('application/json');
    expect(await post.json()).toEqual({ exercise_id: 'lc.q', version: 3, answer: { choices: ['b'] } });
  });
});
