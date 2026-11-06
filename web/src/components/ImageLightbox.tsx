import { useEffect, useRef } from 'react';

/** Full-size view of a lesson image. Closes with Escape, a click anywhere, or the button. */
export function ImageLightbox({ src, alt, onClose }: { src: string; alt: string; onClose: () => void }) {
  const close = useRef<HTMLButtonElement>(null);
  useEffect(() => {
    const before = document.activeElement as HTMLElement | null;
    close.current?.focus();
    const onKey = (e: KeyboardEvent) => e.key === 'Escape' && onClose();
    document.addEventListener('keydown', onKey);
    return () => {
      document.removeEventListener('keydown', onKey);
      before?.focus?.();
    };
  }, [onClose]);
  return (
    <div role="dialog" aria-modal="true" aria-label={alt || 'Image'} className="fixed inset-0 z-40 flex items-center justify-center p-4" onMouseDown={onClose}>
      <div aria-hidden="true" className="absolute inset-0 bg-bg opacity-90" />
      <img src={src} alt={alt} className="relative max-h-full max-w-full rounded-md border border-border bg-surface-raised p-2" />
      <button
        ref={close}
        type="button"
        onClick={onClose}
        className="text-control absolute right-4 top-4 rounded-md border border-border bg-surface px-3 py-1.5 text-text"
      >
        Close
      </button>
    </div>
  );
}
