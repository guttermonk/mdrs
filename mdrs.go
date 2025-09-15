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

	if len(os.Args) >= 2 && (os.Args[1] == "--init-config") {
		initConfig()
		return
	}

	if len(os.Args) >= 2 && (os.Args[1] == "--config-path") {
		fmt.Printf("Config file location: %s\n", getConfigPath())
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

func initConfig() {
	config := DefaultConfig()
	configPath := getConfigPath()
	
	// Check if config already exists
	if _, err := os.Stat(configPath); err == nil {
		fmt.Printf("Config file already exists at: %s\n", configPath)
		fmt.Println("To regenerate, please delete the existing file first.")
		return
	}
	
	// Save the default config
	if err := config.Save(); err != nil {
		exitError(fmt.Errorf("failed to create config file: %w", err))
	}
	
	fmt.Printf("Created default config file at: %s\n", configPath)
	fmt.Println("You can now edit this file to customize colors.")
	fmt.Println("\nExample color values:")
	fmt.Println("  \"#ff0000\" - Red")
	fmt.Println("  \"#00ff00\" - Green")
	fmt.Println("  \"#0000ff\" - Blue")
	fmt.Println("  \"#ffff00\" - Yellow")
	fmt.Println("  \"#ff00ff\" - Magenta")
	fmt.Println("  \"#00ffff\" - Cyan")
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
	
	// configuration
	config          *Config
	
	// help popup
	help            *helpPopup
}

func newUi(g *gocui.Gui) (*ui, error) {
	config, err := LoadConfig()
	if err != nil {
		// Use default config if loading fails
		config = DefaultConfig()
	}
	
	result := &ui{
		width:  -1,
		search: NewSearchState(config),
		config: config,
		help:   newHelpPopup(config),
	}

	g.SetManagerFunc(result.layout)
	g.Cursor = false
	g.InputEsc = true

	// Build keybindings from config
	result.keybindings = []keybinding{}
	
	// Navigation keybindings
	result.keybindings = append(result.keybindings, createKeybindingsFromStrings(renderView, result.config.Keybindings.ScrollUp, result.up)...)
	result.keybindings = append(result.keybindings, createKeybindingsFromStrings(renderView, result.config.Keybindings.ScrollDown, result.down)...)
	result.keybindings = append(result.keybindings, createKeybindingsFromStrings(renderView, result.config.Keybindings.ScrollLeft, result.left)...)
	result.keybindings = append(result.keybindings, createKeybindingsFromStrings(renderView, result.config.Keybindings.ScrollRight, result.right)...)
	result.keybindings = append(result.keybindings, createKeybindingsFromStrings(renderView, result.config.Keybindings.PageUp, result.pageUp)...)
	result.keybindings = append(result.keybindings, createKeybindingsFromStrings(renderView, result.config.Keybindings.PageDown, result.pageDown)...)
	result.keybindings = append(result.keybindings, createKeybindingsFromStrings(renderView, result.config.Keybindings.GoToTop, result.goToTop)...)
	result.keybindings = append(result.keybindings, createKeybindingsFromStrings(renderView, result.config.Keybindings.GoToBottom, result.goToBottom)...)
	
	// Search keybindings
	result.keybindings = append(result.keybindings, createKeybindingsFromStrings(renderView, result.config.Keybindings.StartSearch, result.startSearch)...)
	result.keybindings = append(result.keybindings, createKeybindingsFromStrings(renderView, result.config.Keybindings.NextMatch, result.nextMatch)...)
	result.keybindings = append(result.keybindings, createKeybindingsFromStrings(renderView, result.config.Keybindings.PrevMatch, result.prevMatch)...)
	result.keybindings = append(result.keybindings, createKeybindingsFromStrings(renderView, result.config.Keybindings.ClearSearch, result.clearSearch)...)
	
	// General keybindings
	result.keybindings = append(result.keybindings, createKeybindingsFromStrings(renderView, result.config.Keybindings.Quit, result.quit)...)
	result.keybindings = append(result.keybindings, createKeybindingsFromStrings("", result.config.Keybindings.Quit, result.quit)...)  // Global quit
	result.keybindings = append(result.keybindings, createKeybindingsFromStrings(renderView, result.config.Keybindings.ShowHelp, result.showHelp)...)
	
	// Search view specific keybindings (always fixed)
	result.keybindings = append(result.keybindings, keybinding{searchView, gocui.KeyEnter, gocui.ModNone, result.executeSearch})
	result.keybindings = append(result.keybindings, keybinding{searchView, gocui.KeyEsc, gocui.ModNone, result.cancelSearch})
	result.keybindings = append(result.keybindings, keybinding{searchView, gocui.KeyCtrlC, gocui.ModNone, result.cancelSearch})
	result.keybindings = append(result.keybindings, keybinding{searchView, gocui.KeyCtrlG, gocui.ModNone, result.cancelSearch})
	
	// Register help popup keybindings
	if err := result.help.keybindings(g); err != nil {
		return nil, err
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
	
	// Handle help popup first (it should be on top)
	if err := ui.help.layout(g); err != nil {
		return err
	}
	
	// If help is active, don't update other views
	if ui.help.isActive() {
		return nil
	}

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
		g.Cursor = true
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
		g.Cursor = false
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

	// Get options from config, plus required options
	opts := ui.config.GetMarkdownOptions()

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

func (ui *ui) clearSearch(g *gocui.Gui, v *gocui.View) error {
	ui.searchActive = false
	ui.search.Clear()
	
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

func (ui *ui) showHelp(g *gocui.Gui, v *gocui.View) error {
	ui.help.show()
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