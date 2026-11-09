import { useState } from 'react';
import { Button } from './Button';

/** A shell command with a Copy button. */
export function CopyCommand({ command }: { command: string }) {
  const [copied, setCopied] = useState(false);
  return (
    <div className="flex items-center justify-between gap-3 rounded-md border border-border bg-surface-raised px-3 py-2">
      <code className="text-code-sm min-w-0 font-mono text-text [overflow-wrap:anywhere]">{command}</code>
      <Button
        onClick={() => {
          navigator.clipboard.writeText(command).then(
            () => setCopied(true),
            () => setCopied(false),
          );
        }}
      >
        {copied ? 'Copied' : 'Copy'}
      </Button>
    </div>
  );
}
