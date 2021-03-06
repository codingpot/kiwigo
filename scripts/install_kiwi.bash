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

echo "set OS env to ${OS:?}"
echo "installing Kiwi version ${KIWI_VERSION:?}"

wget -O kiwi.tgz "https://github.com/bab2min/Kiwi/releases/download/${KIWI_VERSION}/kiwi_${OS}_x86_64_${KIWI_VERSION}.tgz" &&
  tar xzvf kiwi.tgz &&
  sudo mv build/libkiwi* /usr/local/lib/ &&
  [[ "$(uname)" == "Linux" ]] && sudo ldconfig || echo 'skip' &&
  rm -rf kiwi.tgz build &&
  wget -O source.tgz https://github.com/bab2min/Kiwi/archive/refs/tags/${KIWI_VERSION}.tar.gz &&
  tar xzvf source.tgz &&
  sudo cp -r Kiwi-${KIWI_VERSION/v/}/include/kiwi /usr/local/include/ &&
  rm -rf source.tgz Kiwi-*
