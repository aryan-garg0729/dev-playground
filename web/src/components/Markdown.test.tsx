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

describe('Markdown lists', () => {
  it('renders bullets and numbers with markers, keeping the start number', () => {
    const { container } = render(<Markdown>{'- one\n- two\n\nText\n\n3. three\n4. four'}</Markdown>);
    const ul = container.querySelector('ul');
    const ol = container.querySelector('ol');
    expect(ul).toHaveClass('list-disc');
    expect(ul?.querySelectorAll('li')).toHaveLength(2);
    expect(ol).toHaveClass('list-decimal');
    expect(ol).toHaveAttribute('start', '3');
  });

  it('keeps code blocks inside numbered steps', () => {
    const { container } = render(<Markdown>{'1. Run this:\n\n   ```sh\n   git status\n   ```\n\n2. Submit.'}</Markdown>);
    const items = container.querySelectorAll('ol > li');
    expect(items).toHaveLength(2);
    expect(items[0]?.querySelector('pre')).toHaveTextContent('git status');
  });

  it('shows task lists without bullets', () => {
    const { container } = render(<Markdown>{'- [x] done\n- [ ] todo'}</Markdown>);
    expect(container.querySelector('ul')).not.toHaveClass('list-disc');
    expect(container.querySelectorAll('input[type="checkbox"]')).toHaveLength(2);
  });
});

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

describe('Markdown tables', () => {
  it('renders GitHub-flavoured tables, with inline code in cells', () => {
    show('| Signal | Number |\n| --- | --- |\n| SIGKILL | `9` |\n| SIGTERM | 15 |');
    const table = screen.getByRole('table');
    expect(screen.getAllByRole('columnheader').map((h) => h.textContent)).toEqual(['Signal', 'Number']);
    expect(screen.getAllByRole('row')).toHaveLength(3);
    expect(table).toHaveTextContent('SIGKILL9');
    expect(screen.getByText('9').tagName).toBe('CODE');
  });
});
