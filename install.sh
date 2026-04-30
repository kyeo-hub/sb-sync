#!/bin/bash

set -e

# Detect OS and Architecture
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Get latest release tag from GitHub API
REPO="kyeo-hub/sb-sync"
LATEST_TAG=$(curl -s https://api.github.com/repos/$REPO/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_TAG" ]; then
    echo "Failed to fetch latest release tag."
    exit 1
fi

EXTENSION="tar.gz"
FILENAME="sb-sync-${LATEST_TAG}-${OS}-${ARCH}.${EXTENSION}"
URL="https://github.com/${REPO}/releases/download/${LATEST_TAG}/${FILENAME}"

echo "Downloading sb-sync ${LATEST_TAG} for ${OS}/${ARCH}..."
curl -L "$URL" -o "sb-sync.tar.gz"

echo "Extracting..."
tar -xzf "sb-sync.tar.gz"
chmod +x sb-sync

# Install to /usr/local/bin
echo "Installing to /usr/local/bin (may require sudo)..."
if [ -w "/usr/local/bin" ]; then
    mv sb-sync /usr/local/bin/
else
    sudo mv sb-sync /usr/local/bin/
fi

rm sb-sync.tar.gz

echo "Successfully installed sb-sync!"
sb-sync --version
