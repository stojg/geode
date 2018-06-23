.PHONY: check
check:
	@echo "gometalinter"
	@! gometalinter --vendor ./...
	@echo "gofmt (simplify)"
	@! gofmt -s -d -l . 2>&1 | grep -vE '^\.git/'
	@echo "goimports"
	@! goimports -l . | grep -vF 'No Exceptions'
