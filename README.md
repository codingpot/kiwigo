# kiwigo

[![CI](https://github.com/codingpot/kiwigo/actions/workflows/ci.yaml/badge.svg)](https://github.com/codingpot/kiwigo/actions/workflows/ci.yaml)
[![Go Reference](https://pkg.go.dev/badge/github.com/codingpot/kiwigo.svg)](https://pkg.go.dev/github.com/codingpot/kiwigo)

https://github.com/bab2min/Kiwi for go

## Dependencies

- Install Kiwi

check out how to install Kiwi [here](https://github.com/bab2min/Kiwi#%EC%84%A4%EC%B9%98)

### Example script (install using the helper script)

```bash
make install-kiwi
```

### Example script (install from source)

```bash
git clone git@github.com:bab2min/Kiwi.git
cd Kiwi
git submodule sync
git submodule update --init --recursive
mkdir build
cd build
cmake -DCMAKE_BUILD_TYPE=Release ../
make
make install
ldconfig
```
