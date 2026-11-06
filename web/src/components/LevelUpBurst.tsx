import { Confetti } from '../effects/Confetti';

/**
 * Level-up moment inside the lesson-clear dialog: a "LEVEL UP" banner, the
 * level badge popping from the old level to the new one, and one confetti
 * burst over the dialog (none under reduced motion). The fanfare is played by
 * the dialog, so a level up never stacks two sounds.
 */
export function LevelUpBurst({ level }: { level: number }) {
  return (
    <>
      <Confetti />
      <div role="status" className="flex flex-col items-center gap-2">
        <p className="cq-pop text-game text-[22px] leading-[26px] tracking-[0.15em] text-accent-text">LEVEL UP</p>
        <p className="flex items-center gap-2" aria-label={`Level ${level - 1} to level ${level}`}>
          <span className="text-game-sm rounded-sm bg-surface-raised px-[7px] py-[2px] text-text-muted">Lv {level - 1}</span>
          <span aria-hidden="true" className="text-text-muted">
            →
          </span>
          <span className="cq-pop text-game rounded-sm bg-level-bg px-[9px] py-[3px] text-level-ink [animation-delay:200ms]">Lv {level}</span>
        </p>
      </div>
    </>
  );
}
