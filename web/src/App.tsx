import { lazy, type ReactNode, Suspense, useEffect } from 'react';
import { Navigate, Route, Routes, useLocation } from 'react-router';
import { track } from './analytics';
import { ApiError } from './api/client';
import { useMe } from './api/queries';
import { RequireAuth } from './components/RequireAuth';
import { LandingPage } from './pages/LandingPage';
import { LoadError, Loading } from './pages/states';

// Everything but the landing page is split out, so a first visit loads only
// what the landing page needs (Lighthouse budget, SPEC Phase 6).
const HomePage = lazy(() => import('./pages/HomePage').then((m) => ({ default: m.HomePage })));
const CoursesPage = lazy(() => import('./pages/CoursesPage').then((m) => ({ default: m.CoursesPage })));
const CoursePage = lazy(() => import('./pages/CoursePage').then((m) => ({ default: m.CoursePage })));
const ExercisePage = lazy(() => import('./pages/ExercisePage').then((m) => ({ default: m.ExercisePage })));
const LeaguePage = lazy(() => import('./pages/LeaguePage').then((m) => ({ default: m.LeaguePage })));
const ShopPage = lazy(() => import('./pages/ShopPage').then((m) => ({ default: m.ShopPage })));
const ProfilePage = lazy(() => import('./pages/ProfilePage').then((m) => ({ default: m.ProfilePage })));
const SignInPage = lazy(() => import('./pages/SignInPage').then((m) => ({ default: m.SignInPage })));
const CliLoginPage = lazy(() => import('./pages/CliLoginPage').then((m) => ({ default: m.CliLoginPage })));
const OnboardingPage = lazy(() => import('./pages/OnboardingPage').then((m) => ({ default: m.OnboardingPage })));
// Dev only: import.meta.env.DEV is false in production builds, so this branch
// and its chunk are dropped from the bundle entirely.
const EffectsLab = import.meta.env.DEV ? lazy(() => import('./pages/dev/EffectsLab').then((m) => ({ default: m.EffectsLab }))) : null;

const auth = (el: ReactNode) => <RequireAuth>{el}</RequireAuth>;

/** "/" is the landing page for visitors, onboarding for new learners, and Home after that. */
function Root() {
  const me = useMe();
  if (me.isPending) return <Loading label="CodeQuest" />;
  if (me.error instanceof ApiError && me.error.status === 401) return <LandingPage />;
  if (me.error) return <LoadError error={me.error} />;
  if (!me.data.user.onboarded) return <Navigate to="/welcome" replace />;
  return <HomePage />;
}

function PageViews() {
  const { pathname } = useLocation();
  useEffect(() => {
    track('page_view');
  }, [pathname]);
  return null;
}

export function App() {
  return (
    <>
      <PageViews />
      <Suspense fallback={<Loading label="page" />}>
        <Routes>
          <Route path="/sign-in" element={<SignInPage />} />
          <Route path="/" element={<Root />} />
          <Route path="/welcome" element={auth(<OnboardingPage />)} />
          <Route path="/cli/login" element={auth(<CliLoginPage />)} />
          <Route path="/courses" element={auth(<CoursesPage />)} />
          <Route path="/courses/:courseId" element={auth(<CoursePage />)} />
          <Route path="/learn/:courseId/:exerciseId" element={auth(<ExercisePage />)} />
          <Route path="/league" element={auth(<LeaguePage />)} />
          <Route path="/shop" element={auth(<ShopPage />)} />
          <Route path="/profile/:tab?" element={auth(<ProfilePage />)} />
          {EffectsLab && <Route path="/dev/effects" element={<EffectsLab />} />}
          <Route path="*" element={<LoadError error={new Error('Page not found')} />} />
        </Routes>
      </Suspense>
    </>
  );
}
