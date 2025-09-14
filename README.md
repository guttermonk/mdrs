# mdrs - Markdown Renderer & Search

[![GitHub license](https://img.shields.io/github/license/guttermonk/mdrs.svg?style=for-the-badge)](https://github.com/guttermonk/mdrs/blob/master/LICENSE)
[![GitHub stars](https://img.shields.io/github/stars/guttermonk/mdrs?style=for-the-badge)](https://github.com/guttermonk/mdrs/stargazers)

A standalone Markdown renderer for the terminal with integrated search functionality.

## Features

- 📖 Beautiful Markdown rendering in your terminal
- 🔍 **Full-text search** with highlighting (Ctrl+F)
- ⌨️ Vim-like keybindings with Colemak-DH support
- ❓ **Interactive help popup** - Press `?` to see all keybindings
- 🎨 Syntax highlighting for code blocks
- 🎨 **Customizable colors** via configuration file
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
mdrs --init-config              # Create default config file
```

## Keybindings

Press `?` at any time to display an interactive help popup with all available keybindings.

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
| `?` | Show help popup |
| `q` `Ctrl+C` | Quit |

### Search
| Key | Action |
|-----|--------|
| `Ctrl+F` `/` | Start search |
| `Enter` | Execute search |
| `n` | Next match |
| `N` | Previous match |
| `ESC` | Clear search/Cancel |

Search highlights all matches (current match in bright yellow, others in yellow text) and shows match count in the status bar. Press `ESC` after searching to clear all highlighting and exit search mode.

## Configuration

Customize colors by creating a config file at `~/.config/mdrs/config.json`:

```bash
mdrs --init-config      # Create default config
mdrs --config-path      # Show config location
```

### Color Customization

All colors are specified as hex values (e.g., `#ff0000`). Configurable elements include:
- **Headings**: `heading1` through `heading6`  
- **Text**: `bold`, `italic`, `strikethrough`
- **Code**: `code`, `code_block`, `code_block_bg`
- **Links**: `link`, `link_url`
- **Lists**: `list_marker`, `task_checked`, `task_unchecked`
- **Layout**: `blockquote`, `table_header`, `table_row`, `table_border`
- **Search**: `search_current`, `search_match`

### Pre-built Themes

Copy a theme to your config:
```bash
# Dracula theme
cp themes/dracula.json ~/.config/mdrs/config.json

# Solarized Dark theme  
cp themes/solarized-dark.json ~/.config/mdrs/config.json
```

Example theme structure:
```json
{
  "colors": {
    "heading1": "#bd93f9",
    "bold": "#f8f8f2",
    "code": "#50fa7b",
    "link": "#8be9fd"
  }
}
```

**Note**: Colors are converted to the nearest ANSI 256 color for terminal display.

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

This tool is an offspring of the [mdr](https://github.com/MichaelMure/mdr), which was the offspring of the [git-bug](https://github.com/MichaelMure/git-bug) project.
