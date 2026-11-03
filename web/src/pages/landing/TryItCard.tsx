import { useState } from 'react';

/**
 * The hero's "Try it" card: a lesson-one preview. It is a static illustration
 * (no runtime is loaded on the landing page, to keep it fast); "Run checks"
 * shows the two checks passing.
 */
export function TryItCard() {
  const [ran, setRan] = useState(false);
  return (
    <figure aria-label="Lesson one preview" className="rounded-lg border border-border bg-surface">
      <div className="flex items-center justify-between border-b border-border px-4 py-3">
        <span className="text-caption tracking-[0.08em] text-text-muted">TRY IT · LESSON 1</span>
        <span className="text-game-sm rounded-pill bg-accent-soft px-2 py-0.5 text-accent-text">+20 XP</span>
      </div>
      <div className="flex flex-col gap-3 p-4">
        <p className="text-small text-text-secondary">
          Make <code className="rounded-sm bg-surface-raised px-1 font-mono text-text">greet(name)</code> return a greeting,
          like <code className="rounded-sm bg-surface-raised px-1 font-mono text-text">"Hello, Ada!"</code>
        </p>
        <pre className="text-code-sm overflow-x-auto rounded-md border border-border bg-bg px-3 py-2" aria-label="Code">
          <span className="text-text-faint">1 </span>
          <span className="text-code-keyword">def</span> <span className="text-code-function">greet</span>
          <span className="text-text">(name):</span>
          {'\n'}
          <span className="text-text-faint">2 </span>
          {'    '}
          <span className="text-code-keyword">return</span> <span className="text-success">f"Hello, {'{name}'}!"</span>
        </pre>
        <div className="flex items-end justify-between gap-3">
          <ul className="text-code-sm flex flex-col gap-0.5 text-text-secondary" aria-label="Checks">
            {['greet("Ada") → "Hello, Ada!"', 'greet("Lin") → "Hello, Lin!"'].map((c) => (
              <li key={c} className="flex items-center gap-2">
                <span aria-hidden className={`h-2 w-2 rounded-pill ${ran ? 'bg-success' : 'bg-border-strong'}`} />
                <span>
                  {c}
                  <span className="sr-only">{ran ? ' passed' : ' not run'}</span>
                </span>
              </li>
            ))}
          </ul>
          <button
            type="button"
            onClick={() => setRan(true)}
            className="text-control h-[34px] shrink-0 rounded-md bg-accent px-4 text-accent-ink hover:brightness-105"
          >
            Run checks
          </button>
        </div>
        <p className="text-small min-h-5 text-success" role="status">
          {ran ? "2 of 2 checks passing. That's lesson one." : ''}
        </p>
      </div>
    </figure>
  );
}
