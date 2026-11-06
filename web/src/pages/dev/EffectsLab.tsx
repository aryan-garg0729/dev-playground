import { type ReactNode, useEffect, useRef, useState } from 'react';
import type { Result } from '../../api/client';
import { Button } from '../../components/Button';
import { ChestOpen, type Loot } from '../../components/ChestOpen';
import { PassCelebration } from '../../components/PassCelebration';
import { powerupName } from '../../components/powerups';
import { SoundToggle } from '../../components/SoundToggle';
import { ThemeToggle } from '../../components/ThemeToggle';
import { SOUND_NAMES, play } from '../../effects/sound';
import { useReducedMotion } from '../../effects/useReducedMotion';
import { attackDuration, BossHP, BossResults, useAttack } from '../../exercise-types/shared/Boss';
import { attack, type CelebrationKind, celebrations, POWERUPS } from './scenarios';

function Section({ title, hint, children }: { title: string; hint?: string; children: ReactNode }) {
  return (
    <section aria-label={title} className="flex flex-col gap-3 rounded-lg border border-border bg-surface p-5">
      <h2 className="text-section text-text">{title}</h2>
      {hint && <p className="text-small text-text-muted">{hint}</p>}
      {children}
    </section>
  );
}

const row = 'flex flex-wrap items-center gap-2';

/** Boss attack with any number of checks and hits; a full clear opens the celebration like a real lesson. */
function BossLab({ onDefeat }: { onDefeat: () => void }) {
  const [total, setTotal] = useState(5);
  const [hits, setHits] = useState(3);
  const [result, setResult] = useState<Result>();
  const revealed = useAttack(result, result !== undefined);
  const timer = useRef<ReturnType<typeof setTimeout>>(undefined);
  useEffect(() => () => clearTimeout(timer.current), []);
  const go = (h: number) => {
    const r = attack(total, Math.min(h, total));
    setResult(r);
    clearTimeout(timer.current);
    if (r.status === 'pass') timer.current = setTimeout(onDefeat, attackDuration(r));
  };
  const field = 'text-body w-16 rounded-md border border-border bg-bg px-2 py-1 text-text';
  return (
    <Section title="Boss fight" hint="Checks are revealed 150 ms apart with hit/block sounds; a full clear plays the fanfare, then the celebration.">
      <div className={row}>
        <label className="text-small flex items-center gap-2 text-text-secondary">
          Checks (HP)
          <input type="number" min={1} max={12} value={total} onChange={(e) => setTotal(Math.max(1, Math.min(12, +e.target.value)))} className={field} />
        </label>
        <label className="text-small flex items-center gap-2 text-text-secondary">
          Hits
          <input type="number" min={0} max={total} value={Math.min(hits, total)} onChange={(e) => setHits(Math.max(0, +e.target.value))} className={field} />
        </label>
        <Button variant="primary" onClick={() => go(hits)}>
          Attack
        </Button>
        <Button onClick={() => go(total)}>Killing blow</Button>
        <Button onClick={() => setResult(undefined)}>Reset</Button>
      </div>
      <div className="grid items-start gap-3 md:grid-cols-2">
        <BossHP maxHP={result?.checks.length ?? total} result={result} revealed={result ? revealed : undefined} />
        {result && <BossResults result={result} revealed={revealed} />}
      </div>
    </Section>
  );
}

/** Daily chest: opens like Home's, with a fake server delay (or a failure). */
function ChestLab() {
  const [open, setOpen] = useState(false);
  const [loot, setLoot] = useState<Loot>();
  const [error, setError] = useState<string>();
  const [powerup, setPowerup] = useState(POWERUPS[0]!);
  const start = (delay: number, fail = false) => {
    setLoot(undefined);
    setError(undefined);
    setOpen(true);
    setTimeout(() => (fail ? setError('Chest not ready yet (simulated error).') : setLoot({ bytes: 40 + Math.floor(Math.random() * 41), powerup })), delay);
  };
  return (
    <Section title="Chest opening" hint="Real server replies take 50–300 ms; the chest wobbles at least 600 ms either way.">
      <div className={row}>
        <select value={powerup} onChange={(e) => setPowerup(e.target.value)} className="text-body rounded-md border border-border bg-bg px-2 py-1.5 text-text" aria-label="Power-up in the chest">
          {POWERUPS.map((p) => (
            <option key={p} value={p}>
              {powerupName(p)}
            </option>
          ))}
        </select>
        <Button variant="primary" onClick={() => start(150)}>
          Open (fast server)
        </Button>
        <Button onClick={() => start(2000)}>Open (slow server, 2 s)</Button>
        <Button onClick={() => start(800, true)}>Open (error)</Button>
      </div>
      {open && <ChestOpen title="Daily chest" loot={loot} error={error} onClose={() => setOpen(false)} />}
    </Section>
  );
}

/**
 * Dev-only page (`/dev/effects`, never in production builds) for trying every
 * game moment: sounds, the lesson-clear dialog, level up, boss fight and chests.
 * Everything here is local make-believe; no request is made and nothing is saved.
 */
export function EffectsLab() {
  const [celebrate, setCelebrate] = useState<CelebrationKind>();
  const reduced = useReducedMotion();
  const allPass = attack(5, 5);
  return (
    <main className="mx-auto flex max-w-[960px] flex-col gap-4 p-6">
      <header className="flex flex-wrap items-center justify-between gap-3">
        <div>
          <p className="text-game-sm tracking-[0.2em] text-accent-text">DEV ONLY</p>
          <h1 className="text-page-title text-text">Effects lab</h1>
        </div>
        <div className="flex items-center gap-2">
          <span className="text-small text-text-muted">Reduced motion: {reduced ? 'on' : 'off'}</span>
          <SoundToggle />
          <ThemeToggle />
        </div>
      </header>
      <p className="text-small text-text-muted">
        To check reduced motion, open DevTools › Rendering › “Emulate CSS prefers-reduced-motion”. Muting here mutes the whole app.
      </p>

      <Section title="Sounds">
        <div className={row}>
          {SOUND_NAMES.map((n) => (
            <Button key={n} onClick={() => play(n)}>
              {n}
            </Button>
          ))}
        </div>
      </Section>

      <Section title="Lesson clear dialog" hint="One sound per dialog: the chime, the level-up fanfare, or none after a boss (its fanfare already played).">
        <div className={row}>
          {(Object.keys(celebrations) as CelebrationKind[]).map((k) => (
            <Button key={k} onClick={() => setCelebrate(k)}>
              {celebrations[k].label}
            </Button>
          ))}
        </div>
      </Section>

      <BossLab onDefeat={() => setCelebrate('boss')} />
      <ChestLab />

      {celebrate && (
        <PassCelebration
          outcome={celebrations[celebrate].outcome}
          result={allPass}
          next={{ to: '/dev/effects', title: 'Back to the lab' }}
          courseTo="/dev/effects"
          onClose={() => setCelebrate(undefined)}
        />
      )}
    </main>
  );
}
