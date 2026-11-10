import { useId, useRef, useState, type KeyboardEvent, type ReactNode } from 'react';
import { createPortal } from 'react-dom';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import type { Element } from 'hast';
import { useCodeLang } from '../stores/codeLang';
import { highlight, HIGHLIGHT_LANGS, TOKEN_CLASS } from './highlight';
import { ImageLightbox } from './ImageLightbox';
import { useLessonAsset } from './LessonAssets';
import { remarkCodeTabs, TAB_LANGS, TABS_LANG, type CodeTab } from './remarkCodeTabs';

/** A lesson image: framed so transparent diagrams read in both themes; click to enlarge. */
function LessonImage({ src, alt }: { src: string; alt: string }) {
  const [open, setOpen] = useState(false);
  return (
    <>
      {/* Full width so SVGs with only a viewBox (no width/height) still get a size; the ratio is kept. */}
      <button type="button" onClick={() => setOpen(true)} aria-label={`Enlarge image: ${alt || 'image'}`} className="block w-full">
        <img
          src={src}
          alt={alt}
          loading="lazy"
          decoding="async"
          className="block h-auto max-h-[420px] w-auto max-w-full rounded-md border border-border bg-surface-raised p-2"
        />
      </button>
      {/* Portal: the image sits inside a <p>, which cannot hold the dialog's divs. */}
      {open && createPortal(<ImageLightbox src={src} alt={alt} onClose={() => setOpen(false)} />, document.body)}
    </>
  );
}

/** A fenced code block with a Copy button (copies the block's text exactly). */
function CodeBlock({ children }: { children: ReactNode }) {
  const ref = useRef<HTMLPreElement>(null);
  const [copied, setCopied] = useState(false);
  const copy = () => {
    const text = ref.current?.textContent ?? '';
    navigator.clipboard.writeText(text.replace(/\n$/, '')).then(
      () => {
        setCopied(true);
        setTimeout(() => setCopied(false), 1500);
      },
      () => setCopied(false),
    );
  };
  return (
    <div className="relative">
      <pre
        ref={ref}
        className="text-code-sm overflow-x-auto border-l-2 border-border-strong pl-[14px] pr-16 text-text-secondary [&_code]:bg-transparent [&_code]:p-0 [&_code]:text-text-secondary"
      >
        {children}
      </pre>
      <button
        type="button"
        onClick={copy}
        aria-label={copied ? 'Copied' : 'Copy code'}
        className="text-caption absolute right-0 top-0 rounded-sm border border-border bg-surface-raised px-2 py-1 text-text-muted hover:text-text"
      >
        {copied ? 'Copied' : 'Copy'}
      </button>
    </div>
  );
}

/** The language and text of a fenced block's `<code>` child, if it has a language. */
function fenced(node: Element | undefined): { lang: string; text: string } | null {
  const code = node?.children[0];
  if (code?.type !== 'element' || code.tagName !== 'code') return null;
  const classes = code.properties.className;
  const cls = Array.isArray(classes) ? classes.find((c) => String(c).startsWith('language-')) : undefined;
  const text = code.children[0];
  if (!cls || text?.type !== 'text') return null;
  return { lang: String(cls).slice('language-'.length), text: text.value.replace(/\n$/, '') };
}

/** The tabs of a merged code block (see remarkCodeTabs), or null for an ordinary block. */
function tabsOf(block: { lang: string; text: string } | null): CodeTab[] | null {
  if (block?.lang !== TABS_LANG) return null;
  try {
    const tabs = JSON.parse(block.text) as CodeTab[];
    return Array.isArray(tabs) && tabs.length > 0 ? tabs : null;
  } catch {
    return null;
  }
}

/** Code coloured like the editor (C++ and Python; other languages stay plain). */
function Highlighted({ code, lang }: { code: string; lang: string }) {
  return (
    <code>
      {highlight(code, lang).map((t, i) =>
        TOKEN_CLASS[t.kind] ? (
          <span key={i} className={TOKEN_CLASS[t.kind]}>
            {t.text}
          </span>
        ) : (
          t.text
        ),
      )}
    </code>
  );
}

/**
 * The same code in several languages, one tab each. The picked language is shared by
 * every tabbed block and remembered per device (C++ until the learner picks another).
 */
function CodeTabs({ tabs }: { tabs: CodeTab[] }) {
  const { lang, setLang } = useCodeLang();
  const id = useId();
  const current = tabs.find((t) => t.lang === lang) ?? tabs[0]!;
  const tabRefs = useRef<(HTMLButtonElement | null)[]>([]);
  const onKeyDown = (e: KeyboardEvent) => {
    if (e.key !== 'ArrowRight' && e.key !== 'ArrowLeft') return;
    e.preventDefault();
    const i = tabs.indexOf(current);
    const next = (i + (e.key === 'ArrowRight' ? 1 : tabs.length - 1)) % tabs.length;
    setLang(tabs[next]!.lang);
    tabRefs.current[next]?.focus();
  };
  return (
    <div className="flex flex-col gap-2">
      <div role="tablist" aria-label="Code language" onKeyDown={onKeyDown} className="flex gap-1 border-b border-border">
        {tabs.map((t, i) => {
          const selected = t === current;
          return (
            <button
              key={t.lang}
              ref={(el) => {
                tabRefs.current[i] = el;
              }}
              type="button"
              role="tab"
              id={`${id}-tab-${t.lang}`}
              aria-selected={selected}
              aria-controls={`${id}-panel`}
              tabIndex={selected ? 0 : -1}
              onClick={() => setLang(t.lang)}
              className={`text-caption -mb-px rounded-t-sm border px-3 py-1.5 ${
                selected
                  ? 'border-border border-b-surface-raised bg-surface-raised text-text'
                  : 'border-transparent text-text-muted hover:text-text'
              }`}
            >
              {TAB_LANGS[t.lang] ?? t.lang}
            </button>
          );
        })}
      </div>
      <div role="tabpanel" id={`${id}-panel`} aria-labelledby={`${id}-tab-${current.lang}`}>
        <CodeBlock>
          <Highlighted code={current.code} lang={current.lang} />
        </CodeBlock>
      </div>
    </div>
  );
}

/**
 * Renders author-written markdown (lesson prompts), with GitHub-flavoured tables,
 * strikethrough and task lists. Raw HTML in the source is not
 * rendered (react-markdown's default), and nothing goes through innerHTML.
 * Images load only from the lesson's own assets; any other image (another site,
 * a missing file, or an image in tutor output) shows as its alt text.
 */
export function Markdown({ children, dropTitle = false }: { children: string; dropTitle?: boolean }) {
  const resolve = useLessonAsset();
  // Lesson pages already show the title, so a prompt's leading "# Title" line is dropped there.
  const source = dropTitle ? children.replace(/^\s*# [^\n]*(\n+|$)/, '') : children;
  return (
    <div className="text-body flex flex-col gap-3 text-text-secondary">
      <ReactMarkdown
        remarkPlugins={[remarkGfm, remarkCodeTabs]}
        components={{
          // Wide tables scroll inside their own box instead of widening the pane.
          table: ({ children }) => (
            <div className="overflow-x-auto rounded-md border border-border">
              <table className="text-small w-full border-collapse text-left">{children}</table>
            </div>
          ),
          thead: ({ children }) => <thead className="bg-surface-raised text-text">{children}</thead>,
          th: ({ children, style }) => (
            <th style={style} className="border-b border-border px-3 py-2 font-semibold">
              {children}
            </th>
          ),
          td: ({ children, style }) => (
            <td style={style} className="border-t border-border px-3 py-2 align-top">
              {children}
            </td>
          ),
          // The CSS reset strips list markers and indentation, so lists bring their own.
          // Items can hold code blocks and nested lists, so their children get spacing too.
          ul: ({ children, className }) => (
            <ul
              className={
                className?.includes('contains-task-list')
                  ? 'flex flex-col gap-1.5'
                  : 'flex list-disc flex-col gap-1.5 pl-5 marker:text-text-muted'
              }
            >
              {children}
            </ul>
          ),
          ol: ({ children, start }) => (
            <ol start={start} className="flex list-decimal flex-col gap-1.5 pl-6 marker:text-text-muted">
              {children}
            </ol>
          ),
          li: ({ children, className }) => (
            <li className={`pl-1 [&>*+*]:mt-2 ${className?.includes('task-list-item') ? 'list-none [&_input]:mr-2' : ''}`}>
              {children}
            </li>
          ),
          strong: ({ children }) => <strong className="font-semibold text-text">{children}</strong>,
          pre: ({ children, node }) => {
            const block = fenced(node);
            const tabs = tabsOf(block);
            if (tabs) return <CodeTabs tabs={tabs} />;
            if (block && HIGHLIGHT_LANGS.includes(block.lang)) {
              return (
                <CodeBlock>
                  <Highlighted code={block.text} lang={block.lang} />
                </CodeBlock>
              );
            }
            return <CodeBlock>{children}</CodeBlock>;
          },
          // Links stand out from body text; links to other sites open in a new tab.
          a: ({ children, href }) => {
            const external = typeof href === 'string' && /^https?:\/\//.test(href);
            return (
              <a
                href={href}
                {...(external ? { target: '_blank', rel: 'noopener noreferrer' } : {})}
                className="font-medium text-accent-text underline decoration-accent-line underline-offset-2 hover:decoration-accent-text"
              >
                {children}
              </a>
            );
          },
          // The CSS reset strips heading sizes, so headings bring their own.
          h1: ({ children }) => <h1 className="text-page-title text-text">{children}</h1>,
          h2: ({ children }) => <h2 className="text-[17px] font-semibold leading-6 text-text [&:not(:first-child)]:mt-3">{children}</h2>,
          h3: ({ children }) => <h3 className="text-section text-text [&:not(:first-child)]:mt-2">{children}</h3>,
          h4: ({ children }) => <h4 className="text-section text-text-secondary">{children}</h4>,
          code: ({ children, className }) =>
            className ? (
              <code>{children}</code>
            ) : (
              <code className="rounded-sm bg-surface-raised px-1.5 py-0.5 font-mono text-[13px] text-text">{children}</code>
            ),
          img: ({ src, alt }) => {
            const url = resolve(typeof src === 'string' ? src : undefined);
            return url ? <LessonImage src={url} alt={alt ?? ''} /> : <span className="italic text-text-muted">{alt ? `[${alt}]` : '[image]'}</span>;
          },
        }}
      >
        {source}
      </ReactMarkdown>
    </div>
  );
}
