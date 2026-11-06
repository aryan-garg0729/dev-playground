import { createContext, type ReactNode, useContext } from 'react';

/** The open lesson's images: name in the content ("assets/loop.svg") -> URL. */
const LessonAssets = createContext<Record<string, string>>({});

export function LessonAssetsProvider({ assets, children }: { assets: Record<string, string>; children: ReactNode }) {
  return <LessonAssets.Provider value={assets}>{children}</LessonAssets.Provider>;
}

/**
 * Resolves an image path from lesson content to its URL, or undefined. Only the
 * lesson's own images resolve: anything else (another site, a typo) does not,
 * so the browser never fetches an image from outside the API.
 */
export function useLessonAsset(): (src: string | undefined) => string | undefined {
  const assets = useContext(LessonAssets);
  return (src) => {
    if (!src) return undefined;
    const name = src.replace(/^\.\//, '');
    return Object.prototype.hasOwnProperty.call(assets, name) ? assets[name] : undefined;
  };
}
