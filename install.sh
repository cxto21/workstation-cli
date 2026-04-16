#!/bin/bash
set -e
set -o pipefail

# Workstation CLI Installer
# Usage: curl -fsSL https://raw.githubusercontent.com/reflecterlabs/workstation-cli/main/install.sh | bash

# Override with: WORKSTATION_REPO=owner/repo (legacy MATO_REPO also supported)
REPO="${WORKSTATION_REPO:-${MATO_REPO:-reflecterlabs/workstation-cli}}"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"

# Detect OS and architecture
OS="$(uname -s)"
ARCH="$(uname -m)"

case "$OS" in
    Linux*)
        OS_TYPE="linux"
        ;;
    Darwin*)
        OS_TYPE="macos"
        ;;
    *)
        echo "Unsupported OS: $OS"
        exit 1
        ;;
esac

case "$ARCH" in
    x86_64|amd64)
        ARCH_TYPE="x86_64"
        ;;
    aarch64|arm64)
        ARCH_TYPE="aarch64"
        ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

LEGACY_BINARY_NAME="mato-${OS_TYPE}-${ARCH_TYPE}"
NEW_BINARY_NAME="workstation-cli-${OS_TYPE}-${ARCH_TYPE}"
MUSL_BINARY_NAME="workstation-cli-${OS_TYPE}-${ARCH_TYPE}-musl"
TARGET_BIN="${INSTALL_DIR}/workstation-cli"
LEGACY_BIN="${INSTALL_DIR}/mato"

echo "Installing Workstation CLI for ${OS_TYPE}-${ARCH_TYPE}..."

# Preflight: detect existing workstation-cli in PATH and warn about PATH precedence.
EXISTING_WORKSTATION_CLI="$(command -v workstation-cli 2>/dev/null || true)"
if [ -n "$EXISTING_WORKSTATION_CLI" ] && [ "$EXISTING_WORKSTATION_CLI" != "$TARGET_BIN" ]; then
    echo "⚠️  Existing 'workstation-cli' found at: $EXISTING_WORKSTATION_CLI"
    echo "    This installer will place workstation-cli at: $TARGET_BIN"
    echo "    If PATH prefers another location, that version will continue to run."
    echo "    Tip: run 'which -a workstation-cli' after install to verify command priority."
    echo ""
fi

# Get latest release metadata
RELEASE_JSON=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest")
LATEST_RELEASE=$(printf '%s' "$RELEASE_JSON" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_RELEASE" ]; then
    echo "Failed to get latest release for ${REPO}"
    exit 1
fi

echo "Latest version: $LATEST_RELEASE"

# Resolve download URL from release assets.
# Supports both legacy and new naming schemes.
# On Linux we prefer a musl binary first to avoid glibc compatibility issues.
DOWNLOAD_URL=""
CANDIDATE_ASSETS=""
if [ "$OS_TYPE" = "linux" ]; then
    CANDIDATE_ASSETS="${MUSL_BINARY_NAME} ${NEW_BINARY_NAME} ${LEGACY_BINARY_NAME}"
else
    CANDIDATE_ASSETS="${NEW_BINARY_NAME} ${LEGACY_BINARY_NAME}"
fi

for ASSET_PREFIX in $CANDIDATE_ASSETS; do
    DOWNLOAD_URL=$(
        printf '%s' "$RELEASE_JSON" \
        | grep '"browser_download_url":' \
        | sed -E 's/.*"([^"]+)".*/\1/' \
        | grep -E "/(${ASSET_PREFIX})(-[^/]+)?\\.tar\\.gz$" \
        | head -n1
    )

    if [ -n "$DOWNLOAD_URL" ]; then
        break
    fi
done

if [ -z "$DOWNLOAD_URL" ]; then
    echo "No matching binary asset found for ${NEW_BINARY_NAME}, ${MUSL_BINARY_NAME} or ${LEGACY_BINARY_NAME} in release ${LATEST_RELEASE}"
    echo "Available assets:"
    printf '%s' "$RELEASE_JSON" \
      | grep '"name":' \
      | sed -E 's/.*"name": "([^"]+)".*/  - \1/' \
      || true
    exit 1
fi

echo "Downloading from: $DOWNLOAD_URL"

# Create temp directory
TMP_DIR=$(mktemp -d)
trap "rm -rf $TMP_DIR" EXIT

# Download and extract
cd "$TMP_DIR"
curl -fsSL "$DOWNLOAD_URL" -o mato.tar.gz
tar xzf mato.tar.gz

# Install
mkdir -p "$INSTALL_DIR"
if [ -f workstation-cli ]; then
    mv workstation-cli "$TARGET_BIN"
elif [ -f mato ]; then
    mv mato "$TARGET_BIN"
else
    echo "No executable found in archive (expected 'workstation-cli' or 'mato')."
    exit 1
fi
chmod +x "$TARGET_BIN"

# Backward compatibility alias
ln -sf "$TARGET_BIN" "$LEGACY_BIN"

echo ""
echo "✅ Workstation CLI installed successfully to $TARGET_BIN"
echo ""

# Check if in PATH
if echo "$PATH" | grep -q "$INSTALL_DIR"; then
    echo "You can now run: workstation-cli"
else
    echo "⚠️  Add $INSTALL_DIR to your PATH:"
    echo ""
    # Detect shell config file
    SHELL_NAME=$(basename "$SHELL" 2>/dev/null || echo "bash")
    case "$SHELL_NAME" in
        zsh)
            RC_FILE="$HOME/.zshrc"
            ;;
        bash)
            RC_FILE="$HOME/.bashrc"
            ;;
        fish)
            RC_FILE="$HOME/.config/fish/config.fish"
            ;;
        *)
            RC_FILE="$HOME/.profile"
            ;;
    esac

    echo "    echo 'export PATH=\"\$PATH:$INSTALL_DIR\"' >> $RC_FILE"
    echo "    source $RC_FILE"
    echo ""
    echo "Or run directly: $TARGET_BIN"
fi

ACTIVE_WORKSTATION_CLI="$(command -v workstation-cli 2>/dev/null || true)"
if [ -n "$ACTIVE_WORKSTATION_CLI" ] && [ "$ACTIVE_WORKSTATION_CLI" != "$TARGET_BIN" ]; then
    echo ""
    echo "⚠️  PATH priority notice:"
    echo "    'workstation-cli' currently resolves to: $ACTIVE_WORKSTATION_CLI"
    echo "    Newly installed binary is at: $TARGET_BIN"
    echo "    Run 'which -a workstation-cli' and adjust PATH order if needed."
fi

echo ""
echo "📚 Documentation: https://github.com/${REPO}"
echo "🐛 Report issues: https://github.com/${REPO}/issues"
