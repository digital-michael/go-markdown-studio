
# Go M5arkdown Studio

A simple Markdown editor with a live preview, written in Go using the Fyne toolkit.

## Features
- Edit and preview Markdown files
- Tree view of files from configured directories
- Light, dark, and system themes
- Configuration screen for directories and themes

## Preparing for first build
```sh
sudo dnf install libXxf86vm-devel
go mod tidy
```

## Optonal Security
```sh
go install github.com/securego/gosec/v2/cmd/gosec@latest

```

## Building

```sh
make build
```

## Running

```sh
./build/markdown-studio
```

## Testing

```sh
make test
```
