import { screen } from '@testing-library/react';
import { renderRoute } from '../test/render';
import { CliInstallPage } from './CliInstallPage';

describe('CliInstallPage', () => {
  it('builds the install commands and downloads from this site', () => {
    renderRoute('/cli/install', '/cli/install', <CliInstallPage />);
    const origin = window.location.origin;
    expect(screen.getByText(`curl -fsSL ${origin}/install.sh | sh`)).toBeInTheDocument();
    expect(screen.getByRole('link', { name: 'codequest.exe (Intel/AMD)' })).toHaveAttribute('href', `${origin}/cli/codequest-windows-amd64.exe`);
    expect(screen.getByText('codequest login')).toBeInTheDocument();
  });
});
