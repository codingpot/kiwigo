#!/usr/bin/env bash

# e.g. v0.10.2
KIWI_VERSION="$1"

if [ "$(uname)" == "Linux" ]; then
  OS='lnx'
elif [ "$(uname)" == "Darwin" ]; then
  OS='mac'
elif [ "$(uname)" == "Windows" ]; then
  OS='win'
fi

if [ "$(uname -m)" == "arm64" ]; then
  ARCH="arm64"
else
  ARCH="x86_64"
fi

echo "set OS env to ${OS:?}"
echo "installing Kiwi version ${KIWI_VERSION:?}"

wget -O kiwi.tgz "https://github.com/bab2min/Kiwi/releases/download/${KIWI_VERSION}/kiwi_${OS}_${ARCH}_${KIWI_VERSION}.tgz" &&
  sudo mkdir -p /usr/local/kiwi &&
  sudo tar xzvf kiwi.tgz &&
  sudo cp lib/libkiwi* /usr/local/lib &&
  sudo cp -rf include/kiwi /usr/local/include &&
  [[ "$(uname)" == "Linux" ]] && sudo ldconfig || echo 'skip' &&
  rm -rf kiwi.tgz bin lib include
