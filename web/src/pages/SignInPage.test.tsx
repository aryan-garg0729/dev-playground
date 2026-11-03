import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter, Route, Routes, useLocation } from 'react-router';
import { RequireAuth } from '../components/RequireAuth';
import { me } from '../test/fixtures';
import { renderRoute, stubApi } from '../test/render';
import { SignInPage, safeNext } from './SignInPage';
import { VerifyPage } from './VerifyPage';

const providers = (o: { github?: boolean; dev?: boolean } = {}) =>
  [200, { magic_link: true, github: !!o.github, dev: !!o.dev }] as [number, unknown];

describe('SignInPage', () => {
  it('emails a link and keeps the page to return to', async () => {
    const spy = stubApi({ 'GET /api/v1/auth/providers': providers(), 'POST /api/v1/auth/magic-link': [202, { status: 'sent' }] });
    renderRoute('/sign-in?next=%2Flearn%2Fc%2Fx', '/sign-in', <SignInPage />);
    await userEvent.type(screen.getByLabelText('Email'), 'ada@example.com');
    await userEvent.click(screen.getByRole('button', { name: 'Email me a link' }));
    expect(await screen.findByRole('heading', { name: 'Check your email' })).toBeInTheDocument();
    const req = spy.mock.calls.map(([r]) => r as Request).find((r) => r.url.endsWith('/auth/magic-link'))!;
    expect(await req.json()).toEqual({ email: 'ada@example.com', next: '/learn/c/x' });
  });

  it('offers GitHub and dev sign-in only when enabled', async () => {
    stubApi({ 'GET /api/v1/auth/providers': providers({ github: true }) });
    renderRoute('/sign-in', '/sign-in', <SignInPage />);
    expect(await screen.findByRole('link', { name: 'Continue with GitHub' })).toHaveAttribute('href', '/api/v1/auth/github/start?next=%2F');
    expect(screen.queryByRole('button', { name: /Dev sign-in/ })).not.toBeInTheDocument();
  });

  it('shows a GitHub error', async () => {
    stubApi({ 'GET /api/v1/auth/providers': providers() });
    renderRoute('/sign-in?error=github', '/sign-in', <SignInPage />);
    expect(screen.getByRole('alert')).toHaveTextContent('GitHub sign-in did not complete');
  });

  it('only follows same-site next paths', () => {
    expect(safeNext('/learn/a')).toBe('/learn/a');
    expect(safeNext('//evil.example')).toBe('/');
    expect(safeNext('https://evil.example')).toBe('/');
    expect(safeNext('/\\evil.example')).toBe('/');
    expect(safeNext(null)).toBe('/');
  });
});

describe('VerifyPage', () => {
  it('redeems the token from the fragment after a click and removes it from the URL', async () => {
    window.history.replaceState(null, '', '/sign-in/verify#token=abc123');
    const spy = stubApi({ 'POST /api/v1/auth/magic-link/verify': [200, { user: me.user, next: '/' }] });
    renderRoute('/sign-in/verify', '/sign-in/verify', <VerifyPage />);
    expect(window.location.hash).toBe('');
    expect(spy).not.toHaveBeenCalled(); // nothing is redeemed without the click
    await userEvent.click(screen.getByRole('button', { name: 'Continue' }));
    const req = spy.mock.calls.map(([r]) => r as Request).find((r) => r.url.endsWith('/verify'))!;
    expect(await req.json()).toMatchObject({ token: 'abc123' });
  });
});

describe('RequireAuth', () => {
  it('sends signed-out users to sign in, with the page to return to', async () => {
    stubApi({ 'GET /api/v1/me': [401, { error: { code: 'unauthorized', message: 'sign in first' } }] });
    render(
      <QueryClientProvider client={new QueryClient({ defaultOptions: { queries: { retry: false } } })}>
        <MemoryRouter initialEntries={['/league']}>
          <Routes>
            <Route path="/league" element={<RequireAuth>secret</RequireAuth>} />
            <Route path="/sign-in" element={<SearchEcho />} />
          </Routes>
        </MemoryRouter>
      </QueryClientProvider>,
    );
    expect(await screen.findByText('next=%2Fleague')).toBeInTheDocument();
    expect(screen.queryByText('secret')).not.toBeInTheDocument();
  });
});

function SearchEcho() {
  return <p>{useLocation().search.slice(1)}</p>;
}
