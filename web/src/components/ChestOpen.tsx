import { useEffect, useRef, useState } from 'react';
import { CountUp } from '../effects/CountUp';
import { play } from '../effects/sound';
import { useReducedMotion } from '../effects/useReducedMotion';
import { Button } from './Button';
import { ChestSprite } from './ChestSprite';
import { PowerUpIcon, powerupName } from './powerups';

export interface Loot {
  bytes: number;
  powerup: string;
}

interface Props {
  /** "Daily chest" or "Boss chest". */
  title: string;
  /** The contents; undefined while the server is still deciding. */
  loot: Loot | undefined;
  /** Shown instead of the loot when opening failed. */
  error?: string;
  onClose: () => void;
}

/** Shortest wobble before the lid opens, so a fast response still reads as "opening" (ms). */
const MIN_WOBBLE_MS = 600;

/**
 * Chest opening (daily chest and boss chest): the closed chest wobbles while
 * the loot is pending, the lid opens, then Bytes count up and the power-up
 * rises. The loot is decided and granted by the server; this only shows it.
 * Skippable with Escape, the backdrop or the button; focus returns on close.
 */
export function ChestOpen({ title, loot, error, onClose }: Props) {
  const reduced = useReducedMotion();
  const [waited, setWaited] = useState(reduced);
  const [counted, setCounted] = useState(false);
  const dialog = useRef<HTMLDivElement>(null);
  const opened = waited && loot !== undefined && !error;

  useEffect(() => {
    if (reduced) return;
    const t = setTimeout(() => setWaited(true), MIN_WOBBLE_MS);
    return () => clearTimeout(t);
  }, [reduced]);

  const played = useRef(false); // once, also under StrictMode's double effects
  useEffect(() => {
    if (!opened || played.current) return;
    played.current = true;
    play('chestOpen');
  }, [opened]);

  useEffect(() => {
    const before = document.activeElement as HTMLElement | null;
    dialog.current?.querySelector('button')?.focus();
    const onKey = (e: KeyboardEvent) => e.key === 'Escape' && onClose();
    document.addEventListener('keydown', onKey);
    return () => {
      document.removeEventListener('keydown', onKey);
      before?.focus?.();
    };
  }, [onClose]);

  return (
    <div className="fixed inset-0 z-40 flex items-center justify-center p-4">
      <div aria-hidden="true" data-testid="chest-backdrop" className="absolute inset-0 bg-bg opacity-80" onMouseDown={onClose} />
      <div
        ref={dialog}
        role="dialog"
        aria-modal="true"
        aria-labelledby="chest-title"
        className="relative flex w-full max-w-[360px] flex-col items-center gap-4 rounded-lg border border-border bg-surface p-6 text-center"
      >
        <p id="chest-title" className="text-game-sm tracking-[0.2em] text-accent-text">
          {title.toUpperCase()}
        </p>
        <ChestSprite open={opened} wobble={!opened && !error} />
        <div aria-live="polite" className="flex min-h-[88px] flex-col items-center gap-2">
          {error && (
            <p role="alert" className="text-small text-danger">
              {error}
            </p>
          )}
          {!error && !opened && <p className="text-small text-text-muted">Opening…</p>}
          {opened && (
            <>
              <p className="cq-pop text-game text-[28px] leading-[32px] text-accent-text">
                +<CountUp to={loot.bytes} duration={900} onDone={() => setCounted(true)} /> Bytes
              </p>
              {(counted || reduced) && (
                <p className="cq-rise flex items-center gap-2 rounded-md border border-border bg-surface-raised px-3 py-2">
                  <PowerUpIcon id={loot.powerup} className="text-accent-text" />
                  <span className="text-section text-text">{powerupName(loot.powerup)}</span>
                </p>
              )}
            </>
          )}
        </div>
        <Button className="h-10 w-full" onClick={onClose}>
          {opened ? 'Collect' : 'Skip'}
        </Button>
      </div>
    </div>
  );
}
