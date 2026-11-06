/**
 * The chest drawn in token colors: body and lid are separate groups so the lid
 * can swing open (effects.css: .cq-lid / .cq-lid-open). Decorative.
 */
export function ChestSprite({ open, wobble }: { open: boolean; wobble: boolean }) {
  return (
    <svg viewBox="0 0 120 110" width="132" height="121" aria-hidden="true" className={wobble ? 'cq-chest-wobble' : ''}>
      {open && <circle className="cq-glow" cx="60" cy="52" r="44" fill="var(--accent-soft)" />}
      {open && (
        <g className="cq-glow" stroke="var(--accent)" strokeWidth="3" strokeLinecap="round">
          <path d="M60 14v10M30 26l6 7M90 26l-6 7M18 50h9M93 50h9" />
        </g>
      )}
      {/* body */}
      <g>
        <rect x="16" y="56" width="88" height="44" rx="5" fill="var(--accent-line)" stroke="var(--accent-text)" strokeWidth="3" />
        <path d="M16 70h88" stroke="var(--accent-text)" strokeWidth="3" />
        <rect x="52" y="64" width="16" height="16" rx="2" fill="var(--accent)" stroke="var(--accent-text)" strokeWidth="3" />
      </g>
      {/* lid, hinged at its back-left corner */}
      <g className={`cq-lid ${open ? 'cq-lid-open' : ''}`}>
        <path d="M16 56V44a14 14 0 0 1 14-14h60a14 14 0 0 1 14 14v12Z" fill="var(--accent)" stroke="var(--accent-text)" strokeWidth="3" strokeLinejoin="round" />
        <path d="M40 30v26M80 30v26" stroke="var(--accent-text)" strokeWidth="3" />
      </g>
    </svg>
  );
}
