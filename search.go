package main

import (
	"fmt"
	"strings"

	"github.com/awesome-gocui/gocui"
)

// SearchState manages the state of the search feature
type SearchState struct {
	active        bool
	term          string
	matches       []SearchMatch
	currentIndex  int
	caseSensitive bool
}

// SearchMatch represents a single search match
type SearchMatch struct {
	lineNumber int
	column     int
	text       string
}

// NewSearchState creates a new search state
func NewSearchState() *SearchState {
	return &SearchState{
		active:        false,
		term:          "",
		matches:       []SearchMatch{},
		currentIndex:  -1,
		caseSensitive: false,
	}
}

// Clear resets the search state
func (s *SearchState) Clear() {
	s.active = false
	s.term = ""
	s.matches = []SearchMatch{}
	s.currentIndex = -1
}

// SetTerm sets the search term and performs the search
func (s *SearchState) SetTerm(term string, content string) {
	s.term = term
	s.findAllMatches(content)
	if len(s.matches) > 0 {
		s.currentIndex = 0
	}
}

// findAllMatches finds all matches in the content
func (s *SearchState) findAllMatches(content string) {
	s.matches = []SearchMatch{}
	if s.term == "" {
		return
	}

	lines := strings.Split(content, "\n")
	searchTerm := s.term
	
	if !s.caseSensitive {
		searchTerm = strings.ToLower(searchTerm)
	}

	for lineNum, line := range lines {
		searchLine := line
		if !s.caseSensitive {
			searchLine = strings.ToLower(line)
		}

		// Find all occurrences in this line
		index := 0
		for {
			pos := strings.Index(searchLine[index:], searchTerm)
			if pos == -1 {
				break
			}
			
			actualPos := index + pos
			s.matches = append(s.matches, SearchMatch{
				lineNumber: lineNum,
				column:     actualPos,
				text:       line[actualPos : actualPos+len(s.term)],
			})
			
			index = actualPos + len(searchTerm)
		}
	}
}

// NextMatch moves to the next match
func (s *SearchState) NextMatch() (SearchMatch, bool) {
	if len(s.matches) == 0 {
		return SearchMatch{}, false
	}

	s.currentIndex = (s.currentIndex + 1) % len(s.matches)
	return s.matches[s.currentIndex], true
}

// PrevMatch moves to the previous match
func (s *SearchState) PrevMatch() (SearchMatch, bool) {
	if len(s.matches) == 0 {
		return SearchMatch{}, false
	}

	s.currentIndex--
	if s.currentIndex < 0 {
		s.currentIndex = len(s.matches) - 1
	}
	return s.matches[s.currentIndex], true
}

// GetCurrentMatch returns the current match
func (s *SearchState) GetCurrentMatch() (SearchMatch, bool) {
	if s.currentIndex < 0 || s.currentIndex >= len(s.matches) {
		return SearchMatch{}, false
	}
	return s.matches[s.currentIndex], true
}

// GetMatchCount returns the total number of matches
func (s *SearchState) GetMatchCount() int {
	return len(s.matches)
}

// GetStatusText returns status text for the search
func (s *SearchState) GetStatusText() string {
	if s.term == "" {
		return ""
	}
	
	if len(s.matches) == 0 {
		return fmt.Sprintf("No matches for: %s", s.term)
	}
	
	return fmt.Sprintf("Match %d of %d: %s", s.currentIndex+1, len(s.matches), s.term)
}

// HighlightContent highlights search matches in the content
func (s *SearchState) HighlightContent(content []byte) []byte {
	if s.term == "" || len(s.matches) == 0 {
		return content
	}

	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")
	
	// Create a map of line numbers to matches for efficient lookup
	lineMatches := make(map[int][]SearchMatch)
	for _, match := range s.matches {
		lineMatches[match.lineNumber] = append(lineMatches[match.lineNumber], match)
	}
	
	// Process each line that has matches
	for lineNum, matches := range lineMatches {
		if lineNum >= len(lines) {
			continue
		}
		
		line := lines[lineNum]
		var newLine strings.Builder
		lastEnd := 0
		
		// Sort matches by column position for this line
		for i, match := range matches {
			isCurrentMatch := false
			// Check if this match is the current one
			for j, m := range s.matches {
				if j == s.currentIndex && m.lineNumber == match.lineNumber && m.column == match.column {
					isCurrentMatch = true
					break
				}
			}
			
			// Add text before the match
			if match.column > lastEnd {
				newLine.WriteString(line[lastEnd:match.column])
			}
			
			// Add highlighted match
			if isCurrentMatch {
				// Current match - bright yellow background with black text
				newLine.WriteString("\033[43;30m")
				newLine.WriteString(line[match.column : match.column+len(s.term)])
				newLine.WriteString("\033[0m")
			} else {
				// Other matches - yellow text
				newLine.WriteString("\033[33m")
				newLine.WriteString(line[match.column : match.column+len(s.term)])
				newLine.WriteString("\033[0m")
			}
			
			lastEnd = match.column + len(s.term)
			
			// If this is the last match and there's text after it
			if i == len(matches)-1 && lastEnd < len(line) {
				newLine.WriteString(line[lastEnd:])
			}
		}
		
		// If there were no matches starting from position 0, add the remaining text
		if len(matches) > 0 && lastEnd == 0 {
			newLine.WriteString(line)
		}
		
		lines[lineNum] = newLine.String()
	}
	
	return []byte(strings.Join(lines, "\n"))
}

// HandleSearchInput handles input in the search view
func (s *SearchState) HandleSearchInput(g *gocui.Gui, v *gocui.View) error {
	// The view is already editable, so text input is handled automatically
	return nil
}

// ToggleCaseSensitive toggles case-sensitive search
func (s *SearchState) ToggleCaseSensitive() {
	s.caseSensitive = !s.caseSensitive
}