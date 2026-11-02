// Generates src/styles/tokens.css and tailwind.config.ts from docs/design/tokens.json.
//   node scripts/gen-tokens.ts           write both files
//   node scripts/gen-tokens.ts --check   exit 1 if either file differs from what would be generated
// Runs on Node's built-in type stripping, so it uses erasable TypeScript only.
import { readFileSync, writeFileSync } from 'node:fs';
import { dirname, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';

interface ColorToken { name: string; value: Record<string, string> }
interface TypeStyle { name: string; fontSize: string; lineHeight: string; fontWeight: number }
interface Tokens {
  color: { themes: { id: string }[]; tokens: ColorToken[] };
  type: { families: Record<string, string>; groups: { family: string; styles: TypeStyle[] }[] };
  spacing: { tokens: { name: string; value: string }[] };
  radius: { tokens: { name: string; value: string }[] };
  layout: { tokens: { name: string; value: string }[] };
}

const webDir = resolve(dirname(fileURLToPath(import.meta.url)), '..');
const tokensPath = process.env.TOKENS_JSON ?? resolve(webDir, '../docs/design/tokens.json');
const cssPath = resolve(webDir, 'src/styles/tokens.css');
const tailwindPath = resolve(webDir, 'tailwind.config.ts');

/** "{text}" -> "var(--text)"; anything else is a literal. */
function resolveValue(v: string): string {
  const alias = /^\{([a-z0-9-]+)\}$/.exec(v);
  return alias ? `var(--${alias[1]})` : v;
}

export function renderCss(t: Tokens): string {
  const out: string[] = ['/* CodeQuest design tokens, generated from docs/design/tokens.json. Do not edit by hand. */'];
  t.color.themes.forEach((theme, i) => {
    out.push(i === 0 ? `:root, [data-theme="${theme.id}"] {` : `[data-theme="${theme.id}"] {`);
    for (const tok of t.color.tokens) {
      const v = tok.value[theme.id];
      if (v === undefined) throw new Error(`color ${tok.name} has no value for theme ${theme.id}`);
      out.push(`  --${tok.name}: ${resolveValue(v)};`);
    }
    out.push(`  color-scheme: ${theme.id};`, '}', '');
  });
  out.push(':root {');
  for (const tok of [...t.spacing.tokens, ...t.radius.tokens, ...t.layout.tokens]) {
    out.push(`  --${tok.name}: ${tok.value};`);
  }
  for (const [name, stack] of Object.entries(t.type.families)) {
    out.push(`  --font-${name}: ${stack};`);
  }
  out.push('}', '');
  for (const group of t.type.groups) {
    for (const s of group.styles) {
      out.push(
        `.text-${s.name} { font-family: var(--font-${group.family}); font-size: ${s.fontSize}; ` +
          `line-height: ${s.lineHeight}; font-weight: ${s.fontWeight}; }`,
      );
    }
  }
  return out.join('\n') + '\n';
}

function stripPrefix(name: string, prefix: string): string {
  return name.startsWith(prefix) ? name.slice(prefix.length) : name;
}

function entries(pairs: [string, string][]): string {
  return pairs.map(([k, v]) => `        '${k}': '${v}'`).join(',\n');
}

export function renderTailwind(t: Tokens): string {
  const colors = entries(t.color.tokens.map((c) => [c.name, `var(--${c.name})`]));
  const spacing = entries(t.spacing.tokens.map((s) => [stripPrefix(s.name, 'space-'), `var(--${s.name})`]));
  const radius = entries(t.radius.tokens.map((r) => [stripPrefix(r.name, 'radius-'), `var(--${r.name})`]));
  const fonts = Object.entries(t.type.families)
    .map(([name, stack]) => `        ${name}: [${stack.split(',').map((f) => `'${f.trim()}'`).join(', ')}],`)
    .join('\n');
  return `// Generated from docs/design/tokens.json. Colors, spacing and radii point at CSS variables in src/styles/tokens.css,
// so the dark/light switch (data-theme on <html>) needs no Tailwind dark: variants.
import type { Config } from 'tailwindcss';

export default {
  content: ['./index.html', './src/**/*.{ts,tsx}'],
  theme: {
    extend: {
      colors: {
${colors}
      },
      spacing: {
${spacing}
      },
      borderRadius: {
${radius}
      },
      fontFamily: {
${fonts}
      },
    },
  },
} satisfies Config;
`;
}

function main(): void {
  const tokens = JSON.parse(readFileSync(tokensPath, 'utf8')) as Tokens;
  const files: [string, string][] = [
    [cssPath, renderCss(tokens)],
    [tailwindPath, renderTailwind(tokens)],
  ];
  if (process.argv.includes('--check')) {
    const stale = files.filter(([p, body]) => readFileSync(p, 'utf8') !== body).map(([p]) => p);
    if (stale.length > 0) {
      console.error(`Out of date with docs/design/tokens.json (run \`make tokens\`):\n  ${stale.join('\n  ')}`);
      process.exit(1);
    }
    console.log('tokens: up to date');
    return;
  }
  for (const [p, body] of files) writeFileSync(p, body);
  console.log('tokens: wrote', files.map(([p]) => p).join(', '));
}

main();
