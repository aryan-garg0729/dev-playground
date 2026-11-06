import { vi } from 'vitest';
import { useSound } from '../stores/sound';
import { play, resetSoundForTests } from './sound';

/** A fake AudioContext that counts oscillators started. */
function fakeAudio() {
  const started = vi.fn();
  const param = () => ({ value: 0, setValueAtTime: vi.fn(), exponentialRampToValueAtTime: vi.fn() });
  const node = () => ({ connect: vi.fn((n: unknown) => n), gain: param() });
  class Ctx {
    state = 'running';
    currentTime = 0;
    destination = {};
    resume = vi.fn();
    createGain = node;
    createOscillator = () => ({ ...node(), type: 'square', frequency: param(), start: started, stop: vi.fn() });
  }
  vi.stubGlobal('AudioContext', Ctx);
  resetSoundForTests();
  return started;
}

describe('sound', () => {
  it('plays notes when unmuted', () => {
    const started = fakeAudio();
    play('levelUp');
    expect(started).toHaveBeenCalled();
  });

  it('is silent when muted', () => {
    const started = fakeAudio();
    useSound.getState().toggleMuted();
    play('chestOpen');
    expect(started).not.toHaveBeenCalled();
  });

  it('does nothing (and does not throw) without Web Audio', () => {
    vi.stubGlobal('AudioContext', undefined);
    resetSoundForTests();
    expect(() => play('success')).not.toThrow();
  });

  it('remembers mute per device', () => {
    useSound.getState().toggleMuted();
    expect(JSON.parse(localStorage.getItem('cq-sound')!)).toEqual({ muted: true, volume: 0.6 });
  });
});
