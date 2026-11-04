import { screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { renderRoute, stubApi } from '../../test/render';
import { CliExercise } from './CliExercise';
import type { CliPublicView } from './steps';

const view: CliPublicView = {
  lesson_id: 'cl.local-server',
  prompt: 'Start a **server**.',
  base_url: 'http://localhost:8000',
  steps: [
    { description: 'Run it', cli_command: { command: 'python3 hello.py', visible_tests: [{ exit_code: 0 }, { stdout_contains_all: ['Hi'] }] } },
    {
      description: 'Fetch it',
      http_request: {
        request: { method: 'GET', url: '${baseURL}/info.json' },
        visible_tests: [{ status_code: 200 }, { json_body: { path: '.lessons', op: 'gt', value: 1 } }, { header: { name: 'Content-Type', contains: 'json' } }],
      },
    },
  ],
};

const exercise = (result: unknown) => ({
  id: 'cl.local-server', version: 1,
  last_submission: result ? { submission_id: 's1', exercise_version: 1, answer: {}, result, created_at: '2026-11-04T00:00:00Z' } : null,
});

const failing = {
  status: 'fail', score: 0.5, stdout: '', runtime_ms: 0,
  checks: [
    { name: 'Step 1: exit code is 0', passed: true, visible: true },
    { name: 'Step 2: status code is 200', passed: false, visible: true, got: '404', expected: '200' },
  ],
};

function show() {
  renderRoute('/l', '/l', <CliExercise header={<h1>Serve a file</h1>} publicView={view} draft={undefined} result={undefined} onChange={vi.fn()} onSubmit={vi.fn()} />);
}

describe('CliExercise', () => {
  it('shows the prompt, the commands to run and every test, with no editor or Submit', async () => {
    stubApi({ 'GET /api/v1/exercises/cl.local-server': [200, exercise(null)] });
    show();
    expect(screen.getByText('server').tagName).toBe('STRONG');
    expect(screen.getByText('codequest run cl.local-server')).toBeInTheDocument();
    expect(screen.getByText('codequest submit cl.local-server')).toBeInTheDocument();
    expect(screen.getByText('$ python3 hello.py')).toBeInTheDocument();
    expect(screen.getByText('> GET ${baseURL}/info.json')).toBeInTheDocument();
    expect(screen.getByText('output contains "Hi"')).toBeInTheDocument();
    expect(screen.getByText('JSON .lessons is greater than 1')).toBeInTheDocument();
    expect(screen.getByText('header Content-Type contains "json"')).toBeInTheDocument();
    expect(screen.queryByRole('button', { name: 'Submit' })).not.toBeInTheDocument();
    expect(await screen.findByText(/No submission yet/)).toBeInTheDocument();
  });

  it('picks up a result submitted from the terminal', async () => {
    stubApi({ 'GET /api/v1/exercises/cl.local-server': [[200, exercise(null)], [200, exercise(failing)]] });
    show();
    await userEvent.click(await screen.findByRole('button', { name: 'Check for new results' }));
    expect(await screen.findByText('1 of 2 checks failing')).toBeInTheDocument();
    expect(screen.getByText('404')).toBeInTheDocument();
  });
});
