import '@testing-library/jest-dom/vitest';
import { cleanup } from '@testing-library/react';
import { afterEach, beforeEach, vi } from 'vitest';
import { resetSoundForTests } from '../effects/sound';
import { useSound } from '../stores/sound';
import { setReducedMotion } from './motion';

// Tests see final states at once (no count-ups or staged reveals) unless a test
// turns motion on with setReducedMotion(false).
beforeEach(() => setReducedMotion(true));

afterEach(() => {
  cleanup();
  vi.restoreAllMocks();
  vi.unstubAllGlobals();
  localStorage.clear();
  document.documentElement.removeAttribute('data-theme');
  useSound.setState({ muted: false, volume: 0.6 });
  resetSoundForTests();
});
