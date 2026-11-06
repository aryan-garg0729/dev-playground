import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { SoundToggle } from './SoundToggle';

describe('SoundToggle', () => {
  it('starts unmuted and toggles mute', async () => {
    render(<SoundToggle />);
    const btn = screen.getByRole('button', { name: 'Mute sounds' });
    expect(btn).toHaveAttribute('aria-pressed', 'false');
    await userEvent.click(btn);
    expect(screen.getByRole('button', { name: 'Unmute sounds' })).toHaveAttribute('aria-pressed', 'true');
  });
});
