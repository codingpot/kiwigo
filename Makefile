.PHONY: test
test: get-model
	go test ./...

get-model:
	if [ ! -f ./ModelGenerator/default.dict ]; then \
		mkdir -p ./libs; \
		curl -L https://github.com/bab2min/Kiwi/releases/download/v0.10.1/kiwi_model_v0.10.1.tgz --output ./libs/model.tgz; \
		tar -xzvf ./libs/model.tgz; \
		rm -rf ./libs; \
	fi

.PHONY: clean
clean:
	rm -rf ./ModelGenerator