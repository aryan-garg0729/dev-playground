import { describe, expect, it } from 'vitest';
import { monacoLanguage } from './language';

describe('monacoLanguage', () => {
  it('maps lesson languages to Monaco ids', () => {
    expect(monacoLanguage('js')).toBe('javascript');
    expect(monacoLanguage('go')).toBe('go');
    expect(monacoLanguage('python')).toBe('python');
  });
});
