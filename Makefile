.PHONY: test
test: ModelGenerator/default.dict
	go test ./...

ModelGenerator/default.dict:
	mkdir -p ./libs; \
	curl -L https://github.com/bab2min/Kiwi/releases/download/v0.10.1/kiwi_model_v0.10.1.tgz --output ./libs/model.tgz; \
	tar -xzvf ./libs/model.tgz; \
	rm -rf ./libs; \

.PHONY: clean
clean:
	rm -rf ./ModelGenerator