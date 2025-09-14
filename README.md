# mdrs : MarkDown Renderer & Search

[![Build Status](https://travis-ci.org/MichaelMure/mdrs.svg?branch=master)](https://travis-ci.org/MichaelMure/mdrs)
[![Go Report Card](https://goreportcard.com/badge/github.com/guttermonk/mdrs)](https://goreportcard.com/report/github.com/guttermonk/mdrs)
[![GitHub license](https://img.shields.io/github/license/guttermonk/mdrs.svg)](https://github.com/guttermonk/mdrs/blob/master/LICENSE)

`mdrs` is a standalone Markdown renderer for the terminal with integrated search functionality.

Note: Markdown being originally designed to render as HTML, rendering in a terminal is occasionally challenging and some adaptation had to be made. 

## Features

- 📖 Render Markdown files beautifully in your terminal
- 🔍 **Search functionality** - Press `Ctrl+F` to search through your documents
- ⌨️ Vim-like keybindings for navigation
- 🎨 Syntax highlighting for code blocks
- 📊 Table rendering support
- 🚀 Fast and lightweight

## Examples

![rendered markdown](examples/markdown.png)
![rendered table](examples/table.png)
![rendered code](examples/code.png)

## Installation

You can grab a [pre-compiled binary](https://github.com/guttermonk/mdrs/releases/latest).

## Keybindings

### Navigation

| Action | Key |
|--------|-----|
| Quit | <kbd>ctrl+C</kbd>, <kbd>Q</kbd>|
| Up | <kbd>↑</kbd>, <kbd>K</kbd>, <kbd>ctrl+P</kbd>|
| Down | <kbd>↓</kbd>, <kbd>J</kbd>, <kbd>ctrl+N</kbd> |
| Left | <kbd>←</kbd>, <kbd>H</kbd> |
| Right | <kbd>→</kbd>, <kbd>L</kbd> |
| Page Up | <kbd>⇞</kbd> |
| Page Down | <kbd>⇟</kbd>, <kbd>space</kbd> |
| Go to Top | <kbd>G</kbd> |
| Go to Bottom | <kbd>shift+G</kbd> |

### Search

| Action | Key |
|--------|-----|
| Start Search | <kbd>ctrl+F</kbd>, <kbd>/</kbd> |
| Next Match | <kbd>N</kbd> |
| Previous Match | <kbd>shift+N</kbd> |
| Cancel Search | <kbd>ESC</kbd>, <kbd>ctrl+C</kbd> |

## Search Feature

The search feature allows you to quickly find text within your markdown documents:

1. Press `Ctrl+F` or `/` to open the search bar
2. Type your search term and press `Enter`
3. All matches will be highlighted:
   - Current match: Bright yellow background
   - Other matches: Yellow text
4. Navigate between matches with `n` (next) and `N` (previous)
5. The status bar shows the current match number and total matches

## Origin

This tool is an offspring of the [git-bug](https://github.com/MichaelMure/git-bug) project.

## Contribute

PRs accepted.

## License

MIT