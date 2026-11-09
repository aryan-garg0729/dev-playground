import { useRef, useState, type ReactNode } from 'react';
import { createPortal } from 'react-dom';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import { ImageLightbox } from './ImageLightbox';
import { useLessonAsset } from './LessonAssets';

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
        remarkPlugins={[remarkGfm]}
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
          pre: ({ children }) => <CodeBlock>{children}</CodeBlock>,
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
