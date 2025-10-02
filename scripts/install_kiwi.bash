#!/usr/bin/env bash

# e.g. v0.10.2
KIWI_VERSION="$1"

if [ "$(uname)" == "Linux" ]; then
  OS='lnx'
  ARCH='x86_64'
elif [ "$(uname)" == "Darwin" ]; then
  OS='mac'
  # For v0.10.3, only x86_64 build is available for macOS
  # ARM64 Macs will use x86_64 build with Rosetta translation
  ARCH='x86_64'
elif [ "$(uname)" == "Windows" ]; then
  OS='win'
  ARCH='x86_64'
fi

echo "set OS env to ${OS:?}"
echo "set ARCH env to ${ARCH:?}"
echo "installing Kiwi version ${KIWI_VERSION:?}"

wget -O kiwi.tgz "https://github.com/bab2min/Kiwi/releases/download/${KIWI_VERSION}/kiwi_${OS}_${ARCH}_${KIWI_VERSION}.tgz" &&
  tar xzvf kiwi.tgz &&
  sudo mkdir -p /usr/local/lib /usr/local/include &&
  sudo mv build/libkiwi* /usr/local/lib/ &&
  [[ "$(uname)" == "Linux" ]] && sudo ldconfig || echo 'skip' &&
  rm -rf kiwi.tgz build &&
  wget -O source.tgz https://github.com/bab2min/Kiwi/archive/refs/tags/${KIWI_VERSION}.tar.gz &&
  tar xzvf source.tgz &&
  sudo cp -r Kiwi-${KIWI_VERSION/v/}/include/kiwi /usr/local/include/ &&
  rm -rf source.tgz Kiwi-*
