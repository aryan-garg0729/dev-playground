// Builds Monaco themes from the design tokens at runtime. Monaco only accepts
// hex colors, so each CSS variable is resolved by the browser and converted;
// no color values live in this file.
import type { editor } from 'monaco-editor/editor';

/** Resolves a CSS custom property (e.g. "--surface") to #rrggbb or #rrggbbaa. */
export function tokenHex(name: string, root: Element = document.documentElement): string {
  const probe = document.createElement('span');
  probe.style.color = `var(${name})`;
  probe.style.display = 'none';
  root.appendChild(probe);
  const rgb = getComputedStyle(probe).color;
  probe.remove();
  return rgbToHex(rgb);
}

export function rgbToHex(rgb: string): string {
  const m = /rgba?\(\s*([\d.]+)[,\s]+([\d.]+)[,\s]+([\d.]+)(?:[,\s/]+([\d.]+%?))?\s*\)/.exec(rgb);
  if (!m) return rgb.startsWith('#') ? rgb : '#00000000';
  const h = (n: number) => Math.round(Math.max(0, Math.min(255, n))).toString(16).padStart(2, '0');
  let a = '';
  if (m[4] !== undefined) {
    const v = m[4].endsWith('%') ? parseFloat(m[4]) / 100 : parseFloat(m[4]);
    if (v < 1) a = h(v * 255);
  }
  const [r, g, b] = [m[1], m[2], m[3]].map((v) => Number(v ?? 0));
  return `#${h(r ?? 0)}${h(g ?? 0)}${h(b ?? 0)}${a}`;
}

const bare = (hex: string) => hex.slice(1);

export function buildTheme(base: 'vs' | 'vs-dark'): editor.IStandaloneThemeData {
  const t = (n: string) => tokenHex(`--${n}`);
  return {
    base,
    inherit: true,
    rules: [
      { token: '', foreground: bare(t('text')) },
      { token: 'keyword', foreground: bare(t('code-keyword')) },
      { token: 'number', foreground: bare(t('code-number')) },
      { token: 'string', foreground: bare(t('text-secondary')) },
      { token: 'comment', foreground: bare(t('text-muted')), fontStyle: 'italic' },
      { token: 'identifier', foreground: bare(t('text')) },
      { token: 'type.identifier', foreground: bare(t('code-function')) },
      { token: 'delimiter', foreground: bare(t('text-secondary')) },
    ],
    colors: {
      'editor.background': t('surface'),
      'editor.foreground': t('text'),
      'editorLineNumber.foreground': t('text-faint'),
      'editorLineNumber.activeForeground': t('text-muted'),
      'editorCursor.foreground': t('accent'),
      'editor.lineHighlightBackground': t('surface'),
      'editor.lineHighlightBorder': t('surface'),
      'editor.selectionBackground': t('border-strong'),
      'editor.inactiveSelectionBackground': t('surface-raised'),
      'editorIndentGuide.background1': t('border'),
      'editorWhitespace.foreground': t('border'),
      'editorWidget.background': t('surface-raised'),
      'editorWidget.border': t('border'),
      'scrollbarSlider.background': t('border'),
      'scrollbarSlider.hoverBackground': t('border-strong'),
    },
  };
}
