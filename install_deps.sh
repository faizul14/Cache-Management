#!/bin/bash

set -e

echo "==> Memeriksa dan menginstall dependency..."

install_debian() {
    echo "Deteksi Debian/Ubuntu"
    sudo apt update
    sudo apt install -y figlet fzf
}

install_fedora() {
    echo "Deteksi Fedora"
    sudo dnf install -y figlet fzf
}

install_arch() {
    echo "Deteksi Arch Linux"
    sudo pacman -Sy --noconfirm figlet fzf
}

install_mac() {
    echo "Deteksi MacOS"
    if ! command -v brew &>/dev/null; then
        echo "Homebrew belum terinstall. Silakan install Homebrew dulu: https://brew.sh"
        exit 1
    fi
    brew install figlet fzf
}

if [ "$(uname)" = "Darwin" ]; then
    install_mac
elif [ -f /etc/debian_version ]; then
    install_debian
elif [ -f /etc/fedora-release ]; then
    install_fedora
elif [ -f /etc/arch-release ]; then
    install_arch
else
    echo "Distro Linux tidak dikenali, silakan install figlet dan fzf secara manual."
    exit 1
fi

echo "==> Dependency sudah terinstall."
