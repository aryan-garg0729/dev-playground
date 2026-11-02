import { rgbToHex } from './monacoTheme';

describe('rgbToHex', () => {
  it.each([
    ['rgb(31, 31, 35)', '#1f1f23'],
    ['rgba(245, 165, 36, 0.5)', '#f5a52480'],
    ['rgb(255 255 255 / 100%)', '#ffffff'],
    ['#abcdef', '#abcdef'],
  ])('%s -> %s', (input, want) => {
    expect(rgbToHex(input)).toBe(want);
  });
});
