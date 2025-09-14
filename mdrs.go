package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/MichaelMure/go-term-markdown"
	"github.com/awesome-gocui/gocui"
	"github.com/mattn/go-isatty"
	"github.com/pkg/errors"
)

const padding = 4

func main() {
	if len(os.Args) >= 2 && (os.Args[1] == "version" || os.Args[1] == "--version") {
		printVersion()
		return
	}

	var content []byte

	switch len(os.Args) {
	case 1:
		if isatty.IsTerminal(os.Stdin.Fd()) {
			exitError(fmt.Errorf("usage: %s <file.md>", os.Args[0]))
		}
		data, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			exitError(errors.Wrap(err, "error while reading STDIN"))
		}
		content = data
	case 2:
		data, err := ioutil.ReadFile(os.Args[1])
		if err != nil {
			exitError(errors.Wrap(err, "error while reading file"))
		}
		err = os.Chdir(path.Dir(os.Args[1]))
		if err != nil {
			exitError(err)
		}
		content = data

	default:
		exitError(fmt.Errorf("only one file is supported"))
	}

	g, err := gocui.NewGui(gocui.OutputNormal, false)
	if err != nil {
		exitError(errors.Wrap(err, "error starting the interactive UI"))
	}
	defer g.Close()

	ui, err := newUi(g)
	if err != nil {
		exitError(err)
	}

	ui.setContent(content)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		exitError(err)
	}
}

func exitError(err error) {
	_, _ = fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

const renderView = "render"
const searchView = "search"
const statusView = "status"

type ui struct {
	keybindings []keybinding

	raw string
	// current width of the view
	width   int
	XOffset int
	YOffset int

	// number of lines in the rendered markdown
	lines int

	// search state
	search          *SearchState
	renderedContent []byte
	searchActive    bool
}

func newUi(g *gocui.Gui) (*ui, error) {
	result := &ui{
		width:  -1,
		search: NewSearchState(),
	}

	g.SetManagerFunc(result.layout)
	g.Cursor = true
	g.InputEsc = true

	result.keybindings = []keybinding{
		// General keybindings
		{"", gocui.KeyCtrlC, gocui.ModNone, result.quit},
		{renderView, 'q', gocui.ModNone, result.quit},
		
		// Navigation keybindings
		{renderView, 'k', gocui.ModNone, result.up},
		{renderView, gocui.KeyCtrlP, gocui.ModNone, result.up},
		{renderView, gocui.KeyArrowUp, gocui.ModNone, result.up},
		{renderView, 'j', gocui.ModNone, result.down},
		{renderView, gocui.KeyCtrlN, gocui.ModNone, result.down},
		{renderView, gocui.KeyArrowDown, gocui.ModNone, result.down},
		{renderView, 'h', gocui.ModNone, result.left},
		{renderView, gocui.KeyArrowLeft, gocui.ModNone, result.left},
		{renderView, 'l', gocui.ModNone, result.right},
		{renderView, gocui.KeyArrowRight, gocui.ModNone, result.right},
		{renderView, gocui.KeyPgup, gocui.ModNone, result.pageUp},
		{renderView, gocui.KeyPgdn, gocui.ModNone, result.pageDown},
		{renderView, gocui.KeySpace, gocui.ModNone, result.pageDown},
		{renderView, 'g', gocui.ModNone, result.goToTop},
		{renderView, 'G', gocui.ModNone, result.goToBottom},
		
		// Search keybindings
		{renderView, gocui.KeyCtrlF, gocui.ModNone, result.startSearch},
		{renderView, '/', gocui.ModNone, result.startSearch},
		{renderView, 'n', gocui.ModNone, result.nextMatch},
		{renderView, 'N', gocui.ModNone, result.prevMatch},
		{searchView, gocui.KeyEnter, gocui.ModNone, result.executeSearch},
		{searchView, gocui.KeyEsc, gocui.ModNone, result.cancelSearch},
		{searchView, gocui.KeyCtrlC, gocui.ModNone, result.cancelSearch},
		{searchView, gocui.KeyCtrlG, gocui.ModNone, result.cancelSearch},
	}

	for _, kb := range result.keybindings {
		err := kb.Register(g)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (ui *ui) setContent(content []byte) {
	ui.raw = string(content)
	ui.width = -1
}

func (ui *ui) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	// Status bar at the bottom
	statusY := maxY - 1
	if ui.searchActive {
		statusY = maxY - 3
	}

	// Main render view
	v, err := g.SetView(renderView, ui.XOffset, -ui.YOffset, maxX, statusY, 0)
	if err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}

		v.Frame = false
		v.Wrap = false
	}

	if len(ui.raw) > 0 && ui.width != maxX {
		ui.width = maxX
		v.Clear()
		ui.renderedContent = ui.render(g)
		
		// Apply search highlighting if search is active
		if ui.search.term != "" {
			_, _ = v.Write(ui.search.HighlightContent(ui.renderedContent))
		} else {
			_, _ = v.Write(ui.renderedContent)
		}
	} else if ui.search.term != "" {
		// Update highlighting even if width hasn't changed
		v.Clear()
		_, _ = v.Write(ui.search.HighlightContent(ui.renderedContent))
	}

	// Status bar
	if ui.search.term != "" || ui.searchActive {
		sv, err := g.SetView(statusView, 0, statusY, maxX-1, statusY+2, 0)
		if err != nil {
			if !gocui.IsUnknownView(err) {
				return err
			}
			sv.Frame = false
			sv.Wrap = false
		}
		sv.Clear()
		
		statusText := ui.search.GetStatusText()
		if statusText != "" {
			fmt.Fprintf(sv, " %s", statusText)
		}
	} else {
		g.DeleteView(statusView)
	}

	// Search input view
	if ui.searchActive {
		searchY := maxY - 2
		sv, err := g.SetView(searchView, 0, searchY, maxX-1, searchY+2, 0)
		if err != nil {
			if !gocui.IsUnknownView(err) {
				return err
			}
			sv.Frame = true
			sv.Title = " Search (Enter to search, ESC to cancel) "
			sv.Editable = true
			sv.Wrap = false
		}

		_, err = g.SetCurrentView(searchView)
		if err != nil {
			return err
		}
		
		// Set cursor to end of search term
		sv.SetCursor(len(sv.Buffer()), 0)
	} else {
		g.DeleteView(searchView)
		_, err = g.SetCurrentView(renderView)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ui *ui) render(g *gocui.Gui) []byte {
	maxX, _ := g.Size()

	opts := []markdown.Options{
		// needed when going through gocui
		markdown.WithImageDithering(markdown.DitheringWithBlocks),
	}

	rendered := markdown.Render(ui.raw, maxX-1-padding, padding, opts...)
	ui.lines = 0
	for _, b := range rendered {
		if b == '\n' {
			ui.lines++
		}
	}
	return rendered
}

func (ui *ui) startSearch(g *gocui.Gui, v *gocui.View) error {
	ui.searchActive = true
	return nil
}

func (ui *ui) executeSearch(g *gocui.Gui, v *gocui.View) error {
	searchText := strings.TrimSpace(v.Buffer())
	if searchText == "" {
		return ui.cancelSearch(g, v)
	}

	// Perform the search
	ui.search.SetTerm(searchText, string(ui.renderedContent))
	
	// If we found matches, scroll to the first one
	if match, ok := ui.search.GetCurrentMatch(); ok {
		ui.scrollToLine(g, match.lineNumber)
	}

	ui.searchActive = false
	v.Clear()
	v.SetCursor(0, 0)
	
	// Force a re-render to show highlights
	ui.width = -1
	
	return nil
}

func (ui *ui) cancelSearch(g *gocui.Gui, v *gocui.View) error {
	ui.searchActive = false
	ui.search.Clear()
	v.Clear()
	v.SetCursor(0, 0)
	
	// Force a re-render to clear highlights
	ui.width = -1
	
	return nil
}

func (ui *ui) nextMatch(g *gocui.Gui, v *gocui.View) error {
	if ui.search.term == "" {
		return nil
	}
	
	if match, ok := ui.search.NextMatch(); ok {
		ui.scrollToLine(g, match.lineNumber)
		// Force re-render to update highlight
		ui.width = -1
	}
	
	return nil
}

func (ui *ui) prevMatch(g *gocui.Gui, v *gocui.View) error {
	if ui.search.term == "" {
		return nil
	}
	
	if match, ok := ui.search.PrevMatch(); ok {
		ui.scrollToLine(g, match.lineNumber)
		// Force re-render to update highlight
		ui.width = -1
	}
	
	return nil
}

func (ui *ui) scrollToLine(g *gocui.Gui, lineNumber int) {
	_, maxY := g.Size()
	
	// Try to center the match on screen
	targetOffset := lineNumber - maxY/2
	
	// Clamp to valid range
	ui.YOffset = max(0, min(targetOffset, ui.lines-maxY+1))
}

func (ui *ui) quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func (ui *ui) up(g *gocui.Gui, v *gocui.View) error {
	ui.YOffset -= 1
	ui.YOffset = max(ui.YOffset, 0)
	return nil
}

func (ui *ui) down(g *gocui.Gui, v *gocui.View) error {
	_, maxY := g.Size()
	ui.YOffset += 1
	ui.YOffset = min(ui.YOffset, ui.lines-maxY+1)
	ui.YOffset = max(ui.YOffset, 0)
	return nil
}

func (ui *ui) left(g *gocui.Gui, v *gocui.View) error {
	ui.XOffset -= 1
	ui.XOffset = max(ui.XOffset, 0)
	return nil
}

func (ui *ui) right(g *gocui.Gui, v *gocui.View) error {
	ui.XOffset += 1
	return nil
}

func (ui *ui) pageUp(g *gocui.Gui, v *gocui.View) error {
	_, maxY := g.Size()
	ui.YOffset -= maxY / 2
	ui.YOffset = max(ui.YOffset, 0)
	return nil
}

func (ui *ui) pageDown(g *gocui.Gui, v *gocui.View) error {
	_, maxY := g.Size()
	ui.YOffset += maxY / 2
	ui.YOffset = min(ui.YOffset, ui.lines-maxY+1)
	ui.YOffset = max(ui.YOffset, 0)
	return nil
}

func (ui *ui) goToTop(g *gocui.Gui, v *gocui.View) error {
	ui.YOffset = 0
	return nil
}

func (ui *ui) goToBottom(g *gocui.Gui, v *gocui.View) error {
	_, maxY := g.Size()
	ui.YOffset = ui.lines - maxY + 1
	ui.YOffset = max(ui.YOffset, 0)
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}