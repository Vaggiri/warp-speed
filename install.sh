#!/usr/bin/env bash
set -e

# warp-speed install script for macOS and Linux

BINARY_NAME="warp-speed"
INSTALL_DIR="$HOME/.local/bin"
BINARY_PATH="$INSTALL_DIR/$BINARY_NAME"

echo -e "\033[1;36m🚀 Starting installation of warp-speed...\033[0m"

# 1. Check if binary exists (assuming they might have compiled it without the .exe extension for unix)
if [ ! -f "./$BINARY_NAME" ]; then
    echo -e "\033[1;33m[INFO] Could not find $BINARY_NAME locally. Attempting to download latest release from GitHub...\033[0m"
    
    # Detect OS & Arch
    OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
    ARCH="$(uname -m)"
    if [ "$ARCH" = "x86_64" ]; then ARCH="amd64"; fi
    if [ "$ARCH" = "aarch64" ] || [ "$ARCH" = "arm64" ]; then ARCH="arm64"; fi
    
    DOWNLOAD_URL="https://github.com/Vaggiri/warp-speed/releases/latest/download/warp-speed-${OS}-${ARCH}"
    
    if curl -sL --fail "$DOWNLOAD_URL" -o "./$BINARY_NAME"; then
        echo -e "\033[1;32m[SUCCESS] Downloaded warp-speed-${OS}-${ARCH}\033[0m"
    else
        echo -e "\033[1;31m[ERROR] Failed to download binary. Make sure you have published the GitHub Release!\033[0m"
        exit 1
    fi
fi

# 2. Create the local bin directory if it doesn't exist
if [ ! -d "$INSTALL_DIR" ]; then
    echo -e "\033[1;34m📁 Creating directory: $INSTALL_DIR\033[0m"
    mkdir -p "$INSTALL_DIR"
fi

# 3. Copy the binary
echo -e "\033[1;34m📦 Copying binary to $INSTALL_DIR...\033[0m"
cp "./$BINARY_NAME" "$BINARY_PATH"
chmod +x "$BINARY_PATH"

# 4. PATH detection and injection
SHELL_RC=""
if [[ "$SHELL" == *"zsh"* ]]; then
    SHELL_RC="$HOME/.zshrc"
elif [[ "$SHELL" == *"bash"* ]]; then
    if [ "$(uname)" == "Darwin" ]; then
        SHELL_RC="$HOME/.bash_profile"
    else
        SHELL_RC="$HOME/.bashrc"
    fi
fi

if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo -e "\033[1;34m🔧 Adding $INSTALL_DIR to your PATH via $SHELL_RC...\033[0m"
    if [ -n "$SHELL_RC" ] && [ -f "$SHELL_RC" ]; then
        echo -e "\n# warp-speed CLI" >> "$SHELL_RC"
        echo "export PATH=\"\$PATH:$INSTALL_DIR\"" >> "$SHELL_RC"
        echo -e "\033[1;33m⚠️ PATH updated in $SHELL_RC. Run 'source $SHELL_RC' or restart your terminal.\033[0m"
    else
        echo -e "\033[1;33m⚠️ Could not automatically update your PATH. Please add the following line to your shell config manually:\033[0m"
        echo -e "export PATH=\"\$PATH:$INSTALL_DIR\""
    fi
else
    echo -e "\033[1;32m✅ $INSTALL_DIR is already in your PATH.\033[0m"
fi

echo -e "\033[1;32m🎉 warp-speed has been successfully installed!\033[0m"
echo -e "Try running 'warp-speed' in your terminal."
