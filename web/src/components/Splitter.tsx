import { type KeyboardEvent, type PointerEvent, useRef } from 'react';

interface Props {
  /** "vertical" sits between side-by-side panes (drag left/right); "horizontal" between stacked ones. */
  orientation: 'vertical' | 'horizontal';
  label: string;
  /** Size of the pane this handle controls, in px. */
  value: number;
  min: number;
  max: number;
  /** +1 when dragging right/down grows the pane, -1 when it shrinks it. */
  direction: 1 | -1;
  onChange: (px: number) => void;
  onReset: () => void;
}

const STEP = 16;

/**
 * Drag handle between two panes (WAI-ARIA window splitter): pointer drag,
 * arrow keys (Shift for bigger steps), Home/End, and Enter or double-click to
 * reset. Hidden below 900px, where the panes stack.
 */
export function Splitter({ orientation, label, value, min, max, direction, onChange, onReset }: Props) {
  const drag = useRef<{ start: number; size: number } | null>(null);
  const vertical = orientation === 'vertical';

  function onPointerDown(e: PointerEvent<HTMLDivElement>) {
    if (e.button !== 0) return;
    e.preventDefault();
    e.currentTarget.setPointerCapture?.(e.pointerId);
    drag.current = { start: vertical ? e.clientX : e.clientY, size: value };
    document.body.style.cursor = vertical ? 'col-resize' : 'row-resize';
    document.body.style.userSelect = 'none';
  }
  function onPointerMove(e: PointerEvent<HTMLDivElement>) {
    const d = drag.current;
    if (!d) return;
    const delta = (vertical ? e.clientX : e.clientY) - d.start;
    onChange(d.size + direction * delta);
  }
  function end(e: PointerEvent<HTMLDivElement>) {
    if (!drag.current) return;
    drag.current = null;
    if (e.currentTarget.hasPointerCapture?.(e.pointerId)) e.currentTarget.releasePointerCapture(e.pointerId);
    document.body.style.cursor = '';
    document.body.style.userSelect = '';
  }
  function onKeyDown(e: KeyboardEvent<HTMLDivElement>) {
    const step = e.shiftKey ? STEP * 4 : STEP;
    const grow = vertical ? { ArrowRight: 1, ArrowLeft: -1 } : { ArrowDown: 1, ArrowUp: -1 };
    const dir = grow[e.key as keyof typeof grow];
    if (dir !== undefined) onChange(value + dir * direction * step);
    else if (e.key === 'Home') onChange(min);
    else if (e.key === 'End') onChange(max);
    else if (e.key === 'Enter') onReset();
    else return;
    e.preventDefault();
  }

  return (
    <div
      role="separator"
      aria-orientation={orientation}
      aria-label={label}
      aria-valuenow={Math.round(value)}
      aria-valuemin={min}
      aria-valuemax={max}
      tabIndex={0}
      title={`Drag to resize · double-click to reset`}
      onPointerDown={onPointerDown}
      onPointerMove={onPointerMove}
      onPointerUp={end}
      onPointerCancel={end}
      onDoubleClick={onReset}
      onKeyDown={onKeyDown}
      className={`group flex shrink-0 touch-none items-center justify-center outline-none max-[900px]:hidden ${
        vertical ? 'w-2 cursor-col-resize self-stretch' : 'h-2 cursor-row-resize'
      }`}
    >
      <span
        aria-hidden="true"
        className={`rounded-pill bg-border transition-colors group-hover:bg-accent group-focus-visible:bg-accent group-active:bg-accent ${
          vertical ? 'h-10 w-1' : 'h-1 w-10'
        }`}
      />
    </div>
  );
}
