import { useState } from 'react';
import { createPortal } from 'react-dom';
import ReactMarkdown from 'react-markdown';
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
 * Renders author-written markdown (lesson prompts). Raw HTML in the source is not
 * rendered (react-markdown's default), and nothing goes through innerHTML.
 * Images load only from the lesson's own assets; any other image (another site,
 * a missing file, or an image in tutor output) shows as its alt text.
 */
export function Markdown({ children }: { children: string }) {
  const resolve = useLessonAsset();
  return (
    <div className="text-body flex flex-col gap-3 text-text-secondary">
      <ReactMarkdown
        components={{
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
