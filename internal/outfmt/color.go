package outfmt

import (
	"os"

	"golang.org/x/term"
)

// ColorMode represents color output setting
type ColorMode string

const (
	ColorAuto   ColorMode = "auto"
	ColorAlways ColorMode = "always"
	ColorNever  ColorMode = "never"
)

// ShouldColorize returns true if output should be colorized
func ShouldColorize(mode string) bool {
	switch ColorMode(mode) {
	case ColorAlways:
		return true
	case ColorNever:
		return false
	default: // auto
		return term.IsTerminal(int(os.Stdout.Fd()))
	}
}

// ANSI color codes
const (
	Reset   = "\033[0m"
	Bold    = "\033[1m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	Gray    = "\033[90m"
)

// Colorize wraps text in ANSI color codes if colors are enabled
func Colorize(text, color string, enabled bool) string {
	if !enabled {
		return text
	}
	return color + text + Reset
}
