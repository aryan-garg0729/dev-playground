import { Component, type ErrorInfo, type ReactNode } from 'react';
import { reportError } from '../monitoring';

interface State {
  error?: Error;
}

/** Reports a render error and shows a way back instead of a blank page. */
export class ErrorBoundary extends Component<{ children: ReactNode }, State> {
  state: State = {};

  static getDerivedStateFromError(error: Error): State {
    return { error };
  }

  componentDidCatch(error: Error, info: ErrorInfo) {
    reportError(error, { componentStack: info.componentStack ?? undefined });
  }

  render() {
    if (!this.state.error) return this.props.children;
    return (
      <main className="flex min-h-dvh flex-col items-start gap-3 bg-bg p-8" role="alert">
        <h1 className="text-section text-text">Something went wrong</h1>
        <p className="text-small text-text-muted">The error was reported. Reloading usually fixes it.</p>
        <button
          type="button"
          onClick={() => window.location.reload()}
          className="text-control h-[34px] rounded-md bg-surface-raised px-4 text-text"
        >
          Reload
        </button>
      </main>
    );
  }
}
