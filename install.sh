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

# Get latest release tag
REPO="kyeo-hub/sb-sync"
GH_PROXY="${GH_PROXY:-}" 

echo "Fetching latest version..."
# 1. Try reading the VERSION file from raw.githubusercontent.com (most proxy-friendly)
LATEST_TAG=$(curl -fsSL "${GH_PROXY}https://raw.githubusercontent.com/${REPO}/main/VERSION" | xargs)

# 2. Fallback to redirect tracking
if [ -z "$LATEST_TAG" ]; then
    LATEST_TAG=$(curl -sI "${GH_PROXY}https://github.com/${REPO}/releases/latest" | grep -i location | awk -F'/' '{print $NF}' | tr -d '\r' | xargs)
fi

# 3. Fallback to API
if [ -z "$LATEST_TAG" ] || [ "$LATEST_TAG" = "latest" ]; then
    LATEST_TAG=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/' | xargs)
fi

if [ -z "$LATEST_TAG" ]; then
    echo "Error: Could not detect the latest version. Please check your network or GH_PROXY."
    exit 1
fi

echo "Latest version detected: $LATEST_TAG"

EXTENSION="tar.gz"
FILENAME="sb-sync-${LATEST_TAG}-${OS}-${ARCH}.${EXTENSION}"
URL="${GH_PROXY}https://github.com/${REPO}/releases/download/${LATEST_TAG}/${FILENAME}"

echo "Downloading from: $URL"
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
