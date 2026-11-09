#!/bin/sh
# Installs the CodeQuest CLI from __SITE__ (the deploy image fills in the URL).
# Usage: curl -fsSL __SITE__/install.sh | sh
# CODEQUEST_INSTALL_DIR overrides the target folder (default ~/.local/bin).
set -eu

site="__SITE__"

os=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$os" in
  darwin | linux) ;;
  *) echo "codequest: unsupported system '$os'. On Windows, see $site/cli/install" >&2; exit 1 ;;
esac

case "$(uname -m)" in
  x86_64 | amd64) arch=amd64 ;;
  arm64 | aarch64) arch=arm64 ;;
  *) echo "codequest: unsupported CPU '$(uname -m)'" >&2; exit 1 ;;
esac

dir="${CODEQUEST_INSTALL_DIR:-$HOME/.local/bin}"
mkdir -p "$dir"
tmp=$(mktemp "${TMPDIR:-/tmp}/codequest.XXXXXX")
trap 'rm -f "$tmp"' EXIT

echo "Downloading codequest for $os/$arch..."
curl -fsSL "$site/cli/codequest-$os-$arch" -o "$tmp"
chmod 0755 "$tmp"
mv "$tmp" "$dir/codequest"
trap - EXIT
echo "Installed $dir/codequest"

case ":$PATH:" in
  *":$dir:"*) ;;
  *)
    echo
    echo "$dir is not on your PATH. Add this line to your ~/.zshrc or ~/.bashrc, then open a new terminal:"
    echo "  export PATH=\"$dir:\$PATH\""
    ;;
esac
echo
echo "Next: codequest login"
