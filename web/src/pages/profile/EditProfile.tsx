import { type FormEvent, useState } from 'react';
import { useUpdateProfile } from '../../api/queries';
import { Button } from '../../components/Button';

/** Inline form for the display name and the worn title (only earned titles are offered). */
export function EditProfile({ name, title, titles, onDone }: { name: string; title: string; titles: string[]; onDone: () => void }) {
  const [draft, setDraft] = useState({ name, title });
  const save = useUpdateProfile();
  const submit = (e: FormEvent) => {
    e.preventDefault();
    save.mutate({ name: draft.name.trim(), title: draft.title }, { onSuccess: onDone });
  };
  const field = 'text-control h-[34px] rounded-md border border-border bg-bg px-2 text-text';
  return (
    <form aria-label="Edit profile" onSubmit={submit} className="flex flex-wrap items-end gap-4 rounded-lg border border-border bg-surface p-5">
      <label className="text-small flex flex-col gap-1 text-text-muted">
        Name
        <input value={draft.name} maxLength={80} required onChange={(e) => setDraft({ ...draft, name: e.target.value })} className={field} />
      </label>
      <label className="text-small flex flex-col gap-1 text-text-muted">
        Title
        <select value={draft.title} onChange={(e) => setDraft({ ...draft, title: e.target.value })} className={field}>
          <option value="">None</option>
          {titles.map((t) => (
            <option key={t} value={t}>
              {t}
            </option>
          ))}
        </select>
      </label>
      <div className="flex gap-2">
        <Button variant="primary" type="submit" disabled={save.isPending}>
          {save.isPending ? 'Saving…' : 'Save'}
        </Button>
        <Button onClick={onDone}>Cancel</Button>
      </div>
      {save.error && (
        <p role="alert" className="text-small w-full text-danger">
          {save.error.message}
        </p>
      )}
    </form>
  );
}
