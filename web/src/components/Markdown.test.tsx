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

describe('Markdown headings and code blocks', () => {
  it('renders ## as a styled heading', () => {
    render(<Markdown>{'# Title\n\nText\n\n## Assignment\n\nDo it.'}</Markdown>);
    expect(screen.getByRole('heading', { level: 1, name: 'Title' })).toHaveClass('text-page-title');
    expect(screen.getByRole('heading', { level: 2, name: 'Assignment' })).toHaveClass('font-semibold');
  });

  it('copies a code block without the trailing newline', async () => {
    const user = userEvent.setup();
    const writeText = vi.spyOn(navigator.clipboard, 'writeText').mockResolvedValue();
    render(<Markdown>{'```sh\ndocker run hello-world\ndocker ps\n```'}</Markdown>);
    await user.click(screen.getByRole('button', { name: 'Copy code' }));
    expect(writeText).toHaveBeenCalledWith('docker run hello-world\ndocker ps');
    expect(await screen.findByRole('button', { name: 'Copied' })).toBeInTheDocument();
  });

  it('gives inline code no copy button', () => {
    render(<Markdown>{'Run `docker ps` now.'}</Markdown>);
    expect(screen.queryByRole('button', { name: 'Copy code' })).toBeNull();
  });
});

describe('Markdown dropTitle', () => {
  it('drops only a leading # title', () => {
    render(<Markdown dropTitle>{'# Image Explorer\n\nBoss time!\n\n## Assignment\n\nGo.'}</Markdown>);
    expect(screen.queryByRole('heading', { level: 1 })).toBeNull();
    expect(screen.getByText('Boss time!')).toBeInTheDocument();
    expect(screen.getByRole('heading', { level: 2, name: 'Assignment' })).toBeInTheDocument();
  });

  it('keeps a prompt that starts with text or ##', () => {
    render(<Markdown dropTitle>{'## Setup\n\nText'}</Markdown>);
    expect(screen.getByRole('heading', { level: 2, name: 'Setup' })).toBeInTheDocument();
  });
});

describe('Markdown links', () => {
  it('styles links and opens other sites in a new tab', () => {
    render(<Markdown>{'Solve [Towers](https://cses.fi/problemset/task/1073) or see [notes](#notes).'}</Markdown>);
    const ext = screen.getByRole('link', { name: 'Towers' });
    expect(ext).toHaveAttribute('target', '_blank');
    expect(ext).toHaveAttribute('rel', 'noopener noreferrer');
    expect(ext).toHaveClass('text-accent-text', 'underline');
    expect(screen.getByRole('link', { name: 'notes' })).not.toHaveAttribute('target');
  });
});

describe('Markdown code tabs', () => {
  const md = '```python\nprint(1)\n```\n\n```cpp\ncout << 1;\n```';

  it('merges back-to-back cpp and python blocks into tabs, C++ first and selected', () => {
    const { container } = render(<Markdown>{md}</Markdown>);
    const tabs = screen.getAllByRole('tab');
    expect(tabs.map((t) => t.textContent)).toEqual(['C++', 'Python']);
    expect(tabs[0]).toHaveAttribute('aria-selected', 'true');
    expect(container.querySelectorAll('pre')).toHaveLength(1);
    expect(screen.getByRole('tabpanel')).toHaveTextContent('cout << 1;');
  });

  it('switches every tabbed block when one tab is picked, and remembers it', async () => {
    render(<Markdown>{`${md}\n\nText\n\n${md}`}</Markdown>);
    await userEvent.click(screen.getAllByRole('tab', { name: 'Python' })[0]!);
    for (const panel of screen.getAllByRole('tabpanel')) expect(panel).toHaveTextContent('print(1)');
    expect(localStorage.getItem('cq-code-lang')).toBe('python');
  });

  it('moves between tabs with the arrow keys', async () => {
    render(<Markdown>{md}</Markdown>);
    screen.getByRole('tab', { name: 'C++' }).focus();
    await userEvent.keyboard('{ArrowRight}');
    expect(screen.getByRole('tab', { name: 'Python' })).toHaveFocus();
    expect(screen.getByRole('tabpanel')).toHaveTextContent('print(1)');
  });

  it('colours C++ and Python like the editor, and leaves other languages plain', () => {
    const { container } = render(<Markdown>{'```python\nreturn 1\n```\n\n```sh\necho return\n```'}</Markdown>);
    const [py, sh] = container.querySelectorAll('pre');
    expect(py?.querySelector('.text-code-keyword')).toHaveTextContent('return');
    expect(py?.querySelector('.text-code-number')).toHaveTextContent('1');
    expect(sh?.querySelector('span')).toBeNull();
  });

  it('leaves single blocks and other languages alone', () => {
    render(<Markdown>{'```python\nx = 1\n```\n\n```sh\nls\n```'}</Markdown>);
    expect(screen.queryByRole('tab')).not.toBeInTheDocument();
  });
});
