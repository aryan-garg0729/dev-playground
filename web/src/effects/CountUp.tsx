import { useEffect, useRef, useState } from 'react';
import { useReducedMotion } from './useReducedMotion';

interface Props {
  to: number;
  from?: number;
  /** Milliseconds; matches the `motion-celebrate` token by default. */
  duration?: number;
  /** Called once the final value shows. */
  onDone?: () => void;
}

/** A number that counts up to `to` (XP, Bytes). Shows `to` at once under reduced motion. */
export function CountUp({ to, from = 0, duration = 1200, onDone }: Props) {
  const reduced = useReducedMotion();
  const [value, setValue] = useState(reduced ? to : from);
  const done = useRef(onDone);
  done.current = onDone;
  useEffect(() => {
    if (reduced || to === from) {
      setValue(to);
      done.current?.();
      return;
    }
    let raf = 0;
    const start = performance.now();
    const tick = (now: number) => {
      const p = Math.min(1, (now - start) / duration);
      const eased = 1 - Math.pow(1 - p, 3);
      setValue(Math.round(from + (to - from) * eased));
      if (p < 1) raf = requestAnimationFrame(tick);
      else done.current?.();
    };
    raf = requestAnimationFrame(tick);
    return () => cancelAnimationFrame(raf);
  }, [to, from, duration, reduced]);
  return <>{value}</>;
}
