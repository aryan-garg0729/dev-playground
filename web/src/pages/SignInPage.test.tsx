import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { render, screen } from '@testing-library/react';
import { MemoryRouter, Route, Routes, useLocation } from 'react-router';
import { RequireAuth } from '../components/RequireAuth';
import { renderRoute, stubApi } from '../test/render';
import { SignInPage, safeNext } from './SignInPage';

const providers = (o: { github?: boolean; google?: boolean; dev?: boolean } = {}) =>
  [200, { github: !!o.github, google: !!o.google, dev: !!o.dev }] as [number, unknown];

describe('SignInPage', () => {
it('offers Google and GitHub with the return path and timezone, and dev sign-in only when enabled', async () => {
  stubApi({ 'GET /api/v1/auth/providers': providers({ github: true, google: true }) });
  renderRoute('/sign-in?next=%2Flearn%2Fc%2Fx', '/sign-in', <SignInPage />);
  const google = new URL((await screen.findByRole('link', { name: 'Continue with Google' })).getAttribute('href')!, 'http://x');
  expect(google.pathname).toBe('/api/v1/auth/google/start');
  expect(google.searchParams.get('next')).toBe('/learn/c/x');
  expect(google.searchParams.get('tz')).toBeTruthy();
  const github = screen.getByRole('link', { name: 'Continue with GitHub' }).getAttribute('href')!;
  expect(github.startsWith('/api/v1/auth/github/start?')).toBe(true);
  expect(screen.queryByRole('button', { name: /Dev sign-in/ })).not.toBeInTheDocument();
  expect(screen.queryByLabelText('Email')).not.toBeInTheDocument();
});

it('hides providers that are not configured', async () => {
  stubApi({ 'GET /api/v1/auth/providers': providers({ github: true }) });
  renderRoute('/sign-in', '/sign-in', <SignInPage />);
  await screen.findByRole('link', { name: 'Continue with GitHub' });
  expect(screen.queryByRole('link', { name: 'Continue with Google' })).not.toBeInTheDocument();
});

  it('shows a provider error', async () => {
    stubApi({ 'GET /api/v1/auth/providers': providers() });
    renderRoute('/sign-in?error=google', '/sign-in', <SignInPage />);
    expect(screen.getByRole('alert')).toHaveTextContent('Google sign-in did not complete');
  });

  it('only follows same-site next paths', () => {
    expect(safeNext('/learn/a')).toBe('/learn/a');
    expect(safeNext('//evil.example')).toBe('/');
    expect(safeNext('https://evil.example')).toBe('/');
    expect(safeNext('/\\evil.example')).toBe('/');
    expect(safeNext(null)).toBe('/');
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
