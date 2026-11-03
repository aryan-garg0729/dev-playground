import { create } from 'zustand';

const STORAGE_KEY = 'cq-panes';

/** Lesson pane sizes in px. */
export const PANE_DEFAULTS = { problemWidth: 320, resultsHeight: 240 } as const;
export const PANE_LIMITS = {
  problemWidth: { min: 240, max: 720 },
  resultsHeight: { min: 96, max: 720 },
} as const;

type Pane = keyof typeof PANE_DEFAULTS;
type Sizes = Record<Pane, number>;

const clamp = (pane: Pane, v: number) => Math.round(Math.min(PANE_LIMITS[pane].max, Math.max(PANE_LIMITS[pane].min, v)));

function initial(): Sizes {
  try {
    const saved = JSON.parse(localStorage.getItem(STORAGE_KEY) ?? '{}') as Partial<Sizes>;
    const pick = (p: Pane) => (typeof saved[p] === 'number' && Number.isFinite(saved[p]) ? clamp(p, saved[p]) : PANE_DEFAULTS[p]);
    return { problemWidth: pick('problemWidth'), resultsHeight: pick('resultsHeight') };
  } catch {
    return { ...PANE_DEFAULTS };
  }
}

interface PanesState extends Sizes {
  setSize: (pane: Pane, px: number) => void;
  reset: (pane: Pane) => void;
}

function save(s: Sizes) {
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(s));
  } catch {
    // storage unavailable: sizes last for this page only
  }
}

/** Resizable lesson panes (problem width, checks height); remembered per device. */
export const usePanes = create<PanesState>((set, get) => ({
  ...initial(),
  setSize: (pane, px) => {
    set({ [pane]: clamp(pane, px) } as Pick<Sizes, Pane>);
    const { problemWidth, resultsHeight } = get();
    save({ problemWidth, resultsHeight });
  },
  reset: (pane) => get().setSize(pane, PANE_DEFAULTS[pane]),
}));
