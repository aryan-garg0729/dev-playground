/** The project_cli public view (api/internal/exercise/projectcli). Every test is visible. */
export interface CliPublicView {
  lesson_id: string;
  prompt: string;
  base_url?: string;
  allowed_os?: string[];
  steps: CliStep[];
}

export interface CliStep {
  description: string;
  cli_command?: { command: string; captures?: Capture[]; visible_tests: CliTest[] };
  http_request?: {
    request: { method: string; url: string; headers?: Record<string, string>; body_json?: unknown };
    captures?: Capture[];
    visible_tests: HttpTest[];
  };
}

interface Capture {
  name: string;
  regex?: string;
  json_path?: string;
}

export interface CliTest {
  exit_code?: number;
  stdout_contains_all?: string[];
  stdout_contains_none?: string[];
  stdout_lines_gt?: number;
}

export interface HttpTest {
  status_code?: number;
  body_contains?: string;
  json_body?: { path: string; op: 'eq' | 'ne' | 'gt' | 'lt' | 'contains'; value: unknown };
  header?: { name: string; contains: string };
}

const quote = (s: string) => JSON.stringify(s);

/** Same wording as the server's check names (projectcli.DescribeCLITest). */
export function describeCliTest(t: CliTest): string {
  if (t.exit_code !== undefined) return `exit code is ${t.exit_code}`;
  if (t.stdout_contains_all) return `output contains ${t.stdout_contains_all.map(quote).join(', ')}`;
  if (t.stdout_contains_none) return `output does not contain ${t.stdout_contains_none.map(quote).join(', ')}`;
  return `output has more than ${t.stdout_lines_gt ?? 0} lines`;
}

const ops = { eq: 'is', ne: 'is not', gt: 'is greater than', lt: 'is less than', contains: 'contains' } as const;

/** Same wording as projectcli.DescribeHTTPTest. */
export function describeHttpTest(t: HttpTest): string {
  if (t.status_code !== undefined) return `status code is ${t.status_code}`;
  if (t.body_contains !== undefined) return `body contains ${quote(t.body_contains)}`;
  if (t.json_body) return `JSON ${t.json_body.path} ${ops[t.json_body.op]} ${JSON.stringify(t.json_body.value)}`;
  return `header ${t.header?.name ?? ''} contains ${quote(t.header?.contains ?? '')}`;
}

/** The command or request line shown for a step. */
export function stepAction(s: CliStep): string {
  if (s.cli_command) return `$ ${s.cli_command.command}`;
  const r = s.http_request?.request;
  return r ? `> ${r.method} ${r.url}` : '';
}

export function stepTests(s: CliStep): string[] {
  if (s.cli_command) return s.cli_command.visible_tests.map(describeCliTest);
  return (s.http_request?.visible_tests ?? []).map(describeHttpTest);
}
