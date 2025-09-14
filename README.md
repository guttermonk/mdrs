# mdrs - Markdown Renderer & Search

[![GitHub license](https://img.shields.io/github/license/guttermonk/mdrs.svg?style=for-the-badge)](https://github.com/guttermonk/mdrs/blob/master/LICENSE)
[![GitHub stars](https://img.shields.io/github/stars/guttermonk/mdrs?style=for-the-badge)](https://github.com/guttermonk/mdrs/stargazers)

A standalone Markdown renderer for the terminal with integrated search functionality.

## Features

- 📖 Beautiful Markdown rendering in your terminal
- 🔍 **Full-text search** with highlighting (Ctrl+F)
- ⌨️ Vim-like keybindings with Colemak-DH support
- 🎨 Syntax highlighting for code blocks
- 📊 Table rendering support
- ❄️ Native NixOS support with flakes

## Installation

### Binary Release
Download a [pre-compiled binary](https://github.com/guttermonk/mdrs/releases/latest) for your platform.

### NixOS / Nix
```bash
# Run directly
nix run github:guttermonk/mdrs -- README.md

# Install to profile
nix profile install github:guttermonk/mdrs

# Build locally
git clone https://github.com/guttermonk/mdrs
cd mdrs
nix build
./result/bin/mdrs README.md
```

### From Source
```bash
git clone https://github.com/guttermonk/mdrs
cd mdrs
go build
./mdrs README.md
```

## Usage

```bash
mdrs README.md                  # Render a markdown file
mdrs < file.md                  # Read from stdin
curl example.com/file.md | mdrs # Pipe from network
```

## Keybindings

### Navigation
| Key | Action |
|-----|--------|
| `↑` `k` `i` | Scroll up |
| `↓` `j` `e` | Scroll down |
| `←` `h` | Scroll left |
| `→` `l` `o` | Scroll right |
| `PgUp` | Page up |
| `PgDn` `Space` | Page down |
| `g` | Go to top |
| `G` | Go to bottom |
| `q` `Ctrl+C` | Quit |

### Search
| Key | Action |
|-----|--------|
| `Ctrl+F` `/` | Start search |
| `Enter` | Execute search |
| `n` | Next match |
| `N` | Previous match |
| `ESC` | Cancel search |

Search highlights all matches (current match in bright yellow, others in yellow text) and shows match count in the status bar.

## Development

### Nix Development Shell
```bash
nix develop  # Or use direnv with the included .envrc
go build
go test ./...
```

### Traditional Development
```bash
go mod download
make build
```

The development environment includes Go, gopls, golangci-lint, and other useful tools.

## Examples

![rendered markdown](examples/markdown.png)
![rendered table](examples/table.png)
![rendered code](examples/code.png)

## Contributing

PRs accepted. When contributing:
1. Use the development shell for consistent tooling
2. Run tests with `go test ./...`
3. Update vendor hash in `flake.nix` if dependencies change

## License

MIT

## Origin

This tool is an offspring of the [git-bug](https://github.com/MichaelMure/git-bug) project.
