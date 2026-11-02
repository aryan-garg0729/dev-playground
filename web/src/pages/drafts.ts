// Unsent answers, kept per exercise and version in this browser so a reload
// does not lose work. (Server-side drafts, PUT /exercises/{id}/draft, come later.)
const key = (id: string, version: number) => `cq-draft:${id}:v${version}`;

export function loadDraft(id: string, version: number): unknown {
  try {
    const raw = localStorage.getItem(key(id, version));
    return raw === null ? undefined : (JSON.parse(raw) as unknown);
  } catch {
    return undefined;
  }
}

export function saveDraft(id: string, version: number, answer: unknown): void {
  try {
    localStorage.setItem(key(id, version), JSON.stringify(answer));
  } catch {
    // storage full or unavailable: the draft lives only in memory
  }
}
