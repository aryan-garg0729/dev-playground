// Guard for the "rare moments only" rule: game sounds may be triggered only from
// the lesson-clear dialog, the boss fight and the chest. Routine submits, runs
// and fails stay silent, so no other module may import the sound engine.
import { readdirSync, readFileSync, statSync } from 'node:fs';
import { join, relative } from 'node:path';

const SRC = join(__dirname, '..');
const ALLOWED = [
  'components/ChestOpen.tsx',
  'components/PassCelebration.tsx',
  'exercise-types/shared/Boss.tsx',
  'pages/dev/EffectsLab.tsx', // dev-only effects lab (/dev/effects), not in production builds
  'test/setup.ts',
];

function files(dir: string): string[] {
  return readdirSync(dir).flatMap((f) => {
    const p = join(dir, f);
    return statSync(p).isDirectory() ? files(p) : /\.tsx?$/.test(f) && !/\.test\.tsx?$/.test(f) ? [p] : [];
  });
}

it('only game moments import the sound engine', () => {
  const users = files(SRC)
    .filter((p) => /from '[./]+\/?(effects\/)?sound'/.test(readFileSync(p, 'utf8')) && !p.endsWith(join('effects', 'sound.ts')))
    .map((p) => relative(SRC, p).split('\\').join('/'))
    .filter((p) => !p.startsWith('stores/') && p !== 'components/SoundToggle.tsx');
  expect(users.sort()).toEqual(ALLOWED.sort());
});
