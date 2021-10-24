#!/usr/bin/env bash

wget -O kiwi.tgz https://github.com/bab2min/Kiwi/releases/download/v0.10.1/kiwi_lnx_x86_64_v0.10.1.tgz &&
  tar xzvf kiwi.tgz &&
  sudo mv build/libkiwi* /usr/local/lib/ &&
  sudo ldconfig &&
  rm -rf kiwi.tgz build &&
  wget -O source.tgz https://github.com/bab2min/Kiwi/archive/refs/tags/v0.10.1.tar.gz &&
  tar xzvf source.tgz &&
  mkdir kiwi &&
  sudo cp -r Kiwi-0.10.1/include/kiwi /usr/local/include/ &&
  rm -rf source.tgz Kiwi-*
