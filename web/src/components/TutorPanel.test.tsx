import { screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { vi } from 'vitest';
import { parseSSE } from '../api/tutor';
import { me } from '../test/fixtures';
import { renderRoute } from '../test/render';
import { TutorPanel } from './TutorPanel';

const json = (status: number, body: unknown) =>
  new Response(JSON.stringify(body), { status, headers: { 'Content-Type': 'application/json' } });
const sse = (text: string) => new Response(text, { status: 200, headers: { 'Content-Type': 'text/event-stream' } });

/** Stubs the tutor endpoints; `hint` answers each POST /hint in turn. */
function stubTutor(tutorOn: boolean, hint: Response[] = [], thread = { daily_limit: 5, remaining: 5, messages: [] as unknown[] }) {
  const posts: { url: string; body: unknown }[] = [];
  vi.spyOn(globalThis, 'fetch').mockImplementation(async (input: RequestInfo | URL, init?: RequestInit) => {
    const url = input instanceof Request ? input.url : String(input);
    const method = input instanceof Request ? input.method : (init?.method ?? 'GET');
    if (url.endsWith('/api/v1/me')) return json(200, { ...me, features: { tutor: tutorOn } });
    if (method === 'GET' && url.endsWith('/tutor')) {
      const saved = posts.length > 0 && hint.length === 0 && thread.messages.length === 0
        ? { daily_limit: 5, remaining: 4, messages: [
          { role: 'user', content: 'why does it fail?', created_at: '2026-11-03T10:00:00Z' },
          { role: 'assistant', content: 'What does line 2 return?', created_at: '2026-11-03T10:00:01Z' }] }
        : thread;
      return json(200, saved);
    }
    if (method === 'POST' && url.endsWith('/hint')) {
      posts.push({ url, body: JSON.parse(String(init?.body)) });
      return hint.shift() ?? json(500, {});
    }
    return json(404, { error: { code: 'not_found', message: 'not found' } });
  });
  return posts;
}

describe('TutorPanel', () => {
  it('is a disabled "Coming soon" button when the server has no tutor', async () => {
    stubTutor(false);
    renderRoute('/', '/', <TutorPanel exerciseId="ex1" code="x = 1" />);
    const btn = await screen.findByRole('button', { name: 'Ask the tutor' });
    expect(btn).toBeDisabled();
    expect(btn).toHaveAttribute('title', 'Coming soon');
  });

  it('streams a reply with the learner code and shows hints left', async () => {
    const reply = 'event: delta\ndata: {"text":"What does "}\n\nevent: delta\ndata: {"text":"line 2 return?"}\n\nevent: done\ndata: {"remaining":4,"used_hint":true}\n\n';
    const posts = stubTutor(true, [sse(reply)]);
    renderRoute('/', '/', <TutorPanel exerciseId="ex1" code={'def f():\n    return 0'} />);
    await userEvent.click(await screen.findByRole('button', { name: 'Ask the tutor' }));
    expect(await screen.findByText('5 of 5 hints left today')).toBeInTheDocument();
    expect(screen.getByText(/removes this lesson's no-hint bonus/)).toBeInTheDocument();

    await userEvent.type(screen.getByRole('textbox', { name: 'Ask the tutor' }), 'why does it fail?{Enter}');
    await waitFor(() => expect(posts).toHaveLength(1));
    expect(posts[0]!.url).toMatch(/\/api\/v1\/exercises\/ex1\/hint$/);
    expect(posts[0]!.body).toEqual({ message: 'why does it fail?', code: 'def f():\n    return 0' });
    expect(await screen.findByText('What does line 2 return?')).toBeInTheDocument();
    expect(await screen.findByText('4 of 5 hints left today')).toBeInTheDocument();
  });

  it('shows the daily limit and keeps the box locked', async () => {
    stubTutor(true, [json(429, { error: { code: 'limit_reached', message: "You've used today's 5 tutor hints." } })]);
    renderRoute('/', '/', <TutorPanel exerciseId="ex1" code="" />);
    await userEvent.click(await screen.findByRole('button', { name: 'Ask the tutor' }));
    await userEvent.type(await screen.findByRole('textbox', { name: 'Ask the tutor' }), 'hint{Enter}');
    expect(await screen.findByRole('alert')).toHaveTextContent("today's 5 tutor hints");
  });
});

describe('parseSSE', () => {
  it('splits complete events and keeps the unfinished tail', () => {
    const { events, rest } = parseSSE('event: delta\ndata: {"text":"a"}\n\nevent: done\ndata: {"rem');
    expect(events).toEqual([{ event: 'delta', data: '{"text":"a"}' }]);
    expect(rest).toBe('event: done\ndata: {"rem');
  });
});
