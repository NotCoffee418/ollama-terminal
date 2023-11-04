#!/bin/bash
sudo_check() {
    if [[ $EUID -ne 0 ]]; then
        echo "This script must be run as root (use sudo)." >&2
        exit 1
    fi
}

get_arch() {
    ARCH=$(uname -m)
    case $ARCH in
        x86_64) ARCH="x86_64";;
        armv6*) ARCH="armv6";;
        armv7*) ARCH="armv7";;
        aarch64) ARCH="aarch64";;
        *) echo "Unsupported architecture: $ARCH"; exit 1;;
    esac
    echo $ARCH
}

download_and_install() {
    echo "Fetching latest release..."
    REPO="NotCoffee418/ollama-terminal"
    TAG=$(curl -s https://api.github.com/repos/$REPO/releases/latest | jq -r '.tag_name')
    ARCH=$(get_arch)
    ASSET_NAME="ollama-terminal-linux-${ARCH}"
    ASSET_URL=$(curl -s https://api.github.com/repos/$REPO/releases/latest | jq -r --arg ASSET_NAME "$ASSET_NAME" '.assets[] | select(.name | contains($ASSET_NAME)).browser_download_url')
    echo "Downloading $TAG..."
    curl -L -s --output /usr/local/bin/ollama-terminal $ASSET_URL
    chmod +x /usr/local/bin/ollama-terminal
    echo "Installing..."
    /usr/local/bin/ollama-terminal install
}
sudo_check
download_and_install
echo "Complete"
