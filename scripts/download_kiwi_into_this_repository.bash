#!/usr/bin/env bash

set -e
set -o pipefail

# e.g. v0.10.2
KIWI_VERSION="$1"

# kiwigo/scripts
SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)
PROJECT_DIR=$(dirname "$SCRIPT_DIR")

OS_LIST=(
  'Linux'
  'Darwin'
)

function convert() {
  if [ "$1" == "Linux" ]; then
    echo 'lnx'
  elif [ "$1" == "Darwin" ]; then
    echo 'mac'
  fi
}

function install_library() {
  wget -qO kiwi.tgz "https://github.com/bab2min/Kiwi/releases/download/${KIWI_VERSION}/kiwi_${KIWI_OS_NAME}_x86_64_${KIWI_VERSION}.tgz" &&
    tar xzvf kiwi.tgz &&
    mkdir -p "$LIBRARY_PATH" &&
    mv build/libkiwi* "$LIBRARY_PATH/" &&
    rm -rf kiwi.tgz build

  return $?
}

function install_header() {
  wget -qO source.tgz "https://github.com/bab2min/Kiwi/archive/refs/tags/${KIWI_VERSION}.tar.gz" &&
    tar xzvf source.tgz &&
    cp -r "Kiwi-${KIWI_VERSION/v/}/include/kiwi" "$PROJECT_DIR/include/" &&
    rm -rf source.tgz Kiwi-*

  return $?
}

function install_library_for_windows() {
  LIBRARY_PATH="${PROJECT_DIR}/libs/Windows_x86_64"

  wget -qO kiwi.zip "https://github.com/bab2min/Kiwi/releases/download/${KIWI_VERSION}/kiwi_win_x64_${KIWI_VERSION}.zip" &&
    mkdir -p "$LIBRARY_PATH" &&
    unzip kiwi.zip -d "$LIBRARY_PATH" &&
    rm -rf kiwi.zip build

  return $?
}

function main() {

  echo "Installing Kiwi version ${KIWI_VERSION:?}"

  for OS in "${OS_LIST[@]}"; do
    echo "Downloading library for ${OS:?}"

    KIWI_OS_NAME=$(convert "$OS")
    LIBRARY_PATH="${PROJECT_DIR}/libs/${OS}_x86_64"

    install_library
  done

  install_library_for_windows

  install_header
}

main
