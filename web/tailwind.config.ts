// Generated from docs/design/tokens.json. Colors, spacing and radii point at CSS variables in src/styles/tokens.css,
// so the dark/light switch (data-theme on <html>) needs no Tailwind dark: variants.
import type { Config } from 'tailwindcss';

export default {
  content: ['./index.html', './src/**/*.{ts,tsx}'],
  theme: {
    extend: {
      colors: {
        'bg': 'var(--bg)',
        'surface': 'var(--surface)',
        'surface-raised': 'var(--surface-raised)',
        'border': 'var(--border)',
        'border-strong': 'var(--border-strong)',
        'text': 'var(--text)',
        'text-secondary': 'var(--text-secondary)',
        'text-muted': 'var(--text-muted)',
        'text-faint': 'var(--text-faint)',
        'accent': 'var(--accent)',
        'accent-ink': 'var(--accent-ink)',
        'accent-text': 'var(--accent-text)',
        'accent-soft': 'var(--accent-soft)',
        'accent-line': 'var(--accent-line)',
        'progress-current': 'var(--progress-current)',
        'success': 'var(--success)',
        'success-soft': 'var(--success-soft)',
        'danger': 'var(--danger)',
        'danger-soft': 'var(--danger-soft)',
        'level-bg': 'var(--level-bg)',
        'level-ink': 'var(--level-ink)',
        'code-keyword': 'var(--code-keyword)',
        'code-function': 'var(--code-function)',
        'code-number': 'var(--code-number)',
        'medal-bronze': 'var(--medal-bronze)',
        'medal-silver': 'var(--medal-silver)',
        'medal-gold': 'var(--medal-gold)',
        'medal-diamond': 'var(--medal-diamond)',
        'boss-hp': 'var(--boss-hp)'
      },
      spacing: {
        '1': 'var(--space-1)',
        '2': 'var(--space-2)',
        '3': 'var(--space-3)',
        '4': 'var(--space-4)',
        '5': 'var(--space-5)',
        '6': 'var(--space-6)'
      },
      borderRadius: {
        'sm': 'var(--radius-sm)',
        'md': 'var(--radius-md)',
        'lg': 'var(--radius-lg)',
        'pill': 'var(--radius-pill)'
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', 'sans-serif'],
        mono: ['"JetBrains Mono"', 'ui-monospace', 'monospace'],
        pixel: ['"Pixelify Sans"', 'Inter', 'sans-serif'],
      },
    },
  },
} satisfies Config;
