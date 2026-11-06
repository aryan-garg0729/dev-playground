// Synthesized chiptune sound effects (Web Audio API, no audio files).
// Only for rare, earned moments: lesson clear, level up, boss fight, chest.
// Routine submits, runs and fails never call play().
import { useSound } from '../stores/sound';

export const SOUND_NAMES = ['success', 'levelUp', 'bossHit', 'bossBlock', 'bossDefeat', 'chestOpen'] as const;
export type SoundName = (typeof SOUND_NAMES)[number];

type Wave = OscillatorType;
/** One note: start offset and length in seconds, frequency in Hz (optionally sliding to `to`). */
interface Note {
  at: number;
  len: number;
  freq: number;
  to?: number;
  wave?: Wave;
  gain?: number;
}

// Note frequencies (equal temperament).
const C5 = 523.25, E5 = 659.25, G5 = 783.99, C6 = 1046.5, E6 = 1318.5, G4 = 392.0, A4 = 440.0;

const RECIPES: Record<SoundName, Note[]> = {
  // Soft two-note chime.
  success: [
    { at: 0, len: 0.12, freq: E5, wave: 'triangle' },
    { at: 0.09, len: 0.28, freq: G5, wave: 'triangle' },
  ],
  // Rising arpeggio with a held top note.
  levelUp: [
    { at: 0, len: 0.1, freq: C5, wave: 'square', gain: 0.5 },
    { at: 0.09, len: 0.1, freq: E5, wave: 'square', gain: 0.5 },
    { at: 0.18, len: 0.1, freq: G5, wave: 'square', gain: 0.5 },
    { at: 0.27, len: 0.45, freq: C6, wave: 'square', gain: 0.5 },
    { at: 0.27, len: 0.45, freq: E6, wave: 'triangle', gain: 0.4 },
  ],
  // Short downward blip.
  bossHit: [{ at: 0, len: 0.09, freq: 880, to: 330, wave: 'square', gain: 0.35 }],
  // Dull low thud.
  bossBlock: [{ at: 0, len: 0.12, freq: 160, to: 110, wave: 'triangle', gain: 0.6 }],
  // Short victory fanfare.
  bossDefeat: [
    { at: 0, len: 0.12, freq: G4, wave: 'square', gain: 0.45 },
    { at: 0.12, len: 0.12, freq: C5, wave: 'square', gain: 0.45 },
    { at: 0.24, len: 0.12, freq: E5, wave: 'square', gain: 0.45 },
    { at: 0.36, len: 0.5, freq: G5, wave: 'square', gain: 0.45 },
    { at: 0.36, len: 0.5, freq: C5, wave: 'triangle', gain: 0.4 },
  ],
  // Creak up, then sparkle.
  chestOpen: [
    { at: 0, len: 0.16, freq: A4, to: 880, wave: 'triangle', gain: 0.5 },
    { at: 0.16, len: 0.08, freq: C6, wave: 'square', gain: 0.25 },
    { at: 0.24, len: 0.08, freq: E6, wave: 'square', gain: 0.25 },
    { at: 0.32, len: 0.3, freq: G5 * 2, wave: 'triangle', gain: 0.3 },
  ],
};

/** Overall loudness before the user's volume; effects stay quiet. */
const MASTER = 0.4;

let ctx: AudioContext | null | undefined;

function context(): AudioContext | null {
  if (ctx !== undefined) return ctx;
  const Ctor = typeof window !== 'undefined' ? window.AudioContext : undefined;
  try {
    ctx = Ctor ? new Ctor() : null;
  } catch {
    ctx = null;
  }
  return ctx;
}

/**
 * Plays a named effect unless sound is muted. Every call follows a click
 * (Submit, Attack, Open), so the browser allows audio. Any failure is silent.
 */
export function play(name: SoundName): void {
  const { muted, volume } = useSound.getState();
  if (muted || volume <= 0) return;
  const ac = context();
  if (!ac) return;
  try {
    if (ac.state === 'suspended') void ac.resume();
    const master = ac.createGain();
    master.gain.value = MASTER * volume;
    master.connect(ac.destination);
    const t0 = ac.currentTime + 0.01;
    for (const n of RECIPES[name]) {
      const osc = ac.createOscillator();
      const env = ac.createGain();
      osc.type = n.wave ?? 'square';
      const start = t0 + n.at;
      const end = start + n.len;
      osc.frequency.setValueAtTime(n.freq, start);
      if (n.to) osc.frequency.exponentialRampToValueAtTime(n.to, end);
      const peak = n.gain ?? 0.6;
      env.gain.setValueAtTime(0.0001, start);
      env.gain.exponentialRampToValueAtTime(peak, start + 0.01);
      env.gain.exponentialRampToValueAtTime(0.0001, end);
      osc.connect(env).connect(master);
      osc.start(start);
      osc.stop(end + 0.02);
    }
  } catch {
    // audio is a nicety; never let it break the page
  }
}

/** Test hook: forget the cached AudioContext. */
export function resetSoundForTests(): void {
  ctx = undefined;
}
