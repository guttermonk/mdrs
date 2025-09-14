package main

import (
	"fmt"
	"strings"

	"github.com/awesome-gocui/gocui"
)

const helpPopupView = "helpPopup"

type helpPopup struct {
	active bool
}

func newHelpPopup() *helpPopup {
	return &helpPopup{
		active: false,
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

	// Navigation section
	sb.WriteString(" NAVIGATION\n")
	sb.WriteString(" ═══════════════════════════════════════════════\n")
	sb.WriteString("  k, i, ↑, Ctrl+P    Move up\n")
	sb.WriteString("  j, e, ↓, Ctrl+N    Move down\n")
	sb.WriteString("  h, ←               Move left\n")
	sb.WriteString("  l, o, →            Move right\n")
	sb.WriteString("  PgUp               Page up\n")
	sb.WriteString("  PgDn, Space        Page down\n")
	sb.WriteString("  g                  Go to top\n")
	sb.WriteString("  G                  Go to bottom\n")
	sb.WriteString("\n")

	// Search section
	sb.WriteString(" SEARCH\n")
	sb.WriteString(" ═══════════════════════════════════════════════\n")
	sb.WriteString("  /, Ctrl+F          Start search\n")
	sb.WriteString("  n                  Next match\n")
	sb.WriteString("  N                  Previous match\n")
	sb.WriteString("  ESC                Clear search\n")
	sb.WriteString("\n")

	// General section
	sb.WriteString(" GENERAL\n")
	sb.WriteString(" ═══════════════════════════════════════════════\n")
	sb.WriteString("  ?                  Show this help\n")
	sb.WriteString("  q, Ctrl+C          Quit\n")
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