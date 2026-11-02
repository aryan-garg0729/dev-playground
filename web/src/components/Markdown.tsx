import ReactMarkdown from 'react-markdown';

/**
 * Renders author-written markdown (lesson prompts). Raw HTML in the source is not
 * rendered (react-markdown's default), and nothing goes through innerHTML.
 */
export function Markdown({ children }: { children: string }) {
  return (
    <div className="text-body flex flex-col gap-3 text-text-secondary">
      <ReactMarkdown
        components={{
          strong: ({ children }) => <strong className="font-semibold text-text">{children}</strong>,
          pre: ({ children }) => (
            <pre className="text-code-sm overflow-x-auto border-l-2 border-border-strong pl-[14px] text-text-secondary [&_code]:bg-transparent [&_code]:p-0 [&_code]:text-text-secondary">{children}</pre>
          ),
          code: ({ children, className }) =>
            className ? (
              <code>{children}</code>
            ) : (
              <code className="rounded-sm bg-surface-raised px-1.5 py-0.5 font-mono text-[13px] text-text">{children}</code>
            ),
        }}
      >
        {children}
      </ReactMarkdown>
    </div>
  );
}
