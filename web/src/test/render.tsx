import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { render } from '@testing-library/react';
import type { ReactElement } from 'react';
import { MemoryRouter, Route, Routes } from 'react-router';
import { vi } from 'vitest';

type Reply = [number, unknown];

/**
 * Stubs fetch with "METHOD /path" -> [status, body]. A list of replies is served
 * in order, repeating the last one. Unknown paths return 404.
 */
export function stubApi(routes: Record<string, Reply | Reply[]>) {
  const served: Record<string, number> = {};
  return vi.spyOn(globalThis, 'fetch').mockImplementation(async (input: RequestInfo | URL) => {
    const url = new URL(input instanceof Request ? input.url : String(input));
    const key = `${input instanceof Request ? input.method : 'GET'} ${url.pathname}`;
    const route = routes[key];
    let hit: Reply | undefined;
    if (route && Array.isArray(route[0])) {
      const list = route as Reply[];
      const i = served[key] ?? 0;
      served[key] = i + 1;
      hit = list[Math.min(i, list.length - 1)];
    } else {
      hit = route as Reply | undefined;
    }
    const [status, body] = hit ?? [404, { error: { code: 'not_found', message: 'not found' } }];
    return new Response(JSON.stringify(body), { status, headers: { 'Content-Type': 'application/json' } });
  });
}

export function renderRoute(path: string, pattern: string, element: ReactElement) {
  const qc = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <QueryClientProvider client={qc}>
      <MemoryRouter initialEntries={[path]}>
        <Routes>
          <Route path={pattern} element={element} />
        </Routes>
      </MemoryRouter>
    </QueryClientProvider>,
  );
}
