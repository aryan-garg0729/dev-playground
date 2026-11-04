import { screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { renderRoute, stubApi } from '../test/render';
import { CliLoginPage } from './CliLoginPage';

describe('CliLoginPage', () => {
  it('shows a one-time code once, with its expiry', async () => {
    const spy = stubApi({ 'POST /api/v1/cli/codes': [200, { code: 'ABCDE-FGHJK', expires_at: new Date(Date.now() + 600_000).toISOString() }] });
    renderRoute('/cli/login', '/cli/login', <CliLoginPage />);
    expect(await screen.findByTestId('cli-code')).toHaveTextContent('ABCDE-FGHJK');
    expect(screen.getByText(/Expires in (10:00|9:5\d) · works once/)).toBeInTheDocument();
    const posts = spy.mock.calls.filter(([r]) => (r as Request).method === 'POST');
    expect(posts).toHaveLength(1);
  });

  it('offers a new code after it expires', async () => {
    const spy = stubApi({
      'POST /api/v1/cli/codes': [
        [200, { code: 'OLDDD-CODEE', expires_at: new Date(Date.now() - 1000).toISOString() }],
        [200, { code: 'NEWWW-CODEE', expires_at: new Date(Date.now() + 600_000).toISOString() }],
      ],
    });
    renderRoute('/cli/login', '/cli/login', <CliLoginPage />);
    await userEvent.click(await screen.findByRole('button', { name: 'Get a new code' }));
    expect(await screen.findByTestId('cli-code')).toHaveTextContent('NEWWW-CODEE');
    await waitFor(() => expect(spy.mock.calls.filter(([r]) => (r as Request).method === 'POST')).toHaveLength(2));
  });

  it('shows the error when the code cannot be created', async () => {
    stubApi({ 'POST /api/v1/cli/codes': [500, { error: { code: 'internal', message: 'something went wrong' } }] });
    renderRoute('/cli/login', '/cli/login', <CliLoginPage />);
    expect(await screen.findByRole('alert')).toHaveTextContent('something went wrong');
  });
});
