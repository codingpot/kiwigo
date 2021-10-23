.PHONY: test
test: libs/kiwi/build/libkiwi.so
	LD_LIBRARY_PATH=libs/kiwi/build go test ./...

libs/kiwi/build/libkiwi.so:
	git submodule update --recursive --init
	mkdir -p libs/kiwi/build && cd libs/kiwi/build && cmake .. && make

.PHONY: clean
clean:
	rm -rf libs/kiwi/build
