KIWI_VERSION := v0.21.0

.PHONY: test
test: base/default.dict
	go test ./...

base/default.dict:
	curl -L https://github.com/bab2min/Kiwi/releases/download/$(KIWI_VERSION)/kiwi_model_$(KIWI_VERSION)_base.tgz --output model.tgz
	tar --no-same-owner -xzvf model.tgz
	rm -f model.tgz


.PHONY: install-kiwi
install-kiwi:
	bash scripts/install_kiwi.bash $(KIWI_VERSION)

.PHONY: clean
clean:
	rm -f model.tgz
	rm -rf ./base

.PHONY: format
format:
	# go install mvdan.cc/gofumpt@latest
	gofumpt -l -w .

