import { heatmapWeeks, level, totals } from './Heatmap';

describe('heatmapWeeks', () => {
  const activity = { from: '2025-10-30', to: '2026-11-04', longest_streak: 3, days: [{ date: '2026-11-02', lessons: 2 }, { date: '2026-10-26', lessons: 5 }] };

  it('ends with the week of today, Monday first, and hides future days', () => {
    const cols = heatmapWeeks(activity, 2);
    expect(cols.map((c) => c[0]?.date)).toEqual(['2026-10-26', '2026-11-02']);
    expect(cols[1]?.map((c) => c?.date ?? null)).toEqual(['2026-11-02', '2026-11-03', '2026-11-04', null, null, null, null]);
    expect(totals(cols)).toEqual({ lessons: 7, days: 2 });
    expect(totals(heatmapWeeks(activity, 1))).toEqual({ lessons: 2, days: 1 });
  });

  it('keeps dates stable across timezones and steps colors up to 4+', () => {
    expect(heatmapWeeks({ ...activity, to: '2026-03-30' }, 1)[0]?.[0]?.date).toBe('2026-03-30'); // a Monday after a DST change
    expect(level(0)).toContain('surface-raised');
    expect(level(9)).toBe(level(4));
  });
});

describe('Heatmap', () => {
  it('labels the current month at the right edge', async () => {
    const { render, screen } = await import('@testing-library/react');
    const { Heatmap } = await import('./Heatmap');
    render(<Heatmap activity={{ from: '2025-10-30', to: '2026-11-04', longest_streak: 0, days: [] }} weeks={6} label="x" />);
    expect(screen.getByRole('img').textContent).toContain('Nov');
  });
});
