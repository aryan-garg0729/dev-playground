import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter } from 'react-router';
import { vi } from 'vitest';
import type { Result } from '../api/client';
import { firstPass } from '../test/fixtures';
import { PassCelebration } from './PassCelebration';

const sounds = vi.hoisted(() => vi.fn());
vi.mock('../effects/sound', () => ({ play: sounds, resetSoundForTests: () => {} }));

const result: Result = {
  status: 'pass', score: 1, stdout: '', runtime_ms: 41,
  checks: [1, 2, 3].map((i) => ({ name: `c${i}`, passed: true, visible: true })),
};

const findTheMax = { to: '/learn/c/max', title: 'Find the max' };

function show(onClose = vi.fn(), next: { to: string; title: string } | null = findTheMax) {
  render(
    <MemoryRouter>
      <PassCelebration outcome={firstPass} result={result} next={next ?? undefined} courseTo="/courses/c" onClose={onClose} />
    </MemoryRouter>,
  );
  return onClose;
}

describe('PassCelebration', () => {
  it('shows the XP breakdown, level progress and streak', () => {
    show();
    const dialog = screen.getByRole('dialog', { name: 'All checks passed' });
    expect(dialog).toHaveTextContent('3 of 3 checks · 41 ms · first try');
    expect(dialog).toHaveTextContent('Lesson complete+20');
    expect(dialog).toHaveTextContent('First-try bonus+10');
    expect(dialog).toHaveTextContent('No hints used+5');
    expect(dialog).toHaveTextContent('+35 XP');
    // 1727 - 1470 = 257 of 382 in level 7; 125 to go.
    expect(screen.getByRole('img', { name: '257 of 382 XP in level 7' })).toBeInTheDocument();
    expect(dialog).toHaveTextContent('125 XP to Lv 8');
    expect(dialog).toHaveTextContent('13-day streak');
  });

  it('focuses Next, which links to the next lesson', () => {
    show();
    const next = screen.getByRole('link', { name: 'Next: Find the max' });
    expect(next).toHaveFocus();
    expect(next).toHaveAttribute('href', '/learn/c/max');
  });

  it('is skippable with Escape, the backdrop and Keep working', async () => {
    const onClose = show();
    await userEvent.keyboard('{Escape}');
    await userEvent.click(screen.getByTestId('backdrop'));
    await userEvent.click(screen.getByRole('button', { name: 'Keep working' }));
    expect(onClose).toHaveBeenCalledTimes(3);
  });

  it('goes back to the course after the last lesson', () => {
    show(vi.fn(), null);
    expect(screen.getByRole('link', { name: 'Back to the course' })).toHaveAttribute('href', '/courses/c');
  });

  it('animates the XP bar only with motion allowed (motion-safe classes)', () => {
    show();
    const gained = screen.getByRole('img', { name: /XP in level/ }).lastElementChild!;
    expect(gained.className).toMatch(/motion-safe:transition-\[width\]/);
    expect(gained.className).not.toMatch(/(^|\s)transition/);
  });
});

describe('PassCelebration effects', () => {
  beforeEach(() => sounds.mockClear());

  const render1 = (outcome: typeof firstPass) =>
    render(
      <MemoryRouter>
        <PassCelebration outcome={outcome} result={result} courseTo="/courses/c" onClose={vi.fn()} />
      </MemoryRouter>,
    );

  it('plays one soft chime on a plain lesson clear, with no level-up burst', () => {
    render1(firstPass);
    expect(sounds.mock.calls).toEqual([['success']]);
    expect(screen.queryByText('LEVEL UP')).not.toBeInTheDocument();
  });

  it('plays the fanfare instead of the chime on a level up and shows old → new level', () => {
    render1({ ...firstPass, level_up: true, level: 8 });
    expect(sounds.mock.calls).toEqual([['levelUp']]);
    expect(screen.getByText('LEVEL UP')).toBeInTheDocument();
    expect(screen.getByLabelText('Level 7 to level 8')).toBeInTheDocument();
  });

  it('has no confetti under reduced motion', () => {
    render1({ ...firstPass, level_up: true, level: 8 });
    expect(screen.queryByTestId('confetti')).not.toBeInTheDocument();
  });

  it('after a boss: no extra sound, and the boss chest opens in place; Escape closes only the chest', async () => {
    const onClose = vi.fn();
    const boss = {
      ...firstPass,
      rewards: {
        bytes: [], medals: [], quests_done: [], chest_ready: false, chapter_complete: true,
        boss: { max_hp: 3, hp: 0, defeated: true, chest: { bytes: 300, powerup: 'peek' }, unlocked_chapter: 'ch2' },
      },
    };
    render(
      <MemoryRouter>
        <PassCelebration outcome={boss} result={result} courseTo="/courses/c" onClose={onClose} />
      </MemoryRouter>,
    );
    expect(sounds).not.toHaveBeenCalled();
    expect(screen.getByText('BOSS DEFEATED')).toBeInTheDocument();
    await userEvent.click(screen.getByRole('button', { name: 'Open boss chest' }));
    expect(screen.getByRole('dialog', { name: 'BOSS CHEST' })).toHaveTextContent('+300 Bytes');
    expect(sounds).toHaveBeenCalledWith('chestOpen');
    await userEvent.keyboard('{Escape}');
    expect(screen.queryByRole('dialog', { name: 'BOSS CHEST' })).not.toBeInTheDocument();
    expect(onClose).not.toHaveBeenCalled();
  });
});
