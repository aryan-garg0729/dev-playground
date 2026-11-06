import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { LessonAssetsProvider } from './LessonAssets';
import { Markdown } from './Markdown';

const assets = { 'assets/loop.svg': '/api/v1/assets/abc.svg' };

function show(md: string, a: Record<string, string> = assets) {
  return render(
    <LessonAssetsProvider assets={a}>
      <Markdown>{md}</Markdown>
    </LessonAssetsProvider>,
  );
}

describe('Markdown images', () => {
  it('resolves a lesson image (also written ./assets/...)', () => {
    show('![How the loop runs](./assets/loop.svg)');
    const img = screen.getByRole('img', { name: 'How the loop runs' });
    expect(img).toHaveAttribute('src', '/api/v1/assets/abc.svg');
    expect(img).toHaveAttribute('loading', 'lazy');
  });

  it('shows alt text instead of external or unknown images', () => {
    show('![Remote](https://example.com/x.png) ![Missing](assets/nope.png) ![](//cdn.test/y.png)');
    expect(screen.queryAllByRole('img')).toHaveLength(0);
    expect(screen.getByText('[Remote]')).toBeInTheDocument();
    expect(screen.getByText('[Missing]')).toBeInTheDocument();
    expect(screen.getByText('[image]')).toBeInTheDocument();
  });

  it('loads no images at all outside a lesson (tutor replies)', () => {
    render(<Markdown>{'![x](assets/loop.svg)'}</Markdown>);
    expect(screen.queryByRole('img')).not.toBeInTheDocument();
  });

  it('opens a full-size view that closes with Escape', async () => {
    show('![Diagram](assets/loop.svg)');
    await userEvent.click(screen.getByRole('button', { name: 'Enlarge image: Diagram' }));
    expect(screen.getByRole('dialog', { name: 'Diagram' })).toBeInTheDocument();
    await userEvent.keyboard('{Escape}');
    expect(screen.queryByRole('dialog')).not.toBeInTheDocument();
  });

  it('never renders raw HTML from the source', () => {
    const { container } = show('<img src="assets/loop.svg" onerror="alert(1)">');
    expect(container.querySelector('img')).toBeNull();
  });
});
