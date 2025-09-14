# Search Feature Documentation

## Overview
The markdown reader (mdrs) now includes a powerful search feature that allows you to quickly find and navigate through text in your markdown documents.

## How to Use

### Starting a Search
- Press `Ctrl+F` or `/` to open the search bar
- The search bar will appear at the bottom of the screen
- Type your search term and press `Enter` to execute the search

### Search Behavior
- **Case-insensitive**: By default, searches are case-insensitive
- **Real-time highlighting**: All matches are highlighted in the document
  - Current match: Highlighted with a bright yellow background
  - Other matches: Highlighted with yellow text
- **Match counter**: The status bar shows the current match number and total matches

### Navigation
Once you've performed a search, you can navigate between matches:
- `n` - Jump to the next match
- `N` - Jump to the previous match
- The view automatically scrolls to center the current match on screen

### Canceling Search
- Press `ESC` or `Ctrl+C` while in the search bar to cancel the search
- Press `Ctrl+G` to cancel and clear the search

## Additional Navigation Keys
The markdown reader also supports these navigation commands:
- `h` / `←` - Scroll left
- `l` / `→` - Scroll right
- `j` / `↓` - Scroll down
- `k` / `↑` - Scroll up
- `g` - Jump to top of document
- `G` - Jump to bottom of document
- `PgUp` - Page up
- `PgDn` / `Space` - Page down
- `q` / `Ctrl+C` - Quit the application

## Search Features

### Smart Highlighting
The search feature uses ANSI color codes to highlight matches:
- The current match is displayed with a bright yellow background for easy visibility
- All other matches are shown with yellow text
- Navigation between matches updates the highlighting in real-time

### Status Bar
When a search is active, a status bar appears showing:
- The current match number (e.g., "Match 2 of 5")
- The search term being used
- "No matches" message if the search term isn't found

### Multiple Matches per Line
The search system can find and highlight multiple occurrences of the search term on the same line, making it easy to spot all instances of your search query.

## Implementation Details

The search feature is implemented in two main components:

1. **search.go**: Contains the `SearchState` struct and all search-related logic
   - Manages search state and match tracking
   - Handles highlighting with ANSI escape codes
   - Provides match navigation methods

2. **mdrs.go**: Integrates the search functionality into the main UI
   - Handles keyboard shortcuts
   - Manages the search input view
   - Updates the display with highlighted content

## Building the Project

To build the project with the new search feature:

```bash
# If you have Go installed:
go build

# Or using make:
make build

# To install:
make install
```

## Example Usage

```bash
# Read a markdown file with search capability
./mdrs README.md

# Then press Ctrl+F to search
# Type your search term and press Enter
# Use n/N to navigate between matches
```

## Future Enhancements

Potential improvements for the search feature:
- Regular expression support
- Case-sensitive search toggle
- Search history
- Replace functionality
- Word boundary matching option
- Search result count in real-time as you type