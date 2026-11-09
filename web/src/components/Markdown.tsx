import { useState } from 'react';
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

/**
 * Renders author-written markdown (lesson prompts), with GitHub-flavoured tables,
 * strikethrough and task lists. Raw HTML in the source is not
 * rendered (react-markdown's default), and nothing goes through innerHTML.
 * Images load only from the lesson's own assets; any other image (another site,
 * a missing file, or an image in tutor output) shows as its alt text.
 */
export function Markdown({ children }: { children: string }) {
  const resolve = useLessonAsset();
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
          pre: ({ children }) => (
            <pre className="text-code-sm overflow-x-auto border-l-2 border-border-strong pl-[14px] text-text-secondary [&_code]:bg-transparent [&_code]:p-0 [&_code]:text-text-secondary">{children}</pre>
          ),
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
        {children}
      </ReactMarkdown>
    </div>
  );
}
