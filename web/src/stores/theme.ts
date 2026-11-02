import { create } from 'zustand';

export type Theme = 'dark' | 'light';
const STORAGE_KEY = 'cq-theme';

function current(): Theme {
  return document.documentElement.getAttribute('data-theme') === 'light' ? 'light' : 'dark';
}

interface ThemeState {
  theme: Theme;
  setTheme: (t: Theme) => void;
  toggle: () => void;
}

// index.html sets data-theme before first paint (saved choice, else prefers-color-scheme);
// this store reads it and keeps <html> and localStorage in sync.
export const useTheme = create<ThemeState>((set, get) => ({
  theme: current(),
  setTheme: (theme) => {
    document.documentElement.setAttribute('data-theme', theme);
    try {
      localStorage.setItem(STORAGE_KEY, theme);
    } catch {
      // storage unavailable (private mode): the choice lasts for this page only
    }
    set({ theme });
  },
  toggle: () => get().setTheme(get().theme === 'dark' ? 'light' : 'dark'),
}));
