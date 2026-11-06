import type { Profile } from '../../api/client';
import { ProjectRow } from './ProjectList';

/** Every chapter boss, grouped by course. */
export function ProjectsTab({ profile, icons }: { profile: Profile; icons: Record<string, string> }) {
  const courses = [...new Set(profile.projects.map((p) => p.course_id))];
  if (courses.length === 0) return <p className="text-body text-text-muted">No projects yet.</p>;
  return (
    <>
      {courses.map((id) => {
        const list = profile.projects.filter((p) => p.course_id === id);
        const done = list.filter((p) => p.status === 'defeated').length;
        return (
          <section key={id} aria-labelledby={`proj-${id}`} className="flex flex-col gap-2 rounded-lg border border-border bg-surface px-6 py-4">
            <h2 id={`proj-${id}`} className="text-section flex items-baseline gap-2 text-text">
              {list[0]?.course_title} <span className="text-small font-normal text-text-muted">{done} of {list.length} shipped</span>
            </h2>
            <ul>
              {list.map((p) => (
                <ProjectRow key={p.exercise_id} p={p} icon={icons[p.course_id] ?? ''} />
              ))}
            </ul>
          </section>
        );
      })}
    </>
  );
}
