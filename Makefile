KIWI_VERSION := v0.23.2

.PHONY: test
test: base/default.dict
	go test ./...

base/default.dict:
	curl -L https://github.com/bab2min/Kiwi/releases/download/$(KIWI_VERSION)/kiwi_model_$(KIWI_VERSION)_base.tgz --output model.tgz
	tar --no-same-owner -xzvf model.tgz
	mv models/cong/base ./base
	rm -rf models model.tgz


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

.PHONY: sync-postypes
sync-postypes:
	@echo "Extracting POS tags from Kiwi $(KIWI_VERSION)..."
	source .venv/bin/activate && python scripts/extract_postags.py $(KIWI_VERSION)
	@echo "Generated postype_generated.go"
	@echo "Comparing with current postype.go..."
	@diff -u postype.go postype_generated.go || true
	@echo "To apply changes, run: mv postype_generated.go postype.go"

.PHONY: check-postypes
check-postypes:
	@source .venv/bin/activate && python scripts/extract_postags.py $(KIWI_VERSION)
	@if diff -q postype.go postype_generated.go > /dev/null 2>&1; then \
		echo "POS types are in sync with Kiwi $(KIWI_VERSION)"; \
	else \
		echo "POS types are OUT OF SYNC with Kiwi $(KIWI_VERSION)"; \
		echo "Run 'make sync-postypes' to update"; \
		exit 1; \
	fi

