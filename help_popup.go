package main

import (
	"fmt"
	"strings"

	"github.com/awesome-gocui/gocui"
)

const helpPopupView = "helpPopup"

type helpPopup struct {
	active bool
	config *Config
}

func newHelpPopup(config *Config) *helpPopup {
	return &helpPopup{
		active: false,
		config: config,
	}
}

func (hp *helpPopup) keybindings(g *gocui.Gui) error {
	// Keys to close the help popup
	closeKeys := []interface{}{
		gocui.KeyEsc,
		gocui.KeyEnter,
		gocui.KeySpace,
		'q',
		'?',
	}

	for _, key := range closeKeys {
		if err := g.SetKeybinding(helpPopupView, key, gocui.ModNone, hp.close); err != nil {
			return err
		}
	}

	return nil
}

func (hp *helpPopup) layout(g *gocui.Gui) error {
	if !hp.active {
		return nil
	}

	maxX, maxY := g.Size()

	// Calculate popup dimensions
	width := 60
	height := 30

	// Ensure popup fits on screen
	if width > maxX-4 {
		width = maxX - 4
	}
	if height > maxY-4 {
		height = maxY - 4
	}

	// Center the popup
	x0 := (maxX - width) / 2
	y0 := (maxY - height) / 2

	v, err := g.SetView(helpPopupView, x0, y0, x0+width, y0+height, 0)
	if err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}

		v.Frame = true
		v.Title = " Keybindings (Press any key to close) "
		v.Autoscroll = false
		v.Wrap = false
	}

	v.Clear()

	// Build the help content
	helpContent := hp.buildHelpContent()
	fmt.Fprint(v, helpContent)

	// Set this view as current
	if _, err := g.SetCurrentView(helpPopupView); err != nil {
		return err
	}

	return nil
}

func (hp *helpPopup) buildHelpContent() string {
	var sb strings.Builder

	// Helper function to format key list
	formatKeys := func(keys []string) string {
		return strings.Join(keys, ", ")
	}

	// Navigation section
	sb.WriteString(" NAVIGATION\n")
	sb.WriteString(" ═══════════════════════════════════════════════\n")
	sb.WriteString(fmt.Sprintf("  %-20s Move up\n", formatKeys(hp.config.Keybindings.ScrollUp)))
	sb.WriteString(fmt.Sprintf("  %-20s Move down\n", formatKeys(hp.config.Keybindings.ScrollDown)))
	sb.WriteString(fmt.Sprintf("  %-20s Move left\n", formatKeys(hp.config.Keybindings.ScrollLeft)))
	sb.WriteString(fmt.Sprintf("  %-20s Move right\n", formatKeys(hp.config.Keybindings.ScrollRight)))
	sb.WriteString(fmt.Sprintf("  %-20s Page up\n", formatKeys(hp.config.Keybindings.PageUp)))
	sb.WriteString(fmt.Sprintf("  %-20s Page down\n", formatKeys(hp.config.Keybindings.PageDown)))
	sb.WriteString(fmt.Sprintf("  %-20s Go to top\n", formatKeys(hp.config.Keybindings.GoToTop)))
	sb.WriteString(fmt.Sprintf("  %-20s Go to bottom\n", formatKeys(hp.config.Keybindings.GoToBottom)))
	sb.WriteString("\n")

	// Search section
	sb.WriteString(" SEARCH\n")
	sb.WriteString(" ═══════════════════════════════════════════════\n")
	sb.WriteString(fmt.Sprintf("  %-20s Start search\n", formatKeys(hp.config.Keybindings.StartSearch)))
	sb.WriteString(fmt.Sprintf("  %-20s Next match\n", formatKeys(hp.config.Keybindings.NextMatch)))
	sb.WriteString(fmt.Sprintf("  %-20s Previous match\n", formatKeys(hp.config.Keybindings.PrevMatch)))
	sb.WriteString(fmt.Sprintf("  %-20s Clear search\n", formatKeys(hp.config.Keybindings.ClearSearch)))
	sb.WriteString("\n")

	// General section
	sb.WriteString(" GENERAL\n")
	sb.WriteString(" ═══════════════════════════════════════════════\n")
	sb.WriteString(fmt.Sprintf("  %-20s Show this help\n", formatKeys(hp.config.Keybindings.ShowHelp)))
	sb.WriteString(fmt.Sprintf("  %-20s Quit\n", formatKeys(hp.config.Keybindings.Quit)))
	sb.WriteString("\n")

	// Notes section
	sb.WriteString(" NOTES\n")
	sb.WriteString(" ═══════════════════════════════════════════════\n")
	sb.WriteString("  • Vim navigation keys are supported\n")
	sb.WriteString("  • Colemak-DH layout is also supported\n")
	sb.WriteString("    (i=up, e=down, o=right)\n")
	sb.WriteString("  • Search is case-insensitive\n")
	sb.WriteString("  • While searching:\n")
	sb.WriteString("    - Enter to execute search\n")
	sb.WriteString("    - ESC or Ctrl+C to cancel\n")

	return sb.String()
}

func (hp *helpPopup) show() {
	hp.active = true
}

func (hp *helpPopup) close(g *gocui.Gui, v *gocui.View) error {
	hp.active = false
	if err := g.DeleteView(helpPopupView); err != nil {
		return err
	}
	// Return focus to the main render view
	if _, err := g.SetCurrentView(renderView); err != nil {
		return err
	}
	return nil
}

func (hp *helpPopup) isActive() bool {
	return hp.active
}