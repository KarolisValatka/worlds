.PHONY: build test validate demo

build:
	go build -o worlds ./cmd/worlds

test:
	go test ./...

validate: build
	./worlds validate ./example-vault

demo: validate
	./worlds sync --vault ./example-vault --target pi
	./worlds sync --vault ./example-vault --target claude
	./worlds sync --vault ./example-vault --target cursor
	./worlds sync --vault ./example-vault --target codex
