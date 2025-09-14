package main

import (
	"strings"
	
	"github.com/awesome-gocui/gocui"
)

type keybinding struct {
	viewName string
	key      interface{}
	mod      gocui.Modifier
	handler  func(*gocui.Gui, *gocui.View) error
}

func (k *keybinding) Register(g *gocui.Gui) error {
	return g.SetKeybinding(k.viewName, k.key, k.mod, k.handler)
}

// parseKey converts a string key representation to gocui key and modifier
func parseKey(keyStr string) (interface{}, gocui.Modifier) {
	// Handle control key combinations
	if strings.HasPrefix(keyStr, "C-") || strings.HasPrefix(keyStr, "Ctrl-") {
		keyStr = strings.TrimPrefix(keyStr, "C-")
		keyStr = strings.TrimPrefix(keyStr, "Ctrl-")
		
		// Handle special control combinations
		switch strings.ToLower(keyStr) {
		case "c":
			return gocui.KeyCtrlC, gocui.ModNone
		case "f":
			return gocui.KeyCtrlF, gocui.ModNone
		case "n":
			return gocui.KeyCtrlN, gocui.ModNone
		case "p":
			return gocui.KeyCtrlP, gocui.ModNone
		case "g":
			return gocui.KeyCtrlG, gocui.ModNone
		}
	}
	
	// Handle special keys
	switch keyStr {
	case "Up", "ArrowUp":
		return gocui.KeyArrowUp, gocui.ModNone
	case "Down", "ArrowDown":
		return gocui.KeyArrowDown, gocui.ModNone
	case "Left", "ArrowLeft":
		return gocui.KeyArrowLeft, gocui.ModNone
	case "Right", "ArrowRight":
		return gocui.KeyArrowRight, gocui.ModNone
	case "PageUp", "PgUp":
		return gocui.KeyPgup, gocui.ModNone
	case "PageDown", "PgDn", "PageDn":
		return gocui.KeyPgdn, gocui.ModNone
	case "Space", " ":
		return gocui.KeySpace, gocui.ModNone
	case "Enter", "Return":
		return gocui.KeyEnter, gocui.ModNone
	case "Escape", "Esc":
		return gocui.KeyEsc, gocui.ModNone
	case "Tab":
		return gocui.KeyTab, gocui.ModNone
	case "Backspace":
		return gocui.KeyBackspace, gocui.ModNone
	case "Delete", "Del":
		return gocui.KeyDelete, gocui.ModNone
	case "Home":
		return gocui.KeyHome, gocui.ModNone
	case "End":
		return gocui.KeyEnd, gocui.ModNone
	}
	
	// Handle single character keys
	if len(keyStr) == 1 {
		return rune(keyStr[0]), gocui.ModNone
	}
	
	// Default fallback
	return nil, gocui.ModNone
}

// createKeybindingFromStrings creates keybindings from string representations
func createKeybindingsFromStrings(viewName string, keys []string, handler func(*gocui.Gui, *gocui.View) error) []keybinding {
	var bindings []keybinding
	
	for _, keyStr := range keys {
		key, mod := parseKey(keyStr)
		if key != nil {
			bindings = append(bindings, keybinding{
				viewName: viewName,
				key:      key,
				mod:      mod,
				handler:  handler,
			})
		}
	}
	
	return bindings
}
