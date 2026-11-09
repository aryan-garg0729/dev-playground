/**
 * A ring spinner for waits (page data, the editor chunk). The label is for
 * screen readers; the ring is decorative. Under prefers-reduced-motion the
 * ring stays still.
 */
export function Spinner({ label, size = 'md' }: { label: string; size?: 'sm' | 'md' }) {
  const box = size === 'sm' ? 'h-4 w-4 border-2' : 'h-8 w-8 border-[3px]';
  return (
    <span role="status" className="inline-flex items-center">
      <span aria-hidden="true" className={`cq-spin inline-block rounded-full border-border-strong border-t-accent ${box}`} />
      <span className="sr-only">{label}</span>
    </span>
  );
}
