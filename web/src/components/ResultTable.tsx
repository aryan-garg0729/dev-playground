import type { Table } from '../runner-client/protocol';

const SHOWN = 200;

/** A query result. Cells are rendered as text. */
export function ResultTable({ table }: { table: Table }) {
  const rows = table.rows.slice(0, SHOWN);
  return (
    <div className="overflow-x-auto rounded-lg border border-border bg-surface">
      <table className="text-code-sm w-full border-collapse text-left">
        <thead>
          <tr className="border-b border-border">
            {table.columns.map((c, i) => (
              <th key={i} scope="col" className="px-4 py-2 font-medium text-text">
                {c}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {rows.map((r, i) => (
            <tr key={i} className="border-b border-border last:border-b-0">
              {r.map((cell, j) => (
                <td key={j} className={`px-4 py-1.5 ${cell === null ? 'text-text-muted' : 'text-text-secondary'}`}>
                  {cell === null ? 'NULL' : String(cell)}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
      <p className="text-caption border-t border-border px-4 py-2 text-text-muted">
        {table.rows.length} row{table.rows.length === 1 ? '' : 's'}
        {(table.truncated || table.rows.length > SHOWN) && ` · showing the first ${rows.length}`}
      </p>
    </div>
  );
}
