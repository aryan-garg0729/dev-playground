import { Link } from 'react-router';
import { CopyCommand } from '../components/CopyCommand';
import { AuthShell } from './SignInPage';

/**
 * How to install the `codequest` CLI. The deployed site serves the binaries
 * (built for this site, so they talk to it) and an install script; this page
 * builds the commands from its own origin. Public: no sign-in needed.
 */
export function CliInstallPage() {
  const origin = window.location.origin;
  const exe = (arch: string) => `${origin}/cli/codequest-windows-${arch}.exe`;
  return (
    <AuthShell>
      <div className="flex flex-col gap-2">
        <h1 className="text-page-title text-text">Install the CodeQuest CLI</h1>
        <p className="text-body text-text-secondary">
          Command-line lessons run on your computer. The <code className="font-mono text-text">codequest</code> command runs each
          lesson's checks and sends the result here.
        </p>
      </div>

      <section className="flex flex-col gap-3" aria-label="macOS and Linux">
        <h2 className="text-section text-text">macOS and Linux</h2>
        <p className="text-small text-text-secondary">
          Installs into <code className="font-mono text-text">~/.local/bin</code>. No admin rights needed.
        </p>
        <CopyCommand command={`curl -fsSL ${origin}/install.sh | sh`} />
      </section>

      <section className="flex flex-col gap-3" aria-label="Windows">
        <h2 className="text-section text-text">Windows</h2>
        <p className="text-small text-text-secondary">
          Download{' '}
          <a href={exe('amd64')} className="text-accent-text underline">
            codequest.exe (Intel/AMD)
          </a>{' '}
          or{' '}
          <a href={exe('arm64')} className="text-accent-text underline">
            codequest.exe (ARM)
          </a>
          , rename it to <code className="font-mono text-text">codequest.exe</code> and put it in a folder on your PATH.
        </p>
      </section>

      <section className="flex flex-col gap-3" aria-label="Sign in">
        <h2 className="text-section text-text">Then sign in</h2>
        <CopyCommand command="codequest login" />
        <p className="text-small text-text-muted">
          It opens this site and shows a code to paste into the terminal. Then start a CLI lesson from{' '}
          <Link to="/courses" className="text-accent-text underline">
            Courses
          </Link>
          .
        </p>
      </section>
    </AuthShell>
  );
}
