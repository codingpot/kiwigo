.PHONY: test
test: ModelGenerator/default.dict
	go test ./...

ModelGenerator/default.dict:
	curl -L https://github.com/bab2min/Kiwi/releases/download/v0.10.2/kiwi_model_v0.10.2.tgz --output model.tgz
	tar -xzvf model.tgz
	rm -f model.tgz

.PHONY: clean
clean:
	rm -f model.tgz
	rm -rf ./ModelGenerator

.PHONY: format
format:
	# go install mvdan.cc/gofumpt@latest
	gofumpt -l -w .
