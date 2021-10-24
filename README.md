# kiwigo

https://github.com/bab2min/Kiwi for go

## Dependencies

- Install Kiwi

check out how to install Kiwi [here](https://github.com/bab2min/Kiwi#%EC%84%A4%EC%B9%98)

### example script

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
