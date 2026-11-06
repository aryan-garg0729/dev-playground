import { act, render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter } from 'react-router';
import { vi } from 'vitest';
import { EffectsLab } from './EffectsLab';

const play = vi.hoisted(() => vi.fn());
vi.mock('../../effects/sound', () => ({
  play,
  resetSoundForTests: () => {},
  SOUND_NAMES: ['success', 'levelUp', 'bossHit', 'bossBlock', 'bossDefeat', 'chestOpen'],
}));

const show = () =>
  render(
    <MemoryRouter>
      <EffectsLab />
    </MemoryRouter>,
  );

describe('EffectsLab (dev only)', () => {
  beforeEach(() => play.mockClear());

  it('plays every sound on demand', async () => {
    show();
    await userEvent.click(screen.getByRole('button', { name: 'bossDefeat' }));
    expect(play).toHaveBeenCalledWith('bossDefeat');
  });

  it('opens each lesson-clear scenario', async () => {
    show();
    await userEvent.click(screen.getByRole('button', { name: 'Boss defeated + level up + chest' }));
    const dialog = screen.getByRole('dialog', { name: 'All checks passed' });
    expect(dialog).toHaveTextContent('BOSS DEFEATED');
    expect(dialog).toHaveTextContent('LEVEL UP');
    await userEvent.click(screen.getByRole('button', { name: 'Open boss chest' }));
    expect(screen.getByRole('dialog', { name: 'BOSS CHEST' })).toHaveTextContent('+300 Bytes');
  });

  it('runs a boss attack with the chosen hits, and a killing blow ends in the celebration', async () => {
    show();
    await userEvent.click(screen.getByRole('button', { name: 'Attack' }));
    expect(screen.getByRole('img', { name: '2 of 5 HP left' })).toBeInTheDocument();
    await userEvent.click(screen.getByRole('button', { name: 'Killing blow' }));
    expect(screen.getByRole('img', { name: '0 of 5 HP left' })).toBeInTheDocument();
    expect(await screen.findByRole('dialog', { name: 'All checks passed' })).toBeInTheDocument();
  });

  it('opens the daily chest, and shows a simulated error', async () => {
    vi.useFakeTimers({ shouldAdvanceTime: true });
    show();
    await userEvent.click(screen.getByRole('button', { name: 'Open (error)' }));
    expect(screen.getByRole('dialog', { name: 'DAILY CHEST' })).toHaveTextContent('Opening…');
    act(() => vi.advanceTimersByTime(900));
    expect(screen.getByRole('alert')).toHaveTextContent('simulated error');
    vi.useRealTimers();
  });
});
