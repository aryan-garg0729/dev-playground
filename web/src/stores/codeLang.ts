import { create } from 'zustand';

const STORAGE_KEY = 'cq-code-lang';
/** Language shown first in tabbed code blocks until the learner picks another. */
export const DEFAULT_CODE_LANG = 'cpp';

function initial(): string {
  try {
    return localStorage.getItem(STORAGE_KEY) ?? DEFAULT_CODE_LANG;
  } catch {
    return DEFAULT_CODE_LANG;
  }
}

interface CodeLangState {
  lang: string;
  setLang: (lang: string) => void;
}

/** The language picked in tabbed code blocks; shared by every block and remembered per device. */
export const useCodeLang = create<CodeLangState>((set) => ({
  lang: initial(),
  setLang: (lang) => {
    try {
      localStorage.setItem(STORAGE_KEY, lang);
    } catch {
      // storage unavailable: remember for this page only
    }
    set({ lang });
  },
}));
