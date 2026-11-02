import { create } from 'zustand';

const STORAGE_KEY = 'cq-rail-open';

function initial(): boolean {
  try {
    return localStorage.getItem(STORAGE_KEY) === '1';
  } catch {
    return false;
  }
}

interface RailState {
  open: boolean;
  setOpen: (open: boolean) => void;
}

/** Whether the lesson outline rail is open; remembered per device (closed by default). */
export const useRail = create<RailState>((set) => ({
  open: initial(),
  setOpen: (open) => {
    try {
      localStorage.setItem(STORAGE_KEY, open ? '1' : '0');
    } catch {
      // storage unavailable: remember for this page only
    }
    set({ open });
  },
}));
