package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/MichaelMure/go-term-markdown"
)

// Config holds the configuration for mdrs
type Config struct {
	Colors ColorConfig `json:"colors"`
}

// ColorConfig holds color settings for markdown elements
type ColorConfig struct {
	// Headings
	Heading1       string `json:"heading1"`
	Heading2       string `json:"heading2"`
	Heading3       string `json:"heading3"`
	Heading4       string `json:"heading4"`
	Heading5       string `json:"heading5"`
	Heading6       string `json:"heading6"`
	
	// Text elements
	Bold           string `json:"bold"`
	Italic         string `json:"italic"`
	Strikethrough  string `json:"strikethrough"`
	Link           string `json:"link"`
	LinkURL        string `json:"link_url"`
	
	// Code
	Code           string `json:"code"`
	CodeBlock      string `json:"code_block"`
	CodeBlockBg    string `json:"code_block_bg"`
	
	// Lists
	ListMarker     string `json:"list_marker"`
	TaskChecked    string `json:"task_checked"`
	TaskUnchecked  string `json:"task_unchecked"`
	
	// Quotes and tables
	BlockQuote     string `json:"blockquote"`
	TableHeader    string `json:"table_header"`
	TableRow       string `json:"table_row"`
	TableBorder    string `json:"table_border"`
	
	// Search highlighting (for our search feature)
	SearchCurrent  string `json:"search_current"`
	SearchMatch    string `json:"search_match"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Colors: ColorConfig{
			// Headings - blue shades
			Heading1:       "#00d7ff",
			Heading2:       "#00afff",
			Heading3:       "#0087ff",
			Heading4:       "#005fff",
			Heading5:       "#0037ff",
			Heading6:       "#001fff",
			
			// Text elements
			Bold:           "#ffffff",
			Italic:         "#87ff00",
			Strikethrough:  "#808080",
			Link:           "#00ffff",
			LinkURL:        "#0087af",
			
			// Code
			Code:           "#ffff00",
			CodeBlock:      "#d7ff00",
			CodeBlockBg:    "#262626",
			
			// Lists
			ListMarker:     "#ff8700",
			TaskChecked:    "#00ff00",
			TaskUnchecked:  "#ff0000",
			
			// Quotes and tables
			BlockQuote:     "#808080",
			TableHeader:    "#ffff00",
			TableRow:       "#ffffff",
			TableBorder:    "#808080",
			
			// Search
			SearchCurrent:  "#ffff00",
			SearchMatch:    "#ff8700",
		},
	}
}

// LoadConfig loads configuration from the config file
func LoadConfig() (*Config, error) {
	configPath := getConfigPath()
	
	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create default config file
		config := DefaultConfig()
		if err := config.Save(); err != nil {
			return config, nil // Return default config even if save fails
		}
		return config, nil
	}
	
	// Read config file
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return DefaultConfig(), fmt.Errorf("failed to read config file: %w", err)
	}
	
	// Parse JSON config
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return DefaultConfig(), fmt.Errorf("failed to parse config file: %w", err)
	}
	
	// Fill in any missing values with defaults
	config.fillDefaults()
	
	return &config, nil
}

// Save saves the configuration to file
func (c *Config) Save() error {
	configPath := getConfigPath()
	configDir := filepath.Dir(configPath)
	
	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	
	// Marshal config to JSON with indentation
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	
	// Write config file
	if err := ioutil.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	
	return nil
}

// fillDefaults fills in any missing configuration values with defaults
func (c *Config) fillDefaults() {
	defaults := DefaultConfig()
	
	if c.Colors.Heading1 == "" { c.Colors.Heading1 = defaults.Colors.Heading1 }
	if c.Colors.Heading2 == "" { c.Colors.Heading2 = defaults.Colors.Heading2 }
	if c.Colors.Heading3 == "" { c.Colors.Heading3 = defaults.Colors.Heading3 }
	if c.Colors.Heading4 == "" { c.Colors.Heading4 = defaults.Colors.Heading4 }
	if c.Colors.Heading5 == "" { c.Colors.Heading5 = defaults.Colors.Heading5 }
	if c.Colors.Heading6 == "" { c.Colors.Heading6 = defaults.Colors.Heading6 }
	
	if c.Colors.Bold == "" { c.Colors.Bold = defaults.Colors.Bold }
	if c.Colors.Italic == "" { c.Colors.Italic = defaults.Colors.Italic }
	if c.Colors.Strikethrough == "" { c.Colors.Strikethrough = defaults.Colors.Strikethrough }
	if c.Colors.Link == "" { c.Colors.Link = defaults.Colors.Link }
	if c.Colors.LinkURL == "" { c.Colors.LinkURL = defaults.Colors.LinkURL }
	
	if c.Colors.Code == "" { c.Colors.Code = defaults.Colors.Code }
	if c.Colors.CodeBlock == "" { c.Colors.CodeBlock = defaults.Colors.CodeBlock }
	if c.Colors.CodeBlockBg == "" { c.Colors.CodeBlockBg = defaults.Colors.CodeBlockBg }
	
	if c.Colors.ListMarker == "" { c.Colors.ListMarker = defaults.Colors.ListMarker }
	if c.Colors.TaskChecked == "" { c.Colors.TaskChecked = defaults.Colors.TaskChecked }
	if c.Colors.TaskUnchecked == "" { c.Colors.TaskUnchecked = defaults.Colors.TaskUnchecked }
	
	if c.Colors.BlockQuote == "" { c.Colors.BlockQuote = defaults.Colors.BlockQuote }
	if c.Colors.TableHeader == "" { c.Colors.TableHeader = defaults.Colors.TableHeader }
	if c.Colors.TableRow == "" { c.Colors.TableRow = defaults.Colors.TableRow }
	if c.Colors.TableBorder == "" { c.Colors.TableBorder = defaults.Colors.TableBorder }
	
	if c.Colors.SearchCurrent == "" { c.Colors.SearchCurrent = defaults.Colors.SearchCurrent }
	if c.Colors.SearchMatch == "" { c.Colors.SearchMatch = defaults.Colors.SearchMatch }
}

// getConfigPath returns the path to the config file
func getConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory
		return "mdrs-config.json"
	}
	return filepath.Join(homeDir, ".config", "mdrs", "config.json")
}

// hexToANSI converts a hex color to the nearest ANSI 256 color
func hexToANSI(hex string) (int, error) {
	// Remove # prefix if present
	hex = strings.TrimPrefix(hex, "#")
	
	// Parse hex color
	if len(hex) != 6 {
		return 0, fmt.Errorf("invalid hex color: %s", hex)
	}
	
	r, err := strconv.ParseInt(hex[0:2], 16, 64)
	if err != nil {
		return 0, err
	}
	g, err := strconv.ParseInt(hex[2:4], 16, 64)
	if err != nil {
		return 0, err
	}
	b, err := strconv.ParseInt(hex[4:6], 16, 64)
	if err != nil {
		return 0, err
	}
	
	// Convert to ANSI 256 color
	// This is a simplified conversion - more sophisticated algorithms exist
	if r == g && g == b {
		// Grayscale
		if r < 8 {
			return 16, nil
		}
		if r > 248 {
			return 231, nil
		}
		return int(232 + ((r-8)/10)), nil
	}
	
	// Color cube (6x6x6)
	r = (r * 5) / 255
	g = (g * 5) / 255
	b = (b * 5) / 255
	
	return int(16 + (36 * r) + (6 * g) + b), nil
}

// GetANSIColor returns ANSI escape code for a hex color
func (c *ColorConfig) GetANSIColor(hex string) string {
	if hex == "" {
		return ""
	}
	
	colorCode, err := hexToANSI(hex)
	if err != nil {
		// Return default color on error
		return "\033[0m"
	}
	
	return fmt.Sprintf("\033[38;5;%dm", colorCode)
}

// GetANSIBackground returns ANSI escape code for background color
func (c *ColorConfig) GetANSIBackground(hex string) string {
	if hex == "" {
		return ""
	}
	
	colorCode, err := hexToANSI(hex)
	if err != nil {
		return ""
	}
	
	return fmt.Sprintf("\033[48;5;%dm", colorCode)
}

// GetMarkdownOptions returns markdown rendering options based on config
func (c *Config) GetMarkdownOptions() []markdown.Options {
	opts := []markdown.Options{
		markdown.WithImageDithering(markdown.DitheringWithBlocks),
	}
	
	// Note: The go-term-markdown library may not support all color customizations
	// We'll need to work with what's available in the API
	// This is a placeholder for when we integrate with the actual markdown renderer
	
	return opts
}

// ApplySearchHighlight applies search highlighting colors to text
func (c *Config) ApplySearchHighlight(text string, isCurrent bool) string {
	if isCurrent {
		// Current match - use background color
		bgColor := c.Colors.GetANSIBackground(c.Colors.SearchCurrent)
		return fmt.Sprintf("%s\033[30m%s\033[0m", bgColor, text)
	} else {
		// Other matches - use foreground color
		fgColor := c.Colors.GetANSIColor(c.Colors.SearchMatch)
		return fmt.Sprintf("%s%s\033[0m", fgColor, text)
	}
}