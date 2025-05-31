
.PHONY: build the markdown-studio binary
build:
	go build -o build/markdown-studio ./cmd/markdown-studio


.PHONY: rebuild everything including downloaded dependencies and be verbose
rebuild: clean
	go build -v -a -o build/markdown-studio ./cmd/markdown-studio


.PHONY: security check for security issues
security:
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...  \
	else \
		echo "gosec not installed. Skipping security check." ; \
	fi    


.PHONY: build rebuild
test: security
	go test ./...
	

.PHONY: build rebuild test clean
clean:
	rm -rf build
	mkdir -p build



