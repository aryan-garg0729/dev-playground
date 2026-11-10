// A tiny syntax highlighter for read-only code in lesson prompts. It colours the same
// token kinds as the editor theme (monacoTheme.ts) using the design tokens, and
// returns plain token objects that the caller renders as React spans (no HTML strings).

export type TokenKind = 'keyword' | 'function' | 'constant' | 'string' | 'comment' | 'ident' | 'plain';

export interface Token {
  kind: TokenKind;
  text: string;
}

interface Grammar {
  keywords: Set<string>;
  constants: Set<string>;
  /** Comment and string patterns, tried in order before anything else. */
  patterns: [TokenKind, RegExp][];
}

const words = (s: string) => new Set(s.split(/\s+/).filter(Boolean));

const GRAMMARS: Record<string, Grammar> = {
  python: {
    keywords: words(`and as assert async await break class continue def del elif else except finally for from
      global if import in is lambda nonlocal not or pass raise return try while with yield`),
    constants: words('True False None'),
    patterns: [
      ['comment', /#[^\n]*/y],
      ['string', /[rRbBfFuU]{0,2}("""[\s\S]*?"""|'''[\s\S]*?'''|"(?:\\.|[^"\\\n])*"|'(?:\\.|[^'\\\n])*')/y],
    ],
  },
  cpp: {
    keywords: words(`auto bool break case char class const constexpr continue default delete do double else enum
      explicit float for friend if inline int long mutable namespace new operator private protected public
      return short signed size_t sizeof static static_cast struct switch template this typedef typename
      unsigned using virtual void volatile while`),
    constants: words('true false nullptr NULL'),
    patterns: [
      ['comment', /\/\/[^\n]*|\/\*[\s\S]*?\*\//y],
      ['keyword', /#\s*[a-z]+(?:[ \t]*<[^>\n]*>)?/y],
      ['string', /"(?:\\.|[^"\\\n])*"|'(?:\\.|[^'\\\n])*'/y],
    ],
  },
};

/** Languages the highlighter knows; other fenced blocks stay plain. */
export const HIGHLIGHT_LANGS = Object.keys(GRAMMARS);

const NUMBER = /(?:0[xX][0-9a-fA-F']+|\d[\d_']*(?:\.\d+)?(?:[eE][+-]?\d+)?)[uUlLfFjJ]*/y;
const IDENT = /[A-Za-z_]\w*/y;
const CALL = /\s*\(/y;

function at(re: RegExp, src: string, i: number): string | null {
  re.lastIndex = i;
  const m = re.exec(src);
  return m ? m[0] : null;
}

/** True right after `class` / `struct` (and whitespace): the next name is a type. */
function afterTypeKeyword(out: Token[]): boolean {
  const n = out.length;
  const prev = n > 0 && out[n - 1]!.kind === 'plain' && /^\s+$/.test(out[n - 1]!.text) ? out[n - 2] : undefined;
  return prev?.kind === 'keyword' && (prev.text === 'class' || prev.text === 'struct');
}

/** Splits `src` into coloured tokens. Unknown languages come back as one plain token. */
export function highlight(src: string, lang: string): Token[] {
  const g = GRAMMARS[lang];
  if (!g) return [{ kind: 'plain', text: src }];
  const out: Token[] = [];
  const push = (kind: TokenKind, text: string) => {
    const last = out[out.length - 1];
    if (last && last.kind === kind) last.text += text;
    else out.push({ kind, text });
  };
  let i = 0;
  outer: while (i < src.length) {
    for (const [kind, re] of g.patterns) {
      const m = at(re, src, i);
      if (m) {
        push(kind, m);
        i += m.length;
        continue outer;
      }
    }
    const ident = at(IDENT, src, i);
    if (ident) {
      const kind: TokenKind = g.keywords.has(ident)
        ? 'keyword'
        : g.constants.has(ident)
          ? 'constant'
          : at(CALL, src, i + ident.length) !== null || /^[A-Z][a-z]/.test(ident) || afterTypeKeyword(out)
            ? 'function'
            : 'ident';
      push(kind, ident);
      i += ident.length;
      continue;
    }
    const num = at(NUMBER, src, i);
    if (num) {
      push('constant', num);
      i += num.length;
      continue;
    }
    push('plain', src[i]!);
    i++;
  }
  return out;
}

/** Tailwind classes per token kind, mirroring the editor theme. */
export const TOKEN_CLASS: Record<TokenKind, string | undefined> = {
  keyword: 'text-code-keyword',
  function: 'text-code-function',
  constant: 'text-code-number',
  string: 'text-text-secondary',
  comment: 'italic text-text-muted',
  ident: 'text-text',
  plain: undefined,
};
