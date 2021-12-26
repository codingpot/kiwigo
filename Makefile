KIWI_VERSION := "v0.10.3"

.PHONY: test
test: ModelGenerator/default.dict
	go test ./...

ModelGenerator/default.dict:
	curl -L https://github.com/bab2min/Kiwi/releases/download/$(KIWI_VERSION)/kiwi_model_$(KIWI_VERSION).tgz --output model.tgz
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
