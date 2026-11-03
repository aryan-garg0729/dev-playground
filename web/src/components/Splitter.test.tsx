import { fireEvent, render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { useState } from 'react';
import { vi } from 'vitest';
import { PANE_DEFAULTS, usePanes } from '../stores/panes';
import { Splitter } from './Splitter';

function Harness({ direction = 1 as 1 | -1, orientation = 'vertical' as 'vertical' | 'horizontal', onReset = vi.fn() }) {
  const [v, setV] = useState(300);
  return (
    <Splitter
      orientation={orientation}
      label="Resize"
      value={v}
      min={200}
      max={500}
      direction={direction}
      onChange={(px) => setV(Math.min(500, Math.max(200, px)))}
      onReset={onReset}
    />
  );
}

const sep = () => screen.getByRole('separator', { name: 'Resize' });

describe('Splitter', () => {
  it('exposes its size to assistive tech', () => {
    render(<Harness />);
    expect(sep()).toHaveAttribute('aria-orientation', 'vertical');
    expect(sep()).toHaveAttribute('aria-valuenow', '300');
    expect(sep()).toHaveAttribute('aria-valuemin', '200');
    expect(sep()).toHaveAttribute('aria-valuemax', '500');
  });

  it('resizes with the keyboard', async () => {
    const onReset = vi.fn();
    render(<Harness onReset={onReset} />);
    sep().focus();
    await userEvent.keyboard('{ArrowRight}');
    expect(sep()).toHaveAttribute('aria-valuenow', '316');
    await userEvent.keyboard('{Shift>}{ArrowLeft}{/Shift}');
    expect(sep()).toHaveAttribute('aria-valuenow', '252');
    await userEvent.keyboard('{End}');
    expect(sep()).toHaveAttribute('aria-valuenow', '500');
    await userEvent.keyboard('{Home}');
    expect(sep()).toHaveAttribute('aria-valuenow', '200');
    await userEvent.keyboard('{Enter}');
    expect(onReset).toHaveBeenCalledOnce();
  });

  it('follows a pointer drag; a pane below its handle grows when dragged up', () => {
    render(<Harness orientation="horizontal" direction={-1} />);
    const el = sep();
    el.setPointerCapture = vi.fn();
    el.hasPointerCapture = () => false;
    fireEvent.pointerDown(el, { button: 0, clientY: 600, pointerId: 1 });
    fireEvent.pointerMove(el, { clientY: 540, pointerId: 1 });
    expect(el).toHaveAttribute('aria-valuenow', '360');
    fireEvent.pointerUp(el, { pointerId: 1 });
    fireEvent.pointerMove(el, { clientY: 400, pointerId: 1 }); // released: no change
    expect(el).toHaveAttribute('aria-valuenow', '360');
    expect(document.body.style.userSelect).toBe('');
  });

  it('resets on double-click', async () => {
    const onReset = vi.fn();
    render(<Harness onReset={onReset} />);
    await userEvent.dblClick(sep());
    expect(onReset).toHaveBeenCalled();
  });
});

describe('usePanes', () => {
  it('clamps, remembers and resets sizes', () => {
    const { setSize, reset } = usePanes.getState();
    setSize('problemWidth', 10_000);
    expect(usePanes.getState().problemWidth).toBe(720);
    setSize('resultsHeight', 5);
    expect(usePanes.getState().resultsHeight).toBe(96);
    expect(JSON.parse(localStorage.getItem('cq-panes')!)).toEqual({ problemWidth: 720, resultsHeight: 96 });
    reset('problemWidth');
    reset('resultsHeight');
    expect(usePanes.getState()).toMatchObject(PANE_DEFAULTS);
  });
});
