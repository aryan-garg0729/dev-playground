import { act, render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { vi } from 'vitest';
import { setReducedMotion } from '../test/motion';
import { ChestOpen } from './ChestOpen';

const play = vi.hoisted(() => vi.fn());
vi.mock('../effects/sound', () => ({ play, resetSoundForTests: () => {} }));

describe('ChestOpen', () => {
  it('wobbles while pending, then shows the loot with one sound', () => {
    play.mockClear();
    const { rerender } = render(<ChestOpen title="Daily chest" loot={undefined} onClose={vi.fn()} />);
    expect(screen.getByRole('dialog', { name: 'DAILY CHEST' })).toHaveTextContent('Opening…');
    expect(play).not.toHaveBeenCalled();
    rerender(<ChestOpen title="Daily chest" loot={{ bytes: 64, powerup: 'peek' }} onClose={vi.fn()} />);
    expect(screen.getByRole('dialog')).toHaveTextContent('+64 Bytes');
    expect(screen.getByRole('dialog')).toHaveTextContent('Peek');
    expect(play).toHaveBeenCalledTimes(1);
    expect(play).toHaveBeenCalledWith('chestOpen');
  });

  it('keeps wobbling for a moment with motion allowed, even if the loot is known', () => {
    vi.useFakeTimers();
    setReducedMotion(false);
    render(<ChestOpen title="Boss chest" loot={{ bytes: 300, powerup: 'focus_boost' }} onClose={vi.fn()} />);
    expect(screen.getByRole('dialog')).toHaveTextContent('Opening…');
    act(() => vi.advanceTimersByTime(700));
    expect(screen.getByRole('button', { name: 'Collect' })).toBeInTheDocument();
    vi.useRealTimers();
  });

  it('says so when the chest held no power-up', () => {
    render(<ChestOpen title="Daily chest" loot={{ bytes: 52, powerup: '' }} onClose={vi.fn()} />);
    expect(screen.getByRole('dialog')).toHaveTextContent('+52 Bytes');
    expect(screen.getByRole('dialog')).toHaveTextContent('No power-up this time');
  });

  it('shows an error instead of loot', () => {
    render(<ChestOpen title="Daily chest" loot={undefined} error="Chest not ready" onClose={vi.fn()} />);
    expect(screen.getByRole('alert')).toHaveTextContent('Chest not ready');
  });

  it('is skippable with Escape, the backdrop and the button, and returns focus', async () => {
    const onClose = vi.fn();
    const opener = document.createElement('button');
    document.body.appendChild(opener);
    opener.focus();
    const { unmount } = render(<ChestOpen title="Daily chest" loot={undefined} onClose={onClose} />);
    expect(screen.getByRole('button', { name: 'Skip' })).toHaveFocus();
    await userEvent.keyboard('{Escape}');
    await userEvent.click(screen.getByTestId('chest-backdrop'));
    await userEvent.click(screen.getByRole('button', { name: 'Skip' }));
    expect(onClose).toHaveBeenCalledTimes(3);
    unmount();
    expect(opener).toHaveFocus();
    opener.remove();
  });
});
