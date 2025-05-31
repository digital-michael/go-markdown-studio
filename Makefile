
.PHONY: build rebuild security test clean

configure:
	sudo dnf install libXxf86vm-devel -y
	go mod tidy
	go install github.com/securego/gosec/v2/cmd/gosec@latest

build:
	go fmt ./...
	go build -o build/markdown-studio ./cmd/markdown-studio

rebuild: clean
	go clean -modcache
	go mod tidy
	go build -v -a -o build/markdown-studio ./cmd/markdown-studio

security:
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./... ; \
	else \
		echo "gosec not installed. Skipping security check." ; \
	fi

test: security
	go test ./...

clean:
	rm -rf build
	mkdir -p build



