.PHONY: run
run:
	go run cmd/tagalign/main.go ./...

.PHONY: lint
lint:
	golangci-lint run ./...