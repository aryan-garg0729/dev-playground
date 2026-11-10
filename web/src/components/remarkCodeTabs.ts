import type { Code, Parent, Root, RootContent } from 'mdast';

/** Languages that can share a tabbed block, in tab order. */
export const TAB_LANGS: Record<string, string> = { cpp: 'C++', python: 'Python' };

/** Marker language of the merged node; its value is JSON of {@link CodeTab}[]. */
export const TABS_LANG = 'cq-tabs';

export interface CodeTab {
  lang: string;
  code: string;
}

const tabbable = (n: RootContent | undefined): n is Code => n?.type === 'code' && !!n.lang && n.lang in TAB_LANGS;
const order = Object.keys(TAB_LANGS);

/**
 * Remark plugin: back-to-back fenced blocks in different tab languages (```cpp then ```python)
 * become one block with a tab per language. A lone block, or a run that repeats a language,
 * stays as it is.
 */
export function remarkCodeTabs() {
  return (tree: Root) => walk(tree);
}

function walk(node: Parent) {
  const out: RootContent[] = [];
  const kids = node.children as RootContent[];
  for (let i = 0; i < kids.length; i++) {
    const n = kids[i];
    if (!tabbable(n)) {
      if (n && 'children' in n) walk(n);
      if (n) out.push(n);
      continue;
    }
    const run: Code[] = [n];
    while (tabbable(kids[i + 1]) && !run.some((c) => c.lang === (kids[i + 1] as Code).lang)) {
      run.push(kids[i + 1] as Code);
      i++;
    }
    if (run.length === 1) {
      out.push(n);
      continue;
    }
    const tabs: CodeTab[] = run
      .map((c) => ({ lang: c.lang as string, code: c.value }))
      .sort((a, b) => order.indexOf(a.lang) - order.indexOf(b.lang));
    out.push({ type: 'code', lang: TABS_LANG, value: JSON.stringify(tabs), position: n.position });
  }
  node.children = out as typeof node.children;
}
