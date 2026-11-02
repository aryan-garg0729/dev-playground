import { useState } from 'react';
import { ApiError } from '../../api/client';
import { Button } from '../../components/Button';
import { CheckIcon, CrossIcon } from '../../components/icons';
import { Markdown } from '../../components/Markdown';
import type { ExerciseRendererProps } from '../types';

export interface McqPublicView {
  question: string;
  multiple: boolean;
  options: { id: string; text: string }[];
}

export interface McqAnswer {
  choices: string[];
}

export function McqExercise({ publicView, draft, onChange, onSubmit, result }: ExerciseRendererProps<McqPublicView, McqAnswer>) {
  const { question, multiple, options } = publicView;
  const choices = draft?.choices ?? [];
  const [pending, setPending] = useState(false);
  const [error, setError] = useState<string>();

  function toggle(id: string) {
    setError(undefined);
    if (!multiple) return onChange({ choices: [id] });
    onChange({ choices: choices.includes(id) ? choices.filter((c) => c !== id) : [...choices, id] });
  }

  async function submit() {
    setPending(true);
    setError(undefined);
    try {
      await onSubmit({ choices });
    } catch (e) {
      setError(e instanceof ApiError ? e.message : 'Could not submit. Check your connection and try again.');
    } finally {
      setPending(false);
    }
  }

  const passed = result?.result.status === 'pass';
  const feedback = result?.result.checks.map((c) => c.feedback).filter(Boolean);

  return (
    <div className="flex flex-col gap-4">
      <Markdown>{question}</Markdown>

      <fieldset className="flex flex-col gap-2">
        <legend className="text-section mb-2 text-text">{multiple ? 'Select all that apply' : 'Choose one answer'}</legend>
        {options.map((o) => {
          const checked = choices.includes(o.id);
          return (
            <label
              key={o.id}
              className={`text-body flex min-h-[var(--touch-min)] cursor-pointer items-center gap-3 rounded-md border px-3 py-2 ${
                checked ? 'border-border-strong bg-surface-raised text-text' : 'border-border text-text-secondary hover:bg-surface-raised'
              }`}
            >
              <input
                type={multiple ? 'checkbox' : 'radio'}
                name="mcq-choice"
                value={o.id}
                checked={checked}
                onChange={() => toggle(o.id)}
                className="h-4 w-4 shrink-0"
              />
              <span className="font-mono text-[13px] leading-[23px]">{o.text}</span>
            </label>
          );
        })}
      </fieldset>

      <div className="flex items-center gap-3">
        <Button variant="primary" onClick={submit} disabled={pending || choices.length === 0}>
          {pending ? 'Checking…' : 'Submit'}
        </Button>
        {error && (
          <p role="alert" className="text-small text-danger">
            {error}
          </p>
        )}
      </div>

      <div aria-live="polite">
        {result && (
          <section className={`rounded-lg border border-border px-4 py-[14px] ${passed ? 'bg-success-soft' : 'bg-danger-soft'}`}>
            <p className={`text-section flex items-center gap-2 ${passed ? 'text-success' : 'text-danger'}`}>
              {passed ? <CheckIcon /> : <CrossIcon />}
              {passed ? 'Correct' : 'Not quite'}
            </p>
            {feedback?.map((f) => (
              <p key={f} className="text-body mt-2 text-text-secondary">
                {f}
              </p>
            ))}
            {!passed && <p className="text-small mt-2 text-text-muted">Change your answer and submit again.</p>}
          </section>
        )}
      </div>
    </div>
  );
}
