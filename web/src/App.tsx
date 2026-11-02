import { Route, Routes } from 'react-router';
import { CoursePage } from './pages/CoursePage';
import { CoursesPage } from './pages/CoursesPage';
import { ExercisePage } from './pages/ExercisePage';
import { LoadError } from './pages/states';

export function App() {
  return (
    <Routes>
      <Route path="/" element={<CoursesPage />} />
      <Route path="/courses/:courseId" element={<CoursePage />} />
      <Route path="/learn/:courseId/:exerciseId" element={<ExercisePage />} />
      <Route path="*" element={<LoadError error={new Error('Page not found')} />} />
    </Routes>
  );
}
