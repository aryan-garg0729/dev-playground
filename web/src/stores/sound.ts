import { create } from 'zustand';

const STORAGE_KEY = 'cq-sound';

interface Saved {
  muted: boolean;
  volume: number;
}

const DEFAULTS: Saved = { muted: false, volume: 0.6 };

function load(): Saved {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (!raw) return DEFAULTS;
    const v = JSON.parse(raw) as Partial<Saved>;
    return {
      muted: typeof v.muted === 'boolean' ? v.muted : DEFAULTS.muted,
      volume: typeof v.volume === 'number' && v.volume >= 0 && v.volume <= 1 ? v.volume : DEFAULTS.volume,
    };
  } catch {
    return DEFAULTS;
  }
}

function save(s: Saved): void {
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(s));
  } catch {
    // storage unavailable (private mode): the choice lasts for this page only
  }
}

interface SoundState extends Saved {
  toggleMuted: () => void;
  setVolume: (v: number) => void;
}

/** Sound settings per device (`cq-sound`). On by default; the header toggle mutes. */
export const useSound = create<SoundState>((set, get) => ({
  ...load(),
  toggleMuted: () => {
    const next = { muted: !get().muted, volume: get().volume };
    save(next);
    set(next);
  },
  setVolume: (volume) => {
    const next = { muted: get().muted, volume: Math.min(1, Math.max(0, volume)) };
    save(next);
    set(next);
  },
}));
