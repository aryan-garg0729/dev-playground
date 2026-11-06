import { useEffect, useRef } from 'react';
import { useReducedMotion } from './useReducedMotion';

/** Token names the particles use; read from CSS at runtime so both themes work and no hex lives here. */
const COLOR_VARS = ['--accent', '--success', '--medal-gold', '--level-bg', '--accent-line'];

interface Particle {
  x: number;
  y: number;
  vx: number;
  vy: number;
  rot: number;
  vr: number;
  size: number;
  color: string;
}

/**
 * A single short confetti burst over its parent (level up only). Decorative:
 * aria-hidden, no pointer events, and not drawn at all under reduced motion.
 */
export function Confetti({ count = 80, duration = 1200 }: { count?: number; duration?: number }) {
  const ref = useRef<HTMLCanvasElement>(null);
  const reduced = useReducedMotion();
  useEffect(() => {
    const canvas = ref.current;
    const g = canvas?.getContext?.('2d');
    if (reduced || !canvas || !g) return;
    const w = (canvas.width = canvas.offsetWidth || 420);
    const h = (canvas.height = canvas.offsetHeight || 420);
    const css = getComputedStyle(document.documentElement);
    const colors = COLOR_VARS.map((v) => css.getPropertyValue(v).trim()).filter(Boolean);
    const parts: Particle[] = Array.from({ length: count }, () => ({
      x: w / 2,
      y: Math.min(h * 0.35, 160), // near the banner, even in a tall dialog
      vx: (Math.random() - 0.5) * 14,
      vy: -Math.random() * 9 - 4,
      rot: Math.random() * Math.PI,
      vr: (Math.random() - 0.5) * 0.4,
      size: 4 + Math.random() * 4,
      color: colors[Math.floor(Math.random() * colors.length)] ?? 'currentColor',
    }));
    const start = performance.now();
    let raf = 0;
    const frame = (now: number) => {
      const t = (now - start) / duration;
      g.clearRect(0, 0, w, h);
      if (t >= 1) return;
      g.globalAlpha = 1 - t * t;
      for (const p of parts) {
        p.vy += 0.3;
        p.x += p.vx;
        p.y += p.vy;
        p.rot += p.vr;
        g.save();
        g.translate(p.x, p.y);
        g.rotate(p.rot);
        g.fillStyle = p.color;
        g.fillRect(-p.size / 2, -p.size / 4, p.size, p.size / 2);
        g.restore();
      }
      raf = requestAnimationFrame(frame);
    };
    raf = requestAnimationFrame(frame);
    return () => cancelAnimationFrame(raf);
  }, [count, duration, reduced]);
  if (reduced) return null;
  return <canvas ref={ref} aria-hidden="true" data-testid="confetti" className="pointer-events-none absolute inset-0 h-full w-full" />;
}
