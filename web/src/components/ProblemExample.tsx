export interface Example {
  input: string;
  output: string;
  explanation?: string;
}

function Line({ label, value }: { label: string; value: string }) {
  return (
    <p>
      <span className="font-medium text-text">{label}</span> {value}
    </p>
  );
}

/** Examples and constraints in the problem-site convention (see ProblemExample README). */
export function ProblemExamples({ examples, constraints }: { examples: Example[]; constraints: string[] }) {
  return (
    <div className="flex flex-col gap-4">
      {examples.map((ex, i) => (
        <section key={i} className="flex flex-col gap-2">
          <h3 className="text-section text-text">Example {i + 1}:</h3>
          <div className="text-code-sm border-l-2 border-border-strong pl-[14px] text-text-secondary [overflow-wrap:anywhere]">
            <Line label="Input:" value={ex.input} />
            <Line label="Output:" value={ex.output} />
            {ex.explanation && <Line label="Explanation:" value={ex.explanation} />}
          </div>
        </section>
      ))}
      {constraints.length > 0 && (
        <section className="flex flex-col gap-2">
          <h3 className="text-section text-text">Constraints:</h3>
          <ul className="text-code-sm list-disc pl-[18px] text-text-secondary [overflow-wrap:anywhere]">
            {constraints.map((c) => (
              <li key={c}>{c}</li>
            ))}
          </ul>
        </section>
      )}
    </div>
  );
}
