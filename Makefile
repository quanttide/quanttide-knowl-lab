.PHONY: build test lint clean

build:
	go build ./...

test:
	go test ./... -v

lint:
	go vet ./...
	@which staticcheck >/dev/null 2>&1 && staticcheck ./... || echo "staticcheck not installed, skipping"

clean:
	rm -rf parsed/ triples.jsonl errors.jsonl generated/ sample.md
	go clean -cache

extractor:
	go run ./cmd/extractor -input ./docs -output ./parsed

genre-detect:
	go run ./cmd/genre detect -input _examples/bylaw

genre-validate:
	go run ./cmd/genre validate -input _examples/bylaw

genre-transcode:
	go run ./cmd/genre transcode -input _examples/bylaw -output _examples/bylaw/archive
