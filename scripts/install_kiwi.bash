#!/usr/bin/env bash

# e.g. v0.10.2
KIWI_VERSION="$1"

if [ "$(uname)" == "Linux" ]; then
  OS='lnx'
  ARCH='x86_64'
elif [ "$(uname)" == "Darwin" ]; then
  OS='mac'
  # Check if this is Apple Silicon (ARM64) or Intel (x86_64)
  if [ "$(uname -m)" == "arm64" ]; then
    ARCH='arm64'
  else
    ARCH='x86_64'
  fi
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
