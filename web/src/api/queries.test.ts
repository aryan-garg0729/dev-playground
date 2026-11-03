import { stubApi } from '../test/render';
import { submitAndWait } from './queries';
import { firstPass } from '../test/fixtures';

const result = { status: 'pass', score: 1, stdout: '', runtime_ms: 5, checks: [{ name: 'a', passed: true, visible: true }] };

describe('submitAndWait', () => {
  it('returns a 200 result directly', async () => {
    stubApi({ 'POST /api/v1/exercises/e/submit': [200, { submission_id: 's', result, xp_awarded: 0 }] });
    await expect(submitAndWait('e', 1, { choices: ['a'] }, 1)).resolves.toMatchObject({ result: { status: 'pass' } });
  });

  it('polls the submission after a 202 until it is graded', async () => {
    const spy = stubApi({
      'POST /api/v1/exercises/e/submit': [202, { submission_id: 's', status: 'pending' }],
      'GET /api/v1/submissions/s': [
        [200, { submission_id: 's', exercise_id: 'e', exercise_version: 1, status: 'pending', result: null, xp_awarded: 0, created_at: '' }],
        [200, { submission_id: 's', exercise_id: 'e', exercise_version: 1, status: 'pass', result, xp_awarded: 35, outcome: firstPass, created_at: '' }],
      ],
    });
    const r = await submitAndWait('e', 1, { files: { 'main.py': '' } }, 1);
    expect(r).toEqual({ submission_id: 's', result, xp_awarded: 35, outcome: firstPass });
    expect(spy.mock.calls.filter(([req]) => (req as Request).url.endsWith('/submissions/s'))).toHaveLength(2);
  });

  it('surfaces API errors with their code', async () => {
    stubApi({ 'POST /api/v1/exercises/e/submit': [429, { error: { code: 'rate_limited', message: 'slow down' } }] });
    await expect(submitAndWait('e', 1, {}, 1)).rejects.toMatchObject({ code: 'rate_limited', status: 429 });
  });
});
