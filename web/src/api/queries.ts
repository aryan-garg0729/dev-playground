import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { api, unwrap } from './client';

export const keys = {
  me: ['me'] as const,
  courses: ['courses'] as const,
  course: (id: string) => ['course', id] as const,
  exercise: (id: string) => ['exercise', id] as const,
};

export function useMe() {
  return useQuery({
    queryKey: keys.me,
    queryFn: () => unwrap(api.GET('/me')),
    retry: false,
  });
}

export function useCourses() {
  return useQuery({
    queryKey: keys.courses,
    queryFn: () => unwrap(api.GET('/courses')),
  });
}

export function useCourse(courseId: string) {
  return useQuery({
    queryKey: keys.course(courseId),
    queryFn: () => unwrap(api.GET('/courses/{courseId}', { params: { path: { courseId } } })),
  });
}

export function useExercise(exerciseId: string) {
  return useQuery({
    queryKey: keys.exercise(exerciseId),
    queryFn: () => unwrap(api.GET('/exercises/{exerciseId}', { params: { path: { exerciseId } } })),
  });
}

export function useSubmit(exerciseId: string, courseId: string | undefined) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (body: { version: number; answer: unknown }) =>
      unwrap(
        api.POST('/exercises/{exerciseId}/submit', {
          params: { path: { exerciseId } },
          body: { exercise_id: exerciseId, version: body.version, answer: body.answer },
        }),
      ),
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: keys.courses });
      if (courseId) void qc.invalidateQueries({ queryKey: keys.course(courseId) });
    },
  });
}
