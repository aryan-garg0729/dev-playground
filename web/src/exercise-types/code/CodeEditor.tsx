// Monaco, self-hosted (no CDN) and loaded lazily: import this module only via React.lazy.
import Editor, { loader, type OnMount } from '@monaco-editor/react';
import * as monaco from 'monaco-editor/editor';
import 'monaco-editor/languages/definitions/python/register';
import EditorWorker from 'monaco-editor/editor/editor.worker?worker';
import { useEffect, useRef } from 'react';
import { useTheme } from '../../stores/theme';
import { buildTheme } from './monacoTheme';

self.MonacoEnvironment = { getWorker: () => new EditorWorker() };
loader.config({ monaco: monaco as unknown as Parameters<typeof loader.config>[0]['monaco'] });

export interface CodeEditorProps {
  value: string;
  onChange: (value: string) => void;
  language: string;
  /** 1-based line to tint as the failing line, if any. */
  errorLine?: number;
  ariaLabel: string;
}

function applyTheme(theme: 'dark' | 'light') {
  const name = `cq-${theme}`;
  monaco.editor.defineTheme(name, buildTheme(theme === 'dark' ? 'vs-dark' : 'vs'));
  monaco.editor.setTheme(name);
}

export default function CodeEditor({ value, onChange, language, errorLine, ariaLabel }: CodeEditorProps) {
  const theme = useTheme((s) => s.theme);
  const editorRef = useRef<monaco.editor.IStandaloneCodeEditor | null>(null);
  const decorations = useRef<monaco.editor.IEditorDecorationsCollection | null>(null);

  useEffect(() => applyTheme(theme), [theme]);

  useEffect(() => {
    const ed = editorRef.current;
    if (!ed) return;
    decorations.current?.clear();
    if (errorLine && errorLine <= (ed.getModel()?.getLineCount() ?? 0)) {
      decorations.current = ed.createDecorationsCollection([
        {
          range: new monaco.Range(errorLine, 1, errorLine, 1),
          options: { isWholeLine: true, className: 'cq-error-line', lineNumberClassName: 'cq-error-number' },
        },
      ]);
    }
  }, [errorLine]);

  const onMount: OnMount = (ed) => {
    editorRef.current = ed as unknown as monaco.editor.IStandaloneCodeEditor;
    applyTheme(theme);
  };

  return (
    <Editor
      value={value}
      language={language}
      theme={`cq-${theme}`}
      onChange={(v) => onChange(v ?? '')}
      onMount={onMount}
      loading={<p className="text-small p-4 text-text-muted">Loading editor…</p>}
      options={{
        ariaLabel,
        fontFamily: 'JetBrains Mono, ui-monospace, monospace',
        fontSize: 14,
        lineHeight: 24,
        minimap: { enabled: false },
        scrollBeyondLastLine: false,
        renderLineHighlight: 'none',
        overviewRulerLanes: 0,
        hideCursorInOverviewRuler: true,
        padding: { top: 16, bottom: 16 },
        tabSize: 4,
        insertSpaces: true,
        automaticLayout: true,
        fixedOverflowWidgets: true,
        contextmenu: false,
        quickSuggestions: false,
      }}
    />
  );
}
