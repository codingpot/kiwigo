KIWI_VERSION := "v0.21.0"

.PHONY: test
test: ModelGenerator/default.dict
	go test -count=1 ./...

ModelGenerator/default.dict:
	curl -L https://github.com/bab2min/Kiwi/releases/download/$(KIWI_VERSION)/kiwi_model_$(KIWI_VERSION)_base.tgz --output model.tgz
	tar -xzvf model.tgz
	rm -f model.tgz


.PHONY: install-kiwi
install-kiwi:
	bash scripts/install_kiwi.bash $(KIWI_VERSION)

.PHONY: clean
clean:
	rm -f model.tgz
	rm -rf ./ModelGenerator

.PHONY: format
format:
	# go install mvdan.cc/gofumpt@latest
	gofumpt -l -w .
