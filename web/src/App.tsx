import type { ReactNode } from 'react';
import { Route, Routes } from 'react-router';
import { RequireAuth } from './components/RequireAuth';
import { CoursePage } from './pages/CoursePage';
import { CoursesPage } from './pages/CoursesPage';
import { ExercisePage } from './pages/ExercisePage';
import { HomePage } from './pages/HomePage';
import { LeaguePage } from './pages/LeaguePage';
import { ProfilePage } from './pages/ProfilePage';
import { ShopPage } from './pages/ShopPage';
import { SignInPage } from './pages/SignInPage';
import { LoadError } from './pages/states';
import { VerifyPage } from './pages/VerifyPage';

const auth = (el: ReactNode) => <RequireAuth>{el}</RequireAuth>;

export function App() {
  return (
    <Routes>
      <Route path="/sign-in" element={<SignInPage />} />
      <Route path="/sign-in/verify" element={<VerifyPage />} />
      <Route path="/" element={auth(<HomePage />)} />
      <Route path="/courses" element={auth(<CoursesPage />)} />
      <Route path="/courses/:courseId" element={auth(<CoursePage />)} />
      <Route path="/learn/:courseId/:exerciseId" element={auth(<ExercisePage />)} />
      <Route path="/league" element={auth(<LeaguePage />)} />
      <Route path="/shop" element={auth(<ShopPage />)} />
      <Route path="/profile" element={auth(<ProfilePage />)} />
      <Route path="*" element={<LoadError error={new Error('Page not found')} />} />
    </Routes>
  );
}
