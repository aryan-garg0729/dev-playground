import { highlight, type TokenKind } from './highlight';

const kinds = (src: string, lang: string) =>
  highlight(src, lang)
    .filter((t) => t.kind !== 'plain')
    .map((t) => [t.kind, t.text] as [TokenKind, string]);

describe('highlight', () => {
  it('colours Python keywords, calls, numbers, constants, strings and comments', () => {
    expect(kinds('def f(x):\n    return max(x, 10) or None  # done\ns = "a # b"', 'python')).toEqual([
      ['keyword', 'def'],
      ['function', 'f'],
      ['ident', 'x'],
      ['keyword', 'return'],
      ['function', 'max'],
      ['ident', 'x'],
      ['constant', '10'],
      ['keyword', 'or'],
      ['constant', 'None'],
      ['comment', '# done'],
      ['ident', 's'],
      ['string', '"a # b"'],
    ]);
  });

  it('colours C++ types, preprocessor lines, class names and both comment styles', () => {
    expect(kinds('#include <vector>\nstruct DSU { int n = 1e9; }; // x\n/* y */ vector<int> v;', 'cpp')).toEqual([
      ['keyword', '#include <vector>'],
      ['keyword', 'struct'],
      ['function', 'DSU'],
      ['keyword', 'int'],
      ['ident', 'n'],
      ['constant', '1e9'],
      ['comment', '// x'],
      ['comment', '/* y */'],
      ['ident', 'vector'],
      ['keyword', 'int'],
      ['ident', 'v'],
    ]);
  });

  it('keeps every character, in order', () => {
    const src = 'for (int i = 0; i < n; i++) s += "a\\"b";\n';
    expect(highlight(src, 'cpp').map((t) => t.text).join('')).toBe(src);
  });

  it('returns unknown languages as one plain token', () => {
    expect(highlight('ls -la', 'sh')).toEqual([{ kind: 'plain', text: 'ls -la' }]);
  });
});
